package models

import (
	"time"
)

type Token struct {
	TokenName   string    `json:"tokenName" bson:"tokenName"`
	TokenString string    `json:"token" bson:"-"` // Token is not stored in the database
	TokenHash   [32]byte  `json:"tokenHash" bson:"tokenHash"`
	UserName    string    `json:"userName" bson:"userName"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt" bson:"expiresAt"`
}
