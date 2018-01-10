package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type ansibleRoleDefinition struct {
	Name    string `yaml:"name"`
	Src     string `yaml:"src"`
	Version string `yaml:"version"`
}

func patchRoleFile(rolesFile string, updates map[string]string) error {
	var (
		inFileContent []byte
		err           error
	)
	if inFileContent, err = ioutil.ReadFile(rolesFile); err != nil {
		return fmt.Errorf("Roles file not found: %s", err)
	}

	in := []ansibleRoleDefinition{}
	if err = yaml.Unmarshal(inFileContent, &in); err != nil {
		return fmt.Errorf("Unable to parse roles file: %s", err)
	}

	for roleName, roleVersion := range updates {
		for i := range in {
			if in[i].Name == roleName {
				in[i].Version = roleVersion
			}
		}
	}

	if inFileContent, err = yaml.Marshal(in); err != nil {
		return fmt.Errorf("Unable to marshal roles file: %s", err)
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte("---\n\n"))
	buf.Write(inFileContent)
	buf.Write([]byte("\n...\n"))

	if err = ioutil.WriteFile(rolesFile, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("Unable to write roles file: %s", err)
	}

	return nil
}
