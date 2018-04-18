package connector

import (
	"io"
)

// Response is returned by some structs which implements Connector interface.
type Response io.ReadCloser

// Connector means the something to communicate with Stanford CoreNLP.
type Connector interface {
	Run(string) (Response, error)
}
