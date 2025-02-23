package middleware

import (
	"context"
	"fmt"
)

type CustomClaims struct {
	Sub string `json:"sub"`
}

func (claims *CustomClaims) Validate(ctx context.Context) error {
	if claims.Sub == "" {
		return fmt.Errorf("sub is required")
	}
	return nil
}
