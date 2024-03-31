package sshconfig

import (
	"fmt"
	"os"
	"strings"
)

func ExampleConfig_GetParam() {
	sshConfigExample := `
# global configuration
VisualHostKey yes

# host-based configuration

# dev
Host dev
  HostName 127.0.0.1
  User ubuntu
  Port 22

Host *.google.com *.yahoo.com
  User root
`

	config, err := Parse(strings.NewReader(sshConfigExample))

	if err != nil {
		panic(err)
	}

	visualHostKey := config.GetParam(VisualHostKeyKeyword)

	fmt.Println(visualHostKey)
	// Output: VisualHostKey yes
}

func ExampleConfig_GetHost() {
	getHostBlock := `
Host get-host
  HostName 127.0.0.1
  User ubuntu
  Port 22
`

	config, err := Parse(strings.NewReader(getHostBlock))

	if err != nil {
		panic(err)
	}

	devHost := config.GetHost("get-host")

	fmt.Println(devHost)
	// Output:
	// Host get-host
	//   HostName 127.0.0.1
	//   User ubuntu
	//   Port 22
}

func ExampleConfig_AddHost() {
	addHostBlock := `
Host dev
  HostName 127.0.0.1
  User ubuntu
  Port 22
`

	config, err := Parse(strings.NewReader(addHostBlock))

	if err != nil {
		panic(err)
	}

	gitExampleHost := NewHost([]string{"git.example.com"}, []string{"My cool git server"})

	config.AddHost(gitExampleHost)

	config.WriteTo(os.Stdout)
	// Output:
	// # global configuration
	//
	// # host-based configuration
	//
	// Host dev
	//   HostName 127.0.0.1
	//   User ubuntu
	//   Port 22
	//
	// # My cool git server
	// Host git.example.com
}

func ExampleHost_AddParam() {
	addParamBlock := `
Host dev
  HostName 127.0.0.1
  User ubuntu
  Port 22
`

	config, err := Parse(strings.NewReader(addParamBlock))

	if err != nil {
		panic(err)
	}

	devHost := config.GetHost("dev")

	knownHostsParam := NewParam(UserKnownHostsFileKeyword, []string{"/dev/null"}, nil)

	devHost.AddParam(knownHostsParam)

	fmt.Println(devHost)
	// Output:
	// Host dev
	//   HostName 127.0.0.1
	//   User ubuntu
	//   Port 22
	//   UserKnownHostsFile /dev/null
}

func ExampleHost_AddParam_withcomment() {
	sshConfigExample := `
Host dev
`

	config, err := Parse(strings.NewReader(sshConfigExample))

	if err != nil {
		panic(err)
	}

	devHost := config.GetHost("dev")

	knownHostsParam := NewParam(UserKnownHostsFileKeyword, []string{"/dev/null"}, []string{"MITM dont scare me"})

	devHost.AddParam(knownHostsParam)

	fmt.Println(devHost)
	// Output:
	// Host dev
	//   # MITM dont scare me
	//   UserKnownHostsFile /dev/null
}

func ExampleConfig_FindByHostname() {
	sshConfigExample := `
Host dev
  HostName dev.example.com
  User ubuntu
  Port 22
`

	config, err := Parse(strings.NewReader(sshConfigExample))

	if err != nil {
		panic(err)
	}

	devExampleHost := config.FindByHostname("dev.example.com")

	fmt.Println(devExampleHost)
	// Output:
	// Host dev
	//   HostName dev.example.com
	//   User ubuntu
	//   Port 22
}

func ExampleHost_GetParam() {

	sshDevHostBlock := `
# dev
Host dev
  HostName 127.0.0.1
  User ubuntu
  Port 22
`

	config, err := Parse(strings.NewReader(sshDevHostBlock))

	if err != nil {
		panic(err)
	}

	host := config.GetHost("dev")

	hostname := host.GetParam(HostNameKeyword)

	fmt.Println(hostname)
	// Output: HostName 127.0.0.1
}

