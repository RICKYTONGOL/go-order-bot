// middleware package implements middlewares
package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dongle/go-order-bot/internal/types"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		var ret types.Middleware

		if len(token) > 0 {
			user, err := ParseJWT(token)
			if err != nil {
				rw.Header().Set("Content-Type", "application/json")

				ret.Success = false

				ret.Error.Msg = "Authorization failed, because of token parsing error"
				json.NewEncoder(rw).Encode(ret)

				return
			}

			if user.ExpiredAt < time.Now().Unix() {
				// rw.WriteHeader(http.StatusUnauthorized)
				rw.Header().Set("Content-Type", "application/json")

				ret.Success = false
				ret.Error.UserId = user.UserId
				ret.Error.Email = user.Email
				ret.Error.ExpiredAt = user.ExpiredAt

				ret.Error.Msg = "Authorization failed, because of expired Date!"
				json.NewEncoder(rw).Encode(ret)

				return
			}

			user_id := r.PostFormValue("user_id")

			if user.UserId != user_id {
				// rw.WriteHeader(http.StatusUnauthorized)
				rw.Header().Set("Content-Type", "application/json")

				ret.Success = false

				ret.Error.Msg = "JwtToken is not correct"
				json.NewEncoder(rw).Encode(ret)

				return
			}
		} else {
			// rw.WriteHeader(http.StatusUnauthorized)
			rw.Header().Set("Content-Type", "application/json")

			ret.Success = false
			ret.Error.Msg = "Authorization failed"
			json.NewEncoder(rw).Encode(ret)

			return
		}
		next.ServeHTTP(rw, r)
	})
}
