package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pusher "github.com/pusher/pusher-http-go"
)

var client = pusher.Client{
	AppID:   "PUSHER_APP_ID",
	Key:     "PUSHER_APP_KEY",
	Secret:  "PUSHER_APP_SECRET",
	Cluster: "PUSHER_APP_CLUSTER",
	Secure:  true,
}

type user struct {
	Username string `json:"username",xml:"username",form:"username",query:"username"`
	Email    string `json:"email",xml:"email",form:"email",query:"email"`
}

var loggedInUser user

func isUserLoggedIn(rw http.ResponseWriter, req *http.Request) {
	if loggedInUser.Username != "" && loggedInUser.Email != "" {
		json.NewEncoder(rw).Encode(loggedInUser)
	} else {
		json.NewEncoder(rw).Encode("false")
	}
}

func NewUser(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &loggedInUser)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(rw).Encode(loggedInUser)
}

func pusherAuth(res http.ResponseWriter, req *http.Request) {
	params, _ := ioutil.ReadAll(req.Body)
	data := pusher.MemberData{
		UserID: loggedInUser.Username,
		UserInfo: map[string]string{
			"email": loggedInUser.Email,
		},
	}
	response, err := client.AuthenticatePresenceChannel(params, data)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(res, string(response))
}

func main() {

	//Define our routes
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/isLoggedIn", isUserLoggedIn)
	http.HandleFunc("/new/user", NewUser)
	http.HandleFunc("/pusher/auth", pusherAuth)

	//Start executing the application pn port 8090
	log.Fatal(http.ListenAndServe(":8090", nil))
}
