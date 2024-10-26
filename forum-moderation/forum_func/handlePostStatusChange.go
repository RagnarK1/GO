package forum_func

import (
	"net/http"
)

func HandlePostStatusChange(w http.ResponseWriter, r *http.Request) {
	var errorMes ErrorMessage
	switch r.Method { //get request method
	case "GET": // GET = Web server got requested to show first page
		errorMes = testMethod(r, "GET", "/poststatus")
	case "POST": // POST = web server got request to generate something additional from user filled form
		errorMes = testMethod(r, "POST", "/poststatus")
		topicID := r.FormValue("id")
		topicStatus := r.FormValue("activate" + topicID)
		comment := r.FormValue("area" + topicID)
		sqlScript := "update topic_header set status='" + topicStatus + "', status_text='" + comment + "', status_date = CURRENT_TIMESTAMP where topic_id =" + topicID
		ExecSQLScript(DatabaseName, sqlScript)
	}
	if errorMes.Key == 200 {
		http.Redirect(w, r, "/myHome", http.StatusFound)
	} else {
		HandleError(w, r, errorMes)
	}
}
