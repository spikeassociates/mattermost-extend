package main

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"regexp"
)

type MMPlugin struct {
	plugin.MattermostPlugin
}

func main() {
	plugin.ClientMain(&MMPlugin{})
}

func (p *MMPlugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {

	r, _ := regexp.Compile("^#(show|edit|view|create) (\\w+)$")
	matches := r.FindStringSubmatch(post.Message)

	if len(matches) > 0 {
		action := matches[1]
		parameter := matches[2]

		broadcast := &model.WebsocketBroadcast{UserId: post.UserId}

		payloadData := map[string]interface{}{
			"action":    action,
			"parameter": parameter,
		}
		p.API.PublishWebSocketEvent("corebos", payloadData, broadcast)
	}

}
