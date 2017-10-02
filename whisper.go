package main

import (
	"flag"
	"fmt"
	"go.mozilla.org/sops/decrypt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// get options: -f <file name> <command>
	file := flag.String("f", "", "path to SOPS file")
	flag.Parse()

	command := flag.Args()

	fmt.Println("file = ", *file)
	fmt.Println("command = ", strings.Join(command, " "))

	if *file == "" || len(command) == 0 {
		fmt.Println("usage: whisper -f <file> <command>")
		os.Exit(1)
	}

	path, err := filepath.Abs(*file)

	if err != nil {
		fmt.Println("Couldn't determine absolute path to file", err)
		os.Exit(1)
	}

	yamlBytes, err := decrypt.File(path, "yaml")

	if err != nil {
		fmt.Println("Couldn't decrypt file!")
		fmt.Println(err)
		os.Exit(1)
	}

	yamlMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal(yamlBytes, &yamlMap)

	if err != nil {
		fmt.Println("Couldn't figure out those yaml bytes")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%v", yamlMap)

	// get all of the values under environment key
	// pass them to os.Exec?
}
