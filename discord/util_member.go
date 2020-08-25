package discord

import "github.com/andersfylling/disgord"

func hasRole(member *disgord.Member, roleId disgord.Snowflake) bool {
    for _, r := range member.Roles {
        if r.String() == roleId.String() {
            return true
        }
    }

    return false
}
