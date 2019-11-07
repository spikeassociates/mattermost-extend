package helper

import "mattermost-extend/configuration"

type Config struct {
	ChatWithMeToken         string
	ChatWithMeExtensionUrl  string
	MatterMostHost          string
	MatterMostAdminUsername string
	MatterMostAdminPassword string
}

func (c *Config) UpdateConfigurations() {
	configuration.ChatWithMeToken = c.ChatWithMeToken
	configuration.ChatWithMeExtensionUrl = c.ChatWithMeExtensionUrl
	configuration.MatterMostHost = c.MatterMostHost
	configuration.MatterMostAdminUsername = c.MatterMostAdminUsername
	configuration.MatterMostAdminPassword = c.MatterMostAdminPassword
}
