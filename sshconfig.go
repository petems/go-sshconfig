package sshconfig

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	// Config struct for the entire SSH config file
	Config struct {
		Source  []byte
		Globals []*Param
		Hosts   []*Host
	}
	// Host struct for host entries
	Host struct {
		Comments  []string
		Hostnames []string
		Params    []*Param
	}
	// Param struct for parameters for configuration
	Param struct {
		Comments []string
		Keyword  string
		Args     []string
	}
)

// All SSH config configuration options
// http://man7.org/linux/man-pages/man5/ssh_config.5.html#top_of_page
const (
	HostKeyword                             = "Host"
	MatchKeyword                            = "Match"
	AddressFamilyKeyword                    = "AddressFamily"
	BatchModeKeyword                        = "BatchMode"
	BindAddressKeyword                      = "BindAddress"
	CanonicalDomainsKeyword                 = "CanonicalDomains"
	CanonicalizeFallbackLocalKeyword        = "CanonicalizeFallbackLocal"
	CanonicalizeHostnameKeyword             = "CanonicalizeHostname"
	CanonicalizeMaxDotsKeyword              = "CanonicalizeMaxDots"
	CanonicalizePermittedCNAMEsKeyword      = "CanonicalizePermittedCNAMEs"
	ChallengeResponseAuthenticationKeyword  = "ChallengeResponseAuthentication"
	CheckHostIPKeyword                      = "CheckHostIP"
	CipherKeyword                           = "Cipher"
	CiphersKeyword                          = "Ciphers"
	ClearAllForwardingsKeyword              = "ClearAllForwardings"
	CompressionKeyword                      = "Compression"
	CompressionLevelKeyword                 = "CompressionLevel"
	ConnectionAttemptsKeyword               = "ConnectionAttempts"
	ConnectTimeoutKeyword                   = "ConnectTimeout"
	ControlMasterKeyword                    = "ControlMaster"
	ControlPathKeyword                      = "ControlPath"
	ControlPersistKeyword                   = "ControlPersist"
	DynamicForwardKeyword                   = "DynamicForward"
	EnableSSHKeysignKeyword                 = "EnableSSHKeysign"
	EscapeCharKeyword                       = "EscapeChar"
	ExitOnForwardFailureKeyword             = "ExitOnForwardFailure"
	FingerprintHashKeyword                  = "FingerprintHash"
	ForwardAgentKeyword                     = "ForwardAgent"
	ForwardX11Keyword                       = "ForwardX11"
	ForwardX11TimeoutKeyword                = "ForwardX11Timeout"
	ForwardX11TrustedKeyword                = "ForwardX11Trusted"
	GatewayPortsKeyword                     = "GatewayPorts"
	GlobalKnownHostsFileKeyword             = "GlobalKnownHostsFile"
	GSSAPIAuthenticationKeyword             = "GSSAPIAuthentication"
	GSSAPIDelegateCredentialsKeyword        = "GSSAPIDelegateCredentials"
	HashKnownHostsKeyword                   = "HashKnownHosts"
	HostbasedAuthenticationKeyword          = "HostbasedAuthentication"
	HostbasedKeyTypesKeyword                = "HostbasedKeyTypes"
	HostKeyAlgorithmsKeyword                = "HostKeyAlgorithms"
	HostKeyAliasKeyword                     = "HostKeyAlias"
	HostNameKeyword                         = "HostName"
	IdentitiesOnlyKeyword                   = "IdentitiesOnly"
	IdentityFileKeyword                     = "IdentityFile"
	IgnoreUnknownKeyword                    = "IgnoreUnknown"
	IPQoSKeyword                            = "IPQoS"
	KbdInteractiveAuthenticationKeyword     = "KbdInteractiveAuthentication"
	KbdInteractiveDevicesKeyword            = "KbdInteractiveDevices"
	KexAlgorithmsKeyword                    = "KexAlgorithms"
	LocalCommandKeyword                     = "LocalCommand"
	LocalForwardKeyword                     = "LocalForward"
	LogLevelKeyword                         = "LogLevel"
	MACsKeyword                             = "MACs"
	NoHostAuthenticationForLocalhostKeyword = "NoHostAuthenticationForLocalhost"
	NumberOfPasswordPromptsKeyword          = "NumberOfPasswordPrompts"
	PasswordAuthenticationKeyword           = "PasswordAuthentication"
	PermitLocalCommandKeyword               = "PermitLocalCommand"
	PKCS11ProviderKeyword                   = "PKCS11Provider"
	PortKeyword                             = "Port"
	PreferredAuthenticationsKeyword         = "PreferredAuthentications"
	ProtocolKeyword                         = "Protocol"
	ProxyCommandKeyword                     = "ProxyCommand"
	ProxyUseFdpassKeyword                   = "ProxyUseFdpass"
	PubkeyAuthenticationKeyword             = "PubkeyAuthentication"
	RekeyLimitKeyword                       = "RekeyLimit"
	RemoteForwardKeyword                    = "RemoteForward"
	RequestTTYKeyword                       = "RequestTTY"
	RevokedHostKeysKeyword                  = "RevokedHostKeys"
	RhostsRSAAuthenticationKeyword          = "RhostsRSAAuthentication"
	RSAAuthenticationKeyword                = "RSAAuthentication"
	SendEnvKeyword                          = "SendEnv"
	ServerAliveCountMaxKeyword              = "ServerAliveCountMax"
	ServerAliveIntervalKeyword              = "ServerAliveInterval"
	StreamLocalBindMaskKeyword              = "StreamLocalBindMask"
	StreamLocalBindUnlinkKeyword            = "StreamLocalBindUnlink"
	StrictHostKeyCheckingKeyword            = "StrictHostKeyChecking"
	TCPKeepAliveKeyword                     = "TCPKeepAlive"
	TunnelKeyword                           = "Tunnel"
	TunnelDeviceKeyword                     = "TunnelDevice"
	UpdateHostKeysKeyword                   = "UpdateHostKeys"
	UsePrivilegedPortKeyword                = "UsePrivilegedPort"
	UserKeyword                             = "User"
	UserKnownHostsFileKeyword               = "UserKnownHostsFile"
	VerifyHostKeyDNSKeyword                 = "VerifyHostKeyDNS"
	VisualHostKeyKeyword                    = "VisualHostKey"
	XAuthLocationKeyword                    = "XAuthLocation"

	GlobalConfigurationHeader = "# global configuration"
	HostConfigurationHeader   = "# host-based configuration"
)

