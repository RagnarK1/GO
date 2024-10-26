package forum_func

import (
	"net/http"
	"text/template"
)

func HandleError(w http.ResponseWriter, r *http.Request, errorMes ErrorMessage) {
	t, err := template.ParseFiles("common/error.htm")
	if err != nil {
		http.Error(w, "Server not reachable - 500 Internal Server Error", http.StatusNotFound)
		return
	}
	t.Execute(w, errorMes)
}
