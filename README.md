# go-corenlp

`go-corenlp` is a Golang wrapper for [Stanford CoreNLP](https://stanfordnlp.github.io/CoreNLP/). 

## Install

Download and install it:

```shell
go get github.com/hironobu-s/go-corenlp
```

Make sure that you can run Stanford CoreNLP on [command line](https://stanfordnlp.github.io/CoreNLP/cmdline.html):

```shell
java -cp "*" edu.stanford.nlp.pipeline.StanfordCoreNLP -h
```

## Usage

A simple code for using `go-corenlp` is:

```go
package main

import (
	"fmt"

	"github.com/hironobu-s/go-corenlp" // exposes "corenlp"
	"github.com/hironobu-s/go-corenlp/connector"
)

func main() {
	// sample text from https://stanfordnlp.github.io/CoreNLP/
	text := `President Xi Jinping of Chaina, on his first state visit to the United States, showed off his familiarity with American history and pop culture on Tuesday night.`

	// LocalExec connector is responsible to run Stanford CoreNLP process.
	c := connector.NewLocalExec(nil)
	c.ClassPath = "./corenlp/*" // set Java class path
	c.Annotators = []string{"tokenize", "ssplit", "pos"}

	// Annotate text
	doc, err := corenlp.Annotate(c, text)
	if err != nil {
		panic(err)
	}

	// Output words and pos
	for _, sentence := range doc.Sentences {
		for _, token := range sentence.Tokens {
			fmt.Printf("%s(%s)%s", token.Word, token.Pos, token.After)
		}
	}
}
	
```

Output:

```text
President(NNP) Xi(NN) Jinping(NN) of(IN) Chaina(NNP),(,) on(IN) his(PRP$) first(JJ) state(NN) visit(NN) to(TO) the(DT) United(NNP) States(NNPS),(,) showed(VBD) off(IN) his(PRP$) familiarity(NN) with(IN) American(JJ) history(NN) and(CC) pop(NN) culture(NN) on(IN) Tuesday(NNP) night(NN).(.)
```

### Handle an annotated documents 

```go
// Annotate text
doc, err := corenlp.Annotate(connector.NewLocalExec(nil), text)
if err != nil {
	panic(err)
}

// First sentence
sentence := doc.Sentences[0]

// RawParse contains text-based result of Parser annotator
fmt.Println(sentence.RawParse) // => (ROOT (S (NP (NP (NNP President)...

// Parse() returns go's struct of Parser annotator
parse, _ := sentence.Parse()
fmt.Printf("%v\n", parse.Pos) // => ROOT

// Tokenizer, PosTagger
for _, token := range sentence.Tokens {
	fmt.Printf("%s(%s)%s", token.Word, token.Pos, token.After)
}

// Dependencies
for _, dep := range sentence.Dependencies {
	fmt.Printf("%s => (%s) => %s\n", dep.GovernorGloss, dep.Dep, dep.DependentGloss)
}
```

### Timeout

`go-corenlp` supports a timeout by using `context.Context`.

```go
ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
defer cancel()

c := connector.NewLocalExec(ctx)
doc, err := corenlp.Annotate(c, text)
```

### Connect to CoreNLP server

To connect [CoreNLP server](https://stanfordnlp.github.io/CoreNLP/corenlp-server.html), You may use `HTTPClient provider`.

```go
ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
defer cancel()

c := connector.NewHTTPClient(ctx, "http://127.0.0.1:9000/")
c.Username = "username"
c.Password = "password"

doc, err := corenlp.Annotate(c, text)
```

### Parse json output 

To use `ParseOutput` method, You can parse the output file which is generated by Stanford CoreNLP.

For example. If you run following command

```shell
java -cp "*" edu.stanford.nlp.pipeline.StanfordCoreNLP -annotators tokenize,ssplit -file input.txt --outputFormat json
```

The output file `input.txt.json` will be generated, So you can parse it as below.

```go
rawjson, err := ioutil.ReadFile("input.txt.json")
if err != nil {
	panic(err)
}
doc, err := ParseOutput(rawjson)

```

## LICENSE

MIT
