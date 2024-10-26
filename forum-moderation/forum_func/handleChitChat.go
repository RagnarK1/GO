package forum_func

import (
	"fmt"
	"net/http"
)

func HandleChitChat(w http.ResponseWriter, r *http.Request) {
	chatID := r.FormValue("reportid")
	chatType := r.FormValue("type")
	chatText := r.FormValue("chitchat")
	_, user := activeSessionToken(r)
	if chatType == "question" {
		sqlScript := "insert into report (user_email, question) values('" + user + "', '" + chatText + "')"
		ExecSQLScript(DatabaseName, sqlScript)
	} else if chatType == "answer" {
		sqlScript := "update report set answer = '" + chatText + "' where report_id = " + chatID
		fmt.Println(sqlScript)
		ExecSQLScript(DatabaseName, sqlScript)
	}
	http.Redirect(w, r, "/myHome", http.StatusFound)
}
