package auth

import (
    "context"
    "errors"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type contextKey string

const userIDKey contextKey = "userID"

func HashPassword(password string) (string, error) {
    value, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(value), err
}

func CheckPassword(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func CreateToken(userID, secret string) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
        "iat": time.Now().Unix(),
    }
    return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func Middleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            header := r.Header.Get("Authorization")
            if !strings.HasPrefix(header, "Bearer ") {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            raw := strings.TrimPrefix(header, "Bearer ")
            token, err := jwt.Parse(raw, func(token *jwt.Token) (any, error) {
                if token.Method != jwt.SigningMethodHS256 {
                    return nil, errors.New("unexpected signing method")
                }
                return []byte(secret), nil
            })
            if err != nil || !token.Valid {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            userID, _ := claims["sub"].(string)
            if userID == "" {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, userID)))
        })
    }
}

func UserID(ctx context.Context) string {
    value, _ := ctx.Value(userIDKey).(string)
    return value
}
