package configuration

import "github.com/mattermost/mattermost-server/model"

const ChatWithMeToken = "your_chatwithme_token"
const ChatWithMeExtensionUrl = "http://your_corebosserver/your_corebos/notifications.php?type=CWM"
const MatterMostHost = "http://your_mattermost_ip:8065"
const MatterMostAdminUsername = "mattermost_admin_username"
const MatterMostAdminPassword = "mattermost_admin_password"

var ChatWithMeTriggerWords = []string{
	"#ayuda",
	"#busca",
	"#muestra",
	"#actualiza",
	"#edita",
	"#crea",
	"#borra",
	"#ver",
	"#iniciacontador",
	"#paracontador",
	"#registratiempo",
	"#avisame",
	"#lista",
	"#help",
	"#find",
	"#show",
	"#update",
	"#edit",
	"#create",
	"#delete",
	"#see",
	"#starttimer",
	"#stoptimer",
	"#logtime",
	"#remindme",
	"#list",
	"#task",
	"#taskfortime",
	"#taskforproject",
	"#sbsavetime",
	"#time",
}

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Position  string `json:"position"`
	Roles     string `json:"roles"`
}

func (u *User) GetMMUser() model.User {
	return model.User{Username: u.Username,
		Password:  u.Password,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Position:  u.Position,
		Roles:     u.Roles}
}

func (u User) GetUser(mmUser *model.User) User {
	return User{Username: mmUser.Username,
		Id:        mmUser.Id,
		Email:     mmUser.Email,
		FirstName: mmUser.FirstName,
		LastName:  mmUser.LastName,
		Position:  mmUser.Position,
		Roles:     mmUser.Roles}
}
