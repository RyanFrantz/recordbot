package main

import (
    "errors"
    "fmt"
    "os/exec"
    "strings"
)

func Uuid() (string, error) {
    out, err := exec.Command("uuidgen").Output()
    if err != nil {
        return "", errors.New(fmt.Sprintf("Unable to generate UUID: %s\n", err))
    }
    uuid := string(out)
    return strings.Trim(uuid, "\n"), nil
}
