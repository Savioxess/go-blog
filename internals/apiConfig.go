package internals

import (
	"net/http"
	"strings"

	"github.com/Savioxess/blog/internals/utils"
	"github.com/golang-jwt/jwt/v5"
)

type APIConfig struct {
	JWT_SECRET string
}

func (cfg *APIConfig) Handler(next Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.Handle(w, r)
	})
}

func (cfg *APIConfig) GetUserIDFromToken(next Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		claims := jwt.RegisteredClaims{}
		jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT_SECRET), nil
		})

		if err != nil {
			response := &utils.Error{
				Error: "Invalid JSON Token",
			}

			responseJSON, err := utils.EncodeJSONResponse(response)

			if err != nil {
				utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
				return
			}

			utils.ClientErrorResponse(400, responseJSON, r.Method, r.URL.Path, w)
			return
		}

		userId, err := jwtToken.Claims.GetSubject()

		if err != nil {
			utils.ServerErrorResponse(500, "Internal Server Error", r.Method, r.URL.Path, w)
			return
		}

		r.Header.Set("UserID", userId)
		next.Handle(w, r)
	})
}
