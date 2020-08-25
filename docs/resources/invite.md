# Discord Invite Resource

A resource to create an Invite for a channel

## Example Usage

```hcl-terraform
resource discord_invite chatting {
    channel_id = var.channel_id
    max_age = 0
}
```

## Argument Reference

* `channel_id` (Required) ID of the channel to create an invite for
* `max_age` (Optional) Age of the invite. 0 for permanent (default 86400)
* `max_uses` (Optional) Max number of uses for the invite. 0 (the default) for unlimited
* `temporary` (Optional) Whether the invite kicks users after the close discord (default false)
* `unique` (Optional) Whether this should create a new invite every time

## Attributes Reference

* `id` / `code` The invite code