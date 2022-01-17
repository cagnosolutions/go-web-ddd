package webapp

import (
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const sessionIDLen = 32

type session struct {
	sid  string
	ts   time.Time
	data url.Values
}

func (s *session) ID() string {
	return s.sid
}

func (s *session) Has(key string) bool {
	return s.data.Has(key)
}

func (s *session) Set(key string, val string) {
	s.data.Set(key, val)
}

func (s *session) Get(key string) string {
	return s.data.Get(key)
}

func (s *session) Del(key string) {
	s.data.Del(key)
}

type sessionStore struct {
	sessCookID    string
	rateInSeconds int64
	sessions      *sync.Map
}

func NewSessionStore(sessCookID string, rateInSeconds int64) *sessionStore {
	ss := &sessionStore{
		sessCookID:    sessCookID,
		rateInSeconds: rateInSeconds,
		sessions:      new(sync.Map),
	}
	go ss.gc()
	return ss
}

func (ss *sessionStore) NewSession(w http.ResponseWriter, r *http.Request) {
	sid := randomN(sessionIDLen)
	ss.sessions.Store(sid, &session{
		sid:  sid,
		ts:   time.Now(),
		data: url.Values{},
	})
	c := NewCookie(w, ss.sessCookID, sid,
		time.Now().Add(time.Duration(ss.rateInSeconds)*time.Second),
		int(ss.rateInSeconds))
	http.SetCookie(w, c)
}

func (ss *sessionStore) GetSession(w http.ResponseWriter, r *http.Request) *session {
	c := GetCookie(r, ss.sessCookID)
	if c == nil {
		ss.NewSession(w, r)
		c = GetCookie(r, ss.sessCookID)
	}
	v, _ := ss.sessions.Load(c.Value)
	return v.(*session)
}

func (ss *sessionStore) UpdateSession(w http.ResponseWriter, r *http.Request) {
	c := GetCookie(r, ss.sessCookID)
	if c == nil {
		return
	}
	v, ok := ss.sessions.Load(c.Value)
	if ok {
		currentTime := time.Now()
		v.(*session).ts = currentTime
		c = NewCookie(w, ss.sessCookID, c.Value,
			currentTime.Add(time.Duration(ss.rateInSeconds)*time.Second),
			int(ss.rateInSeconds))
		http.SetCookie(w, c)
	}
}

func (ss *sessionStore) EndSession(w http.ResponseWriter, r *http.Request) {
	c := GetCookie(r, ss.sessCookID)
	if c == nil {
		return
	}
	ss.sessions.Delete(c.Value)
	c = NewCookie(w, ss.sessCookID, c.Value, time.Now(), -1)
	http.SetCookie(w, c)
}

func (ss *sessionStore) gc() {
	ss.sessions.Range(func(sid, sess interface{}) bool {
		if (sess.(*session).ts.Unix() + ss.rateInSeconds) < time.Now().Unix() {
			ss.sessions.Delete(sid)
		}
		return true
	})
	time.AfterFunc(time.Duration(ss.rateInSeconds)*time.Second, func() {
		ss.gc()
	})
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randomN(n int) string {
	var src = rand.NewSource(time.Now().UnixNano() + int64(rand.Uint64()))
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
