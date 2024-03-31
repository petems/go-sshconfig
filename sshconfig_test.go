package sshconfig

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sshConfigTest = `
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

	devHostBlockTest = `
# dev
Host dev
  HostName 127.0.0.1
  User ubuntu
  Port 22
`

	githubBlockTest = `
# github.com global config
Host github.com
  ControlMaster auto
  ControlPath ~/.ssh/ssh-%r@%h:%p
  ControlPersist yes
  User git
`
)

func TestParseAndWriteTo(t *testing.T) {

	config, err := Parse(strings.NewReader(sshConfigTest))
	if err != nil {
		t.Error(err)
	}

	buf := &bytes.Buffer{}

	config.WriteTo(buf)

	assert.Equal(t, sshConfigTest, buf.String())
}

func TestGetParam(t *testing.T) {

	config, err := Parse(strings.NewReader(sshConfigTest))

	assert.NoError(t, err)

	visualHostKey := config.GetParam(VisualHostKeyKeyword)

	assert.Equal(t, visualHostKey.Value(), "yes")
}

func TestHostParamString(t *testing.T) {

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

	assert.Equal(t, githubhost.String(), githubBlockTest)
}

func TestFindByHostname_Host(t *testing.T) {

	config, err := Parse(strings.NewReader(sshConfigTest))

	assert.NoError(t, err)

	devHost := config.FindByHostname("dev")

	assert.Equal(t, devHost.String(), devHostBlockTest)
}

func TestFindByHostname_InParam(t *testing.T) {

	config, err := Parse(strings.NewReader(sshConfigTest))

	assert.NoError(t, err)

	devHost := config.FindByHostname("127.0.0.1")

	assert.Equal(t, devHost.String(), devHostBlockTest)

}

func TestWriteTo(t *testing.T) {
	config, err := Parse(strings.NewReader(sshConfigTest))

	assert.NoError(t, err)

	var b bytes.Buffer

	writtenCount, err := config.WriteTo(&b)

	assert.NoError(t, err)

	assert.Equal(t, sshConfigTest, b.String())
	assert.Equal(t, writtenCount, int64(156))
}

func TestWriteToWithNewParam(t *testing.T) {
	config, err := Parse(strings.NewReader(sshConfigTest))

	host := config.FindByHostname("dev")
	param := host.GetParam("User")
	param.Args = []string{"ec2-user"}

	assert.NoError(t, err)

	var b bytes.Buffer

	writtenCount, err := config.WriteTo(&b)

	assert.NoError(t, err)

	expected := `
# global configuration
VisualHostKey yes

# host-based configuration

# dev
Host dev
  HostName 127.0.0.1
  User ec2-user
  Port 22

Host *.google.com *.yahoo.com
  User root
`

	assert.Equal(t, expected, b.String())
	assert.Equal(t, writtenCount, int64(158))
}

func TestWriteToFilepath(t *testing.T) {
	config, err := Parse(strings.NewReader(sshConfigTest))

	assert.NoError(t, err)

	err = config.WriteToFilepath("./example_config")

	assert.NoError(t, err)

	exampleConfigContents, err := ioutil.ReadFile("./example_config")

	assert.NoError(t, err)

	assert.Equal(t, sshConfigTest, string(exampleConfigContents))
}
