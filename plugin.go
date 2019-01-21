package main

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/pkg/errors"
	"mattermos-extend/configuration/language"
	"regexp"
	"strings"
)

type MMPlugin struct {
	plugin.MattermostPlugin
}

func main() {
	plugin.ClientMain(&MMPlugin{})
}

func (p *MMPlugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {

	r, _ := regexp.Compile("^#(\\w+) (\\w+)(?: (\\d+))?$")
	matches := r.FindStringSubmatch(strings.TrimSpace(post.Message))

	if len(matches) > 0 {

		if action, ok := language.Command[matches[1]]; ok {

			module := matches[2]

			broadcast := &model.WebsocketBroadcast{UserId: post.UserId}

			payloadData := map[string]interface{}{
				"action": action,
				"module": module,
			}

			if matches[3] != "" {
				payloadData["id"] = matches[3]
			}

			p.API.PublishWebSocketEvent("corebos", payloadData, broadcast)
		}
	}

}

func (p *MMPlugin) OnActivate() error {

	teams, err := p.API.GetTeams()
	if err != nil {
		return err
	}

	if len(teams) == 0 {
		return errors.New("there are no existing teams")
	}

	team := teams[0]
	channel, _ := p.API.GetChannelByNameForTeamName(team.Name, "chatwithme", false)

	if channel == nil {

		channel, err = p.API.CreateChannel(&model.Channel{
			TeamId:      team.Id,
			Type:        model.CHANNEL_OPEN,
			DisplayName: "Chat With Me",
			Name:        "chatwithme",
			Header:      "The channel used by the mattermost-extend plugin.",
			Purpose:     "The channel was created by the mattermost-extend plugin to extend the server functionality.",
		})

		if err != nil {
			return err
		}

	}

	return nil
}
