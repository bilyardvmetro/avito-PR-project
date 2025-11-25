package domain

type UserRepository interface {
	UpsertUser(user *User) error
	GetUser(id string) (*User, error)
	UpdateUser(user *User) error
	GetUsersByTeam(teamName string) ([]*User, error)
}
