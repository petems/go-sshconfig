# go-sshconfig
[![Build Status](https://travis-ci.com/petems/go-sshconfig.svg?branch=master)](https://travis-ci.com/petems/go-sshconfig)[![](https://godoc.org/github.com/petems/go-sshconfig?status.svg)](http://godoc.org/github.com/petems/go-sshconfig)[![Go Report Card](https://goreportcard.com/badge/github.com/petems/go-sshconfig)](https://goreportcard.com/report/github.com/petems/go-sshconfig)

A simple [ssh_config](https://man.openbsd.org/ssh_config) parser/writer library

## Example

An app to add a global config entry to a config file:

```
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

    fmt.Println(config)
```

`$HOME/.ssh/config` contents:

```
Host github.com
  ControlMaster auto
  ControlPath ~/.ssh/ssh-%r@%h:%p
  ControlPersist yes
  User git
``` 

Output:

```
VisualHostKey not found! Adding with value: yes
# global configuration

# Added by the petems/go-sshconfig example app
VisualHostKey yes

# host-based configuration

Host github.com
  ControlMaster auto
  ControlPath ~/.ssh/ssh-%r@%h:%p
  ControlPersist yes
  User git
```

The best way to go deeper is to read the [docs](https://godoc.org/github.com/petems/go-sshconfig).

## Attribution

Forked and refactored from [emptyinterface/sshconfig/](https://github.com/emptyinterface/sshconfig/)

## Development
If you'd like to other features or anything else, check out the contributing guidelines in [CONTRIBUTING.md](CONTRIBUTING.md).