package state

import "time"

const (
	AccountDataType = "run.nathan.{{.Name}}.account"
)

// Account represents an Account object.
type Account struct {
	Name     string
	Email    string
	Password string
	IsActive bool
	Status   AccountStatus
	Node
}

// Accounts represents a collection of Accounts.
type Accounts struct {
	Results []Account
	PageInfo
}

type AccountStatus string

const (
	AccountReader    AccountStatus = "READER"
	AccountModerator               = "MODERATOR"
	AccountSuperuser               = "SUPERUSER"
)

// Node represents an abstract Node.
type Node struct {
	Id       string
	Cursor   int
	Created  time.Time
	Modified time.Time
}

// Edge represents a relationship between two Nodes.
type Edge struct {
	FromId string
	ToId   string
	Kind   string
}

// PageInfo represents paging information.
type PageInfo struct {
	Total       int
	HasNext     bool
	HasPrevious bool
	StartID     string
	EndID       string
}

// Conformance

func (i *Node) IdentifyID() string     { return i.Id }
func (i *Node) ApplyID(id string)      { i.Id = id }
func (i *Node) ApplyCursor(cursor int) { i.Cursor = cursor }
func (i *Node) ApplyTime(created, modified time.Time) {
	i.Created = created
	i.Modified = modified
}

func (i *Account) IdentifyType() string { return AccountDataType }
