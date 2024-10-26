package forum_func

import (
	"net/http"
)

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	deleteID := r.FormValue("Id")
	mainId := r.FormValue("MainId")
	topicType := r.FormValue("TopicType")
	_, user := activeSessionToken(r)
	if len(user) == 0 {
		http.Redirect(w, r, "/signin", http.StatusFound)
	} else if topicType == "Topic" {
		DeleteTopic(DatabaseName, user, deleteID)
		http.Redirect(w, r, "/", http.StatusFound)
	} else if topicType == "Sub" {
		DeleteSubTopic(DatabaseName, user, deleteID)
		http.Redirect(w, r, "/topic?Id="+mainId, http.StatusFound)
	}
}
