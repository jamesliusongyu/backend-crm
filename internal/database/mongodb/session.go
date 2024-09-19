package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Session represents a user session with a unique session ID, the associated user's email, and expiration details.
type Session struct {
	SessionID string    `json:"session_id"`
	Email     string    `json:"email"`
	JWTToken  string    `bson:"jwt_token"`
	Tenant    string    `json:"tenant"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionResponse represents the response structure for a session.
type SessionResponse struct {
	ID        string    `bson:"_id,omitempty"`
	SessionID string    `json:"session_id"`
	Email     string    `json:"email"`
	JWTToken  string    `bson:"jwt_token"`
	Tenant    string    `json:"tenant"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionCollection provides methods to interact with the session collection in MongoDB.
type SessionCollection struct {
	*GenericCollection[Session, SessionResponse]
}

// NewSessionCollection creates a new SessionCollection.
func NewSessionCollection(collection *mongo.Collection) *SessionCollection {
	return &SessionCollection{
		GenericCollection: NewGenericCollection[Session, SessionResponse](collection),
	}
}

// Ensure SessionCollection satisfies the Collection interface.
var _ Collection[Session, SessionResponse] = (*SessionCollection)(nil)

// Create stores a new session in the collection.
func (r *SessionCollection) Create(ctx context.Context, entity Session) (string, error) {
	session := Session{
		SessionID: entity.SessionID,
		Email:     entity.Email,
		JWTToken:  entity.JWTToken,
		Tenant:    entity.Tenant,
		CreatedAt: entity.CreatedAt,
		ExpiresAt: entity.ExpiresAt,
	}
	return r.GenericCollection.Create(ctx, session)
}

func (r *SessionCollection) Update(ctx context.Context, id string, entity Session) error {
	session := Session{
		SessionID: entity.SessionID,
		Email:     entity.Email,
		JWTToken:  entity.JWTToken,
		Tenant:    entity.Tenant,
		CreatedAt: entity.CreatedAt,
		ExpiresAt: entity.ExpiresAt,
	}
	return r.GenericCollection.Update(ctx, id, session)
}

// GetBySessionID retrieves a session by its session ID.
func (r *SessionCollection) GetBySessionID(ctx context.Context, sessionID string, tenant string) (SessionResponse, error) {
	return r.GenericCollection.GetByKeyValue(ctx, "session_id", sessionID, tenant)
}

// Delete removes a session from the collection.
func (r *SessionCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
