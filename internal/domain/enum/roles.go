package enum

import public "github.com/knstch/users-api/public"

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	UnknownRole        Role = "unknown"
	UnverifiedUserRole Role = "unverified_user"
	VerifiedUserRole   Role = "verified_user"
)

func ConvertServiceRoleToPublic(role Role) public.Role {
	switch role {
	case UnknownRole:
		return public.Role_UNKNOWN
	case UnverifiedUserRole:
		return public.Role_UNVERIFIED_USER
	case VerifiedUserRole:
		return public.Role_VERIFIED_USER
	default:
		return public.Role_UNKNOWN
	}
}
