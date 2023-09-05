package auth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type OIDCInfo struct {
	JwksURI string `json:"jwks_uri"`
}

type PublicKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	E   string `json:"e"`
	N   string `json:"n"`
	Use string `json:"use"`
	Key rsa.PublicKey
}

var oidcInfo OIDCInfo
var publicKeysMap = make(map[string]PublicKey)

func retrieveFromURL(url string) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	httpClient := common.GetHTTPClient(true)
	resp, err := httpClient.Do(req)

	if err != nil {
		logger.Log.Errorf("Error in retrieving OIDC info: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Error in reading OIDC info: %v", err)
		return nil, err
	}

	return body, nil

}

func RefreshOIDCInfo() {

	body, err := retrieveFromURL(config.Store.Auth.OidcURL + "/.well-known/openid-configuration")
	if err != nil {
		logger.Log.Errorf("Error in retrieving OIDC info: %v", err)
		return
	}

	err = json.Unmarshal(body, &oidcInfo)
	if err != nil {
		logger.Log.Errorf("Error in parsing OIDC configuration info: %v", err)
		return
	}

	jwksURLBody, err := retrieveFromURL(oidcInfo.JwksURI)
	if err != nil {
		logger.Log.Errorf("Error in retrieving OIDC public keys: %v", err)
		return
	}

	var publicKeys struct {
		Keys []PublicKey `json:"keys"`
	}
	err = json.Unmarshal(jwksURLBody, &publicKeys)
	if err != nil {
		logger.Log.Errorf("Error in parsing OIDC public keys: %v", err)
		return
	}

	if len(publicKeys.Keys) == 0 {
		logger.Log.Info("No OIDC public keys found")
		return
	}

	// Clear the old publicKeysMap
	publicKeysMap = make(map[string]PublicKey)

	for _, key := range publicKeys.Keys {
		if key.Use == "sig" && key.Kty == "RSA" {
			key.Key, err = generateRSAPublicKey(key.N, key.E)
			if err != nil {
				logger.Log.Errorf("Error in generating RSA public key: %v", err)
			} else {
				publicKeysMap[key.Kid] = key
			}
		}
	}

	logger.Log.Info("Refreshed OIDC info")
}

func Middleware(next http.Handler) http.Handler {

	if len(oidcInfo.JwksURI) == 0 || len(publicKeysMap) == 0 {
		logger.Log.Info("No OIDC info found")
		RefreshOIDCInfo()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			//Guest User

			//To allow the user to access the protected routes, comment the following line
			common.RespondWithUnauthorized(w)

			// To stop the user from accessing the protected routes, comment the following line
			// r = setUserNameInReq(r, models.GuestUser)
			// next.ServeHTTP(w, r)

		} else if tokenString == config.Store.SecretToken {

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
	return r.WithContext(context.WithValue(r.Context(), models.UserContextKey, userName))
}

func validateToken(tokenString string, next http.Handler, w http.ResponseWriter, r *http.Request) {
	token, err := parseToken(r.Context(), tokenString)
	if err != nil {
		logger.Log.Error("Error while parsing token")
		common.RespondWithUnauthorized(w)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid && claims["sub"] != nil {
		r = setUserNameInReq(r, claims["sub"].(string))
		next.ServeHTTP(w, r)
	} else {
		common.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}
}

func parseToken(ctx context.Context, tokenString string) (*jwt.Token, error) {

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		pubKey, err := validateOIDCToken(ctx, token, tokenString)
		if err != nil {
			pubKey, err = validateJwtInHouse(ctx, token, tokenString)
		}

		return pubKey, err
	})

	return parsedToken, err
}
