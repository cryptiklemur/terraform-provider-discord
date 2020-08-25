# Discord Member Roles Resource

A resource to manage member roles for a server

## Example Usage

```hcl-terraform
resource discord_member_roles jake {
    user_id = var.user_id
    server_id = var.server_id
    role {
        role_id = var.role_id_to_add
    }
    role {
        role_id = var.role_id_to_always_remove
        has_role = false
    }
}
```

## Argument Reference

* `user_id` (Required) ID of the user to manage roles for
* `server_id` (Required) ID of the server to manage roles in

The **role** blocks have the following arguments:

* `role_id` (Required) The role id to manage
* `has_role` (Optional) Whether the user should have the role

There can be multiple `role` blocks