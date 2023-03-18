package main

// an simple http server that can be used to check username and password

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		log.Println("check request")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		u := &user{}
		if err := json.Unmarshal(body, u); err != nil {
			w.WriteHeader(400)
			return
		}
		log.Println("check request, username: ", u.Username, " password: ", u.Password)
		if u.Username == "foo" && u.Password == "bar" {
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(401)
	})
	log.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
