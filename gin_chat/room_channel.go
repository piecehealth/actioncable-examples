package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/piecehealth/actioncable"
)

var getRoomId = func(c *actioncable.Channel) string {
	params := struct {
		Id string `json:"id"`
	}{}

	json.Unmarshal(c.Params, &params)

	return fmt.Sprintf("room_%s", params.Id)
}

var roomChannel = &actioncable.ChannelDescripion{
	Name: "RoomChannel",
	Subscribed: func(c *actioncable.Channel) {
		roomId := getRoomId(c)
		currentUser := fmt.Sprintf("%v", c.ConnIdentifier)

		if roomId == "forbidden_room" && currentUser != "the chosen one" {
			c.Reject()
			return
		}

		log.Printf("%v just join the %s.", c.ConnIdentifier, roomId)

		c.StreamFrom(roomId)

		m := map[string]string{
			"user_join": currentUser + " just join the room.",
		}
		c.Broadcast(roomId, m)
	},
	Unsubscribed: func(c *actioncable.Channel) {
		roomId := getRoomId(c)

		log.Printf("%v just left the %s.", c.ConnIdentifier, roomId)

		m := map[string]string{
			"user_leave": fmt.Sprintf("%v", c.ConnIdentifier) + " just left the room.",
		}
		cable.Broadcast("RoomChannel", roomId, m)
	},
	PerformAction: func(c *actioncable.Channel, data string) {
		d := struct {
			Action  string `json:"action"`
			Message string `json:"message"`
			Name    string `json:"name"`
		}{}

		json.Unmarshal([]byte(data), &d)

		roomId := getRoomId(c)

		switch d.Action {
		case "send_message":
			currentUser := fmt.Sprintf("%v", c.ConnIdentifier)

			m := map[string]string{
				"send_by": currentUser,
				"message": d.Message,
			}

			c.Broadcast(roomId, m)
		case "kick":
			log.Printf("disconnect %s's connection.", d.Name)
			cable.DisconnectRemoteConnection(d.Name)
		case "stop_stream":
			log.Println("stop all streams.")
			c.StopAllStreams()
		default:
			log.Printf("Unable to process RoomChannel#%s", d.Action)
		}
	},
}
