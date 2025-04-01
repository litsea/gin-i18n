package testdata

import (
	"embed"
)

//go:embed localize/*
var Localize embed.FS
