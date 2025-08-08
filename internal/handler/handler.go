package handler

import (
	"go.uber.org/fx"
)

type Handlers struct {
	Hello *HelloHandler
	User  *UserHandler
}

func NewHandlers(hello *HelloHandler, user *UserHandler) *Handlers {
	return &Handlers{
		Hello: hello,
		User:  user,
	}
}

var Module = fx.Module("handlers",
	fx.Provide(
		NewHelloHandler,
		NewUserHandler,
		NewHandlers,
	),
)
