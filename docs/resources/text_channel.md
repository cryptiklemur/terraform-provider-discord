# Discord Text Channel Resource

A resource to create a text channel

## Example Usage

```hcl-terraform
resource discord_text_channel general {
  name = "general"
  server_id = var.server_id
  position = 0
}
```

## Argument Reference

* `name` (Required) Name of the category
* `server_id` (Required) ID of server this category is in
* `position` (Optional) Position of the channel, 0-indexed
* `topic` (Optional) Topic of the channel
* `nsfw` (Optional) Whether the channel is NSFW
* `category` (Optional) ID of category to place this channel in