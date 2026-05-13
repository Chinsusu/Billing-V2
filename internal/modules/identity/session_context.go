package identity

import "context"

type sessionIdentityContextKey struct{}

func WithSessionIdentity(ctx context.Context, identity SessionIdentity) context.Context {
	return context.WithValue(ctx, sessionIdentityContextKey{}, identity)
}

func SessionIdentityFromContext(ctx context.Context) (SessionIdentity, bool) {
	identity, ok := ctx.Value(sessionIdentityContextKey{}).(SessionIdentity)
	return identity, ok
}
