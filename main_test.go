package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"testing"
)

var content []byte

func init() {
	b, err := os.ReadFile("testdata/mock.yaml")
	if err != nil {
		panic("oops")
	}
	content = b
}

func BenchmarkSplitToMultidoc(b *testing.B) {
	r := bytes.NewReader(content)
	b.ResetTimer()
	if err := splitToMultidoc(r, io.Discard); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkYQ(b *testing.B) {
	r := bytes.NewReader(content)
	b.ResetTimer()
	cmd := exec.Command("yq", ".items[] | split_doc")
	cmd.Stdin = r
	cmd.Stdout = io.Discard
	b.ResetTimer()
	cmd.Run()
}
