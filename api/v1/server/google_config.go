package server

import (
	"context"

	"golang.org/x/oauth2"
)

type IGoogleConfig interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	Exchange(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}
