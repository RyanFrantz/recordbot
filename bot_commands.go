package main

import (
    "fmt"
    "regexp"
)

func is_bot_command(message string) (is_command bool, bot_command string, event_name string, err error) {
    // "@recordbot start outage event"
    // "@recordbot stop"
    // "@recordbot status"
    // "@recordbot wutang   "
    regex_str := `\<@\w+\>\s+(?P<bot_command>\w+)\s*(?P<event_name>.*)?`

    re, err := regexp.Compile(regex_str)
    if err != nil {
        fmt.Printf("regexp: Could not compile regex '%s'\n", regex_str)
        return false, "", "", err // Add detail from the fmt.Printf statement.
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
        return true, matches["bot_command"], matches["event_name"], nil
    } else {
        fmt.Printf("regexp: Failed to match!\n")
        return false, "", "", nil
    }
}
