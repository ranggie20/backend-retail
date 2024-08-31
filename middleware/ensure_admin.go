package middleware

import (
	"database/sql"
	"net/http"
	// "github.com/google/uuid"
	// "github.com/online-bnsp/backend/middleware/auth"
	// "github.com/online-bnsp/backend/util"
)

func EnsureAdmin(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// user := auth.GetClaim(r.Context())

			// uid, err := uuid.Parse(user.UserID)
			// if err != nil {
			// 	resp := util.NewResponse(http.StatusUnauthorized, http.StatusUnauthorized, "Invalid user ID", nil)
			// 	resp.WriteResponse(w, r)
			// 	return

			// }

			// _, err = repo.New(db).GetAdminUser(r.Context(), uid)
			// if err != nil {
			// 	resp := util.NewResponse(http.StatusUnauthorized, http.StatusUnauthorized, "Not an admin user", nil)
			// 	resp.WriteResponse(w, r)
			// 	return

			// }

			next.ServeHTTP(w, r)
		})
	}
}
