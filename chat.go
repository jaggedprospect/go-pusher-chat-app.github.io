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
        AppID:   "996271",
        Key:     "de95eab4d71f410e9e3d",
        Secret:  "1a26d99bc2365bb2ee52",
        Cluster: "us2",
        Secure:  true,
    }

    type user struct {
        Name  string `json:"name" xml:"name" form:"name" query:"name"`
        Email string `json:"email" xml:"email" form:"email" query:"email"`
	}

	/*
	* Trigger a Pusher event on public channel to send details
	* to subscribed clients
	*/
	func registerNewUser(rw http.ResponseWriter, req *http.Request) {
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
            panic(err)
        }

        var newUser user

        err = json.Unmarshal(body, &newUser)
        if err != nil {
            panic(err)
        }

        client.Trigger("update", "new-user", newUser) // public channel, Pusher event, userc

        json.NewEncoder(rw).Encode(newUser)
    }

    func pusherAuth(res http.ResponseWriter, req *http.Request) {
        params, _ := ioutil.ReadAll(req.Body)
        response, err := client.AuthenticatePrivateChannel(params)
        if err != nil {
            panic(err)
        }

        fmt.Fprintf(res, string(response))
    }

    func main(){
        http.Handle("/", http.FileServer(http.Dir("./public"))) // serves static files for server

        http.HandleFunc("/new/user", registerNewUser) // handles new users
		http.HandleFunc("/pusher/auth", pusherAuth) // authorizes users from client-side
		
        log.Fatal(http.ListenAndServe(":8090", nil))
    }