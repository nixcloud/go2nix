// All go2nix CLI related stuff
package main

//go:generate go-bindata -o assets.go templates/

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jawher/mow.cli"
)

// go2nix entry-point
func main() {
	go2nix := cli.App("go2nix", "Nix derivations for Go packages")
	go2nix.Version("v version", "go2nix "+version)

	go2nix.Command("save", "Saves dependecies for cwd within GOPATH", func(cmd *cli.Cmd) {
		outputFile := cmd.StringOpt("o output", "default.nix",
			"Write the resulting nix file to the named output file")
		depsFile := cmd.StringOpt("d deps-file", "deps.nix",
			"Write the resulting dependencies file to the named output file")
		testImports := cmd.BoolOpt("t test-imports", false,
			"Include test imports.")
		buildTags := cmd.StringOpt("tags", "",
			"the dependencies will be generated with the specified build tags")

		cmd.Action = func() {
			goPath := os.Getenv("GOPATH")
			if goPath == "" {
				log.Fatal("No GOPATH set, can't find dependencies")
			}
			currPkg, err := currentPackage(goPath)
			if err != nil {
				log.Fatal(err)
			}
			buildTagsList := strings.Split(*buildTags, ",")
			if err := save(currPkg, goPath, *outputFile, *depsFile,
				*testImports, buildTagsList); err != nil {
				log.Fatal(err)
			}
		}
	})

	go2nix.Run(os.Args)
}

func currentPackage(goPath string) (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for _, goPathDir := range strings.Split(goPath, ":") {
		if strings.HasPrefix(currDir, goPathDir+"/src/") {
			return strings.TrimPrefix(currDir, goPathDir+"/src/"), nil
		}
	}

	return "", fmt.Errorf("Current dir %v is outside of GOPATH(%v). "+
		"Can't get current package name", currDir, goPath)
}
