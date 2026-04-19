package context

import "context"

func NewContext() context.Context {
	ctx := context.Background()

	return ctx
}
