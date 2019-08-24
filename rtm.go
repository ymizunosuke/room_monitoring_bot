package  main

import (
	"fmt"
	"github.com/nlopes/slack"
	"math/rand"
	"os"
	"strings"
	"log"
	"time"
)


const (
	imgDir = "./images/"
	letsMeow = 0

	// patterns
	patternPicture = "写真"
	patternDailyGraph = "日"
	patternHourlyGraph = "時間"

	// messages
	messageMeow = "なぁぁぁぁぁぁぁーーーーーー"
	messageWait = "少し待つにゃ"
	messageTemperature = "室温は %4.2f℃ だにゃ"
)

func main() {

	slackToken := os.Getenv("SLACK_TOKEN")
	botId := os.Getenv("BOT_ID")

	api := slack.New(slackToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			text := ev.Text

			if strings.HasPrefix(text, botId) {
				var tx [] Temperatures

				// meow once in 10 times
				rand.Seed(time.Now().UnixNano())
				flgMeow := rand.Intn(10)
				if flgMeow == letsMeow {
					rtm.SendMessage(rtm.NewOutgoingMessage(messageMeow, ev.Channel))

				} else if strings.Contains(text, patternHourlyGraph) {
					// post Hourly Graph

					tx = GetTemperatureByTerm(hourly)
					MakeTemperatureGraph(&tx)
					file, err := os.Open(imgDir + "graph.png")
					if err != nil {
						log.Fatal(err)
					}
					_, err = api.UploadFile(
						slack.FileUploadParameters{
							Reader:          file,
							Filename:        "hourly graph",
							Channels:        []string{ev.Channel},
						})
					if err != nil {
						log.Fatal(err)
					}
				} else if strings.Contains(text, patternDailyGraph) {
					// post Daily Graph

					tx = GetTemperatureByTerm(daily)
					MakeTemperatureGraph(&tx)
					file, err := os.Open(imgDir + "graph.png")
					if err != nil {
						log.Fatal(err)
					}
					_, err = api.UploadFile(
						slack.FileUploadParameters{
							Reader:          file,
							Filename:        "daily graph",
							Channels:        []string{ev.Channel},
						})
					if err != nil {
						log.Fatal(err)
					}
				} else if strings.Contains(text, patternPicture) {
					// post Picture

					rtm.SendMessage(rtm.NewOutgoingMessage(messageWait, ev.Channel))
					TakePicture()
					file, err := os.Open(imgDir + "picture.jpg")
					if err != nil {
						log.Fatal(err)
					}
					_, err = api.UploadFile(
						slack.FileUploadParameters{
							Reader:          file,
							Filename:        "picture",
							Channels:        []string{ev.Channel},
						})
					if err != nil {
						log.Fatal(err)
					}
				} else {
					// post room temperature

					t, _, _ := GetLatestTemperature()
					msg := fmt.Sprintf(messageTemperature, t)
					rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
				}
			}

		}
	}
}

