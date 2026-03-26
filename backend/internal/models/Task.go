package models

type Task struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
}

func (u User) Greet() string {
	return "Hello, " + u.Name
}
