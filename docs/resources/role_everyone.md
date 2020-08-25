# Discord Role Everyone Resource

A resource to create a role

## Example Usage

```hcl-terraform
resource discord_role moderator {
    server_id = var.server_id
    permissions = data.discord_permission.moderator.allow_bits
}
```

## Argument Reference

* `server_id` (Required) Which server the role will be in
* `permissions` (Optional) The permission bits of the role
