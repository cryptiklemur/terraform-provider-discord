package discord

import (
    "encoding/json"
    "github.com/andersfylling/disgord"
    "time"
)

type UnmappedEmbed struct {
    Title       string                    `json:"title,omitempty"`       // title of embed
    Description string                    `json:"description,omitempty"` // description of embed
    URL         string                    `json:"url,omitempty"`         // url of embed
    Timestamp   string                    `json:"timestamp,omitempty"`   // timestamp	timestamp of embed content
    Color       int                       `json:"color,omitempty"`       // color code of the embed
    Footer      []*disgord.EmbedFooter    `json:"footer,omitempty"`      // embed footer object	footer information
    Image       []*disgord.EmbedImage     `json:"image,omitempty"`       // embed image object	image information
    Thumbnail   []*disgord.EmbedThumbnail `json:"thumbnail,omitempty"`   // embed thumbnail object	thumbnail information
    Video       []*disgord.EmbedVideo     `json:"video,omitempty"`       // embed video object	video information
    Provider    []*disgord.EmbedProvider  `json:"provider,omitempty"`    // embed provider object	provider information
    Author      []*disgord.EmbedAuthor    `json:"author,omitempty"`      // embed author object	author information
    Fields      []*disgord.EmbedField     `json:"fields,omitempty"`      //	array of embed field objects	fields information
}

func buildEmbed(embedList []interface{}) (*disgord.Embed, error) {
    embedMap := embedList[0].(map[string]interface{})

    var time disgord.Time
    if embedMap["timestamp"].(string) != "" {
        err := time.UnmarshalText([]byte(embedMap["timestamp"].(string)))
        if err != nil {
            return nil, err
        }
    }

    embed := &disgord.Embed{
        Title:       embedMap["title"].(string),
        Description: embedMap["description"].(string),
        URL:         embedMap["url"].(string),
        Color:       embedMap["color"].(int),
        Timestamp:   time,
    }

    if len(embedMap["footer"].([]interface{})) > 0 {
        footerMap := embedMap["footer"].([]interface{})[0].(map[string]interface{})
        embed.Footer = &disgord.EmbedFooter{
            Text:    footerMap["text"].(string),
            IconURL: footerMap["icon_url"].(string),
        }
    }

    if len(embedMap["image"].([]interface{})) > 0 {
        imageMap := embedMap["image"].([]interface{})[0].(map[string]interface{})
        embed.Image = &disgord.EmbedImage{
            URL:    imageMap["url"].(string),
            Width:  imageMap["width"].(int),
            Height: imageMap["height"].(int),
        }
    }

    if len(embedMap["thumbnail"].([]interface{})) > 0 {
        thumbnailMap := embedMap["thumbnail"].([]interface{})[0].(map[string]interface{})
        embed.Thumbnail = &disgord.EmbedThumbnail{
            URL:    thumbnailMap["url"].(string),
            Width:  thumbnailMap["width"].(int),
            Height: thumbnailMap["height"].(int),
        }
    }

    if len(embedMap["video"].([]interface{})) > 0 {
        videoMap := embedMap["video"].([]interface{})[0].(map[string]interface{})
        embed.Video = &disgord.EmbedVideo{
            URL:    videoMap["url"].(string),
            Width:  videoMap["width"].(int),
            Height: videoMap["height"].(int),
        }
    }

    if len(embedMap["provider"].([]interface{})) > 0 {
        providerMap := embedMap["provider"].([]interface{})[0].(map[string]interface{})
        embed.Provider = &disgord.EmbedProvider{
            URL:  providerMap["url"].(string),
            Name: providerMap["name"].(string),
        }
    }

    if len(embedMap["author"].([]interface{})) > 0 {
        authorMap := embedMap["author"].([]interface{})[0].(map[string]interface{})
        embed.Author = &disgord.EmbedAuthor{
            Name:    authorMap["name"].(string),
            URL:     authorMap["url"].(string),
            IconURL: authorMap["icon_url"].(string),
        }
    }

    for _, field := range embedMap["fields"].([]interface{}) {
        fieldMap := field.(map[string]interface{})

        embed.Fields = append(embed.Fields, &disgord.EmbedField{
            Name:   fieldMap["name"].(string),
            Value:  fieldMap["value"].(string),
            Inline: fieldMap["inline"].(bool),
        })
    }

    return embed, nil
}

func unbuildEmbed(embed *disgord.Embed) []interface {} {
    var ret interface {}

    var timestamp string
    if !embed.Timestamp.IsZero() {
        timestamp = embed.Timestamp.Format(time.RFC3339)
    }

    e := &UnmappedEmbed{
        Title:       embed.Title,
        Description: embed.Description,
        URL:         embed.URL,
        Timestamp:   timestamp,
        Color:       embed.Color,
        Fields:      embed.Fields,
    }

    if embed.Footer != nil {
        e.Footer = []*disgord.EmbedFooter{embed.Footer}
    }
    if embed.Image != nil {
        e.Image = []*disgord.EmbedImage{embed.Image}
    }
    if embed.Thumbnail != nil {
        e.Thumbnail = []*disgord.EmbedThumbnail{embed.Thumbnail}
    }
    if embed.Video != nil {
        e.Video = []*disgord.EmbedVideo{embed.Video}
    }
    if embed.Provider != nil {
        e.Provider = []*disgord.EmbedProvider{embed.Provider}
    }
    if embed.Author != nil {
        e.Author = []*disgord.EmbedAuthor{embed.Author}
    }

    j, _ := json.MarshalIndent(e, "", "    ")
    _ = json.Unmarshal(j, &ret)

    return []interface{}{ret}
}
