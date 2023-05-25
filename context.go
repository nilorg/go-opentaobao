package opentaobao

import "context"

type sessionKey struct{}

func NewSessionContext(parent context.Context, session string) context.Context {
	return context.WithValue(parent, sessionKey{}, session)
}

func FromSessionContext(ctx context.Context) (session string, ok bool) {
	session, ok = ctx.Value(sessionKey{}).(string)
	return
}