// NewHost creates a new parameter based on the main objects: the hostnames and comments
func NewHost(hostnames []string, comments []string) *Host {
	return &Host{
		Comments:  comments,
		Hostnames: hostnames,
	}
}

func (host *Host) String() string {

	buf := &bytes.Buffer{}

	fmt.Fprintln(buf)
	if len(host.Comments) > 0 {
		for _, comment := range host.Comments {
			if !strings.HasPrefix(comment, "#") {
				comment = "# " + comment
			}
			fmt.Fprintln(buf, comment)
		}
	}

	fmt.Fprintf(buf, "%s %s\n", HostKeyword, strings.Join(host.Hostnames, " "))
	for _, param := range host.Params {
		fmt.Fprint(buf, "  ", param.String())
	}

	return buf.String()

}

// NewParam creates a new parameter based on the main objects: the keyword, the argument and a comment
func NewParam(keyword string, args []string, comments []string) *Param {
	return &Param{
		Comments: comments,
		Keyword:  keyword,
		Args:     args,
	}
}

func (param *Param) String() string {

	buf := &bytes.Buffer{}

	if len(param.Comments) > 0 {
		fmt.Fprintln(buf)
		for _, comment := range param.Comments {
			if !strings.HasPrefix(comment, "#") {
				comment = "# " + comment
			}
			fmt.Fprintln(buf, comment)
		}
	}

	fmt.Fprintf(buf, "%s %s\n", param.Keyword, strings.Join(param.Args, " "))

	return buf.String()

}

