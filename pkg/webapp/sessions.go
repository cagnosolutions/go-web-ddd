package webapp

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

var ErrNilSession = errors.New("session is nil, or not found")

const sessionIDLen = 32

type SessionManager interface {
	// New should create and return a new session
	New() *Session

	// Get should return a cached session
	Get(r *http.Request) (*Session, bool)

	// Save should persist session to the underlying store
	// implementation. Passing a nil session erases it.
	Save(w http.ResponseWriter, r *http.Request, s *Session)
}

func AddTime(t time.Time, duration time.Duration) time.Time {
	return t.Add(duration)
}

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

// SessionStore implements the session manager interface
// and is a basic session manager using cookies.
type SessionStore struct {
	sid      string        // sid is the store id
	timeout  time.Duration // expires is the max idle session time allowed
	sessions *sync.Map
}

// NewSessionStore takes a session id and a make session timeout. The sid
// will be used as the key for all session cookies, and the timeout is the
// maximum allowable idle session time before the session is expired
func NewSessionStore(sid string, timeout time.Duration) *SessionStore {
	ss := &SessionStore{
		sid:      sid,
		timeout:  timeout,
		sessions: new(sync.Map),
	}
	go ss.gc()
	return ss
}

// New creates and returns a new session
func (ss *SessionStore) New() *Session {
	return &Session{
		id:      RandStringN(sessionIDLen),
		data:    make(map[string]interface{}),
		expires: AddTime(time.Now(), ss.timeout),
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
func (ss *SessionStore) Save(w http.ResponseWriter, r *http.Request, session *Session) {
	if session == nil {
		cook := GetCookie(r, ss.sid)
		if cook == nil {
			return
		}
		ss.sessions.Delete(cook.Value)
		http.SetCookie(w, NewCookie(ss.sid, cook.Value, time.Now(), -1))
		return
	}
	session.expires = AddTime(time.Now(), ss.timeout)
	ss.sessions.Store(session.id, session)
	http.SetCookie(w, NewCookie(ss.sid, session.id, session.expires, int(ss.timeout)))
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
	time.AfterFunc(ss.timeout/2, func() { ss.gc() })
}
