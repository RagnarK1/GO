package server

import (
	"fmt"
	"html/template"
	"net/http"
	models "server/models"
	"strconv"
)

func main_page(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound) //future: http.ServeFile(w,r, "templates/err404.html")
			return
		}
	}
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, models.ARtists)
}

func band_page(w http.ResponseWriter, r *http.Request) {
	var id int
	if r.Method == "POST" {
		var err error
		input := r.FormValue("band-ID")
		id, err = strconv.Atoi(input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed) // future http.ServeFile(w,r, "templates/err400.html")
		return
	}
	tmpl, err := template.ParseFiles("views/about.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	tmpl.Execute(w, models.ARtists[id-1])

}

func Main() {
	fmt.Println("link: http://localhost:6969/ ")

	http.Handle("/views/assets/", http.StripPrefix("/views/assets/", http.FileServer(http.Dir("views/assets"))))

	http.HandleFunc("/", main_page)
	http.HandleFunc("/about", band_page)

	http.ListenAndServe(":6969", nil)
}
