/*This file is for the demo Social Network application by May Contain Awesome's Software Design School.
It demonstrates how a simple social network can be coded in the Go programming language.
It is not intended to be used as an actual social network.
It does not feature any form of user authentication.

Copyright (c) May Contain Awesome, 2015*/

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

type Post struct {
	Id      string
	Name    string
	Message string
	Time    time.Time
	Replies []Post
}

var Posts []Post

//This function handles a new post request, either a reply or a new
func newPostHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	nameStr := req.FormValue("Name") //when we create either a new post or a reply, this contains the name of the user

	messageStr := req.FormValue("Post") //when we are creating a new post, this contains the post message

	replyStr := req.FormValue("Reply")     //when we are replying to a post this contains the reply message
	replyIdStr := req.FormValue("ReplyId") //when we are replying to a post, this contains the CSV id of that post

	currentTime := time.Now()
	//Check to see if we are creating a New post
	if nameStr != "" && messageStr != "" {
		//Creating a new post is easy, just generate the post...
		newPost := Post{
			Id:      strconv.Itoa(len(Posts)),
			Name:    nameStr,
			Message: messageStr,
			Time:    currentTime,
		}
		//...then append it to the global Posts array.
		Posts = append(Posts, newPost)
	} else if replyStr != "" && nameStr != "" && replyIdStr != "" { //check if we are creating a Reply post

		//Replying to a post is a little more complicated than creating a new one, as we need to
		//track down the post we're replying to before we append to its .Replies
		//The process, however, can be broken down into a simple 5-step algorithm

		//Step 1: Convert the string of comma separated values into an array of integers
		Ids := strings.Split(replyIdStr, ",")
		IdInts := make([]int, len(Ids))
		var err error
		var IdInt int
		for i, IdStr := range Ids {
			IdInt, err = strconv.Atoi(IdStr)
			if err != nil {
				break
			}
			IdInts[i] = IdInt
		}
		if err == nil {
			//Step 2: find the post we are replying to
			pAddr := &Posts[IdInts[0]]

			for i, Id := range IdInts {
				if i == 0 {
					continue
				}
				pAddr = &pAddr.Replies[Id]
			}

			//Step 3: Generate the ID of the new post
			newId := strconv.Itoa(IdInts[0])

			IdInts = append(IdInts, len(pAddr.Replies))

			for i, Id := range IdInts {
				if i == 0 {
					continue
				}
				newId += "," + strconv.Itoa(Id)
			}

			//Step 4: Generate the new post

			newPost := Post{
				Id:      newId,
				Name:    nameStr,
				Message: replyStr,
				Time:    currentTime,
			}

			//Step 5: Append the new post to the post we're replying to
			pAddr.Replies = append(pAddr.Replies, newPost)
		}
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

//This handler is the home page, and just shows the client the templates that render all the posts.
func handler(w http.ResponseWriter, req *http.Request) {
	if err := templates.Execute(w, Posts); err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func main() {
	fmt.Println("I am a webserver")

	serverAddress := ":8080"

	//An example of a Post tree.
	//Note how the ID progression works:
	firstPost := Post{
		Id:      "0", //Parent ID
		Name:    "Jebediah7",
		Message: "Good afternoon",
		Time:    time.Now(),
		Replies: []Post{
			Post{
				Id:      "0,0", //Parent ID followed by child ID
				Name:    "Franklin24",
				Message: "Good afternoon indeed",
				Time:    time.Now(),
				Replies: []Post{
					Post{
						Id:      "0,0,0", //Parent ID followed by Parent ID followed by child ID
						Name:    "Jebediah7",
						Message: "Indeed, for it is a good afternoon.",
						Time:    time.Now(),
					},
				},
			},
			Post{
				Id:      "0,1", //Parent ID followed by child ID, 1 because 0 already used
				Name:    "Shroomsbury616",
				Message: "I disagree about how good this afternoon is :/",
				Time:    time.Now(),
				Replies: []Post{
					Post{
						Id:      "0,1,0", //Parent ID followed by Parent ID followed by child ID
						Name:    "Jebediah7",
						Message: "Then I am sorry to hear that, Shroomsbury616.",
						Time:    time.Now(),
					},
				},
			},
		},
	}

	//Append our example post to the Posts global array
	Posts = append(Posts, firstPost)

	//Handle public files (eg css)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	//Handlers for different URLs
	http.HandleFunc("/newpost", newPostHandler)
	http.HandleFunc("/", handler)

	//Run server. Print any error messages.
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		fmt.Println("Error:", err.Error())
	}
}
