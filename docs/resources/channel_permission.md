# Discord Channel Permission Resource

A resource to create a Permission Overwrite for a channel

## Example Usage

```hcl-terraform
resource discord_channel_permission chatting {
    channel_id = var.channel_id
    type = "role"
    overwrite_id = var.role_id
    allow = data.discord_permission.chatting.allow_bits
}
```

## Argument Reference

* `type` (Required) Type of the overwrite, `role` or `user`
* `channel_id` (Required) ID of channel for this overwrite
* `overwrite_id` (Required) ID of user or role for this overwrite
* `allow` (Optional) Permission bits for the allowed permissions on this overwrite. At least one of these two (allow, deny) are required
* `deny` (Optional) Permission bits for the denied permissions on this overwrite. At least one of these two (allow, deny) are required

## Attribute Reference

* `id` Hash of the channel id, overwrite id, and type