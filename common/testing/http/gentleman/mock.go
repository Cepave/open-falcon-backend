package mock

import (
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
	"gopkg.in/h2non/gock.v1"
)

// Plugin exports the mock plugin
var MockPlugin = plugin.NewPhasePlugin("before dial", func(ctx *context.Context, h context.Handler) {
	gock.InterceptClient(ctx.Client)
	h.Next(ctx)
})

// Disable disables the registered mocks.
// It's a shorthand to gock.Disable().
func DisableMock() {
	gock.Disable()
}
