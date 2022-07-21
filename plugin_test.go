package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMessageWillBePosted(t *testing.T) {

}
func TestOnActivate(t *testing.T) {
	var teams []model.Team
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
	channel := &model.Channel{
		Id:          "2",
		CreateAt:    34774,
		TeamId:      "2",
		Type:        "D",
		DisplayName: "Test Channel",
		Name:        "chatwithme",
		Purpose:     "test",
	}
	teams = append(teams, *team)
	api := &plugintest.API{}
	api.On("GetTeams").Return(teams, nil)
	api.On("GetChannelByNameForTeamName", team.Name, "chatwithme", false).Return(channel, nil)
	defer api.AssertExpectations(t)
	testPlugin := &MMPlugin{}
	testPlugin.SetAPI(api)
	err := testPlugin.OnActivate()
	require.NoError(t, err)
	assert.Equal(t, nil, err)
}
