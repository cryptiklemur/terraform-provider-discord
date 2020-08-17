package discord

import (
    "fmt"
    "strings"
)

func parseTwoIds(id string) (string, string, error) {
    parts := strings.SplitN(id, ":", 2)

    if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
        return "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
    }

    return parts[0], parts[1], nil
}
