package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/hironobu-s/go-corenlp/document"
	"github.com/hironobu-s/go-corenlp/provider"
)

func Annotate(p provider.Provider, text string) (root *document.Document, err error) {
	response, err := p.Run(text)
	if err != nil {
		return nil, err
	}
	defer response.Close()

	rawjson, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}
	return ParseOutput(rawjson)
}

func ParseOutput(rawjson []byte) (root *document.Document, err error) {
	err = json.Unmarshal(rawjson, &root)
	return root, err
}