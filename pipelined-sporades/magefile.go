//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	Default = Build
)

// Install build dependencies.
func BuildDeps() error {
	err := sh.RunV("protoc", "--version")
	if err != nil {
		return err
	}
	err = sh.RunV("go", "get", "-u", "github.com/golang/protobuf/protoc-gen-go")
	if err != nil {
		return err
	}
	err = sh.RunV("go", "get", "-u", "google.golang.org/grpc")
	if err != nil {
		return err
	}

	return nil
}

// Install dependencies.
func Deps() error {
	err := sh.RunV("go", "mod", "vendor")
	if err != nil {
		return err
	}

	return nil
}

// Generate code.
func Generate() error {
	err := sh.RunV("protoc", "--go_out=./", "./proto/definitions.proto")
	if err != nil {
		return err
	}
	return nil
}

// Run tests.
func Test() error {
	mg.Deps(Generate)
	return sh.RunV("go", "test", "-v", "./...")
}

// Build binary executables.
func Build() error {
	err := sh.RunV("go", "build", "-v", "-o", "./client/bin/client", "./client/")
	err = sh.RunV("go", "build", "-v", "-o", "./replica/bin/replica", "./replica/")
	if err != nil {
		return err
	}

	return nil
}
