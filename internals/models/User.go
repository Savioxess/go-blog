package models

type User struct {
	ID       []byte `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	JoinedOn string	`json:"joined_on"`
}
