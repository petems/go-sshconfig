package sshconfig

import (
	"fmt"
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
	sshConfigExample := `
Host dev
  HostName 127.0.0.1
  User ubuntu
  Port 22
`

	config, err := Parse(strings.NewReader(sshConfigExample))

	if err != nil {
		panic(err)
	}

	devHost := config.GetHost("dev")

	fmt.Println(devHost)
	// Output:
	// Host dev
	//   HostName 127.0.0.1
	//   User ubuntu
	//   Port 22
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

func ExampleNewHost() {

	host := NewHost([]string{"git.example.com"}, []string{"My cool git server"})

	fmt.Println(host)
	// Output:
	// # My cool git server
	// Host git.example.com
}
