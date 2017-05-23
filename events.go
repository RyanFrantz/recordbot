package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Create a dynamic map of channel names to events.
// It's a map of maps (hash of hashes).
// We'll use this to set/lookup if an event is ongoing in a channel and tag
// activity accordingly.
var events = make(map[string]map[string]string)

/* An example
 * events = [
 *   "jurassicpark" = [
 *     "uuid"  = "66D40634-52EB-4379-AA82-F48FE487FEE5",
 *     "name"  = "Loose lizards",
 *     "start" = "1494982840",
 *   ],
 *   "cybertron" = [
 *     "uuid"  = "C0F4D7CB-5B5D-4AB0-B2CD-6CE54CD6C4EF",
 *     "name"  = "Dang Decepticons overprovisioned energon
 *     "start" = "1494982739",
 *   ]
 * ]
 *
 */

/*
for channel := range events {
    fmt.Printf("\n#%s: %s (%s)\n", channel, events[channel]["name"], events[channel]["uuid"])
    fmt.Printf("Started: %s\n", events[channel]["start"])
}
*/

// Shell out to `uuidgen` to generate a UUID.
func Uuid() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to generate UUID: %s\n", err))
	}
	uuid := string(out)
	return strings.Trim(uuid, "\n"), nil
}
