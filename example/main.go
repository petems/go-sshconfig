package main

import (
	"fmt"
	"log"
	"os"

	"github.com/petems/go-sshconfig"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	sshConfigFile := os.ExpandEnv("$HOME/.ssh/config")

	file, err := os.Open(sshConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	config, err := sshconfig.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()

	// modify by reference for existing params
	// or create a new param and append it to global
	if param := config.GetParam(sshconfig.VisualHostKeyKeyword); param != nil {
		fmt.Println("VisualHostKey found! Switching to value: yes")
		param.Args = []string{"yes"}
		param.Comments = []string{"Added by the petems/go-sshconfig example app"}
	} else {
		fmt.Println("VisualHostKey not found! Adding with value: yes")
		param = sshconfig.NewParam(sshconfig.VisualHostKeyKeyword, []string{"yes"}, []string{"Added by the petems/go-sshconfig example app"})
		config.Globals = append(config.Globals, param)
	}

	fmt.Println("Creating example config file with VisualHostKey added in")

	// atomic write to file to ensure config is preserved in
	// the event of an error
	if err := config.WriteToFilepath("./example_config"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created file at ./example_config")

}
