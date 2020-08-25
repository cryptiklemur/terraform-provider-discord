package discord

import (
    "fmt"
    "github.com/andersfylling/disgord"
    "strings"
)

func parseTwoIds(id string) (string, string, error) {
    parts := strings.SplitN(id, ":", 2)

    if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
        return "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
    }

    return parts[0], parts[1], nil
}

func getId(v string) disgord.Snowflake {
    return disgord.ParseSnowflakeString(v)
}

func getMinorId(v interface{}) disgord.Snowflake {
    str := v.(string)
    if strings.Contains(str, ":") {
        _, secondId, _ := parseTwoIds(str)

        return getId(secondId)
    }

    return getId(v.(string))
}

func getMajorId(v interface{}) disgord.Snowflake {
    str := v.(string)
    if strings.Contains(str, ":") {
        firstId, _, _ := parseTwoIds(str)

        return getId(firstId)
    }

    return getId(v.(string))
}

func getBothIds(v interface{}) (disgord.Snowflake, disgord.Snowflake, error) {
    firstId, secondId, err := parseTwoIds(v.(string))
    if err != nil {
        return 0, 0, err
    }

    return disgord.ParseSnowflakeString(firstId), disgord.ParseSnowflakeString(secondId), nil
}
