package stores

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/zoobzio/rocco/session"
)

// Sessions provides database-backed session and OAuth state storage.
type Sessions struct {
	db *sqlx.DB
}

// NewSessions creates a new sessions store.
func NewSessions(db *sqlx.DB) *Sessions {
	return &Sessions{db: db}
}

// CreateState persists a CSRF state token.
func (s *Sessions) CreateState(ctx context.Context, state string) error {
	query := `INSERT INTO oauth_states (state, expires_at, created_at) VALUES ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, state, time.Now().Add(10*time.Minute), time.Now())
	return err
}

// VerifyState checks that a state token exists and is valid, then deletes it.
func (s *Sessions) VerifyState(ctx context.Context, state string) (bool, error) {
	query := `DELETE FROM oauth_states WHERE state = $1 AND expires_at > $2 RETURNING id`
	var id int64
	err := s.db.QueryRowContext(ctx, query, state, time.Now()).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Create persists a new session.
func (s *Sessions) Create(ctx context.Context, id string, data session.Data) error {
	meta, err := json.Marshal(data.Meta)
	if err != nil {
		return err
	}
	query := `INSERT INTO sessions (id, user_id, tenant_id, email, scopes, roles, meta, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = s.db.ExecContext(ctx, query, id, data.UserID, data.TenantID, data.Email,
		pq.StringArray(data.Scopes), pq.StringArray(data.Roles), meta,
		time.Now().Add(24*time.Hour))
	return err
}

// Get retrieves session data by ID.
func (s *Sessions) Get(ctx context.Context, id string) (*session.Data, error) {
	query := `SELECT user_id, tenant_id, email, scopes, roles, meta FROM sessions
		WHERE id = $1 AND expires_at > $2`
	var (
		data   session.Data
		scopes pq.StringArray
		roles  pq.StringArray
		meta   []byte
	)
	err := s.db.QueryRowContext(ctx, query, id, time.Now()).Scan(
		&data.UserID, &data.TenantID, &data.Email, &scopes, &roles, &meta)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("session not found")
	}
	if err != nil {
		return nil, err
	}
	data.Scopes = scopes
	data.Roles = roles
	if len(meta) > 0 {
		_ = json.Unmarshal(meta, &data.Meta)
	}
	return &data, nil
}

// Refresh extends the session's expiry.
func (s *Sessions) Refresh(ctx context.Context, id string) error {
	query := `UPDATE sessions SET expires_at = $1 WHERE id = $2`
	res, err := s.db.ExecContext(ctx, query, time.Now().Add(24*time.Hour), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("session not found")
	}
	return nil
}

// Delete removes a session by ID.
func (s *Sessions) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// Cleanup removes expired states and sessions.
func (s *Sessions) Cleanup(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM oauth_states WHERE expires_at < $1`, time.Now())
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < $1`, time.Now())
	return err
}
