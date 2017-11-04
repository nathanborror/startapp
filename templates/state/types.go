package state

import "time"

const (
	AccountDataType = "run.nathan.{{.Name}}.account"
)

type Node struct {
	Id       string
	Cursor   int
	Created  time.Time
	Modified time.Time
}

func (i *Node) IdentifyID() string     { return i.Id }
func (i *Node) ApplyID(id string)      { i.Id = id }
func (i *Node) ApplyCursor(cursor int) { i.Cursor = cursor }
func (i *Node) ApplyTime(created, modified time.Time) {
	i.Created = created
	i.Modified = modified
}

type Account struct {
	Name     string
	Email    string
	Password string
	IsActive bool
	Status   AccountStatus
	Node
}

func (i *Account) IdentifyType() string { return AccountDataType }

type AccountStatus string

const (
	AccountReader    AccountStatus = "READER"
	AccountModerator               = "MODERATOR"
	AccountSuperuser               = "SUPERUSER"
)

type Session struct {
	Account *Account
	Token   string
}
