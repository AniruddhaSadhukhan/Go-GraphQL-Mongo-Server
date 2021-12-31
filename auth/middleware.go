package auth

import (
	"context"
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			//Guest User

			//To allow the user to access the protected routes, comment the following line
			common.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})

			// To stop the user from accessing the protected routes, comment the following line
			// r = setUserNameInReq(r, models.GuestUser)
			// next.ServeHTTP(w, r)

		} else if tokenString == config.ConfigManager.SecretToken {

			//Internal User (Eg. Other backend services)
			r = setUserNameInReq(r, models.InternalUser)
			next.ServeHTTP(w, r)

		} else {
			//Validate Token
			validateToken(tokenString, next, w, r)
		}

	})
}

func setUserNameInReq(r *http.Request, userName string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "User", userName))
}

func validateToken(tokenString string, next http.Handler, w http.ResponseWriter, r *http.Request) {
	token, err := parseToken(tokenString, r.Context())
	if err != nil {
		logger.Log.Error("Error while parsing token")
		common.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid && claims["user"] != nil {
		r = setUserNameInReq(r, claims["user"].(string))
		next.ServeHTTP(w, r)
	} else {
		common.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}
}

func parseToken(tokenString string, ctx context.Context) (*jwt.Token, error) {
	// Add other methods of token validation here eg. OIDC etc.
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return validateJwtInHouse(token, tokenString, ctx)
	})
	return parsedToken, err
}
