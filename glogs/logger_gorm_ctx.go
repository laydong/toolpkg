package glogs

import (
	"context"
	"net/http"
)

type GormCtx struct {
	context.Context
	request *http.Request
	val     string
}

func (c *GormCtx) String() string {
	return c.val
}

func WithRequest(r *http.Request) context.Context {
	ctx := new(GormCtx)
	ctx.Context = context.Background()
	ctx.request = r
	return ctx
}

func WithValue(value string) context.Context {
	ctx := new(GormCtx)
	ctx.Context = context.Background()
	ctx.val = value
	return ctx
}
