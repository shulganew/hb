package entities

// Type action.
type ActionType string

const (
	SUB   ActionType = "S"  // Subscribe to the user.
	UNSUB ActionType = "UN" // Unsub from the user.
)

func (s *ActionType) String() string {
	return string(*s)
}

type Action struct {
	Type       ActionType `json:"t"`
	ChatID     int64      `json:"c"`
	Subscriber string     `json:"sr"` // User, who will be notified
	Subscribed string     `json:"sd"` // User, which Happy B important for Subscriber
}
