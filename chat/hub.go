// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	DB "chat/db"
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run(repo *DB.Service) {
	for {
		select {
		case client := <-h.register:
			// register the client with the hub
			h.clients[client] = true
			go sendToClientAllPrevMessages(client, repo.Storage.Fetch)
		case client := <-h.unregister:
			// unregister the client with the hub and close his channel
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			go saveMsgToDb(message, repo.Storage.Append)
			for client := range h.clients {
				select {
				case client.send <- message:
				default: // If the send channel is not blocked, the message is successfully sent to the client.
					// However, if the send channel is blocked (indicated by the default case), it means the client is not receiving messages anymore.
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func saveMsgToDb(message []byte, append func(msg string) error) {
	stringMsg := string(message[:])
	fmt.Printf("saving msg to db %s\n\n", stringMsg)
	err := append(stringMsg)
	if err != nil {
		fmt.Printf("saving msg to db error %s\n\n", stringMsg)
	}
}

func sendToClientAllPrevMessages(client *Client, fetch func() (*[]string, error)) {
	messages, err := fetch()
	if err != nil {
		fmt.Printf("sendToClientAllPrevMessages %s\n\n", err)
	}
	for _, message := range *messages {
		client.send <- []byte(message)
	}
}
