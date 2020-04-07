package model

type role int

const (
	guest role = iota
	admin
)

type User struct {
	Name      string
	Email     string
	Password  string
	Lastname  string
	Role      role
	AvatarURL *string
	Token     string
	TokenExp  string
}