package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	g18n "github.com/litsea/gin-i18n"
	"github.com/litsea/i18n"
	"golang.org/x/text/language"
)

//go:embed localize/*
var fs embed.FS

func main() {
	// new gin engine
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	gi := g18n.New(
		g18n.WithOptions(
			i18n.WithLanguages(language.English, language.German),
			i18n.WithLoaders(
				i18n.EmbedLoader(fs, "./localize/"),
			),
		),
	)
	// apply i18n middleware
	r.Use(gi.Localize())

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, g18n.T(ctx, "welcome"))
	})

	r.GET("/:name", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, g18n.T(ctx, "welcomeWithName", map[any]any{
			"name": ctx.Param("name"),
		}))
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
