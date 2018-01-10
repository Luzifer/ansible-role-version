package main

import "github.com/Luzifer/ansible-role-version/cmd"

var version = "dev"

func main() {
	cmd.Execute(version)
}
