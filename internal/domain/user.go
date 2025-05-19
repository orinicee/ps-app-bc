package domain

type Role string

const (
	Editor Role = "Editor"
	Client Role = "client"
)

type User struct {
	ID       string
	Email    string
	Password string
	Role     Role
}
