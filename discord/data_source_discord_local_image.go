package discord

import (
    "context"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/polds/imgbase64"
    "hash/crc32"
    "strconv"
)

func Hashcode(s string) int {
    v := int(crc32.ChecksumIEEE([]byte(s)))
    if v >= 0 {
        return v
    }
    if -v >= 0 {
        return -v
    }
    // v == MinInt
    return 0
}

func dataSourceDiscordLocalImage() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceDiscordLocalImageRead,
        Schema: map[string]*schema.Schema{
            "file": {
                Type:     schema.TypeString,
                Required: true,
            },
            "data_uri": {
                Type: schema.TypeString,
                Computed: true,
            },
        },
    }
}

func dataSourceDiscordLocalImageRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

    img, err := imgbase64.FromLocal(d.Get("file").(string))
    if err != nil {
        return diag.Errorf("Failed to process %s: %s", d.Get("file").(string), err.Error())
    }

    d.Set("data_uri", img)
    d.SetId(strconv.Itoa(Hashcode(d.Get("data_uri").(string))))

    return diags
}
