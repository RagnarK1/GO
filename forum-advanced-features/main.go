package main

import (
	"crypto/tls"
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

const (
	port               = "8081"
	cert               = "forum_func/cert/forum.crt"
	key                = "forum_func/cert/forum.key"
	RateLimitSeconds   = 30
	allowedConnections = 10
)

var (
	RequestAmounts = make(map[string]int)
	cipherSuites   = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}

	tlsCfg = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites:             cipherSuites,
	}

	server = &http.Server{
		Addr:      ":" + port,
		Handler:   nil,
		TLSConfig: tlsCfg,
	}
)

func main() {
	// forum_func.CreateDatabase(forum_func.DatabaseName) // create new database when database not found
	if _, err := os.Stat(forum_func.DbName); err == nil { // control if database exists
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
	http.HandleFunc("/poststatus", accessLimiter(forum_func.HandlePostStatusChange))
	http.HandleFunc("/usertype", accessLimiter(forum_func.HandleUserTypeChange))

	fmt.Println("Server was started on PORT :8081. Please go to your web browser and open https://localhost:8081")
	fmt.Println("To terminate server click CTRL + \"c\"")
	// fmt.Println(http.ListenAndServe(":8080", nil)) //start web server
	// fmt.Println(http.ListenAndServeTLS(":8081", "forum_func/cert/localhost/localhost.crt", "forum_func/cert/localhost/localhost.decrypted.key", nil))
	fmt.Println(server.ListenAndServeTLS(cert, key))
	// fmt.Println(http.ListenAndServeTLS(":8081", cert, key, nil))

}

// Count the amount of requests from one IPaddress
func accessLimiter(handlePage http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorMes forum_func.ErrorMessage
		userIP := r.RemoteAddr // user IP address
		go func() {            // remember IP address for ...
			time.Sleep(RateLimitSeconds * time.Second) // when 30 seconds passed
			delete(RequestAmounts, userIP)             // delete IPaddress from map.
		}()
		if RequestAmounts[userIP] > allowedConnections { // give an error when 10 attemps are done within 30 seconds. Last part deletes IP key from map after 30 seconds.
			errorMes.Key = 429
			errorMes.Text = "Too many requests from current IP address!"
			forum_func.HandleError(w, r, errorMes)
			return
		}
		RequestAmounts[userIP] += 1 // IPaddress has made a connection to our site. write into map that n+1 attempts are done
		handlePage.ServeHTTP(w, r)  // send the handle function to correct route
	})
}
