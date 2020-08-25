# Discord Color Data Source

A simple helper to get the integer representation of a hex or rgb color

## Example Usage

```hcl-terraform
data discord_color blue {
    hex = "#4287f5"
}

data discord_color green {
  rgb = "rgb(46, 204, 113)"
}

resource discord_role blue {
    // ...
    color = data.discord_color.blue.dec
}
resource discord_role green {
    // ...
    color = data.discord_color.green.dec
}
```

## Argument Reference

* `hex` (Optional) The hex color code. One of these must be present
* `rgb` (Optional) The RGB color (format: `rgb(R, G, B)`). One of these must be present

## Attribute Reference

* `dec` The integer representation of the passed color