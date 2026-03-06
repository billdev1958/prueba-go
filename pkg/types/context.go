package types

import "context"

type contextKey string

const (
	ActorKey contextKey = "actor"
)

// GetActor extracts the actor from context or returns "system_unknown"
func GetActor(ctx context.Context) string {
	if actor, ok := ctx.Value(ActorKey).(string); ok {
		return actor
	}
	return "system_unknown"
}
