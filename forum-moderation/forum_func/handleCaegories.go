package forum_func

import (
	"net/http"
)

func HandleCategoriesChange(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue("cattext")
	cattype := r.FormValue("cattype")
	if cattype == "add" {
		sqlScript := "insert into setup (setup_type, setup_value) values('CATEGORY', '" + category + "')"
		ExecSQLScript(DatabaseName, sqlScript)
	} else if cattype == "delete" {
		sqlScript := "delete from setup where setup_type = 'CATEGORY' and setup_value= '" + category + "'"
		ExecSQLScript(DatabaseName, sqlScript)
	}
	http.Redirect(w, r, "/myHome", http.StatusFound)
}
