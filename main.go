/*
 recordbot

 Archives activity. Optionally labels activity to delineate events such as
 maintenance windows or outages for easier recall later.

 Based on https://github.com/nlopes/slack/blob/master/examples/websocket/websocket.go
*/

package main

import (
    "fmt"
    "log"
    "os"
    "regexp"
    "encoding/json"
    "strconv"
    "time"

    "github.com/nlopes/slack"
)

/*
type Config struct {
    Api_key string
}
*/

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

    type ElasticsearchDocument struct {
        SlackTimestamp         string `json:"slack_timestamp"`
        SlackTimestampISO8601  string `json:"slack_timestamp_iso8601"`
        ChannelID              string `json:"channel_id"`
        ChannelName            string `json:"channel_name"`
        UserID                 string `json:"user_id"`
        UserName               string `json:"user_name"`
        Message                string `json:"message"`
        EventUuid              string `json:"event_uuid"`
    }

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
            channelInfo, _ := rtm.GetChannelInfo(ev.Channel)
            userInfo, _ := rtm.GetUserInfo(ev.User)
            if userInfo.ID != bot_id { // We may not want to respond to our own bot and get in a loop.
                re := regexp.MustCompile("^<@" + bot_id + ">\\s+(\\w+)")
                fmt.Printf("Matches bot command? %q\n", re.FindString(ev.Text))
                // Test for an existing event in eventsByChannel.
                event_uuid, event_exists := eventsByChannel[channelInfo.Name]
                if event_exists {
                    fmt.Println("Ongoing event being tracked!")
                } else {
                    // Generate a UUID to tag messages.
                    if len(re.FindString(ev.Text)) > 0 {
                        event_uuid, err = Uuid()
                        if err != nil {
                            log.Fatal(err)
                        }
                        eventsByChannel[channelInfo.Name] = event_uuid
                    }

                }

                // golang's time doesn't parse epoch strings. Convert to int64 and do some magic.
                intSize := 64
                ts, err := strconv.ParseFloat(ev.Timestamp, intSize)
                if err != nil {
                    fmt.Printf("Unable to convert timestamp '%s' to int64: %s\n", ev.Timestamp, err)
                }
                slack_time := time.Unix(int64(ts), 0)
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
                    slack_time.Format(time.RFC3339),  // ISO 8601
                    channelInfo.ID,
                    channelInfo.Name,
                    userInfo.ID,
                    userInfo.Name,
                    ev.Text,
                    event_uuid,
                }
                // TODO: Actually use err here.
                es_json, _ := json.Marshal(edoc)
                fmt.Println(string(es_json))
                rtm.SendMessage(rtm.NewOutgoingMessage("Recorded " + string(es_json), channelInfo.ID)) // DEBUG
                match_command(ev.Text)
                //}
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
