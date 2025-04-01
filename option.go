package i18n

import (
	"github.com/gin-gonic/gin"
	"github.com/litsea/i18n"
)

type Option func(*I18n)

type GetLngHandler func(ctx *gin.Context) string

func WithOptions(opts ...i18n.Option) Option {
	return func(i *I18n) {
		i.options = opts
	}
}

func WithGetLngHandler(h GetLngHandler) Option {
	return func(i *I18n) {
		i.getLngHandler = h
	}
}

func WithLogger(l Logger) Option {
	return func(i *I18n) {
		i.logger = l
	}
}

func WithDefaultLogger() Option {
	return func(i *I18n) {
		i.logger = defaultLogger{}
	}
}
