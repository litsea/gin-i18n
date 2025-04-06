package i18n

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/litsea/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/litsea/gin-i18n/testdata"
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
		{
			name: "en-hello-messageTemplate",
			args: args{
				path: "/messageTemplate/alex",
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
		{
			name: "de-hello-messageTemplate",
			args: args{
				path: "/messageTemplate/alex",
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
		{
			name: "en-hello-messageTemplate",
			args: args{
				path: "/messageTemplate/alex",
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

func newServer(mw ...gin.HandlerFunc) *gin.Engine {
	// TODO: data race warning for gin mode
	// https://github.com/gin-gonic/gin/pull/1580 (not yet released)
	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	if len(mw) > 0 {
		r.Use(mw...)
	}

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, T(ctx, "welcome"))
	})

	r.GET("/messageId/:name", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, T(ctx, "welcomeWithName", map[any]any{
			"name": ctx.Param("name"),
		}))
	})

	r.GET("/messageTemplate/:name", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, T(ctx, "hello {{ .name }}", map[any]any{
			"name": ctx.Param("name"),
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
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, http.NoBody)
	req.Header.Add("Accept-Language", lng.String())

	gi := New(
		WithOptions(
			i18n.WithLanguages(language.English, language.German),
			i18n.WithLoaders(
				i18n.EmbedLoader(testdata.Localize, "./localize/"),
			),
		),
	)

	w := httptest.NewRecorder()
	s := newServer(gi.Localize())
	s.ServeHTTP(w, req)

	return w.Body.String()
}

func TestNoI18nContext(t *testing.T) {
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
		// English (fallback to msgID)
		{
			name: "en-no-i18n-hello",
			args: args{
				path: "/",
				lng:  language.English,
			},
			want: "welcome",
		},
		{
			name: "en-no-i18n-hello-messageId",
			args: args{
				path: "/messageId/alex",
				lng:  language.English,
			},
			want: "welcomeWithName",
		},
		{
			// (fallback to template)
			name: "en-no-i18n-hello-messageTemplate",
			args: args{
				path: "/messageTemplate/alex",
				lng:  language.English,
			},
			want: "hello alex",
		},
		// German (fallback to msgID)
		{
			name: "de-no-i18n-hello",
			args: args{
				path: "/",
				lng:  language.German,
			},
			want: "welcome",
		},
		{
			name: "de-no-i18n-hello-messageId",
			args: args{
				path: "/messageId/alex",
				lng:  language.German,
			},
			want: "welcomeWithName",
		},
		{
			// (fallback to template)
			name: "de-no-i18n-hello-messageTemplate",
			args: args{
				path: "/messageTemplate/alex",
				lng:  language.German,
			},
			want: "hello alex",
		},
		// exist (all not exists)
		{
			name: "no-i18n-lang-not-exist-" + language.English.String(),
			args: args{
				path: fmt.Sprintf("/exist/%s", language.English.String()),
				lng:  language.English,
			},
			want: "false",
		},
		{
			name: "no-i18n-lang-not-exist-" + language.SimplifiedChinese.String(),
			args: args{
				path: fmt.Sprintf("/exist/%s", language.SimplifiedChinese.String()),
				lng:  language.English,
			},
			want: "false",
		},
		// default-lang (always English)
		{
			name: "no-i18n-lang-is-default-" + language.English.String(),
			args: args{
				path: "/lng/default",
				lng:  language.English,
			},
			want: language.English.String(),
		},
		{
			name: "no-i18n-lang-is-not-default-" + language.German.String(),
			args: args{
				path: "/lng/default",
				lng:  language.German,
			},
			want: language.English.String(),
		},
		// current-lang (always English)
		{
			name: "no-i18n-current-lang-" + language.English.String(),
			args: args{
				path: "/lng/current",
				lng:  language.English,
			},
			want: language.English.String(),
		},
		{
			name: "no-i18n-current-lang-" + language.German.String(),
			args: args{
				path: "/lng/current",
				lng:  language.German,
			},
			want: language.English.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := makeRequestNoI18nContext(tt.args.lng, tt.args.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func makeRequestNoI18nContext(lng language.Tag, path string) string {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, http.NoBody)
	req.Header.Add("Accept-Language", lng.String())

	w := httptest.NewRecorder()
	s := newServer()
	s.ServeHTTP(w, req)

	return w.Body.String()
}
