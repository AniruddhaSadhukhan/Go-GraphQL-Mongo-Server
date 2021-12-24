package auth

import (
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/models"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			//Guest User

			//To allow the user to access the protected routes, comment the following line
			// common.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})

			// To stop the user from accessing the protected routes, comment the following line
			setUserName(r, models.GuestUser)
			next.ServeHTTP(w, r)

		} else if token == config.ConfigManager.SecretToken {

			//Internal User (Eg. Other backend services)
			setUserName(r, models.InternalUser)
			next.ServeHTTP(w, r)

		} else {
			//Validate Token

			// If token is valid,set username to the Username extracted from token
			// and proceed to the next middleware
			setUserName(r, models.GuestUser)
			next.ServeHTTP(w, r)
		}

	})
}

func setUserName(r *http.Request, userName string) {
	r.Header.Set("User", userName)
}
