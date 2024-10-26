package forum_func

import (
	"net/http"
	"time"
)

func (s UserSession) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

func HandleSignOut(w http.ResponseWriter, r *http.Request) {
	sessionToken := ""
	c, err := r.Cookie("session_token")
	if err == nil {
		sessionToken = c.Value
	}
	c.Expires = time.Now()
	delete(UserSessions, sessionToken)
	http.Redirect(w, r, "/", http.StatusFound)
}
