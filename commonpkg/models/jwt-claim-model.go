package models

import (
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	AuthHeaderName      = "Authorization"
	ContextKey_BranchId = "branchId"
	ContextKey_Username = "username"
	ContextKey_Roles    = "roles"
)

type JwtClaims struct {
	BranchId string `json:"branchId,omitempty"`
	Username string `json:"username,omitempty"`
	Roles    []int  `json:"roles,omitempty"`
	jwt.StandardClaims
}

func (claims JwtClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true) {
		return nil
	}
	return fmt.Errorf("Token is invalid")
}

func (claims JwtClaims) VerifyAudience(origin string) bool {
	return strings.Compare(claims.Audience, origin) == 0
}
