package discord

import (
    "context"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "gopkg.in/go-playground/colors.v1"
    "strconv"
    "strings"
)

func dataSourceDiscordColor() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceDiscordColorRead,
        Schema: map[string]*schema.Schema{
            "hex": {
                ExactlyOneOf: []string{"hex", "rgb"},
                Type:         schema.TypeString,
                Optional:     true,
            },
            "rgb": {
                ExactlyOneOf: []string{"hex", "rgb"},
                Type:         schema.TypeString,
                Optional:     true,
            },
            "dec": {
                Type:     schema.TypeInt,
                Computed: true,
            },
        },
    }
}

func ConvertToInt(hex string) (int64, error) {
    hex = strings.Replace(hex, "0x", "", 1)
    hex = strings.Replace(hex, "0X", "", 1)
    hex = strings.Replace(hex, "#", "", 1)

    return strconv.ParseInt(hex, 16, 64)
}

func dataSourceDiscordColorRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

    var hex string
    if v, ok := d.GetOk("hex"); ok {
        clr, err := colors.ParseHEX(v.(string))
        if err != nil {
            return diag.Errorf("Failed to parse hex %s: %s", v.(string), err.Error())
        }
        hex = clr.String()
    }
    if v, ok := d.GetOk("rgb"); ok {
        clr, err := colors.ParseRGB(v.(string))
        if err != nil {
            return diag.Errorf("Failed to parse rgb %s: %s", v.(string), err.Error())
        }

        hex = clr.ToHEX().String()
    }

    intColor, err := ConvertToInt(hex)
    if err != nil {
        return diag.Errorf("Failed to parse hex %s: %s", hex, err.Error())
    }

    d.SetId(strconv.Itoa(int(intColor)))
    d.Set("dec", int(intColor))

    return diags
}
