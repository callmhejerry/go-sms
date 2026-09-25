package audit

import "github.com/google/uuid"

type Entry struct {
	TenantID     *uuid.UUID
	UserID       *uuid.UUID
	Action       string
	ResourceType string
	ResourceID   *uuid.UUID
	Metadata     map[string]any
	IPAddress    string
	UserAgent    string
	RequestID    string
}
