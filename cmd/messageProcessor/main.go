package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
	. "twitch_chat_analysis/cmd"
	"twitch_chat_analysis/cmd/messageProcessor/cache"
	"twitch_chat_analysis/config"
)

func main() {
	err := godotenv.Load(".env")
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

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			cfg := config.GetConfig()
			c := cache.NewCache(cfg.Cache.Default).Repository

			var messages []string
			cachedMessages, err := c.Get("message")

			if err != redis.Nil {
				json.Unmarshal([]byte(cachedMessages), &messages)
				messages = append(messages, string(d.Body))
			} else {
				messages = append(messages, string(d.Body))
			}

			mod, _ := json.Marshal(messages)
			err = c.Set("message", mod, 24*time.Hour)
			FailOnError(err, "failed on setting")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
