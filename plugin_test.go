package main

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"mattermost-extend/helper"
	"net/http/httptest"
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
var helperUser = &helper.User{
	Id:        "2",
	Username:  "test",
	Password:  "test",
	Email:     "test@gmail.com",
	FirstName: "test",
	LastName:  "test",
	Position:  "test",
	Roles:     "test",
	TeamNames: "chatwithme",
}
var team = &model.Team{
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
func TestMMPlugin_ServeHTTP(t *testing.T) {
	testPlugin.SetAPI(api)
	defer api.AssertExpectations(t)
	w := httptest.NewRecorder()
	routes := []string{"/syncuser", "/health", "/postmessage"}

	for _, route := range routes {
		// r.Header.body
		switch route {
		case "/syncuser":
			tests := []struct {
				name        string
				expected    interface{}
				body        *helper.User
				description string
			}{
				{
					name:        "When Body Nil",
					expected:    "Error Decoding Json user",
					description: "the function will return nothing in the handler",
					body:        nil,
				},
				{
					name:        "When Body Set",
					expected:    nil,
					description: "the function will return nothing in the handler",
					body:        helperUser,
				},
			}
			for _, test := range tests {
				r := httptest.NewRequest("POST", route, nil)
				r.Body.Read([]byte(fmt.Sprintf("%v", test.body)))
				api.On("GetUserByUsername", helperUser.Username).Return(user, nil)
				api.On("GetUserByEmail", helperUser.Email).Return(user, nil)
				api.On("CreateUser", user).Return(user, nil)
				testPlugin.ServeHTTP(&plugin.Context{}, w, r)
				body, err := ioutil.ReadAll(w.Result().Body)
				require.NoError(t, err)
				assert.NotEqual(t, test.expected, string(body))

			}
		// do something
		case "/health":
		// check health
		case "/postmessage":
		// post
		default:

		}

	}
}

func TestHandleHealth(t *testing.T) {
	testPlugin.SetAPI(api)
	defer api.AssertExpectations(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	testPlugin.handleHealth(w, r)
	body, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	assert.Equal(t, "{\"success\":true,\"data\":{\"information\":\"\",\"message\":\"Spike Mattermost Corebos Server Plugin is running ...\",\"status\":200}}", string(body))
}
func TestAddTeam(t *testing.T) {
	var teams []*model.Team
	teams = append(teams, team)
	testPlugin.SetAPI(api)
	defer api.AssertExpectations(t)
	w := httptest.NewRecorder()
	api.On("GetTeams").Return(teams, nil)
	addTeam(testPlugin, w, *user, *helperUser)
	assert.NoError(t, nil, nil)
}

func TestContains(t *testing.T) {
	list := []string{"create", "list", "update"}
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "positive",
			input:    "create",
			expected: true,
		},
		{
			name:     "negative",
			input:    "nothing",
			expected: false,
		},
	}
	for _, testCase := range testCases {
		result := helper.Contains(list, testCase.input)
		assert.Equal(t, testCase.expected, result)
	}
}
func TestRemoveIfLast(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "positive",
			input:    "e",
			expected: "mileAg",
		},
		{
			name:     "negative",
			input:    "mil",
			expected: "mileAge",
		},
	}
	for _, testCase := range testCases {
		result := helper.RemoveIfISLast("mileAge", testCase.input)
		assert.Equal(t, testCase.expected, result)
	}
}
