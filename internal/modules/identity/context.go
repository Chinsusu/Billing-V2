package identity

import (
	"context"
	"errors"
)

var ErrActorContextMissing = errors.New("actor context missing")

type actorContextKey struct{}

func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, actorContextKey{}, actor)
}

func FromContext(ctx context.Context) (Actor, bool) {
	actor, ok := ctx.Value(actorContextKey{}).(Actor)
	return actor, ok
}

func RequireActor(ctx context.Context) (Actor, error) {
	actor, ok := FromContext(ctx)
	if !ok {
		return Actor{}, ErrActorContextMissing
	}
	if err := actor.Validate(); err != nil {
		return Actor{}, err
	}
	return actor, nil
}
