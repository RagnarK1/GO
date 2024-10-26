package forum_func

import (
	"database/sql"
	"html"
	"net/http"
	"strconv"
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
		signin.LoginType = r.FormValue("logintype")
		signin.User = html.EscapeString(r.FormValue("user"))

		if signin.Email == "" || signin.LoginType == "" || signin.User == "" {
			signin.Email = r.FormValue("hidden_email")
			signin.LoginType = r.FormValue("hidden_loginType")
			signin.User = r.FormValue("hidden_user")
		}

		signin.Password = cryptPW(r.FormValue("password"))

		status := tempAddUserToDB(signin.User, signin.Email, signin.Password, signin.LoginType, w)
		if status == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			t.Execute(w, "User name already exists! Try another one!")
			return
		} else if status == 2 {
			w.WriteHeader(http.StatusUnauthorized)
			t.Execute(w, "User email already exists!")
			return
		}

		// sqlScript := "select user_name from users where upper(user_name) = upper('" + signin.User + "')"
		// value := SelectSQLScript(DatabaseName, sqlScript)
		// var user_name sql.NullString
		// for value.Next() {
		// 	value.Scan(&user_name)
		// }
		// if user_name.String != "" {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	t.Execute(w, "User name already exists! Try another one!")
		// 	return
		// }
		// sqlInsertScript := "insert into users (user_email, user_method, user_name, user_pw) values('" + signin.Email + "','" + signin.LoginType + "','" + signin.User + "','" + signin.Password + "')"
		// _, err := ExecSQLScript(DatabaseName, sqlInsertScript)
		// if err != nil {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	t.Execute(w, "User email already exists!")
		// 	return
		// }
		// sqlLoginTypeScript := "insert into login_method (user_email, method) values('" + signin.Email + "','" + signin.LoginType + "')"
		// ExecSQLScript(DatabaseName, sqlLoginTypeScript)

		statusNumber := LogIn(signin.User, signin.Email, signin.Password, signin.LoginType, w)
		if statusNumber == 1 {
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

func tempAddUserToDB(userName, userEmail, userPassword, userLoginType string, w http.ResponseWriter) int {
	// sqlScript := "select user_name from users where upper(user_name) = upper('" + userName + "')"
	// value := SelectSQLScript(DatabaseName, sqlScript)
	// var user_name sql.NullString
	// for value.Next() {
	// 	value.Scan(&user_name)
	// }
	user_name := checkName(userName)
	if user_name.String != "" {
		if userLoginType == "google" || userLoginType == "github" {
			aNumber := 1
			for {
				user_name = checkName(userName + strconv.Itoa(aNumber))
				if user_name.String == "" {
					userName = userName + strconv.Itoa(aNumber)
					break
				}
				aNumber++
			}
		} else {
			return 1
		}

	}
	sqlInsertScript := "insert into users (user_email, user_method, user_name, user_pw) values('" + userEmail + "','" + userLoginType + "','" + userName + "','" + userPassword + "')"
	_, err := ExecSQLScript(DatabaseName, sqlInsertScript)
	if err != nil {
		return 2
	}
	sqlLoginTypeScript := "insert into login_method (user_email, method) values('" + userEmail + "','" + userLoginType + "')"
	ExecSQLScript(DatabaseName, sqlLoginTypeScript)
	return 0
}

func checkName(userName string) sql.NullString {
	sqlScript := "select user_name from users where upper(user_name) = upper('" + userName + "')"
	value := SelectSQLScript(DatabaseName, sqlScript)
	var user_name sql.NullString
	for value.Next() {
		value.Scan(&user_name)
	}
	return user_name
}
