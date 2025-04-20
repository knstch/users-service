package enum

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	UnverifiedUserRole Role = "unverified_user"
)
