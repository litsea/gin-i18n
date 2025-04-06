package i18n

import (
	"github.com/gin-gonic/gin"
	"github.com/litsea/i18n"
)

const (
	ginI18nContextKey = "litsea.gin-i18n"
	i18nContextKey    = "litsea.i18n"
)

func getGinI18nFromContext(ctx *gin.Context) *I18n {
	v, exists := ctx.Get(ginI18nContextKey)
	if !exists {
		return nil
	}

	i, ok := v.(*I18n)
	if !ok {
		return nil
	}

	return i
}

func getI18nFromContext(ctx *gin.Context) *i18n.I18n {
	v, exists := ctx.Get(i18nContextKey)
	if !exists {
		return nil
	}

	i, ok := v.(*i18n.I18n)
	if !ok {
		return nil
	}

	return i
}
