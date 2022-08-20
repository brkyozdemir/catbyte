package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"sort"
	. "twitch_chat_analysis/cmd"
	"twitch_chat_analysis/cmd/messageProcessor/cache"
	"twitch_chat_analysis/cmd/models"
	"twitch_chat_analysis/config"
)

func main() {
	r := gin.Default()

	r.GET("/message/list", func(ctx *gin.Context) {
		cfg := config.GetConfig()
		c := cache.NewCache(cfg.Cache.Default).Repository
		cachedMessages, err := c.Get("message")

		sender := ctx.Query("sender")
		receiver := ctx.Query("receiver")

		messageModel := models.Message{}
		var messageArray []models.Message

		status := 200

		var messages []string
		if err != redis.Nil {
			err = json.Unmarshal([]byte(cachedMessages), &messages)
			status = FailOnError(err, "unmarshal error")

			for _, message := range messages {
				err = json.Unmarshal([]byte(message), &messageModel)
				status = FailOnError(err, "unmarshal error")
				messageArray = append(messageArray, messageModel)
			}
		}

		var response []models.Message
		for _, message := range messageArray {
			if message.Receiver == receiver && message.Sender == sender {
				response = append(response, message)
			}
		}

		sort.Slice(response, func(i, j int) bool {
			return response[i].CreatedAt > response[j].CreatedAt
		})

		ctx.JSON(status, response)
	})

	r.Run(":8081")
}