// Value returns the current value for a given parameter
func (param *Param) Value() string {
	if len(param.Args) > 0 {
		return param.Args[0]
	}
	return ""
}

// Parse is the main guts of the library
// It reads from a given io.Reader and parses it
// into a ssh config object
func Parse(r io.Reader) (*Config, error) {

	// dat state
	var (
		global = true

		param = &Param{}
		host  *Host
	)

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Source: data,
	}

	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {

		line := strings.TrimSpace(sc.Text())
		if len(line) == 0 {
			continue
		}

		switch line {
		case GlobalConfigurationHeader,
			HostConfigurationHeader:
			continue
		}

		if line[0] == '#' {
			param.Comments = append(param.Comments, line)
			continue
		}

		psc := bufio.NewScanner(strings.NewReader(line))
		psc.Split(bufio.ScanWords)
		if !psc.Scan() {
			continue
		}

		param.Keyword = psc.Text()

		for psc.Scan() {
			param.Args = append(param.Args, psc.Text())
		}

		if param.Keyword == HostKeyword {
			global = false
			if host != nil {
				config.Hosts = append(config.Hosts, host)
			}
			host = &Host{
				Comments:  param.Comments,
				Hostnames: param.Args,
			}
			param = &Param{}
			continue
		} else if global {
			config.Globals = append(config.Globals, param)
			param = &Param{}
			continue
		}

		host.Params = append(host.Params, param)
		param = &Param{}

	}

	if global {
		config.Globals = append(config.Globals, param)
	} else if host != nil {
		config.Hosts = append(config.Hosts, host)
	}

	return config, nil

}

// WriteTo writes to an io.Writer the config file
// This is useful for outputting an SSH config to stdout
func (config *Config) WriteTo(w io.Writer) error {

	fmt.Fprintln(w)
	fmt.Fprintln(w, GlobalConfigurationHeader)

	for _, param := range config.Globals {
		fmt.Fprint(w, param.String())
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, HostConfigurationHeader)

	for _, host := range config.Hosts {
		fmt.Fprint(w, host.String())
	}

	return nil
}

// WriteToFilepath creates a file on disk at a given path from a given sshconfig object
func (config *Config) WriteToFilepath(filePath string) error {

	// create a tmp file in the same path with the same mode
	tmpFilePath := filePath + "." + strconv.FormatInt(time.Now().UnixNano(), 10)

	var mode os.FileMode = 0600
	if stat, err := os.Stat(filePath); err == nil {
		mode = stat.Mode()
	}

	file, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL|os.O_SYNC, mode)
	if err != nil {
		return err
	}

	if err := config.WriteTo(file); err != nil {
		file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpFilePath, filePath); err != nil {
		return err
	}

	return nil

}

// GetParam returns a global parameter from an SSH config file
func (config *Config) GetParam(keyword string) *Param {
	for _, param := range config.Globals {
		if param.Keyword == keyword {
			return param
		}
	}
	return nil
}

// GetHost returns a host from an SSH config file
func (config *Config) GetHost(hostname string) *Host {
	for _, host := range config.Hosts {
		for _, hn := range host.Hostnames {
			if hn == hostname {
				return host
			}
		}
	}
	return nil
}

// GetParam returns a parameter for a specific host
func (host *Host) GetParam(keyword string) *Param {
	for _, param := range host.Params {
		if param.Keyword == keyword {
			return param
		}
	}
	return nil
}

// FindByHostname takes a string argument of a host name
// then searches through the ssh config for that host
// 	config.FindByHostname("github.com")
func (config *Config) FindByHostname(hostname string) *Host {
	for _, host := range config.Hosts {
		for _, hn := range host.Hostnames {
			if hn == hostname {
				return host
			}
		}
		if hns := host.GetParam(HostNameKeyword); hns != nil {
			for _, hn := range hns.Args {
				if hn == hostname {
					return host
				}
			}
		}
	}
	return nil
}
