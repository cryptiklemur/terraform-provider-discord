# Discord Voice Channel Resource

A resource to create a voice channel

## Example Usage

```hcl-terraform
resource discord_voice_channel general {
  name = "General"
  server_id = var.server_id
  position = 0
}
```

## Argument Reference

* `name` (Required) Name of the category
* `server_id` (Required) ID of server this category is in
* `position` (Optional) Position of the channel, 0-indexed
* `bitrate` (Optional) Bitrate of the channel
* `userlimit` (Optional) User Limit of the channel
* `category` (Optional) ID of category to place this channel in