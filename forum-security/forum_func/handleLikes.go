package forum_func

import (
	"encoding/json"
	"net/http"
	"strings"
)

func HandleLikes(w http.ResponseWriter, r *http.Request) {
	var LikesDisLikes LikeDisLike
	subId := r.FormValue("Id")
	mainId := r.FormValue("MainId")
	LikeType := r.FormValue("LikeType")
	LikeType = strings.Title(strings.ToLower(LikeType))
	_, user := activeSessionToken(r)
	if len(user) == 0 {
		http.Redirect(w, r, "/signin", http.StatusFound)
	} else {
		if subId == "null" { //it is Topic like
			LikesDisLikes = TopicLikes(DatabaseName, mainId, LikeType, user)
		} else {
			LikesDisLikes = SubTopicLikes(DatabaseName, mainId, subId, LikeType, user)
		}
		a, _ := json.Marshal(LikesDisLikes)
		w.Write(a)
	}
}
