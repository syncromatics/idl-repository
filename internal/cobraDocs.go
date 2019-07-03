package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	idlRepository "github.com/syncromatics/idl-repository/cmd/idl-repository/cmd"
	idl "github.com/syncromatics/idl-repository/cmd/idl/cmd"
)

func main() {
	err := doc.GenMarkdownTree(idl.RootCmd, "docs/idl")
	if err != nil {
		log.Fatal(err)
	}
	err = doc.GenMarkdownTree(idlRepository.RootCmd, "docs/idl-repository")
	if err != nil {
		log.Fatal(err)
	}
}
