package webserver

import (
	"Back_end/redisdb"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
)

func StartHTTPServer() {

	RedisClient := redisdb.OpenRedis()
	defer RedisClient.Close()

	// create indexes
	redisdb.CreateFetchChatBetweenIndex()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	  }))
	r.Get ("/status" , func (w http.ResponseWriter , r* http.Request) {
		
		fmt.Fprintln(w , "began the server ")
	})

	r.Post("/register", registerHandler)
	r.Post("/login", loginHandler)
	r.Post("/verify-contact", verifyContactHandler)
	r.Get("/chat-history", chatHistoryHandler)
	r.Get("/contact-list", contactListHandler)

	serv := http.Server {
		Addr: viper.GetString("PORT"),
		Handler: r,
	}

	err := serv.ListenAndServe(); if err != nil {
		log.Fatal(err)
	}
}