package i18n

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/litsea/i18n"
	"golang.org/x/text/language"
)

type I18n struct {
	getLngHandler GetLngHandler
	options       []i18n.Option
	logger        Logger
}

func New(opts ...Option) *I18n {
	i := &I18n{
		getLngHandler: defaultGetLngHandler,
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func T(ctx *gin.Context, msgID string, tplData ...map[any]any) string {
	gi := ctx.Value("gin-i18n").(*I18n)

	return gi.T(ctx, msgID, tplData...)
}

func (i *I18n) T(ctx *gin.Context, msgID string, tplData ...map[any]any) string {
	i18 := ctx.Value("i18n").(*i18n.I18n)
	lng := GetCurrentLanguage(ctx)

	msg, err := i18.Translate(lng, msgID, tplData...)
	if err != nil && i.logger != nil {
		i.logger.Warn(fmt.Errorf("translation: %w", err).Error(), "msgID", msgID)
	}

	return msg
}

func (i *I18n) GetCurrentLanguage(ctx *gin.Context) language.Tag {
	return language.Make(i.getLngHandler(ctx))
}

func HasLanguage(ctx *gin.Context, l string) bool {
	i := ctx.MustGet("i18n").(*i18n.I18n)
	lng := language.Make(l)

	if lng == i.GetDefaultLanguage() {
		return true
	}

	return slices.Contains(i.GetLanguages(), lng)
}

func GetDefaultLanguage(ctx *gin.Context) language.Tag {
	i := ctx.MustGet("i18n").(*i18n.I18n)
	return i.GetDefaultLanguage()
}

func GetCurrentLanguage(ctx *gin.Context) language.Tag {
	i := ctx.MustGet("gin-i18n").(*I18n)
	return i.GetCurrentLanguage(ctx)
}

func defaultGetLngHandler(ctx *gin.Context) string {
	if ctx == nil || ctx.Request == nil {
		return language.English.String()
	}

	i := ctx.Value("i18n").(*i18n.I18n)
	ls := i.GetLanguages()
	defaultLng := i.GetDefaultLanguage()

	lng := ctx.Query("lng")
	if lng != "" {
		return lng
	}

	lng = ctx.GetHeader("Accept-Language")
	if lng != "" {
		ts, _, err := language.ParseAcceptLanguage(lng)
		if err != nil || len(ts) == 0 {
			return defaultLng.String()
		}

		for _, t := range ts {
			if slices.Contains(ls, t) {
				return t.String()
			}
		}
	}

	return defaultLng.String()
}
