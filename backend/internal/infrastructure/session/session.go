package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// ErrInvalidSession はセッションが無効な場合のエラー
var ErrInvalidSession = errors.New("invalid session")

// Store はセッションストアのインターフェース
type Store interface {
	Get(r *http.Request, name string) (*Session, error)
	Save(w http.ResponseWriter, r *http.Request, name string, session *Session) error
	Delete(w http.ResponseWriter, name string)
}

// Session はセッションデータを保持する
type Session struct {
	Values  map[string]any
	Options *Options
}

// Options はCookieのオプション
type Options struct {
	Path     string
	MaxAge   int
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
}

// CookieStore は署名付きCookieベースのセッションストア
type CookieStore struct {
	secret []byte
}

// NewCookieStore は新しいCookieStoreを作成する
func NewCookieStore(secret []byte) *CookieStore {
	return &CookieStore{
		secret: secret,
	}
}

// Get はリクエストからセッションを取得する
func (s *CookieStore) Get(r *http.Request, name string) (*Session, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		// Cookieが存在しない場合は新しいセッションを返す
		return &Session{
			Values: make(map[string]any),
			Options: &Options{
				Path:     "/",
				MaxAge:   60 * 60 * 24 * 7,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			},
		}, nil
	}

	// Cookieの値をデコード・検証
	session, err := s.decode(cookie.Value)
	if err != nil {
		// デコードに失敗した場合は新しいセッションを返す
		return &Session{
			Values: make(map[string]any),
			Options: &Options{
				Path:     "/",
				MaxAge:   60 * 60 * 24 * 7,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			},
		}, nil
	}

	return session, nil
}

// Save はセッションをCookieに保存する
func (s *CookieStore) Save(w http.ResponseWriter, r *http.Request, name string, session *Session) error {
	encoded, err := s.encode(session)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     session.Options.Path,
		MaxAge:   session.Options.MaxAge,
		HttpOnly: session.Options.HttpOnly,
		Secure:   session.Options.Secure,
		SameSite: session.Options.SameSite,
	})

	return nil
}

// Delete はセッションCookieを削除する
func (s *CookieStore) Delete(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// encode はセッションをエンコードして署名する
func (s *CookieStore) encode(session *Session) (string, error) {
	// セッションデータをJSON化
	data, err := json.Marshal(session.Values)
	if err != nil {
		return "", err
	}

	// Base64エンコード
	encoded := base64.RawURLEncoding.EncodeToString(data)

	// HMAC署名を生成
	signature := s.sign(encoded)

	// 署名付きの値を返す (署名.データ)
	return signature + "." + encoded, nil
}

// decode は署名を検証してセッションをデコードする
func (s *CookieStore) decode(value string) (*Session, error) {
	// 署名とデータを分離
	parts := strings.SplitN(value, ".", 2)
	if len(parts) != 2 {
		return nil, ErrInvalidSession
	}

	signature := parts[0]
	encoded := parts[1]

	// 署名を検証
	expectedSignature := s.sign(encoded)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, ErrInvalidSession
	}

	// Base64デコード
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// JSONをデコード
	var values map[string]any
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}

	return &Session{
		Values: values,
		Options: &Options{
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 7,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		},
	}, nil
}

// sign はHMAC-SHA256で署名を生成する
func (s *CookieStore) sign(data string) string {
	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// GenerateRandomString はランダムな文字列を生成する
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

// NewSession は新しいセッションを作成する
func NewSession() *Session {
	return &Session{
		Values: make(map[string]any),
		Options: &Options{
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 7,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		},
	}
}

// GetString はセッションから文字列を取得する
func (s *Session) GetString(key string) (string, bool) {
	v, ok := s.Values[key]
	if !ok {
		return "", false
	}
	str, ok := v.(string)
	return str, ok
}

// GetInt64 はセッションからint64を取得する
func (s *Session) GetInt64(key string) (int64, bool) {
	v, ok := s.Values[key]
	if !ok {
		return 0, false
	}
	// JSONデコード後はfloat64になる
	switch val := v.(type) {
	case float64:
		return int64(val), true
	case int64:
		return val, true
	case int:
		return int64(val), true
	}
	return 0, false
}

// Set はセッションに値を設定する
func (s *Session) Set(key string, value any) {
	s.Values[key] = value
}

// Delete はセッションから値を削除する
func (s *Session) Delete(key string) {
	delete(s.Values, key)
}

// IsExpired はセッションが期限切れかどうかを確認する
func (s *Session) IsExpired(expiresAtKey string) bool {
	expiresAt, ok := s.GetInt64(expiresAtKey)
	if !ok {
		return true
	}
	return time.Now().Unix() > expiresAt
}
