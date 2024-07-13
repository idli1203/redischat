package redisdb

import (
	"fmt"
	"time"
)

// finds the user keys
func userSetKey() string {
	return "users"
}

// generates a key for every session ka client 
func sessionKey(client string) string {
	return "session#" + client
}

// returns the chat message ka unique key by making the chat message paired to the time in milliseconds
func chatKey() string {
	return fmt.Sprintf("chat#%d", time.Now().UnixMilli())
}

// 
func chatIndex() string {
	return "idx#chats"
}


func contactListZKey(username string) string {
	return "contacts:" + username 
}