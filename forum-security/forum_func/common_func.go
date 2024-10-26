package forum_func

import (
	"database/sql"
	"encoding/hex"
	"hash/fnv"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// const DatabaseName = "file:sqlite-forum.db?_auth&_auth_user=admin&_auth_pass=admin&_auth_crypt=sha1"
// const DbName = "sqlite-forum.db"
const DatabaseName = "file:sqlite-forum.db?key=admin"
const DbName = "sqlite-forum.db"

type ErrorMessage struct { //Error message definiton for error.htm
	Key  int
	Text string
}
type UserCredentials struct { //User credentials to signin or signup into forum
	Email     string
	User      string
	Password  string
	LoginType string
}

type UserSession struct { //user session user and end time
	Username string
	Email    string
	UserType string
	Expiry   time.Time
}

var UserSessions = map[string]UserSession{} //Collection of all user sessions

type TopicVarable struct {
	Session    UserSession
	Code       []ForumTopic
	Topic      Topic
	Categories []string
	Users      []string
}

type SubTopicVarable struct {
	Session UserSession
	Code    []ForumTopicContent
	Topic   Topic
}

type UserType struct {
	UserName string
	UserType string
}

func activeSessionToken(r *http.Request) (string, string) {
	sessionToken := ""
	user := ""
	c, err := r.Cookie("session_token")
	if err == nil {
		sessionToken = c.Value
		user = UserSessions[sessionToken].Email
	}
	return sessionToken, user
}

func RandStringToken(n int) string {
	b := make([]byte, n)
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func cryptPW(s string) string {
	h := fnv.New64()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func testMethod(r *http.Request, method, link string) ErrorMessage {
	var errorMes ErrorMessage
	errorMes.Key = 200
	if r.URL.Path != link {
		errorMes.Key = http.StatusNotFound
		errorMes.Text = "Bad request - Page Not Found"
	}
	if r.Method != method {
		errorMes.Key = http.StatusMethodNotAllowed
		errorMes.Text = "Bad request - Method Not Allowed"
	}
	return errorMes
}

func LogIn(user, pw, loginType string, w http.ResponseWriter) int {
	for key, data := range UserSessions {
		if data.Email == user {
			delete(UserSessions, key)
			break
		}
	}
	sqlScript := "select user_name, user_type from users where user_email ='" + user + "' and user_pw = '" + pw + "'"
	value := SelectSQLScript(DatabaseName, sqlScript)
	var user_name, user_type sql.NullString
	for value.Next() {
		value.Scan(&user_name, &user_type)
	}
	if !user_name.Valid {
		return 0
	}
	sqlUpdateScript := "update users set login_time_prev=login_time, login_time=CURRENT_TIMESTAMP where user_email = '" + user + "'"
	ExecSQLScript(DatabaseName, sqlUpdateScript)

	sessionToken := (uuid.Must(uuid.NewRandom())).String()
	// sessionToken := RandStringToken(20)
	expiresAt := time.Now().Add(1200 * time.Second)
	UserSessions[sessionToken] = UserSession{
		Username: user_name.String,
		Email:    user,
		UserType: user_type.String,
		Expiry:   expiresAt,
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	return 1
}

func FindUserRights(DatabaseName, user string) []UserType {
	var Types []UserType
	sqlScript := "select user_name, user_type from users where user_email !='" + user + "' and status !='DELETED'"
	value := SelectSQLScript(DatabaseName, sqlScript)
	var user_name, user_type sql.NullString
	for value.Next() {
		value.Scan(&user_name, &user_type)
	}

	mainTypes := UserType{
		UserName: user_name.String,
		UserType: user_type.String,
	}
	Types = append(Types, mainTypes)
	return Types
}
