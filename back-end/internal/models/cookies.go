package models 


import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"
	
)

const sessionCookieName = "session_id"


var (
	
	sessions = map[string]int{}
	
	userSessions = map[int]string{}

	sessionsMutex sync.Mutex
)

type CurrentUser struct {
	ID          int
	Pseudo      string
	Email       string
	CreatedAt   string
	Description string
	HasAvatar   bool
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func CreateSession(w http.ResponseWriter, userID int) error {
	sessionID, err := generateSessionID()
	if err != nil {
		return err
	}

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	
	oldSessionID, exists := userSessions[userID]
	if exists {
		delete(sessions, oldSessionID)
	}

	sessions[sessionID] = userID
	userSessions[userID] = sessionID

	cookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	return nil
}


func GetCurrentUser(db *sql.DB, r *http.Request) (*CurrentUser, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}

	sessionsMutex.Lock()
	userID, exists := sessions[cookie.Value]
	sessionsMutex.Unlock()

	if !exists {
		return nil, fmt.Errorf("session invalide")
	}

	var user CurrentUser

	err = db.QueryRow(`SELECT id, pseudo, email, created_at FROM users WHERE id = ?`, userID).Scan(
		&user.ID,
		&user.Pseudo,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	var avatarCount int
	if err = db.QueryRow(`SELECT COUNT(1) FROM user_avatars WHERE user_id = ?`, userID).Scan(&avatarCount); err == nil {
		user.HasAvatar = avatarCount > 0
	}

	return &user, nil
}


func DeleteSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		sessionsMutex.Lock()

		userID, exists := sessions[cookie.Value]
		if exists {
			delete(userSessions, userID)
		}

		delete(sessions, cookie.Value)

		sessionsMutex.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}