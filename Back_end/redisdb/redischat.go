package redisdb

import (
	"Back_end/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func CreateChat(c *models.Person) (string, error) {

	chatKey := chatKey()
	fmt.Println("chat key", chatKey)

	msg , _ := json.Marshal(c)

	// JSON.SET chat#1661360942123 $ '{"from": "sun", "to":"earth","message":"good morning!"}'
	res, err := RedisClient.Do(
		context.Background(),
		"JSON.SET",
		chatKey,
		"$",
		string(msg),
	).Result()

	if err != nil {
		log.Println("error while setting chat json", err)
		return "", err
	}

	log.Println("chat successfully set", res)

	// add contacts to both user's contact list
	err = UpdateContactList(c.From, c.To)
	if err != nil {
		log.Println("error while updating contact list of", c.From)
	}

	err = UpdateContactList(c.To, c.From)
	if err != nil {
		log.Println("error while updating contact list of", c.To)
	}

	return chatKey, nil
}


func UpdateContactList(username, contact string) error {

	 upd := &redis.Z{Score: float64(time.Now().Unix()), Member: contact}

	
	 err := RedisClient.ZAdd(context.Background(),
		contactListZKey(username),
		*upd,
	).Err()

	if err != nil {
		log.Println("error while updating contact list. username: ",
			username, "contact:", contact, err)
		return err
	}

	return nil
}

func CreateFetchChatBetweenIndex() {
	res, err := RedisClient.Do(context.Background(),
		"FT.CREATE",
		chatIndex(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA", "$.from", "AS", "from", "TAG",
		"$.to", "AS", "to", "TAG",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE",
	).Result()

	fmt.Println(res, err)
}

func FetchChatBetween(username1, username2, fromTS, toTS string) ([]models.Person, error) {


	query := fmt.Sprintf("@from:{%s|%s} @to:{%s|%s} @timestamp:[%s %s]",
		username1, username2, username1, username2, fromTS, toTS)

	res, err := RedisClient.Do(context.Background(),
		"FT.SEARCH",
		chatIndex(),
		query,
		"SORTBY", "timestamp", "DESC",
	).Result()

	if err != nil {
		return nil, err
	}

	// convert redis data to something usable like map 
	data := Deserialise(res)

	// port the map to chat between users
	chats := DeserialiseChat(data)
	return chats, nil
}


func FetchContactList(username string) ([]models.Contact, error) {
	zRangeArg := redis.ZRangeArgs{
		Key:   contactListZKey(username),
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}


	res, err := RedisClient.ZRangeArgsWithScores(context.Background(), zRangeArg).Result()

	if err != nil {
		log.Println("error while fetching contact list. username: ",
			username, err)
		return nil, err
	}

	contactList := DeserialiseContactList(res)

	return contactList, nil
}