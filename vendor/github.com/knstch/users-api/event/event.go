package event

type UserCreated struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type UserResetPassword struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
