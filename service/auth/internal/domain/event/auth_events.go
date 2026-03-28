package event

import "time"

type AuthEvent struct {
	eventType  string
	userID     int64
	operatorID int64
	data       map[string]interface{}
	occurredAt time.Time
}

func NewAuthEvent(eventType string, userID, operatorID int64, data map[string]interface{}) *AuthEvent {
	return &AuthEvent{
		eventType:  eventType,
		userID:     userID,
		operatorID: operatorID,
		data:       data,
		occurredAt: time.Now(),
	}
}

func (e *AuthEvent) EventType() string { return e.eventType }
func (e *AuthEvent) UserID() int64 { return e.userID }
func (e *AuthEvent) OperatorID() int64 { return e.operatorID }
func (e *AuthEvent) Data() map[string]interface{} { return e.data }
func (e *AuthEvent) OccurredAt() time.Time { return e.occurredAt }

const (
	EventTypeUserLoggedIn     = "user.logged_in"
	EventTypeUserLoggedOut    = "user.logged_out"
	EventTypeUserCreated      = "user.created"
	EventTypeUserUpdated      = "user.updated"
	EventTypePasswordChanged  = "user.password_changed"
	EventTypeRoleAssigned     = "role.assigned"
	EventTypeRoleRemoved      = "role.removed"
	EventTypeUserLocked       = "user.locked"
	EventTypeUserUnlocked     = "user.unlocked"
)

func UserLoggedInEvent(userID int64, ip string) *AuthEvent {
	return NewAuthEvent(EventTypeUserLoggedIn, userID, userID, map[string]interface{}{
		"ip": ip,
	})
}

func UserLoggedOutEvent(userID int64) *AuthEvent {
	return NewAuthEvent(EventTypeUserLoggedOut, userID, userID, nil)
}

func UserCreatedEvent(userID, operatorID int64, username string) *AuthEvent {
	return NewAuthEvent(EventTypeUserCreated, userID, operatorID, map[string]interface{}{
		"username": username,
	})
}

func PasswordChangedEvent(userID int64) *AuthEvent {
	return NewAuthEvent(EventTypePasswordChanged, userID, userID, nil)
}

func RoleAssignedEvent(userID, operatorID, roleID int64) *AuthEvent {
	return NewAuthEvent(EventTypeRoleAssigned, userID, operatorID, map[string]interface{}{
		"role_id": roleID,
	})
}

func RoleRemovedEvent(userID, operatorID, roleID int64) *AuthEvent {
	return NewAuthEvent(EventTypeRoleRemoved, userID, operatorID, map[string]interface{}{
		"role_id": roleID,
	})
}