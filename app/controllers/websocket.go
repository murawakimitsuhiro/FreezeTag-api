package controllers

import (
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"github.com/revel/examples/chat/app/chatroom"
	"freezetag/app/models"
	"fmt"
)

type WebSocket struct {
	*revel.Controller
}

func (c WebSocket) RoomSocket(user string, ws *websocket.Conn) revel.Result {
	/*
	type request struct {
		UserLoation models.UserLocation `json: Connection`
	}*/
	fmt.Println("繋がれ")

	// Make sure the websocket is valid.
	if ws == nil {
		return nil
	}

	// Join the room.
	subscription := chatroom.Subscribe()
	defer subscription.Cancel()

	chatroom.Join(user)
	defer chatroom.Leave(user)

	// Send down the archive.
	for _, event := range subscription.Archive {
		if websocket.JSON.Send(ws, &event) != nil {
			// They disconnected
			return nil
		}
	}

	// In order to select between websocket messages and subscription events, we
	// need to stuff websocket events into a channel.
	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(newMessages)
				return
			}
			newMessages <- msg
		}
	}()

	// Now listen for new events from either the websocket or the chatroom.
	for {
		select {
		case event := <-subscription.New:
			if websocket.JSON.Send(ws, &event) != nil {
				// They disconnected.
				return nil
			}
		case msg, ok := <-newMessages:
			// If the channel is closed, they disconnected.
			if !ok {
				return nil
			}

			// Otherwise, say something.
			chatroom.Say(user, msg)
		}
	}
	return nil
}
