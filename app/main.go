package main

import (
	"github.com/kennethklee/ssh-authorized-manager/ssham"

	// Built-in plugins
	_ "github.com/kennethklee/ssh-authorized-manager/plugins/header-auth"
)

func main() {
	ssham.Main()
}
