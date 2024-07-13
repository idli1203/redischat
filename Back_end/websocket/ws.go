package websocket

import (
	"Back_end/models"
	"Back_end/redisdb"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type Client struct {
	Conn     *websocket.Conn
	Username string
}

type Message struct {
	Type string     `json:"type"`
	User string     `json:"user,omitempty"`
	Chat models.Person `json:"chat,omitempty"`
}

var clients = make(map[*Client] bool)
var broadcast = make(chan *models.Person)

// Upgrader upgrades the http connection to a websocket 
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host, r.URL.Query())

	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{Conn: ws}
	// register client
	clients[client] = true
	fmt.Println("clients", len(clients), clients , ws.RemoteAddr())

	// listen indefinitely for new messages coming
	receiver(client)

	fmt.Println("exiting", ws.RemoteAddr().String())
	delete(clients, client)
}

func receiver(client *Client) {
	for {
		// read in a message
		// messageType: 1-> Text Message, 2 -> Binary Message
		messagetype  , p, err := client.Conn.ReadMessage()

		fmt.Println(messagetype)

		if err != nil {
			log.Println(err)
			return
		}

		m := &Message{}

		err = json.Unmarshal(p, m)
		if err != nil {
			log.Println("error while unmarshaling chat", err)
			continue
		}

		fmt.Println("host", client.Conn.RemoteAddr())
		if m.Type == "bootup" {
			client.Username = m.User
			fmt.Println("client successfully mapped", &client, client, client.Username)
		} else {
			fmt.Println("received message", m.Type, m.Chat)
			c := m.Chat
			c.Sendtime = time.Now().Unix()

			// save in redis
			id, err := redisdb.CreateChat(&c)
			if err != nil {
				log.Println("error while saving chat in redis", err)
				return
			}

			c.ID = id
			broadcast <- &c
		}
	}
}

func broadcaster() {
	for {
		message := <-broadcast
		fmt.Println("new message", message)

		for client := range clients {
			// send message only to involved users
			fmt.Println("username:", client.Username,
				"from:", message.From,
				"to:", message.To)

			if client.Username == message.From || client.Username == message.To {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}


func WsRoutes() http.Handler{
	r := chi.NewRouter();

	r.Get("/ws" , serveWs) 


	return r
}

func StartWebsocketServer() {
	RedisClient := redisdb.OpenRedis()
	defer RedisClient.Close()

	go broadcaster()
	
	srv := http.Server {
		 Addr : viper.GetString("WEBSOCKETPORT") ,
		 Handler : WsRoutes(),
	}

	err := srv.ListenAndServe() ; if err != nil {
		log.Fatal(err)
	}
}
