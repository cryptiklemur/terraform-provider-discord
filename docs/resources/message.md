# Discord Message Resource

A resource to create a message

## Example Usage

### Content Example

```hcl-terraform
resource discord_message hello_world {
    channel_id = var.channel_id
    content = "hello world"
}
```

### Embed Example

```hcl-terraform
resource discord_message hello_world {
    channel_id = var.channel_id
    embed {
        title = "Hello World"
        footer {
            text = "I'm awesome"
        }
        
        fields {
            name = "foo"
            value = "bar"
            inline = true
        }
        
        fields {
            name = "bar"
            value = "baz"
            inline = false
        }
    }
}
```

## Argument Reference

* `channel_id` (Required) Which channel the message will be in
* `content` (Optional) Text content of message. Either this or embed (or both) must be set
* `tts` (Optional) Whether this message triggers tts (default false)
* `embed` (Optional) An embed block (detailed below). There can only be one of these. Either this or content (or both) must be set
* `pinned` (Optional) Whether this message is pinned (default false)

The **embed** block has the following arguments:

Details on arguments can be found [here](https://discord.com/developers/docs/resources/channel#message-object)

* `title`
* `description`
* `url`
* `timestamp`
* `color`
* `footer` (only one allowed)
    * `text`
    * `icon_url`
* `image` (only one allowed)
    * `url`
    * `height`
    * `width`
* `thumbnail` (only one allowed)
    * `url`
    * `height`
    * `width`
* `video` (only one allowed)
    * `url`
    * `height`
    * `width`
* `provider` (only one allowed)
    * `name`
    * `url`
* `author` (only one allowed)
    * `name`
    * `url`
    * `icon_url`
* `fields` (multiple allowed)
    * `name`
    * `value`
    * `inline`

## Attribute Reference

* `server_id` ID of the server this message is in
* `author` ID of the user who wrote the message
* `timestamp` When the message was sent
* `edited_timestamp` When the message was edited