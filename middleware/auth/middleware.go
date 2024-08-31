package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/online-bnsp/backend/util"
)

var refreshSecret = []byte("your_refresh_secret_key")
var jwtKey = []byte("your_secret_key") // Replace with your secret key

type Claims struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // Add Role field
	jwt.StandardClaims
}

type claimContextKey struct{}
type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// func GenerateRefreshToken(userID string) (string, error) {
// 	expirationTime := time.Now().Add(7 * 24 * time.Hour)
// 	claims := &RefreshClaims{
// 		UserID: userID,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(refreshSecret)
// }

func GenerateRefreshToken() (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Set refresh token expiration time
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func DecodeRefreshToken(tokenStr string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			util.NewResponse(http.StatusUnauthorized, http.StatusUnauthorized, "Unauthorized!", nil).WriteResponse(w, r)
			return
		}

		// Check if the header starts with "Bearer"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid authorization format", nil).WriteResponse(w, r)
			return
		}

		// Extract the token part after "Bearer "
		token := strings.TrimPrefix(authHeader, "Bearer ")
		identity, err := DecodeJWT(token)
		if err != nil {
			util.NewResponse(http.StatusUnauthorized, http.StatusUnauthorized, "Unauthorized!", nil).WriteResponse(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), claimContextKey{}, identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaim(ctx context.Context) Identity {
	v, _ := ctx.Value(claimContextKey{}).(Identity)
	return v
}

func VerifyRefreshToken(tokenString string) (int32, string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return 0, "", errors.New("invalid token signature")
		}
		return 0, "", errors.New("invalid token")
	}
	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}
	return claims.UserID, claims.Username, nil
}

// Middleware otentikasi
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const ContextKeyRole = "role"
		const ContextUserID = "user_id"
		var ErrInvalidSigningMethod = errors.New("Invalid signing method")

		// Ambil token dari cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Unauthorized - Unable to retrieve token", http.StatusUnauthorized)
			return
		}

		// Ambil token dari nilai cookie
		tokenString := cookie.Value
		if tokenString == "" {
			http.Error(w, "Unauthorized - Invalid token format", http.StatusUnauthorized)
			return
		}

		// Verifikasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing token sesuai
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidSigningMethod
			}
			// Kunci signing untuk validasi token
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		// Ambil klaim dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized - Invalid claims", http.StatusUnauthorized)
			return
		}

		// Ambil role dari klaim dan simpan ke context
		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Unauthorized - Role not found", http.StatusUnauthorized)
			return
		}

		user_id, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Unauthorized - user_id not found", http.StatusUnauthorized)
			return
		}

		// Simpan role ke dalam context
		ctx := context.WithValue(r.Context(), ContextKeyRole, role)
		ctx = context.WithValue(ctx, ContextUserID, int32(user_id))
		r = r.WithContext(ctx)

		// Lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

// Middleware role
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value("role")
			if userRole == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if userRole is present in allowedRoles slice
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					next.ServeHTTP(w, r)
					return // Exit the handler function after successful role check
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}

func ExtractTokenClaims(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			// If there's no token, proceed without claims
			next.ServeHTTP(w, r)
			return
		}

		tokenString := cookie.Value
		var claims Claims
		_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
