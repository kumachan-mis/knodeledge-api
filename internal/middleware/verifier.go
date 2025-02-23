package middleware

import (
	"context"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/mock_$GOFILE -package=$GOPACKAGE

type UserVerifier interface {
	Verify(ctx context.Context, userId string) *Error
}

type userVerifier struct {
}

func NewUserVerifier() UserVerifier {
	return userVerifier{}
}

func (v userVerifier) Verify(ctx context.Context, userId string) *Error {
	claims, ok := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		return Errorf(VerificationFailurepPanic, "failed to get JWT claims from context")
	}

	customClaims, ok := claims.CustomClaims.(*CustomClaims)
	if !ok {
		return Errorf(VerificationFailurepPanic, "failed to get custom claims from JWT claims")

	}

	if userId != customClaims.Sub {
		return Errorf(AuthorizationError, "user id in request does not match the user id in JWT")
	}

	return nil
}
