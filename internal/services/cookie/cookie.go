package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

const (
	AgeSecond = 1
	AgeMinute = AgeSecond * 60
	AgeHour   = AgeMinute * 60
)

const (
	UserIdCookieName = "user_id"
	Path             = "/"
	sep              = ":"
)

type UserIdCookie struct {
	UserId string
	Sign   string
}

func GetUserIdCookie(r *http.Request) (*UserIdCookie, error) {
	c, err := r.Cookie(UserIdCookieName)
	if err != nil {
		return nil, err
	}

	return parseUserIdCookie(c.Value), nil
}

func parseUserIdCookie(source string) *UserIdCookie {
	shards := strings.Split(source, sep)
	if len(shards) < 2 {
		return nil
	}

	return &UserIdCookie{UserId: shards[0], Sign: shards[1]}
}

func SetUserId(w http.ResponseWriter, userId, sign string) {
	http.SetCookie(w, &http.Cookie{
		Name:     UserIdCookieName,
		Value:    userId + sep + sign,
		Path:     Path,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   AgeMinute,
	})
}

type SecureCookie struct {
	hashKey []byte
}

func NewSecureCookie(hashKey string) *SecureCookie {
	return &SecureCookie{
		hashKey: []byte(hashKey),
	}
}

func (s *SecureCookie) EncodeString(val string) []byte {
	h := hmac.New(sha256.New, s.hashKey)
	h.Write([]byte(val))

	return h.Sum(nil)
}

func (s *SecureCookie) IsValidUserIdCookie(u *UserIdCookie) (bool, error) {
	sign, err := hex.DecodeString(u.Sign)
	if err != nil {
		return false, err
	}

	return hmac.Equal(sign, s.EncodeString(u.UserId)), nil
}
