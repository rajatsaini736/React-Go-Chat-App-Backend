package websocket

import (
	"fmt"
)

type Pool struct {
	Register 	chan *Client
	Unregister 	chan *Client
	Clients 	map[*Client]bool
	Broadcast	chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register: 	make(chan *Client), // unbuffered channel
		Unregister: make(chan *Client),
		Clients: 	make(map[*Client]bool),
		Broadcast: 	make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			
			msg := `{"user":"", "msg":"New User Joined..."}`
			for client, _ := range pool.Clients {

				client.Conn.WriteJSON(Message{Type: 1, Body: msg})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))

			msg := `{"user":"","msg":"User Disconnected..."}`
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: msg})
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}