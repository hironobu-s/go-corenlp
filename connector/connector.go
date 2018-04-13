package connector

import (
	"io"
)

type Response io.ReadCloser

type Connector interface {
	Run(string) (Response, error)
}
