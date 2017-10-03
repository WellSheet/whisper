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
	file := flag.String("f", "", "path to SOPS file")
	passThroughEnv := flag.Bool("p", false, "pass the current environment through")
	flag.Parse()

	command := strings.Join(flag.Args(), " ")

	if *file == "" || len(command) == 0 {
		log.Fatal("usage: whisper -f <file> [-p] <command>")
	}

	path, err := filepath.Abs(*file)

	if err != nil {
		log.Fatalf("Couldn't determine absolute path to file: %v", err)
	}

	decryptedEnvironment, err := decryptEnvironmentVars(path)

	if err != nil {
		log.Fatalf("Failed to decrypt environment vars: %v", err)
	}

	cmd := exec.Command("sh", "-c", command)

	if *passThroughEnv {
		cmd.Env = os.Environ()
	}

	cmd.Env = append(cmd.Env, decryptedEnvironment...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func decryptEnvironmentVars(sopsFile string) ([]string, error) {
	yamlBytes, err := decrypt.File(sopsFile, "yaml")

	if err != nil {
		return nil, err
	}

	secrets := make(map[string]interface{})
	err = yaml.Unmarshal(yamlBytes, &secrets)

	if err != nil {
		return nil, err
	}

	environmentStrings := make([]string, 1)

	environmentMap := secrets["environment"].(map[interface{}]interface{})
	for k, v := range environmentMap {
		environmentStrings = append(environmentStrings, fmt.Sprintf("%s=%v", k, v))
	}

	return environmentStrings, nil
}
