package acl

// Action represents a MediaMTX action type
type Action string

const (
	ActionPublish  Action = "publish"
	ActionRead     Action = "read"
	ActionPlayback Action = "playback"
)

// AuthRequest represents an authorization request from MediaMTX
type AuthRequest struct {
	Action   Action
	Path     string
	Protocol string
	IP       string
	User     string
	Password string
	Token    string
}

// AuthResult represents the result of an authorization check
type AuthResult struct {
	Allowed bool
	Reason  string
}
