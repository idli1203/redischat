package redisdb

import (
	"Back_end/models"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type Information struct {
	ID      string `json:"_id"`
	Payload []byte `json:"payload"`
	Total   int64  `json:"total"`
}

func Deserialise(res interface{}) [] Information {


	switch v := res.(type) {
	case []interface{}:

		if len(v) > 1 {
			total := len(v) - 1
			var docs = make([] Information, 0, total/2)

			for i := 1; i <= total; i = i + 2 {
				arrOfValues := v[i+1].([]interface{})
				value := arrOfValues[len(arrOfValues)-1].(string)

				doc := Information{
					ID:      v[i].(string),
					Payload: []byte(value),
					Total:   v[0].(int64),
				}

				docs = append(docs, doc)
			}
			return docs
		}
	default:
		log.Printf("different response type otherthan []interface{}. type: %T", res)
		return nil
	}

	return nil
}

func DeserialiseChat(docs []Information ) []models.Person {
	chats := []models.Person{}
	for _, doc := range docs {
		var c models.Person
		json.Unmarshal(doc.Payload, &c)

		c.ID = doc.ID
		chats = append(chats, c)
	}

	return chats
}

func DeserialiseContactList(contacts [] redis.Z) []models.Contact {
	contactList := make([]models.Contact, 0, len(contacts))

	for _ , contact := range contacts {
		contactList = append(contactList, models.Contact{
			Username:     contact.Member.(string),
			Last_activity: int64(contact.Score),
		})
	}

	return contactList
}