package main

import (
    "fmt"
    "regexp"
)

func match_command(message string) {
    // "bot start outage event"
    // "bot stop"
    // "bot status"
    // "bot wutang   "
    regex_str := `bot\s+(?P<bot_command>\w+)\s*(?P<event_name>.*)?`

    re, err := regexp.Compile(regex_str)
    if err != nil {
        fmt.Printf("regexp: Could not compile regex '%s'\n", regex_str)
    }

    fmt.Printf("regexp: Received message '%s'\n", message)
    if re.MatchString(message) == true {
        fmt.Printf("regexp: Matched '%s'\n", message)
        match_result := re.FindStringSubmatch(message)
        names := re.SubexpNames()
        matches := make(map[string]string)
        for index, name := range names {
            if index != 0 {
                matches[name] = match_result[index]
            }
        }
        fmt.Printf("regexp: Bot command = '%s'\n", matches["bot_command"])
        fmt.Printf("regexp: Event name = '%s'\n", matches["event_name"])
    } else {
        fmt.Printf("regexp: Failed to match!\n")
    }
}
