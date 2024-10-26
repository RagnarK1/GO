package forum_func

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	githubClientID     = "6d8d14d9ff53cb1e9a6c"
	githubClientSecret = "f1cd1d7bfd8cc5975a6b006f309a9a239b147c23"
	githubRedirectURL  = "https://localhost:8081/github/callback"
)

func HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	state := RandStringToken(20)

	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		url.QueryEscape(githubClientID),
		url.QueryEscape(githubRedirectURL),
		url.QueryEscape(state),
		url.QueryEscape("user:email"),
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func HandleGithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := getGithubAccessToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := getGithubUserDetails(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Println("user.id:", user.ID)
	// fmt.Println("user.login:", user.Login)
	// fmt.Println("user.email:", user.Email)
	// var errorMes ErrorMessage
	// t, err := template.ParseFiles("common/signin.htm")
	// if err != nil {
	// 	errorMes.Key = 500
	// 	errorMes.Text = "Server not reachable - 500 Internal Server Error"
	// 	HandleError(w, r, errorMes)
	// 	return
	// }

	//LogIn(user.Email, "", "github", w)
	statusNumber := LogIn(user.Login, user.Email, cryptPW(""), "github", w)
	if statusNumber == 1 { // LogIn(signin.Email, signin.Password, signin.LoginType, w) == 1
		http.Redirect(w, r, "/", http.StatusFound)
	} else if statusNumber == 2 {
		if status := LogIn(user.Login, user.Email, cryptPW(""), "github", w); status == 1 {
			//fmt.Println("here")
			http.Redirect(w, r, "/", http.StatusFound)
		}
		//http.Redirect(w, r, "/signup", http.StatusFound)
		//return
	} // else {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	t.Execute(w, "Incorrect email or password!")
	// 	return
	// }

	//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func getGithubAccessToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", githubClientID)
	data.Set("client_secret", githubClientSecret)
	data.Set("code", code)

	resp, err := http.PostForm("https://github.com/login/oauth/access_token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyStr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	params, err := url.ParseQuery(string(bodyStr))
	if err != nil {
		return "", err
	}

	return params.Get("access_token"), nil
}

func getGithubUserDetails(accessToken string) (*User, error) {
	user, err := getUserFromAPI(accessToken)
	if err != nil {
		return nil, err
	}

	primaryEmail, err := getUserPrimaryEmail(accessToken)
	if err != nil {
		return nil, err
	}

	user.Email = primaryEmail
	return user, nil
}

func getUserFromAPI(accessToken string) (*User, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting user details: %v", string(body))
	}

	user := &User{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func getUserPrimaryEmail(accessToken string) (string, error) {
	userEmailURL := "https://api.github.com/user/emails"
	userEmailReq, err := http.NewRequest("GET", userEmailURL, nil)
	if err != nil {
		return "", err
	}
	userEmailReq.Header.Set("Authorization", "Bearer "+accessToken)
	userEmailResp, err := http.DefaultClient.Do(userEmailReq)
	if err != nil {
		return "", err
	}
	defer userEmailResp.Body.Close()

	userEmails := []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}{}
	err = json.NewDecoder(userEmailResp.Body).Decode(&userEmails)
	if err != nil {
		return "", err
	}

	for _, email := range userEmails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	return "", errors.New("no primary email found for user")
}

type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}
