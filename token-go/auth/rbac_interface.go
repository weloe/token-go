package auth

type RBAC interface {
	GetRole(id string) []string
}
