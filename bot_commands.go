package main

import (
    "errors"
    "fmt"
    "regexp"
)

func is_bot_command(message string) (is_command bool, bot_command string, event_name string, err error) {
    // "@recordbot start"
    // "@recordbot start outage event"
    // "@recordbot stop"
    // "@recordbot status"
    // "@recordbot wutang   "
    pattern := `\<@\w+\>\s+(?P<bot_command>\w+)\s*(?P<event_name>.*)?`

    re, err := regexp.Compile(pattern)
    if err != nil {
        err := errors.New(fmt.Sprintf("bot_commands.go: Could not compile regex pattern '%s': %s\n", pattern, err))
        return false, "", "", err
    }

    if re.MatchString(message) == true {
        match_result := re.FindStringSubmatch(message)
        // Create a hash of the capture group names to their values.
        names := re.SubexpNames()
        matches := make(map[string]string)
        for index, name := range names {
            if index != 0 {
                matches[name] = match_result[index]
            }
        }
        return true, matches["bot_command"], matches["event_name"], nil
    } else {
        // No match, carry on.
        return false, "", "", nil
    }
}
