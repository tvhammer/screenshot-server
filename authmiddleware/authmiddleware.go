package authmiddleware

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	u "os/user"
)

func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !checkUsernameAndPassword(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		handler(w, r)
	}
}

type user struct {
	USERNAME string
	PASSWORD string
}

type users struct {
	USERS []user
}

func checkUsernameAndPassword(username, password string) bool {
	current, err := u.Current()
	data, err := ioutil.ReadFile(current.HomeDir + "/.screencapture.json")
	if err != nil {
		log.Print(err)
	}

	var obj users

	err = json.Unmarshal(data, &obj)
	if err != nil {
		log.Println("error:", err)
	}
	for _, u := range obj.USERS {
		if username == u.USERNAME && password == u.PASSWORD {
			return true
		}
	}

	return false
}
