package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/line/line-bot-sdk-go/linebot"
	)
	
func Handler(w http.ResponseWriter, req *http.Request) {
		bot, err := linebot.New(
			os.Getenv("ChannelSecret"),
			os.Getenv("ChannelAccessToken"),
		)
		if err != nil {
			log.Fatal(err)
		}
	
		// Setup HTTP Server for receiving requests from LINE platform
		
		events, err := bot.ParseRequest(req)
			if err != nil {
				if err == linebot.ErrInvalidSignature {
					w.WriteHeader(400)
				} else {
					w.WriteHeader(500)
				}
				return
			}
			for _, event := range events {
				if event.Type == linebot.EventTypeMessage {
					switch message := event.Message.(type) {
					case *linebot.TextMessage:
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
							log.Print(err)
						}
					case *linebot.StickerMessage:
						replyMessage := fmt.Sprintf(
							"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							log.Print(err)
						}
					}
				}
			}
		
		// This is just sample code.
		// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
		if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
			log.Fatal(err)
		}
	}
