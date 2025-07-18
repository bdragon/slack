package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// Get tokens from environment variables
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		fmt.Println("SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		fmt.Println("SLACK_APP_TOKEN environment variable is required")
		os.Exit(1)
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(appToken),
	)
	client := socketmode.New(api, socketmode.OptionDebug(true))

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)
					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)
				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.FunctionExecutedEvent:
						callbackID := ev.Function.CallbackID
						if callbackID == "sample_function" {
							userId := ev.Inputs["user_id"]
							payload := map[string]string{
								"user_id": userId.(string),
							}

							err := api.FunctionCompleteSuccess(ev.FunctionExecutionID, slack.FunctionCompleteSuccessRequestOptionOutput(payload))
							if err != nil {
								fmt.Printf("failed posting message: %v \n", err)
							}
						}
					}
				default:
					client.Debugf("unsupported Events API event received\n")
				}

			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()
	client.Run()
}
