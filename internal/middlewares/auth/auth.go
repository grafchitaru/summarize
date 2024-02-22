package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/grafchitaru/summarize/internal/config"
	"net/http"
)

func GenerateToken(userID uuid.UUID, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type UserDataID struct {
	Value uuid.UUID
}

func WithUserCookie(ctx config.HandlerContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil && err != http.ErrNoCookie {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			isBatchByUserID := r.Method == http.MethodGet && r.RequestURI != "/api/user/register" && r.RequestURI != "/api/user/login" && r.RequestURI != "/ping"

			if err == http.ErrNoCookie {
				if isBatchByUserID {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			if err != nil || cookie.Value == "" {
				userID := uuid.New()
				token, _ := GenerateToken(userID, ctx.Config.SecretKey)
				setCookieAuthorization(w, r, token)
			} else {
				_, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
					return []byte(ctx.Config.SecretKey), nil
				})

				if err != nil {
					userID := uuid.New()
					token, _ := GenerateToken(userID, ctx.Config.SecretKey)
					setCookieAuthorization(w, r, token)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func setCookieAuthorization(w http.ResponseWriter, r *http.Request, token string) {
	//nolint:exhaustruct
	cook := &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
	w.Header().Add("Authorization", "Bearer "+token)
	r.Header.Add("Authorization", "Bearer "+token)
	http.SetCookie(w, cook)
	r.AddCookie(cook)
}

func GetUserID(req *http.Request, secretKey string) (string, error) {
	cookie, err := req.Cookie("token")
	if err != nil {
		return "", err
	}
	tokenString := cookie.Value

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	userID, _ := claims["user_id"].(string)
	return userID, nil
}
