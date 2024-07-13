package models

type Person struct {
	ID string       `json:"id"`
	From string      `json:"from"`
	To   string      `json:"to"`
	Message string    `json:"message"`
	Sendtime int64 `json:"Sendtime"`
} 

type Contact struct {
	Username string `json:"username"`
	Last_activity int64 `json:"last_active"`
}