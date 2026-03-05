package client

import "embed"

//go:embed static/*
var StaticFiles embed.FS

//go:embed page.html
var PageHtml string
