# Discord Member Data Source

Fetches a member's information from a server.

## Example Usage

```hcl-terraform
data discord_member jake {
    server_id = "81384788765712384"
    user_id   = "103559217914318848"
}

output jakes_username_and_discrim {
    value = "${data.discord_member.jake.username}#${data.discord_member.jake.discriminator}"
}
```

## Argument Reference

* `server_id` (Required) The server id to search for the user in
* `user_id` (Optional) The user id to search for. Required if not searching by username/discriminator
* `username` (Optional) The username to search for. Discriminator is required when using this
* `discriminator` (Optional) The discriminator to search for. Username is required when using this

## Attribute Reference

* `id` The user's id
* `joined_at` The time at which the user joined
* `premium_since` The time at which the user became premium
* `username` The username of the user
* `discriminator` The discriminator (#0000) of the user
* `nick` The current nickname of the user
* `avatar` The avatar hash of the user
* `roles` Array of role ids that the user has
* `in_server` Bool of whether or not the user is in the server
