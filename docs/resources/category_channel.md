# Discord Category Channel Resource

A resource to create a Category channel

## Example Usage

```hcl-terraform
resource discord_category_channel chatting {
  name = "Chatting"
  server_id = var.server_id
  position = 0
}
```

## Argument Reference

* `name` (Required) Name of the category
* `server_id` (Required) ID of server this category is in
* `position` (Optional) Position of the channel, 0-indexed

## Attribute Reference

* `id` The ID of the channel