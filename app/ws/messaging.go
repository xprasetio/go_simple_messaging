package ws

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
)

type MessagePayload struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func ServeWSMessaging(app *fiber.App) {
	var (
		clients   = make(map[*websocket.Conn]bool) // Connected clients
		broadcast = make(chan MessagePayload)      // Broadcast channel
	)

	app.Get("/message/v1/send", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()

		// Register the client
		clients[c] = true

		for {
			var msg MessagePayload
			if err := c.ReadJSON(&msg); err != nil {
				fmt.Println("error reading JSON:", err)
				break // Exit on error
			}

			// Broadcast the message to all clients
			broadcast <- msg
		}
	}))
	go func() {
		for {
			msg := <-broadcast
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					fmt.Println("error writing JSON:", err)
					client.Close()          // Close the connection on error
					delete(clients, client) // Remove the client from the map
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT_SOCKET", "4000"))))
}
