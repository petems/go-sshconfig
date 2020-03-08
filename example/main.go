package main

import (
	"log"
	"os"

	"github.com/petems/sshconfig"
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
		param.Args = []string{"yes"}
	} else {
		param = sshconfig.NewParam(sshconfig.VisualHostKeyKeyword, []string{"yes"}, []string{"good to see you"})
		config.Globals = append(config.Globals, param)
	}

	// atomic write to file to ensure config is preserved in
	// the event of an error
	if err := config.WriteToFilepath("./example_config"); err != nil {
		log.Fatal(err)
	}

}
