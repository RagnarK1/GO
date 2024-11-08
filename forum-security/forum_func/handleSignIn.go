package forum_func

import (
	"net/http"
	"text/template"
)

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	var errorMes ErrorMessage
	signin := &UserCredentials{}
	// errorMes.Key = 200
	t, err := template.ParseFiles("common/signin.htm")
	if err != nil {
		errorMes.Key = 500
		errorMes.Text = "Server not reachable - 500 Internal Server Error"
		HandleError(w, r, errorMes)
		return
	}
	switch r.Method { //get request method
	case "GET": // GET = Web server got requested to show first page
		errorMes = testMethod(r, "GET", "/signin")
	case "POST": // POST = web server got request to generate something additional from user filled form
		errorMes = testMethod(r, "POST", "/signin")
		signin.Email = r.FormValue("email")
		signin.Password = cryptPW(r.FormValue("password"))
		signin.LoginType = r.FormValue("logintype")
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
