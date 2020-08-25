# Discord Role Resource

A resource to create a role

## Example Usage

```hcl-terraform
resource discord_role moderator {
    server_id = var.server_id
    name = "Moderator"
    permissions = data.discord_permission.moderator.allow_bits
    color = data.discord_color.blue.dec
    hoist = true
    mentionable = true
    position = 5
}
```

## Argument Reference

* `server_id` (Required) Which server the role will be in
* `name` (Required) The name of the role
* `permissions` (Optional) The permission bits of the role
* `color` (Optional) The integer representation of the role color
* `hoist` (Optional) Whether the role should be hoisted (default false)
* `mentionable` (Optional) Whether the role should be mentionable (default false)
* `position` (Optional) The position of the role. This is reverse indexed (@everyone is 0)

## Attribute Reference

* `managed` Whether this role is managed by another service
