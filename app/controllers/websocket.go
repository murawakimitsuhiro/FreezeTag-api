package controllers

import (
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"fmt"
	"freezetag/app/models"
	"encoding/json"
)

type WebSocket struct {
	*revel.Controller
}

func (c WebSocket) RoomSocket(user string, ws *websocket.Conn) revel.Result {
	// Make sure the websocket is valid.
	if ws == nil {
		return nil
	}

	// Join the room.
	subscription := Subscribe()//chatroom.Subscribe()
	defer subscription.Cancel()

	/*
	chatroom.Join(user)
	defer chatroom.Leave(user)
	*/
	Join(user)
	defer Leave(user)

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

			jsonByte := ([]byte)(msg)
			userLocation := models.UserLocation{}

			if err := json.Unmarshal(jsonByte, &userLocation); err != nil {
				fmt.Println("JSON Unmarshal error:", err)
				return nil
			}

			UpdateLocation(userLocation)
			//chatroom.UpdateLocation(userLocation)
		}
	}
	return nil
}
