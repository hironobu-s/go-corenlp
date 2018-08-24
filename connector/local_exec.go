package connector

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

// LocalExec connector is responsible to run Stanford CoreNLP process.
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

// NewLocalExec returns a pointer of LocalExec
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

// Run marshals Connector interface implementation.
func (c *LocalExec) Run(text string) (response Response, err error) {
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
	args := c.JavaArgs

	if c.ClassPath != "" {
		args = append(args, "-cp", c.ClassPath)
	}

	if c.Class != "" {
		args = append(args, c.Class)
	}

	args = append(args, c.CoreNlpArgs...)

	if c.Props != "" {
		args = append(args, "-props", c.Props)
	}

	if len(c.Annotators) > 0 {
		args = append(args, "-annotators", strings.Join(c.Annotators, ","))
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
	cmd := exec.CommandContext(c.ctx, c.JavaCmd, args...)
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
