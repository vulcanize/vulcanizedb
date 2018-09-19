package sshconfig

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// Test parsing
func TestParsing(t *testing.T) {
	config := `Host google
  HostName google.se
  User goog
  Port 2222
  ProxyCommand ssh -q pluto nc saturn 22
  HostKeyAlgorithms ssh-dss
  # comment
  IdentityFile ~/.ssh/company

Host face
  HostName facebook.com
  User mark
  Port 22`

	_, err := parse(config)

	if err != nil {
		t.Errorf("unable to parse config: %s", err.Error())
	}

	configCR := strings.Replace(`Host google
  HostName google.se
  User goog
  Port 2222
  ProxyCommand ssh -q pluto nc saturn 22
  HostKeyAlgorithms ssh-dss
  # comment
  IdentityFile ~/.ssh/company

Host face
  HostName facebook.com
  User mark
  Port 22`, "\n", "\r\n", -1)

	_, err = parse(configCR)

	if err != nil {
		t.Errorf("unable to parse config: %s", err.Error())
	}
}

func TestMultipleHost(t *testing.T) {
	config := `Host google google2 aws
  HostName google.se
  User goog
  Port 2222`

	hosts, err := parse(config)

	if err != nil {
		t.Errorf("unable to parse config: %s", err.Error())
	}

	h := hosts[0]
	if ok := reflect.DeepEqual([]string{"google", "google2", "aws"}, h.Host); !ok {
		t.Error("unexpected host mismatch")
	}

}

// TestTrailingSpace ensures the parser does not hang when attempting to parse
// a Host declaration with a trailing space after a pattern
func TestTrailingSpace(t *testing.T) {
	// in the config below, the first line is "Host google \n"
	config := `
Host googlespace 
    HostName google.com
`
	parse(config)
}

func TestIgnoreKeyword(t *testing.T) {
	config := `Host google
  HostName google.se
  User goog
  Port 2222
  ProxyCommand ssh -q pluto nc saturn 22
  HostKeyAlgorithms ssh-dss
  # comment
  IdentityOnly yes
  IdentityFile ~/.ssh/company

Host face
  HostName facebook.com
  User mark
  Port 22`

	expected := []*SSHHost{
		{
			Host:              []string{"google"},
			HostName:          "google.se",
			User:              "goog",
			Port:              2222,
			HostKeyAlgorithms: "ssh-dss",
			ProxyCommand:      "ssh -q pluto nc saturn 22",
			IdentityFile:      "~/.ssh/company",
		},
		{
			Host:              []string{"face"},
			User:              "mark",
			Port:              22,
			HostName:          "facebook.com",
			HostKeyAlgorithms: "",
			ProxyCommand:      "",
			IdentityFile:      "",
		},
	}
	actual, err := parse(config)
	if err != nil {
		t.Errorf("unexpected error parsing config: %s", err.Error())
	}

	compare(t, expected, actual)
}

func compare(t *testing.T, expected, actual []*SSHHost) {
	for i, ac := range actual {
		exMap := toMap(t, expected[i])
		acMap := toMap(t, ac)

		if ok := reflect.DeepEqual(exMap, acMap); !ok {
			t.Errorf("unexpected parsed \n expected: %+v \n actual: %+v", exMap, acMap)
		}
	}
}

func toMap(t *testing.T, a *SSHHost) map[string]interface{} {
	ab, err := json.Marshal(a)
	if err != nil {
		t.Errorf("marshaling expected %s", err)
	}

	var aMap map[string]interface{}
	if err := json.Unmarshal(ab, &aMap); err != nil {
		t.Errorf("unmarshaling expected %s", err)
	}

	return aMap
}
