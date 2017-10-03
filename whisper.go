package main

import (
	"flag"
	"fmt"
	"go.mozilla.org/sops/decrypt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// get options: -f <file name> <command>
	file := flag.String("f", "", "path to SOPS file")
	passThroughEnv := flag.Bool("p", false, "pass the current environment through")
	flag.Parse()

	command := strings.Join(flag.Args(), " ")

	log.Printf("file = %v", *file)
	log.Printf("command = %v", command)

	if *file == "" || len(command) == 0 {
		log.Fatal("usage: whisper -f <file> <command>")
	}

	path, err := filepath.Abs(*file)

	if err != nil {
		log.Fatalf("Couldn't determine absolute path to file: %v", err)
	}

	yamlBytes, err := decrypt.File(path, "yaml")

	if err != nil {
		log.Fatalf("Couldn't decrypt file: %v", err)
	}

	secrets := make(map[string]interface{})
	err = yaml.Unmarshal(yamlBytes, &secrets)

	if err != nil {
		log.Fatalf("Couldn't unmarshal YAML bytes: %v", err)
	}

	cmd := exec.Command("sh", "-c", command)

	if *passThroughEnv {
		cmd.Env = os.Environ()
	}

	environmentMap := secrets["environment"].(map[interface{}]interface{})
	for k, v := range environmentMap {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%v", k, v))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
