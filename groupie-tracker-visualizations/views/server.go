package server

import (
	"fmt"
	"html/template"
	"net/http"
	models "server/models"
	"strconv"
	print2 "server/models"
)

func main_page(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path != "/" {		
			errHandle(http.StatusNotFound, w)
			return
		}
	}
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		fmt.Println(err)
		errHandle(http.StatusInternalServerError, w)
		return
	}
	print2.HandleJson()
	tmpl.Execute(w, models.ARtists)
}

func band_page(w http.ResponseWriter, r *http.Request) {
	var id int
	if r.URL.Path != "/about"{
		errHandle(http.StatusNotFound, w)			
		return
	}
	if r.Method == "POST" {	
		var err error
		input := r.FormValue("band-ID")
		id, err = strconv.Atoi(input)
		if err != nil {
			errHandle(http.StatusInternalServerError, w)
		}
	} else {
		errHandle(http.StatusMethodNotAllowed, w)
		return
	}
	tmpl, err := template.ParseFiles("views/about.html")
	if err != nil {
		errHandle(http.StatusInternalServerError, w)
	}	
	tmpl.Execute(w, models.ARtists[id-1])
}

func Main() {
	fmt.Println("link: http://localhost:6969/ ")

	http.Handle("/views/assets/", http.StripPrefix("/views/assets/", http.FileServer(http.Dir("views/assets"))))
	http.Handle("/about/views/assets/", http.StripPrefix("/about/views/assets/", http.FileServer(http.Dir("views/assets"))))
	http.HandleFunc("/", main_page)
	http.HandleFunc("/about", band_page)

	http.ListenAndServe(":6969", nil)
}

//takes error code and handles it
 func errHandle(statusCode int, w http.ResponseWriter){
	var tmpl *template.Template
	var err error
	switch statusCode {
	case 404:
		tmpl, err = template.ParseFiles("views/err404.html")
	case 405:
		tmpl, err = template.ParseFiles("views/err405.html")
	default:
		tmpl, err = template.ParseFiles("views/err500.html")
	}
	if err != nil {
		var err1 error
		tmpl, err1 = template.ParseFiles("views/err500.html")
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	tmpl.Execute(w, nil)
 }