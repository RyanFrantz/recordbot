package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Create a dynamic map of channel names to event UUID.
// We'll use this to set/lookup if an event is ongoing in a channel and tag
// activity accordingly.
var eventsByChannel = make(map[string]string)

func Uuid() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to generate UUID: %s\n", err))
	}
	uuid := string(out)
	return strings.Trim(uuid, "\n"), nil
}
