package event

type UserCreated struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
