package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// test constant variables
var api = &plugintest.API{}
var testPlugin = &MMPlugin{}
var post = &model.Post{
	Id:         "2",
	CreateAt:   34455,
	UpdateAt:   485858,
	IsPinned:   true,
	UserId:     "2",
	ChannelId:  "2",
	RootId:     "2",
	OriginalId: "2",
	Type:       "D",
	Message:    "open",
}
var channel = &model.Channel{
	Id:          "2",
	CreateAt:    34774,
	TeamId:      "2",
	Type:        "D",
	DisplayName: "Test Channel",
	Name:        "chatwithme",
	Purpose:     "test",
}
var user = &model.User{
	Id:       "2",
	CreateAt: 22344,
}

func TestMMPlugin_MessageWillBePosted(t *testing.T) {
	testPlugin.SetAPI(api)
	result, text := testPlugin.MessageWillBePosted(&plugin.Context{}, post)
	require.Equal(t, text, "")
	assert.NotEqual(t, "Posted Ephemeral Trigger Word", result.Message)
}

func TestMMPlugin_MessageHasBeenPosted(t *testing.T) {
	testPlugin.SetAPI(api)
	defer api.AssertExpectations(t)
	testPlugin.MessageHasBeenPosted(&plugin.Context{}, post)
	assert.NoError(t, nil, nil)
}
func TestMMPlugin_OnActivate(t *testing.T) {
	var teams []*model.Team
	team := &model.Team{
		Id:             "2",
		CreateAt:       34455,
		UpdateAt:       485858,
		DisplayName:    "Test Team",
		Name:           "TEST",
		Description:    "API test mockup team",
		Email:          "test@test.com",
		Type:           "D",
		CompanyName:    "test",
		AllowedDomains: "test.com",
	}

	teams = append(teams, team)
	api.On("GetTeams").Return(teams, nil)
	api.On("GetChannelByNameForTeamName", team.Name, "chatwithme", false).Return(channel, nil)
	defer api.AssertExpectations(t)
	testPlugin.SetAPI(api)
	err := testPlugin.OnActivate()
	require.NoError(t, err)
	assert.Equal(t, nil, err)
}

func TestSendPostToChatWithMeExtension(t *testing.T) {
	testPlugin.SetAPI(api)
	defer api.AssertExpectations(t)
	api.On("GetChannel", post.ChannelId).Return(channel, nil)
	api.On("GetUser", post.UserId).Return(user, nil)
	errorPost := &model.Post{
		UserId:    post.UserId,
		ChannelId: post.ChannelId,
		Message:   ":x::x::x: Connection with super-brain is currently not available, please be patient while the universe reorganizes to get back in touch and try in a little while. Thanks! :milky_way:",
	}
	api.On("CreatePost", errorPost).Return(nil, nil)
	//TODO: configure ChatWithMe extension url
	err := SendPostToChatWithMeExtension(post, "create", testPlugin)
	require.NoError(t, err)
	assert.Equal(t, nil, err)
}
