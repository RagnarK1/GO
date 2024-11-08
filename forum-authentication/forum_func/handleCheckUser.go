package forum_func

import (
	"database/sql"
	"fmt"
	"html"
	"html/template"
	"net/http"
)

func HandleCheckUser(w http.ResponseWriter, r *http.Request) {
	User := html.EscapeString(r.FormValue("user"))
	sqlScript := "select user_name from users where upper(user_name) = upper('" + User + "')"
	value := SelectSQLScript(DatabaseName, sqlScript)
	var user_name sql.NullString
	for value.Next() {
		value.Scan(&user_name)
	}
	fmt.Fprint(w, template.HTML(user_name.String))
}

func HandleCheckEmail(w http.ResponseWriter, r *http.Request) {
	Email := html.EscapeString(r.FormValue("email"))
	sqlScript := "SELECT user_email FROM users WHERE UPPER(user_email) = UPPER('" + Email + "')"
	value := SelectSQLScript(DatabaseName, sqlScript)
	var user_email sql.NullString
	for value.Next() {
		value.Scan(&user_email)
	}
	fmt.Fprint(w, template.HTML(user_email.String))
}
