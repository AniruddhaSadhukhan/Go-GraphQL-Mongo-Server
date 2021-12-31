package auth

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

func GenerateToken(token *models.Token) error {
	token.CreatedAt = time.Now()
	newToken := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"user":      token.UserName,
		"tokenName": token.TokenName,
		"exp":       token.ExpiresAt.Unix(),
		"iat":       token.CreatedAt.Unix(),
	})

	tokenString, err := newToken.SignedString(getPrivateKey())
	if err != nil {
		logger.Log.Error("Error generating token: " + err.Error())
		return err
	}

	token.TokenString = tokenString
	token.TokenHash = sha256.Sum256([]byte(tokenString))
	return nil
}

func getPrivateKey() *rsa.PrivateKey {

	privateKeyString := config.ConfigManager.JWT_PrivateKey

	if privateKeyString == "" {
		logger.Log.Error("JWT Private Key is not configured")
		return &rsa.PrivateKey{}
	}

	block, _ := pem.Decode([]byte(privateKeyString))

	if block == nil || block.Type != "RSA PRIVATE KEY" {
		logger.Log.Error("Can not decode JWT Private Key")
		return &rsa.PrivateKey{}
	}

	// Parse the key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Log.Error("Can not parse PKCS1 JWT Private Key")
	}
	return privateKey

}

func validateJwtInHouse(token *jwt.Token, tokenString string, ctx context.Context) (interface{}, error) {
	logger.Log.Info("Validating JWT token for Inhouse flow")
	pubKey := getPrivateKey().Public()

	// Check for signing method
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		logger.Log.Error("Invalid JWT token")
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	// Verify token in db
	tokenHash := sha256.Sum256([]byte(tokenString))
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && claims["user"] != nil {
		userName := claims["user"].(string)
		if models.IsExist(
			models.TokenCollection,
			bson.M{
				"userName":  userName,
				"tokenHash": tokenHash,
			},
			ctx,
		) {
			return pubKey, nil
		}
	}

	return nil, fmt.Errorf("invalid token")
}
