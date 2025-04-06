package i18n

import (
	"github.com/gin-gonic/gin"
	"github.com/litsea/i18n"
)

func (i *I18n) Localize() gin.HandlerFunc {
	if i.logger != nil {
		i.options = append(i.options, i18n.WithLogger(i.logger))
	}

	i18 := i18n.New(i.options...)

	return func(ctx *gin.Context) {
		ctx.Set(i18nContextKey, i18)
		ctx.Set(ginI18nContextKey, i)

		ctx.Next()
	}
}
