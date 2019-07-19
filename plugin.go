package main

import (
	"mattermost-extend/configuration"
	"mattermost-extend/configuration/language"
	"mattermost-extend/helper"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/pkg/errors"
)

type MMPlugin struct {
	plugin.MattermostPlugin
}

func main() {
	plugin.ClientMain(&MMPlugin{})
}

func (p *MMPlugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {

	//Regular expression used for the replacement logic of incoming and outgoing webhooks
	r, _ := regexp.Compile("^\\S+")
	triggerWord := r.FindString(post.Message)

	if helper.Contains(configuration.ChatWithMeTriggerWords, triggerWord) {
		SendPostToChatWithMeExtension(post, triggerWord, p)
	}

	//Regular expression user for special commands like: open, create, edit, list that
	r, _ = regexp.Compile("^#(\\w+) (\\w+)(?: (\\d+))?$")
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

func SendPostToChatWithMeExtension(post *model.Post, triggerWord string, p *MMPlugin) error {

	cnl, _ := p.API.GetChannel(post.ChannelId)

	formData := url.Values{
		"text":         {post.Message},
		"token":        {configuration.ChatWithMeToken},
		"trigger_word": {triggerWord},
		"user_id":      {post.UserId},
		"chnl_name":    {cnl.Name},
		"chnl_dname":   {cnl.DisplayName},
	}

	resp, err := http.PostForm(configuration.ChatWithMeExtensionUrl, formData)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	incomingWebhookPayload, decodeError := model.IncomingWebhookRequestFromJson(resp.Body)
	if decodeError != nil {
		return decodeError
	}

	if len(incomingWebhookPayload.Text) == 0 && incomingWebhookPayload.Attachments == nil {
		return errors.New("Wrong response format")
	}

	newPost := &model.Post{
		UserId:    post.UserId,
		ChannelId: post.ChannelId,
		Type:      model.POST_SLACK_ATTACHMENT,
	}

	if incomingWebhookPayload.Props != nil {
		newPost.Props = incomingWebhookPayload.Props
	}

	newPost.AddProp("attachments", incomingWebhookPayload.Attachments)

	_, err = p.API.CreatePost(newPost)
	if err != nil {
		return err
	}
	return nil
}
