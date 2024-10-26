package main

import (
	"fmt"
	"forum/forum_func"
	"net/http"
	"os"
	"time"
)

// Had error in the beginning! Repaired it with next commands!
// sudo apt update
// sudo apt install build-essential
// gcc --version

// const DatabaseName = "sqlite-forum.db"
var RequestAmounts = make(map[string]int)

func main() {
	//os.Remove(forum_func.DatabaseName)                          // When you want to create new database // delete current database file.
	if _, err := os.Stat(forum_func.DatabaseName); err == nil { // control if database exists
		fmt.Println("Forum database file found. Executing server start.")
	} else {
		fmt.Println("Forum database file not found. Creating new database file.")
		forum_func.CreateDatabase(forum_func.DatabaseName) // create new database when database not found
	}
	http.Handle("/common/", http.StripPrefix("/common/", http.FileServer(http.Dir("common")))) // give access to common folder so we could access it for extra (htm, css, js) files
	http.HandleFunc("/", accessLimiter(forum_func.HandleMain))
	http.HandleFunc("/signup", accessLimiter(forum_func.HandleSignUp))
	http.HandleFunc("/signout", accessLimiter(forum_func.HandleSignOut))
	http.HandleFunc("/signin", accessLimiter(forum_func.HandleSignIn))
	http.HandleFunc("/uitopic", accessLimiter(forum_func.HandleCreate))
	http.HandleFunc("/dtopic", accessLimiter(forum_func.HandleDelete))
	http.HandleFunc("/topic", accessLimiter(forum_func.HandleTopics))
	http.HandleFunc("/likeme", accessLimiter(forum_func.HandleLikes))
	http.HandleFunc("/myHome", accessLimiter(forum_func.HandleUserPage))
	http.HandleFunc("/checkUser", accessLimiter(forum_func.HandleCheckUser))
	http.HandleFunc("/github/login", accessLimiter(forum_func.HandleGithubLogin))
	http.HandleFunc("/github/callback", accessLimiter(forum_func.HandleGithubCallback))
	http.HandleFunc("/poststatus", accessLimiter(forum_func.HandlePostStatusChange))
	http.HandleFunc("/usertype", accessLimiter(forum_func.HandleUserTypeChange))
	http.HandleFunc("/checkEmail", accessLimiter(forum_func.HandleCheckEmail))

	fmt.Println("Server was started on PORT :8081. Please go to your web browser and open https://localhost:8081")
	fmt.Println("To terminate server click CTRL + \"c\"")
	// fmt.Println(http.ListenAndServe(":8080", nil)) //start web server
	// fmt.Println(http.ListenAndServeTLS(":8081", "forum_func/cert/localhost/localhost.crt", "forum_func/cert/localhost/localhost.decrypted.key", nil))
	fmt.Println(http.ListenAndServeTLS(":8081", "forum_func/cert/forum.crt", "forum_func/cert/forum.key", nil))
	// fmt.Println(err)
}

// Count the amount of requests from one IPaddress
func accessLimiter(handlePage http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorMes forum_func.ErrorMessage
		userIP := r.RemoteAddr // user IP address
		go func() {            // remember IP address for ...
			time.Sleep(30 * time.Second)   // when 30 seconds passed
			delete(RequestAmounts, userIP) // delete IPaddress from map.
		}()
		if RequestAmounts[userIP] > 50 { // give an error when 20 attemps are done within 30 seconds. Last part deletes IP key from map after 30 seconds.
			errorMes.Key = 429
			errorMes.Text = "Too many requests from current IP address!"
			forum_func.HandleError(w, r, errorMes)
			return
		}
		RequestAmounts[userIP] += 1 // IPaddress has made a connection to our site. write into map that n+1 attempts are done
		handlePage.ServeHTTP(w, r)  // send the handle function to correct route
	})
}
