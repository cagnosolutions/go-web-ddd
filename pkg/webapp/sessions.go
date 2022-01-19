package webapp

import (
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const sessionIDLen = 32

type Session struct {
	sid  string
	ts   time.Time
	data url.Values
}

func (s *Session) ID() string {
	return s.sid
}

func (s *Session) Has(key string) bool {
	return s.data.Has(key)
}

func (s *Session) Set(key string, val string) {
	s.data.Set(key, val)
}

func (s *Session) Get(key string) string {
	return s.data.Get(key)
}

func (s *Session) Del(key string) {
	s.data.Del(key)
}

func (ss *SessionStore) Secure(role string, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if u, ok := ss.CurrentUser(r); !ok || !u.Has(role) {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type SessionStore struct {
	sessCookID    string
	rateInSeconds int64
	sessions      *sync.Map
}

func NewSessionStore(sessCookID string, rateInSeconds int64) *SessionStore {
	ss := &SessionStore{
		sessCookID:    sessCookID,
		rateInSeconds: rateInSeconds,
		sessions:      new(sync.Map),
	}
	go ss.gc()
	return ss
}

func (ss *SessionStore) NewSession(w http.ResponseWriter, r *http.Request) {
	sid := randomN(sessionIDLen)
	ss.sessions.Store(sid, &Session{
		sid:  sid,
		ts:   time.Now(),
		data: url.Values{},
	})
	c := NewCookie(w, ss.sessCookID, sid,
		time.Now().Add(time.Duration(ss.rateInSeconds)*time.Second),
		int(ss.rateInSeconds))
	http.SetCookie(w, c)
}

func (ss *SessionStore) GetSession(w http.ResponseWriter, r *http.Request) *Session {
	c := GetCookie(r, ss.sessCookID)
	if c == nil {
		ss.NewSession(w, r)
		c = GetCookie(r, ss.sessCookID)
	}
	v, _ := ss.sessions.Load(c.Value)
	return v.(*Session)
}

func (ss *SessionStore) UpdateSession(w http.ResponseWriter, r *http.Request) {
	c := GetCookie(r, ss.sessCookID)
	if c == nil {
		return
	}
	v, ok := ss.sessions.Load(c.Value)
	if ok {
		currentTime := time.Now()
		v.(*Session).ts = currentTime
		c = NewCookie(w, ss.sessCookID, c.Value,
			currentTime.Add(time.Duration(ss.rateInSeconds)*time.Second),
			int(ss.rateInSeconds))
		http.SetCookie(w, c)
	}
}

func (ss *SessionStore) EndSession(w http.ResponseWriter, r *http.Request) {
	c := GetCookie(r, ss.sessCookID)
	if c == nil {
		return
	}
	ss.sessions.Delete(c.Value)
	c = NewCookie(w, ss.sessCookID, c.Value, time.Now(), -1)
	http.SetCookie(w, c)
}

func (ss *SessionStore) CurrentUser(r *http.Request) (*Session, bool) {
	c := GetCookie(r, ss.sessCookID)
	if c == nil {
		return nil, false
	}
	v, _ := ss.sessions.Load(c.Value)
	return v.(*Session), true
}

func (ss *SessionStore) gc() {
	ss.sessions.Range(func(sid, sess interface{}) bool {
		if (sess.(*Session).ts.Unix() + ss.rateInSeconds) < time.Now().Unix() {
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
