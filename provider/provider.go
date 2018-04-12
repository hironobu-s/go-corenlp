package provider

import (
	"io"
)

type Response io.ReadCloser

type Provider interface {
	Run(string) (Response, error)
}
