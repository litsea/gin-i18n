# gin-i18n

## Usage

### Use Middleware

```golang
import (
	"embed

	"github.com/gin-gonic/gin"
	g18n "github.com/litsea/gin-i18n"
	"github.com/litsea/i18n"
	"golang.org/x/text/language"
)

//go:embed localize/*
var fs embed.FS

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
```

### Translate

```golang
import (
	g18n "github.com/litsea/gin-i18n"
)

g18n.T(ctx, "welcomeWithName", map[any]any{
	"name": ctx.Param("name"),
}))
```