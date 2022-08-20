package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
	. "twitch_chat_analysis/cmd"
	. "twitch_chat_analysis/cmd/models"
)

func main() {
	r := gin.Default()

	r.POST("/message", func(c *gin.Context) {
		message := Message{}
		c.ShouldBindJSON(&message)
		message.CreatedAt = time.Now().Unix()

		conn, err := amqp.Dial("amqp://user:password@localhost:7001/")
		if err != nil {
			log.Panicf("%s: %s", "failed this connection", err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		FailOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"catByte",
			false,
			false,
			false,
			false,
			nil,
		)
		FailOnError(err, "Failed to declare a queue")

		marshaledMessage, _ := json.Marshal(message)
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        marshaledMessage,
			})
		FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", string(marshaledMessage))
	})

	r.Run()
}
