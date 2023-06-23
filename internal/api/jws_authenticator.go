package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/lestrrat-go/jwx/jwt"
)

//const JWTClaimsContextKey = "jwt_claims"

// JWSValidator is used to validate JWS payloads and return a JWT if they're valid.
type JWSValidator interface {
	ValidateJWS(jws string) (jwt.Token, error)
}

// GetJWSFromRequest extracts a JWS string from an Authorization: Bearer <jws> header.
func GetJWSFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no Authorization header found")
	}

	prefix := "BearerAuth "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header is malformed")
	}

	return strings.TrimPrefix(authHeader, prefix), nil
}

func NewAuthenticator(v JWSValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

// Authenticate uses the specified validator to ensure a JWT is valid, then makes
// sure that the claims provided by the JWT match the scopes as required by the API.
func Authenticate(v JWSValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	// Verify security scheme name
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	// Validate JWS
	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	token, err := v.ValidateJWS(jws)
	if err != nil {
		return fmt.Errorf("validating jws: %w", err)
	}

	// Verify claims
	err = CheckTokenClaims(input.Scopes, token)
	if err != nil {
		return fmt.Errorf("token claims don't match: %w", err)
	}

	// TODO: Store token in context, so handlers can use it
	//ctx = context.WithValue(ctx, JWTClaimsContextKey, token) ????

	return nil
}

// GetClaimsFromToken returns a list of claims from the token. We store these
// as a list under the "perms" claim, short for permissions, to keep the token
// shorter.
func GetClaimsFromToken(t jwt.Token) ([]string, error) {
	rawPerms, found := t.Get(PermissionsClaim)
	if !found {
		// The "perms" aren't present in the token, so this means that the token
		// has none. The token is still valid since it passed signature validation.
		return make([]string, 0), nil
	}

	// rawPerms is an untyped JSON list, so it have to be converted to a list of strings
	rawList, ok := rawPerms.([]interface{}) // TODO: not sure what is it?
	if !ok {
		return nil, fmt.Errorf("'%s' claim is unexpected type'", PermissionsClaim)
	}

	claims := make([]string, len(rawList))

	for i, rawClaim := range rawList {
		var ok bool
		claims[i], ok = rawClaim.(string)
		if !ok {
			return nil, fmt.Errorf("%s[%d] is not a string", PermissionsClaim, i)
		}
	}

	return claims, nil
}

func CheckTokenClaims(expectedClaims []string, t jwt.Token) error {
	claims, err := GetClaimsFromToken(t)
	if err != nil {
		return fmt.Errorf("getting claims from token: %w", err)
	}

	claimsMap := make(map[string]bool, len(claims))
	for _, claim := range claims {
		claimsMap[claim] = true
	}

	for _, expected := range expectedClaims {
		if !claimsMap[expected] {
			return errors.New("provided claims do not match expected scopes")
		}
	}

	return nil
}
