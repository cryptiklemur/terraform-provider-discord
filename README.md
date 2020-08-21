# Discord Terraform Provider

## Resources

* discord_channel
* discord_invite
* discord_message
* discord_role
* discord_server
* discord_server_member

## Data

## Todo

* data.discord_permission

    ```hcl-terraform
    data discord_permission allow {
        manage_channel = "allow"
        read_channel = "unset"
        add_reaction = "deny"
    }  
    ```

* resource.discord_channel_permission (Permission Overwrides)
    
    ```hcl-terraform
    resource discord_channel_permission everyone {
        channel_id = discord_channel.rules.id
        type = "role"
        role_id = discord_server.test.everyone_role_id
        allow = data.discord_permissions.allow.bits
        deny = data.discord_permissions.deny.bits
    }
    ```