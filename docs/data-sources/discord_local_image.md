# Discord Local Image Data Source

A simple helper to get data uri of a local image

## Example Usage

```hcl-terraform
data discord_local_image logo {
    file = "logo.png"
}

resource discord_server server {
    // ...
    icon_data_uri = data.discord_local_image.logo.data_uri
}
```

## Argument Reference

* `file` (Required) The path to the file to process

## Attribute Reference

* `data_uri` The data uri of the `file`