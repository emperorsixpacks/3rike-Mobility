package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionTTL  = 72 * time.Hour
	maxSessions = 3
)

// Redis key helpers
func sessionKey(id string) string        { return "session:" + id }
func userSessionsKey(uid uint) string    { return fmt.Sprintf("user_sessions:%d", uid) }

type authService struct {
	users     domain.UserRepository
	jwtSecret string
	rdb       *redis.Client
	cantonParticipantFingerprint string
}

func newAuthService(users domain.UserRepository, jwtSecret string, rdb *redis.Client, fingerprint string) domain.AuthService {
	return &authService{users: users, jwtSecret: jwtSecret, rdb: rdb, cantonParticipantFingerprint: fingerprint}
}

func (s *authService) Register(ctx context.Context, email, password string, role domain.Role) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &domain.User{Email: email, PasswordHash: string(hash), Role: role}
	return u, s.users.Create(ctx, u)
}

// Login creates a session in Redis, enforces max 3 sessions (evicts oldest), returns JWT.
func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Derive Canton party ID from Keycloak sub + participant fingerprint.
	if u.KeycloakSub != "" && s.cantonParticipantFingerprint != "" && u.CantonPartyID == "" {
		u.CantonPartyID = u.KeycloakSub + "::" + s.cantonParticipantFingerprint
		_ = s.users.Update(ctx, u)
	}

	now := time.Now()
	sess := domain.Session{
		ID:        fmt.Sprintf("%d-%d", u.ID, now.UnixNano()),
		UserID:    u.ID,
		Role:      u.Role,
		CreatedAt: now,
		ExpiresAt: now.Add(sessionTTL),
	}

	if s.rdb != nil {
		if err := s.storeSession(ctx, &sess, u.ID); err != nil {
			return "", "", err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":             u.ID,
		"role":            string(u.Role),
		"session_id":      sess.ID,
		"canton_party_id": u.CantonPartyID,
		"exp":             sess.ExpiresAt.Unix(),
	})
	signed, err := token.SignedString([]byte(s.jwtSecret))
	return signed, sess.ID, err
}

func (s *authService) storeSession(ctx context.Context, sess *domain.Session, userID uint) error {
	data, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	listKey := userSessionsKey(userID)

	// Enforce max sessions: trim oldest if at limit
	ids, _ := s.rdb.LRange(ctx, listKey, 0, -1).Result()
	if len(ids) >= maxSessions {
		oldest := ids[0]
		s.rdb.LPop(ctx, listKey)
		s.rdb.Del(ctx, sessionKey(oldest))
	}

	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, sessionKey(sess.ID), data, sessionTTL)
	pipe.RPush(ctx, listKey, sess.ID)
	pipe.Expire(ctx, listKey, sessionTTL)
	_, err = pipe.Exec(ctx)
	return err
}

func (s *authService) Logout(ctx context.Context, sessionID string, userID uint) error {
	if s.rdb == nil {
		return nil
	}
	s.rdb.Del(ctx, sessionKey(sessionID))
	s.rdb.LRem(ctx, userSessionsKey(userID), 0, sessionID)
	return nil
}

func (s *authService) ListSessions(ctx context.Context, userID uint) ([]domain.Session, error) {
	if s.rdb == nil {
		return nil, nil
	}
	ids, err := s.rdb.LRange(ctx, userSessionsKey(userID), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	sessions := make([]domain.Session, 0, len(ids))
	for _, id := range ids {
		data, err := s.rdb.Get(ctx, sessionKey(id)).Bytes()
		if err != nil {
			continue // expired or deleted
		}
		var sess domain.Session
		if err := json.Unmarshal(data, &sess); err == nil {
			sessions = append(sessions, sess)
		}
	}
	return sessions, nil
}

func (s *authService) RevokeSession(ctx context.Context, sessionID string, userID uint) error {
	return s.Logout(ctx, sessionID, userID)
}
