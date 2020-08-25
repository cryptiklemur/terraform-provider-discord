# Discord Role Data Source

Fetches a role's information from a server.

## Example Usage

```hcl-terraform
data discord_role mods_id {
    server_id = "81384788765712384"
    role_id   = "175643578071121920"
}
data discord_role mods_name {
    server_id = "81384788765712384"
    name      = "Mods"
}

output mods_color {
    value = data.discord_role.mods_id.color
}
```

## Argument Reference

* `server_id` (Required) The server id to search for the user in
* `role_id` (Optiona) The user id to search for. Either this or `name` is required
* `name` (Optional) The role name to search for. Either this or `role_id` is required

## Attribute Reference

* `id` The id of the role
* `position` Position of the role. This is reverse-indexed. the `@everyone` role is 0
* `color` The integer representation of the role's color
* `permissions` The permission bits of the role
* `hoist` Whether the role is hoisted
* `mentionable` Whether the role is mentionable
* `managed` Whether the role is managed
