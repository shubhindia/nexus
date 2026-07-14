package types

type Role string

const (
	RoleSystem    Role = "system"
	RoleDeveloper Role = "developer"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

func ParseRole(role string) Role {
	switch role {
	case "system":
		return RoleSystem
	case "developer":
		return RoleDeveloper
	case "assistant":
		return RoleAssistant
	case "tool":
		return RoleTool
	default:
		return RoleUser
	}
}

func (r Role) String() string {
	return string(r)
}
