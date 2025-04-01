package i18n

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/litsea/gin-i18n/testdata"
	"github.com/litsea/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestI18n(t *testing.T) {
	t.Parallel()

	type args struct {
		lng  language.Tag
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// English
		{
			name: "en-hello",
			args: args{
				path: "/",
				lng:  language.English,
			},
			want: "hello",
		},
		{
			name: "en-hello-messageId",
			args: args{
				path: "/messageId/alex",
				lng:  language.English,
			},
			want: "hello alex",
		},
		// German
		{
			name: "de-hello",
			args: args{
				path: "/",
				lng:  language.German,
			},
			want: "hallo",
		},
		{
			name: "de-hello-messageId",
			args: args{
				path: "/messageId/alex",
				lng:  language.German,
			},
			want: "hallo alex",
		},
		// French (fallback)
		{
			name: "fr-hello",
			args: args{
				path: "/",
				lng:  language.French,
			},
			want: "hello",
		},
		{
			name: "fr-hello-messageId",
			args: args{
				path: "/messageId/alex",
				lng:  language.French,
			},
			want: "hello alex",
		},
		// exist
		{
			name: "lang-exist-" + language.English.String(),
			args: args{
				path: fmt.Sprintf("/exist/%s", language.English.String()),
				lng:  language.English,
			},
			want: "true",
		},
		{
			name: "lang-not-exist-" + language.SimplifiedChinese.String(),
			args: args{
				path: fmt.Sprintf("/exist/%s", language.SimplifiedChinese.String()),
				lng:  language.English,
			},
			want: "false",
		},
		// default-lang
		{
			name: "lang-is-default-" + language.English.String(),
			args: args{
				path: "/lng/default",
				lng:  language.English,
			},
			want: language.English.String(),
		},
		{
			name: "lang-is-not-default-" + language.German.String(),
			args: args{
				path: "/lng/default",
				lng:  language.German,
			},
			want: language.English.String(),
		},
		// current-lang
		{
			name: "current-lang-" + language.English.String(),
			args: args{
				path: "/lng/current",
				lng:  language.English,
			},
			want: language.English.String(),
		},
		{
			name: "current-lang-" + language.German.String(),
			args: args{
				path: "/lng/current",
				lng:  language.German,
			},
			want: language.German.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := makeRequest(tt.args.lng, tt.args.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	gi := New(
		WithOptions(
			i18n.WithLanguages(language.English, language.German),
			i18n.WithLoaders(
				i18n.EmbedLoader(testdata.Localize, "./localize/"),
			),
		),
	)

	r.Use(gi.Localize())

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, T(ctx, "welcome"))
	})

	r.GET("/messageId/:name", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, T(ctx, "welcomeWithName", map[any]any{
			"name": ctx.Param("name"),
		}))
	})

	r.GET("/messageIdWithField/:messageId/:field", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, T(ctx, ctx.Param("messageId"), map[any]any{
			"field": ctx.Param("field"),
		}))
	})

	r.GET("/exist/:lng", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%v", HasLanguage(ctx, ctx.Param("lng")))
	})

	r.GET("/lng/default", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", GetDefaultLanguage(ctx).String())
	})

	r.GET("/lng/current", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", GetCurrentLanguage(ctx).String())
	})

	return r
}

func makeRequest(lng language.Tag, path string) string {
	req, _ := http.NewRequestWithContext(context.Background(), "GET", path, nil)
	req.Header.Add("Accept-Language", lng.String())

	w := httptest.NewRecorder()
	s := newServer()
	s.ServeHTTP(w, req)

	return w.Body.String()
}
