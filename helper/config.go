package helper

import (
	"mattermost-extend/configuration"
)

type Config struct {
	ChatWithMeToken                 string
	ChatWithMeExtensionUrl          string
	MatterMostHost                  string
	MatterMostAdminUsername         string
	MatterMostAdminPassword         string
	ChatWithMeTriggerWords          string
	ChatWithMeTriggerWordsEphemeral string
}

func (c *Config) UpdateConfigurations() {
	configuration.ChatWithMeToken = c.ChatWithMeToken
	configuration.ChatWithMeExtensionUrl = c.ChatWithMeExtensionUrl
	configuration.MatterMostHost = RemoveIfISLast(c.MatterMostHost, "/")
	configuration.MatterMostAdminUsername = c.MatterMostAdminUsername
	configuration.MatterMostAdminPassword = c.MatterMostAdminPassword
	configuration.ChatWithMeTriggerWords = ToArray(c.ChatWithMeTriggerWords, ",")
	configuration.ChatWithMeTriggerWordsEphemeral = ToArray(c.ChatWithMeTriggerWordsEphemeral, ",")
}
