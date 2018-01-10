package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/Luzifer/rconfig"
	log "github.com/sirupsen/logrus"
)

var (
	cfg = struct {
		RolesFile      string `flag:"roles-file,f" default:"requirements.yml" description:"File containing the requirements"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

type ansibleRoleDefinition struct {
	Name    string `yaml:"name"`
	Src     string `yaml:"src"`
	Version string `yaml:"version"`
}

func init() {
	if err := rconfig.Parse(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("ansible-role-version %s\n", version)
		os.Exit(0)
	}
}

func main() {
	args := rconfig.Args()[1:]
	if len(args) != 2 {
		log.Fatalf("Usage: ansible-role-version <role-name> <new version>")
	}

	roleName := args[0]
	roleVersion := args[1]

	var (
		inFileContent []byte
		err           error
	)
	if inFileContent, err = ioutil.ReadFile(cfg.RolesFile); err != nil {
		log.WithError(err).Fatalf("Roles file not found")
	}

	in := []ansibleRoleDefinition{}
	if err = yaml.Unmarshal(inFileContent, &in); err != nil {
		log.WithError(err).Fatal("Unable to parse roles file")
	}

	for i := range in {
		if in[i].Name == roleName {
			in[i].Version = roleVersion
		}
	}

	if inFileContent, err = yaml.Marshal(in); err != nil {
		log.WithError(err).Fatal("Unable to marshal roles file")
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte("---\n\n"))
	buf.Write(inFileContent)
	buf.Write([]byte("\n...\n"))

	if err = ioutil.WriteFile(cfg.RolesFile, buf.Bytes(), 0644); err != nil {
		log.WithError(err).Fatal("Unable to write roles file")
	}

	log.Info("Roles file written successfully")
}
