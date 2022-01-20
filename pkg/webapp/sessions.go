package webapp

import (
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var ErrNilSession = errors.New("session is nil, or not found")

type AuthUser interface {
	Register(username, password, role string)
	Authenticate(username, password string) (*SystemUser, bool)
}

type SystemUser struct {
	Username string
	Password string
	Role     string
}

type BasicAuthUser struct {
	users *sync.Map
}

func NewBasicAuthUser() *BasicAuthUser {
	return &BasicAuthUser{
		users: new(sync.Map),
	}
}

func (a *BasicAuthUser) Register(username, password, role string) {
	a.users.Store(username, &SystemUser{
		Username: username,
		Password: password,
		Role:     role,
	})
}

func (a *BasicAuthUser) Authenticate(username, password string) (*SystemUser, bool) {
	su, ok := a.users.Load(username)
	if !ok {
		return nil, false
	}
	if su.(*SystemUser).Password != password {
		return nil, false
	}
	return su.(*SystemUser), true
}

const sessionIDLen = 32

type Session struct {
	id      string
	data    map[string]interface{}
	expires time.Time
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Has(k string) bool {
	_, ok := s.data[k]
	return ok
}

func (s *Session) Set(k string, val interface{}) {
	s.data[k] = val
}

func (s *Session) Get(k string) (interface{}, bool) {
	v, ok := s.data[k]
	return v, ok
}

func (s *Session) Del(k string) {
	delete(s.data, k)
}

func (s *Session) ExpiresIn() int64 {
	return s.expires.Unix() - time.Now().Unix()
}

type SessionStorer interface {
	// New should create and return a new session
	New() *Session

	// Get should return a cached session
	Get(r *http.Request) (*Session, bool)

	// Save should persist session to the underlying store
	// implementation. Passing a nil session erases it.
	Save(w http.ResponseWriter, r *http.Request, s *Session)
}

type SessionStore struct {
	sid      string // sid is the store id
	rate     int64  // rate is the max idle session time in seconds
	sessions *sync.Map
}

func NewSessionStore(sid string, rate int64) *SessionStore {
	ss := &SessionStore{
		sid:      sid,
		rate:     rate,
		sessions: new(sync.Map),
	}
	go ss.gc()
	return ss
}

// New creates and returns a new session
func (ss *SessionStore) New() *Session {
	return &Session{
		id:      randomN(sessionIDLen),
		data:    make(map[string]interface{}),
		expires: time.Now().Add(time.Duration(ss.rate) * time.Second),
	}
}

// Get returns a cached session (if one exists)
func (ss *SessionStore) Get(r *http.Request) (*Session, bool) {
	c := GetCookie(r, ss.sid)
	if c == nil {
		return nil, false
	}
	v, ok := ss.sessions.Load(c.Value)
	if !ok {
		return nil, false
	}
	return v.(*Session), true
}

// Save persists the provided session. If you would like to remove a session, simply
// pass it a nil session, and it will time the cookie out.
func (ss *SessionStore) Save(w http.ResponseWriter, r *http.Request, s *Session) {
	if s == nil {
		c := GetCookie(r, ss.sid)
		if c == nil {
			return
		}
		ss.sessions.Delete(c.Value)
		http.SetCookie(w, NewCookie(ss.sid, c.Value, time.Now(), -1))
		return
	}
	s.expires = time.Now().Add(time.Duration(ss.rate) * time.Second)
	ss.sessions.Store(s.id, s)
	http.SetCookie(w, NewCookie(ss.sid, s.id, s.expires, int(ss.rate)))
}

func (ss *SessionStore) String() string {
	var sessions []string
	ss.sessions.Range(func(id, sess interface{}) bool {
		sessions = append(sessions, id.(string))
		return true
	})
	return strings.Join(sessions, "\n")
}

func (ss *SessionStore) gc() {
	ss.sessions.Range(func(id, sess interface{}) bool {
		if sess.(*Session).ExpiresIn() < 0 {
			ss.sessions.Delete(id)
		}
		return true
	})
	time.AfterFunc(time.Duration(ss.rate/2)*time.Second, func() { ss.gc() })
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
