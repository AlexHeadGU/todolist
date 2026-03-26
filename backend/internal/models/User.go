package models

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
}

func (u User) Greet() string {
	return "Hello, " + u.Name
}
