package forum_func

import (
	"net/http"
	"text/template"
)

func HandleTopics(w http.ResponseWriter, r *http.Request) {
	var errorMes ErrorMessage
	t, err := template.ParseFiles("common/topic_content.htm")
	if err != nil {
		errorMes.Key = 500
		errorMes.Text = "Server not reachable - 500 Internal Server Error"
		HandleError(w, r, errorMes)
		return
	}
	errorMes = testMethod(r, "GET", "/topic")

	if errorMes.Key == 200 {
		topicId := r.FormValue("Id")
		sessionToken, user := activeSessionToken(r)
		homeVarables := SubTopicVarable{
			Session: UserSessions[sessionToken],
			Code:    SubTopics(DatabaseName, user, topicId),
			Topic:   GetTopic(DatabaseName, topicId)}
		t.Execute(w, homeVarables)
	} else {
		HandleError(w, r, errorMes)
	}
}
