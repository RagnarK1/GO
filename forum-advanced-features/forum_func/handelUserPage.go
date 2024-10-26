package forum_func

import (
	"html/template"
	"net/http"
)

type myPageVariable struct {
	Session     UserSession
	MyTopics    []ForumTopic
	Rights      []UserType
	Pending     []ForumTopic
	UserTypes   []string
	TopicStatus []string
	LikeDislike []MyLikeDislikeTopic
	Comment     []MyComment
	Updates     []ChangesAfterLogin
	Categories  []string
}

func HandleUserPage(w http.ResponseWriter, r *http.Request) {
	var errorMes ErrorMessage
	var myPageVariables myPageVariable
	t, err := template.ParseFiles("common/userpage.htm")
	if err != nil {
		errorMes.Key = 500
		errorMes.Text = "Server not reachable - 500 Internal Server Error"
		HandleError(w, r, errorMes)
		return
	}
	switch r.Method { //get request method
	case "GET": // GET = Web server got requested to show first page
		errorMes = testMethod(r, "GET", "/myHome")
		sessionToken, user := activeSessionToken(r)
		if len(user) == 0 {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		// homeVariables = TopicVarable{Session: UserSessions[sessionToken], Code: MainTopics(DatabaseName, user), Categories: ReadSetupTable(DatabaseName, "CATEGORY"), Users: ReadTopicUsers(DatabaseName, user)}
		myPageVariables = myPageVariable{Session: UserSessions[sessionToken],
			MyTopics:    MainTopics(DatabaseName, user, "MyTopics"),
			Rights:      FindUserRights(DatabaseName, user),
			Pending:     MainTopics(DatabaseName, user, "PendingTopics"),
			UserTypes:   ReadSetupTable(DatabaseName, "USER_TYPE"),
			TopicStatus: ReadSetupTable(DatabaseName, "STATUS"),
			LikeDislike: MyLikesDislikes(DatabaseName, user),
			Comment:     MyComments(DatabaseName, user),
			Updates:     ChangesAfterLastLogin(DatabaseName, user),
			Categories:  ReadSetupTable(DatabaseName, "CATEGORY")}
		// case "POST": // POST = web server got request to generate something additional from user filled form
		// 	errorMes = testMethod(r, "POST", "/myHome")
	}
	if errorMes.Key == 200 {
		t.Execute(w, myPageVariables)
	} else {
		HandleError(w, r, errorMes)
	}
	return
}
