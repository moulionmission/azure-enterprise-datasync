package auth

import (
	"errors"
	"strings"
)

// Claims represents claims parsed from an Entra ID (Azure AD) JWT token.
type Claims struct {
	TenantID string
	Roles    []string
	Subject  string
}

// EntraIDValidator enforces Zero-Trust tenant boundaries and RBAC permissions.
type EntraIDValidator struct {
	expectedTenantID string
}

// NewEntraIDValidator creates a validator bound to the enterprise tenant ID.
func NewEntraIDValidator(expectedTenantID string) *EntraIDValidator {
	return &EntraIDValidator{
		expectedTenantID: expectedTenantID,
	}
}

// ValidateToken verifies tenant isolation and role existence.
func (v *EntraIDValidator) ValidateToken(claims *Claims, requiredRole string) error {
	if claims == nil {
		return errors.New("empty claims")
	}

	// Enforce strict tenant boundary
	if !strings.EqualFold(claims.TenantID, v.expectedTenantID) {
		return errors.New("unauthorized tenant boundary violation")
	}

	// Enforce RBAC
	for _, role := range claims.Roles {
		if strings.EqualFold(role, requiredRole) {
			return nil
		}
	}

	return errors.New("insufficient role privileges")
}
