package main

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
