package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/waypoint/sdk/component"
)

// CanAuth returns true if the provided component supports authenticating and
// validating authentication for plugins
func (a *App) CanAuth(comp interface{}) bool {
	_, ok := comp.(component.Authenticator)
	return ok
}

// ValidateAuth validates if the component is properly authenticated. This
// will always return nil if the component doesn't support auth.
func (a *App) ValidateAuth(ctx context.Context, c interface{}) error {
	auth, ok := c.(component.Authenticator)
	if !ok {
		return nil
	}

	_, err := a.callDynamicFunc(ctx,
		a.logger.Named("validate_auth"),
		nil,
		auth,
		auth.ValidateAuthFunc(),
	)
	return err
}

// Auth authenticates a component. This will return an error if the component
// doesn't support auth. If this returns nil, then the auth function succeeded
// but the component itself may still not be authenticated. You must check
// again with ValidateAuth.
func (a *App) Auth(ctx context.Context, c interface{}) (*component.AuthResult, error) {
	auth, ok := c.(component.Authenticator)
	if !ok {
		return nil, fmt.Errorf("does not implement authenticator")
	}

	result, err := a.callDynamicFunc(ctx,
		a.logger.Named("auth"),
		nil,
		auth,
		auth.AuthFunc(),
	)
	if result == nil || err != nil {
		return nil, err
	}

	authresult, ok := result.(*component.AuthResult)
	if !ok {
		authresult = nil
	}

	return authresult, nil
}
