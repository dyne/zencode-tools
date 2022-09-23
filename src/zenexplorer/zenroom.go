package main

import (
	"bytes"
	"io"
	"os/exec"
)

type ZenResult struct {
    Output string;
    Logs string;
}
type ZenInput struct {
	Script string
	Keys string;
	Data string;
}

func ZenroomExec(stdin io.Reader, zenIn ZenInput) (ZenResult, bool) {
	return wrapper(stdin, zenIn)
}

func ZencodeExec(stdin io.Reader, zenIn ZenInput) (ZenResult, bool) {
	return wrapper(stdin, zenIn, "-z")
}

func wrapper(stdin io.Reader, zenIn ZenInput, args ...string) (ZenResult, bool) {
	args = append(args, zenIn.Script)
	if zenIn.Keys != "" {
		args = append(args, "-k", zenIn.Keys)
	}
	if zenIn.Data != "" {
		args = append(args, "-a", zenIn.Data)
	}
	cmd := exec.Command("zenroom", args...)
	cmd.Stdin = stdin

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()

	return ZenResult{Logs: errBuf.String(), Output: outBuf.String()}, err != nil
}
