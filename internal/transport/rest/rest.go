package rest

import (
	"github.com/gorilla/sessions"
	"github.com/kiloMIA/documed/internal/service"
	"go.uber.org/zap"
)

type Transport struct {
	Service      *service.Service
	SessionStore *sessions.CookieStore
	Logger       *zap.Logger
}

func NewTransport(
	service *service.Service,
	sessionStore *sessions.CookieStore,
	logger *zap.Logger,
) *Transport {
	return &Transport{
		Service:      service,
		SessionStore: sessionStore,
		Logger:       logger,
	}
}
