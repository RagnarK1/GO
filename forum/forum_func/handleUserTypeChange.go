package forum_func

import (
	"net/http"
)

func HandleUserTypeChange(w http.ResponseWriter, r *http.Request) {
	var errorMes ErrorMessage
	switch r.Method { //get request method
	case "GET": // GET = Web server got requested to show first page
		errorMes = testMethod(r, "GET", "/usertype")
	case "POST": // POST = web server got request to generate something additional from user filled form
		errorMes = testMethod(r, "POST", "/usertype")
		user := r.FormValue("user")
		status := r.FormValue("type")
		sqlScript := "update users set user_type='" + status + "' where user_name = '" + user + "'"
		ExecSQLScript(DatabaseName, sqlScript)
	}
	if errorMes.Key == 200 {
		http.Redirect(w, r, "/myHome", http.StatusFound)
	} else {
		HandleError(w, r, errorMes)
	}
}
