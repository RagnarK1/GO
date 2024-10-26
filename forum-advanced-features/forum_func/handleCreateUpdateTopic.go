package forum_func

import (
	"html"
	"net/http"
	"strings"
)

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	var errorMes ErrorMessage
	switch r.Method { //get request method
	case "GET": // GET = Web server got requested to show first page
		errorMes.Key = http.StatusNotFound
		errorMes.Text = "Bad request - Page Not Found"
		HandleError(w, r, errorMes)
	case "POST": // POST = web server got request to generate something additional from user filled form
		topicTitle := html.EscapeString(r.FormValue("TopicTitle"))
		topicText := html.EscapeString(r.FormValue("TopicText"))
		r.ParseForm()
		categories := r.Form["TopicCategory"]
		topicId := r.FormValue("TopicId")
		mainId := r.FormValue("MainId")
		TopicType := r.FormValue("TopicType")
		removeImg := r.FormValue("removeImg")
		topicImg, _, _ := uploadImg(w, r)
		_, user := activeSessionToken(r)
		if len(user) == 0 {
			http.Redirect(w, r, "/signin", http.StatusFound)
		} else if (TopicType == "Sub" || len(strings.TrimSpace(topicTitle)) != 0) && len(strings.TrimSpace(topicText)) != 0 {
			if TopicType == "Topic" {
				if len(topicId) != 0 {
					UpdateTopic(DatabaseName, user, topicId, topicTitle, topicText, topicImg, removeImg, categories)
				} else {
					InsertTopic(DatabaseName, user, topicTitle, topicText, topicImg, categories)
				}
				http.Redirect(w, r, "/", http.StatusFound)
			} else if TopicType == "Sub" {
				if len(topicId) != 0 {
					UpdateSubTopic(DatabaseName, user, topicId, topicTitle, topicText)
				} else {
					InsertSubTopic(DatabaseName, user, mainId, topicTitle, topicText)
				}
				http.Redirect(w, r, "/topic?Id="+mainId, http.StatusFound)
			}
		} else {
			errorMes.Key = 403
			errorMes.Text = "Inserted value(s) are in invalid format"
			HandleError(w, r, errorMes)
		}
	}
}
