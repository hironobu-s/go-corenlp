package provider

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type LocalExec struct {
	JavaCmd   string
	JavaArgs  []string
	Class     string
	ClassPath string

	Model       string
	Props       string
	Annotators  []string
	CoreNlpArgs []string

	ctx context.Context
}

func NewLocalExec(ctx context.Context) *LocalExec {
	if ctx == nil {
		ctx = context.Background()
	}
	return &LocalExec{
		JavaCmd:   "java",
		JavaArgs:  []string{},
		Class:     "edu.stanford.nlp.pipeline.StanfordCoreNLP",
		ClassPath: "*",

		Model:       "",
		Props:       "",
		Annotators:  []string{},
		CoreNlpArgs: []string{},

		ctx: ctx,
	}
}

func (p *LocalExec) Run(text string) (response Response, err error) {
	// create tmp file which write the input text
	tmp, err := ioutil.TempFile("", "go-corenlp")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())

	if _, err = tmp.WriteString(text); err != nil {
		return nil, err
	}

	// output file name
	outputDir := os.TempDir()
	outputFile := tmp.Name() + ".json"
	defer os.Remove(outputFile)

	// build arguments
	args := p.JavaArgs

	if p.ClassPath != "" {
		args = append(args, "-cp", p.ClassPath)
	}

	if p.Class != "" {
		args = append(args, p.Class)
	}

	args = append(args, p.CoreNlpArgs...)

	if p.Props != "" {
		args = append(args, "-props", p.Props)
	}

	if len(p.Annotators) > 0 {
		args = append(args, "-annotators", strings.Join(p.Annotators, ","))
	}

	args = append(args,
		"-file",
		tmp.Name(),
		"--outputFormat",
		"json",
		"--outputDirectory",
		outputDir,
		"--outputFile",
		outputFile,
	)

	// execute command
	cmd := exec.CommandContext(p.ctx, p.JavaCmd, args...)
	logrus.Debugf("Run command [%s]", strings.Join(cmd.Args, " "))

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("%s: %s", err.Error(), stderr.String())
		logrus.Warnf("Failed to execute command. [%s]", err.Error())
		return nil, err
	}

	if stdout.Len() > 0 {
		logrus.Debugf("Success to execute command. [%s]", stdout.String())
	} else {
		logrus.Debugf("Success to execute command.")
	}

	response, err = os.Open(outputFile)
	if err != nil {
		return nil, err
	}

	return response, nil
}
