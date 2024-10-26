package forum_func

import (
	"net/http"
	"text/template"
)

func HandleMain(w http.ResponseWriter, r *http.Request) {
	var errorMes ErrorMessage
	errorMes.Key = 400
	errorMes.Text = "Page not found"
	var homeVariables TopicVarable
	t, err := template.ParseFiles("common/index.htm")
	if err != nil {
		errorMes.Key = 500
		errorMes.Text = "Server not reachable - 500 Internal Server Error"
		HandleError(w, r, errorMes)
		return
	}
	if r.URL.Path == "/" {
		switch r.Method { //get request method
		case "GET": // GET = Web server got requested to show first page
			errorMes = testMethod(r, "GET", "/")
			sessionToken, user := activeSessionToken(r)
			homeVariables = TopicVarable{Session: UserSessions[sessionToken], Code: MainTopics(DatabaseName, user, "AllTopics"), Categories: ReadSetupTable(DatabaseName, "CATEGORY"), Users: ReadTopicUsers(DatabaseName, user)}
		case "POST": // POST = web server got request to generate something additional from user filled form
			errorMes = testMethod(r, "POST", "/")
		}
	}
	if errorMes.Key == 200 {
		t.Execute(w, homeVariables)
	} else {
		HandleError(w, r, errorMes)
	}

}
