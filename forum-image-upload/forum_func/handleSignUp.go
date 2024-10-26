package forum_func

import (
	"database/sql"
	"html"
	"net/http"
	"text/template"
)

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	var errorMes ErrorMessage
	signin := &UserCredentials{}
	t, err := template.ParseFiles("common/signup.htm")
	if err != nil {
		errorMes.Key = 500
		errorMes.Text = "Server not reachable - 500 Internal Server Error"
		HandleError(w, r, errorMes)
		return
	}
	switch r.Method { //get request method
	case "GET": // GET = Web server got requested to show first page
		errorMes = testMethod(r, "GET", "/signup")
	case "POST": // POST = web server got request to generate something additional from user filled form
		errorMes = testMethod(r, "POST", "/signup")
		signin.Email = r.FormValue("email")
		signin.User = html.EscapeString(r.FormValue("user"))
		signin.Password = cryptPW(r.FormValue("password"))
		signin.LoginType = r.FormValue("logintype")
		sqlScript := "select user_name from users where upper(user_name) = upper('" + signin.User + "')"
		value := SelectSQLScript(DatabaseName, sqlScript)
		var user_name sql.NullString
		for value.Next() {
			value.Scan(&user_name)
		}
		if user_name.String != "" {
			w.WriteHeader(http.StatusUnauthorized)
			t.Execute(w, "User name already exists! Try another one!")
			return
		}
		sqlInsertScript := "insert into users (user_email, user_name, user_pw) values('" + signin.Email + "','" + signin.User + "','" + signin.Password + "')"
		_, err := ExecSQLScript(DatabaseName, sqlInsertScript)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			t.Execute(w, "User email already exists!")
			return
		}
		sqlLoginTypeScript := "insert into login_method (user_email, method) values('" + signin.Email + "','" + signin.LoginType + "')"
		ExecSQLScript(DatabaseName, sqlLoginTypeScript)

		if LogIn(signin.Email, signin.Password, signin.LoginType, w) == 1 {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			t.Execute(w, "Incorrect email or password!")
			return
		}
	}
	if errorMes.Key == 200 {
		t.Execute(w, "")
	} else {
		HandleError(w, r, errorMes)
	}
}
