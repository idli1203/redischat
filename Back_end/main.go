package main

import (
	"Back_end/webserver"
	"Back_end/websocket"
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig() 
	
	if err != nil {
		panic (err)
	}
}
func main() {

	server := flag.String("server" , "" , "start")
	flag.Parse()

	if *server == "start" {
		fmt.Println("http server is starting on :3000 & websocket server is starting on :3030")
		webserver.StartHTTPServer()
		go websocket.StartWebsocketServer()
	}

	
} 