func ExampleParam_Value() {

	sshConfigExample := `
# global configuration
VisualHostKey yes

# host-based configuration

# dev
Host dev
  HostName 127.0.0.1
  User ubuntu
  Port 22

Host *.google.com *.yahoo.com
  User root
`

	config, err := Parse(strings.NewReader(sshConfigExample))

	if err != nil {
		panic(err)
	}

	visualHostKey := config.GetParam(VisualHostKeyKeyword)

	fmt.Println(visualHostKey.Value())
	// Output: yes
}

func ExampleNewParam() {

	param := NewParam(VisualHostKeyKeyword, []string{"yes"}, []string{"Add ascii art of key, see https://man.openbsd.org/ssh_config.5#VisualHostKey"})

	fmt.Println(param)
	// Output:
	// # Add ascii art of key, see https://man.openbsd.org/ssh_config.5#VisualHostKey
	// VisualHostKey yes
}

func ExampleNewParam_hostwithcomment() {
	sshConfigExample := `
Host dev
	HostName 127.0.0.1
	User ubuntu
`

	config, err := Parse(strings.NewReader(sshConfigExample))

	if err != nil {
		panic(err)
	}

	devHost := config.GetHost("dev")

	knownHostsParam := NewParam(UserKnownHostsFileKeyword, []string{"/dev/null"}, []string{"MITM dont scare me"})

	devHost.Params = append(devHost.Params, knownHostsParam)

	fmt.Println(devHost)
	// Output:
	// Host dev
	//   HostName 127.0.0.1
	//   User ubuntu
	//   # MITM dont scare me
	//   UserKnownHostsFile /dev/null
}

func ExampleNewParam_globalwithcomment() {
	sshConfigExample := `
Host dev
	HostName 127.0.0.1
	User ubuntu
`

	config, err := Parse(strings.NewReader(sshConfigExample))

	if err != nil {
		panic(err)
	}

	param := NewParam(VisualHostKeyKeyword, []string{"yes"}, []string{"Add ascii art of key, see https://man.openbsd.org/ssh_config.5#VisualHostKey"})

	config.Globals = append(config.Globals, param)

	config.WriteTo(os.Stdout)
	// Output:
	// # global configuration
	// # Add ascii art of key, see https://man.openbsd.org/ssh_config.5#VisualHostKey
	// VisualHostKey yes

	// # host-based configuration

	// Host dev
	// 	HostName 127.0.0.1
	// 	User ubuntu
	// want:
	// Host dev
	// 	# MITM dont scare me
	// 	UserKnownHostsFile /dev/null
}

func ExampleNewHost() {

	host := NewHost([]string{"git.example.com"}, []string{"My cool git server"})

	fmt.Println(host)
	// Output:
	// # My cool git server
	// Host git.example.com
}

func ExampleNewHost_githubdotcomexample() {

	githubhost := NewHost([]string{"github.com"}, []string{"github.com global config"})

	controlmasterParam := NewParam(ControlMasterKeyword, []string{"auto"}, []string{})
	controlpathParam := NewParam(ControlPathKeyword, []string{"~/.ssh/ssh-%r@%h:%p"}, []string{})
	controlpersistParam := NewParam(ControlPersistKeyword, []string{"yes"}, []string{})
	userParam := NewParam(UserKeyword, []string{"git"}, []string{})

	newParams := []*Param{
		controlmasterParam,
		controlpathParam,
		controlpersistParam,
		userParam,
	}

	for i := 0; i < len(newParams); i++ {
		githubhost.AddParam(newParams[i])
	}

	fmt.Println(githubhost)
	// Output:
	// # github.com global config
	// Host github.com
	//   ControlMaster auto
	//   ControlPath ~/.ssh/ssh-%r@%h:%p
	//   ControlPersist yes
	//   User git
}
