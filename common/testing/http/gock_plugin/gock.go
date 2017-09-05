// Plugin for using Gock on Gentleman library.
//
// GockPlugin
//
// This variable is a plugin of gentleman, could be fed into "Use(plugin.Plugin)" of
// "*gentleman.Client" or "*gentleman.Request" object.
//
// Why
//
// Because https://github.com/h2non/gentleman-mock is still depending on "gentleman.v1" library,
// so we provide this package to make it support "gentleman.v2".
package gock_plugin

import (
	"gopkg.in/h2non/gock.v1"

	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// Supports Gock in gentleman framework.
//
// 	gentleman.Client.Use(gock_plugin.GockPlugin)
var GockPlugin = plugin.NewPhasePlugin("before dial", func(ctx *context.Context, h context.Handler) {
	gock.InterceptClient(ctx.Client)
	h.Next(ctx)
})
