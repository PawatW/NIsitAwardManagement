package enums

type UserRole string

const (
	Admin            UserRole = "ADMIN"
	Student          UserRole = "STUDENT"
	HeadOfDepartment UserRole = "HEAD_OF_DEPARTMENT"
	ViceDean         UserRole = "VICE_DEAN"
	Dean             UserRole = "DEAN"
	CommitteeChair   UserRole = "COMMITTEE_CHAIR"
)

// ContextKey is the type for context keys
type ContextKey string

const (
	UserIDContextKey    ContextKey = "user_id"
	UserEmailContextKey ContextKey = "user_email"
	UserRoleContextKey  ContextKey = "user_role"
)
