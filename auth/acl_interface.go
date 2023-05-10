package auth

type ACL interface {
	GetPermission(id string) []string
}
