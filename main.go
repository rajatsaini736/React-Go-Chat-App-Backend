package main

import (
	"fmt"
	"net/http"	
	"os"

	"github.com/rajatsaini736/React-Go-Chat-App-Backend/pkg/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Println(w, "%+V\n", err)
	}

	client := &websocket.Client {
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	setupRoutes()
	http.ListenAndServe(":" + port, nil)
}