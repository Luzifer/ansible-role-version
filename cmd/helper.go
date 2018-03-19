package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ansibleRoleDefinition struct {
	Name    string `yaml:"name"`
	Src     string `yaml:"src"`
	Version string `yaml:"version"`
}

func getRoleDefinitions(rolesFile string) ([]ansibleRoleDefinition, error) {
	rf, err := os.Open(rolesFile)
	if err != nil {
		return nil, err
	}
	defer rf.Close()

	def := []ansibleRoleDefinition{}
	return def, yaml.NewDecoder(rf).Decode(&def)
}

func patchRoleFile(rolesFile string, updates map[string]string) error {
	inFile, err := os.Open(rolesFile)
	if err != nil {
		return fmt.Errorf("Unable to open roles files: %s", err)
	}
	defer inFile.Close()

	in := []ansibleRoleDefinition{}
	if err = yaml.NewDecoder(inFile).Decode(&in); err != nil {
		return fmt.Errorf("Unable to parse roles file: %s", err)
	}

	for roleName, roleVersion := range updates {
		for i := range in {
			if in[i].Name == roleName {
				in[i].Version = roleVersion
			}
		}
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte("---\n\n"))

	if err := yaml.NewEncoder(buf).Encode(in); err != nil {
		return fmt.Errorf("Unable to marshal roles file: %s", err)
	}

	buf.Write([]byte("\n...\n"))

	if err = ioutil.WriteFile(rolesFile, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("Unable to write roles file: %s", err)
	}

	return nil
}
