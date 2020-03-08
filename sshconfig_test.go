package sshconfig

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sshConfigExample = `# ssh config generated by some go code (github.com/emptyinterface/ssh_config)

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

	devHostBlock = `
# dev
Host dev
  HostName 127.0.0.1
  User ubuntu
  Port 22
`
)

func TestParseAndWriteTo(t *testing.T) {

	config, err := Parse(strings.NewReader(sshConfigExample))
	if err != nil {
		t.Error(err)
	}

	buf := &bytes.Buffer{}

	config.WriteTo(buf)

	assert.Equal(t, sshConfigExample, buf.String())
}

func TestGetParam(t *testing.T) {

	config, err := Parse(strings.NewReader(sshConfigExample))

	assert.NoError(t, err)

	visualHostKey := config.GetParam(VisualHostKeyKeyword)

	assert.Equal(t, visualHostKey.Value(), "yes")
}

func TestFindByHostname(t *testing.T) {

	config, err := Parse(strings.NewReader(sshConfigExample))

	assert.NoError(t, err)

	devHost := config.FindByHostname("dev")

	assert.Equal(t, devHost.String(), devHostBlock)
}
