package webserver

import (
	"Back_end/redisdb"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &UserReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decodng request object", http.StatusBadRequest)
		return
	}

	res := register(u)
	json.NewEncoder(w).Encode(res)
}

func register(u *UserReq) *Response {
	
	// if user not exist create , otherwise return error as response
	res := &Response{Status: true}

	status := redisdb.ExistUser(u.Username)
	if status {
		res.Status = false
		res.Message = "username already taken. try something else."
		return res
	}

	err := redisdb.RegisterUser(u.Username, u.Password)
	if err != nil {
		res.Status = false
		res.Message = "something went wrong while registering the user. please try again after sometime."
		return res
	}

	return res
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &UserReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoidng request object", http.StatusBadRequest)
		return
	}

	res := login(u)
	json.NewEncoder(w).Encode(res)
}

func login(u *UserReq) *Response {
	// if invalid username and password return error
	// if valid user create new session
	res := &Response{Status: true}

	err := redisdb.IsUserAuthentic(u.Username, u.Password)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return res
	}

	return res
}

func verifyContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &UserReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoidng request object", http.StatusBadRequest)
		return
	}

	res := verifyContact(u.Username)
	json.NewEncoder(w).Encode(res)
}

func verifyContact(username string) *Response {
	

	res := &Response{Status: true}

	status := redisdb.ExistUser(username)
	if !status {
		res.Status = false
		res.Message = "invalid username"
	}

	return res
}

func chatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u1 := r.URL.Query().Get("u1")
	u2 := r.URL.Query().Get("u2")

	fromTS, toTS := "0", "+inf"

	if r.URL.Query().Get("from-ts") != "" && r.URL.Query().Get("to-ts") != "" {
		fromTS = r.URL.Query().Get("from-ts")
		toTS = r.URL.Query().Get("to-ts")
	}

	res := chatHistory(u1, u2, fromTS, toTS)
	json.NewEncoder(w).Encode(res)
}

func chatHistory(username1, username2, fromTS, toTS string) *Response {
	res := &Response{}

	fmt.Println(username1, username2)
	// check if user exists
	if !redisdb.ExistUser(username1) || !redisdb.ExistUser(username2) {
		res.Message = "incorrect username"
		return res
	}
   // fetch the chats
	chats, err := redisdb.FetchChatBetween(username1, username2, fromTS, toTS)
	if err != nil {
		log.Println("error in fetch chat between", err)
		res.Message = "unable to fetch chat history. please try again later."
		return res
	}

	res.Status = true
	res.Data = chats
	res.Total = len(chats)
	return res
}

func contactListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := r.URL.Query().Get("username")

	res := contactList(u)
	json.NewEncoder(w).Encode(res)
}

func contactList(username string) *Response {
	// if invalid username return error
	// if valid users fetch chats
	res := &Response{}

	// check if user exists
	if !redisdb.ExistUser(username) {
		res.Message = "incorrect username"
		return res
	}

	contactList, err := redisdb.FetchContactList(username)
	if err != nil {
		log.Println("error in fetch contact list of username: ", username, err)
		res.Message = "unable to fetch contact list. please try again later."
		return res
	}

	res.Status = true
	res.Data = contactList
	res.Total = len(contactList)
	return res
}

