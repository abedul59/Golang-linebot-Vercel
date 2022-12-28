package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)


func Main() {
	bot, err := linebot.New(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", CallbackHandler) 
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
	
	
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)
	
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
				// Handle only on text message
				case *linebot.TextMessage:
					// GetMessageQuota: Get how many remain free tier push message quota you still have this month. (maximum 500)
					quota, err := bot.GetMessageQuota().Do()
					if err != nil {
						log.Println("Quota err:", err)
					}
					// message.ID: Msg unique ID
					// message.Text: Msg text
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("msg ID:"+message.ID+":"+"Get:"+message.Text+" , \n OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
						log.Print(err)
					}
	
				// Handle only on Sticker message
				case *linebot.StickerMessage:
					var kw string
					for _, k := range message.Keywords {
						kw = kw + "," + k
					}
	
					outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s  text: %s", message.StickerID, message.PackageID, kw, message.Text)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
