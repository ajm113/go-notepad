//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var buildTags = []string{"pango_1_42", "gtk_3_22"}

const binaryName = "notepad"
const binaryDir = "./build"
const version = "v1.0.0"

func Run() error {
	mg.Deps(Build)

	return sh.RunV(path.Join(binaryDir, binaryName))
}

func Build() error {
	fmt.Println("Building binary:", binaryName)
	return build(binaryName, buildTags)
}

func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll(binaryDir)
}

func build(binaryName string, tags []string) error {
	ldflags := "-s -w -X main.version=" + version
	if tags != nil && len(tags) > 0 {
		return sh.RunV("go", "build", "-ldflags", ldflags, "-tags", strings.Join(tags, " "), "-o", path.Join(binaryDir, binaryName), ".")
	}
	return sh.RunV("go", "build", "-ldflags", ldflags, "-o", path.Join(binaryDir, binaryName), ".")
}
