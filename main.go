/*
 recordbot

 Archives activity. Optionally labels activity to delineate events such as
 maintenance windows or outages for easier recall later.

 Based on https://github.com/nlopes/slack/blob/master/examples/websocket/websocket.go
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/nlopes/slack"
)

func main() {
	config, err := ReadConfig("recordbot.json")
	if err != nil {
		log.Fatal(err)
	}

	api := slack.New(config.Api_key)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// We'll find our ID when we connect.
	bot_id := ""

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			fmt.Printf("Received 'Hello'!!\n\n")

		case *slack.ConnectedEvent:
			fmt.Println("Connected to Slack!")
			bot_id = ev.Info.User.ID
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			channel_info, _ := rtm.GetChannelInfo(ev.Channel)
			channel_id := channel_info.ID
			channel_name := channel_info.Name
			// TODO: Handle this case and log accordingly to Elasticsearch.
			// TODO: Break out bits of code below into more manageable functions.
			if ev.SubType == "message_changed" {
				fmt.Printf("Received MessageEvent for message which has been changed/edited! Is hidden? [%s] There is likely no user info available.\n", strconv.FormatBool(ev.Hidden))
			}
			user_info, err := rtm.GetUserInfo(ev.User)
			if err != nil {
				fmt.Printf("Unable to get user information: %s\n", err) // 'user_not_found' when ev.SubType == "message_changed"
			} else {
				if user_info.ID != bot_id { // We may not want to respond to our own bot and get in a loop.
					re_bot_request := regexp.MustCompile("^<@" + bot_id + ">\\s+(\\w+)")
					event_uuid := eventsByChannel[channel_name]
					if re_bot_request.MatchString(ev.Text) == true {
						is_command, bot_command, event_name, err := is_bot_command(ev.Text)
						if err != nil {
							// TODO: Log and return from this block.
							fmt.Printf("Failed to check if '%s' is a bot command: %s\n", ev.Text, err)
						}
						if is_command {
							fmt.Printf("COMMAND: '%s'; EVENT: '%s'\n", bot_command, event_name) // DEBUG
							// TODO: Consider logging the bot command as a separate field.
							switch bot_command {
							case "start":
								// Test for an existing event in eventsByChannel.
								if event_uuid != "" {
									rtm.SendMessage(rtm.NewOutgoingMessage("Already tracking an event in this channel", channel_id))
								} else {
									// Generate a UUID to tag messages.
									event_uuid, err = Uuid()
									if err != nil {
										log.Fatal(err)
									}
									// TODO: Let's track the event name as well.
									eventsByChannel[channel_name] = event_uuid
								}
							case "stop":
								eventsByChannel[channel_name] = "" // Clear any event UUID.
							}

						}
					}

					// golang's time doesn't parse epoch strings. Convert to int64 and do some magic.
					intSize := 64
					ts, err := strconv.ParseFloat(ev.Timestamp, intSize)
					if err != nil {
						fmt.Printf("Unable to convert timestamp '%s' to int64: %s\n", ev.Timestamp, err)
					}
					slack_time := time.Unix(int64(ts), 0) // TODO: Address failures above.
					//now := time.Now()
					// Slack often feeds us the last few messages in a channel when
					// we join. Let's ignore those lest we accidentally duplicate
					// stored messages or worse, enables/disables recording an
					// event because it picked up an old command.
					//if slack_time.Unix() < now.Unix() {
					//    fmt.Printf("Slack time is older than now. Ignoring message! (%s)\n", ev.Text)

					//} else {
					edoc := ElasticsearchDocument{
						ev.Timestamp,
						slack_time.Format(time.RFC3339), // ISO 8601
						channel_id,
						channel_name,
						user_info.ID,
						user_info.Name,
						ev.Text,
						event_uuid,
					}
					// TODO: Actually use err here.
					es_json, _ := json.Marshal(edoc)
					fmt.Println(string(es_json))                                                     // DEBUG
					rtm.SendMessage(rtm.NewOutgoingMessage("Recorded "+string(es_json), channel_id)) // DEBUG
				}
			}

		case *slack.PresenceChangeEvent:
			//fmt.Printf("\nPresence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("\nCurrent latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("\nError: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("\nInvalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
