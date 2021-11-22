package goremote

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

type ExecResult struct {
	Stdout string
	Stderr string
	Err    error
}

type SSHClient struct {
	client *ssh.Client
}

func (c *SSHClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

func (c *SSHClient) Run(cmd string) (*ExecResult, error) {
	s, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}

	var bufStderr, bufStdout bytes.Buffer
	s.Stdout = &bufStdout
	s.Stderr = &bufStderr
	err = s.Run(cmd)
	result := &ExecResult{
		Stdout: bufStdout.String(),
		Stderr: bufStderr.String(),
		Err:    err,
	}

	return result, nil
}

func NewSSHClientBuilder() *SSHClientBuilder {
	return &SSHClientBuilder{}
}

type SSHClientBuilder struct {
	key      []byte
	password []byte
	keyPath  string
	user     string
	host     string
	port     int
}

func (b *SSHClientBuilder) WithPrivateKey(key []byte) *SSHClientBuilder {
	b.key = key
	return b
}

func (b *SSHClientBuilder) WithPrivateKeyPath(keyPath string) *SSHClientBuilder {
	b.keyPath = keyPath
	return b
}

func (b *SSHClientBuilder) WithHost(host string) *SSHClientBuilder {
	b.host = host
	return b
}

func (b *SSHClientBuilder) WithKeyPass(pass string) *SSHClientBuilder {
	b.password = []byte(pass)
	return b
}

func (b *SSHClientBuilder) WithPort(port int) *SSHClientBuilder {
	b.port = port
	return b
}

func (b *SSHClientBuilder) WithUser(user string) *SSHClientBuilder {
	b.user = user
	return b
}

func (b *SSHClientBuilder) Build() (*SSHClient, error) {
	if b.host == "" {
		return nil, fmt.Errorf("ssh host cannot be empty")
	}

	if len(b.key) == 0 && b.keyPath == "" {
		b.keyPath = path.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}

	if b.keyPath != "" {
		key, err := ioutil.ReadFile(b.keyPath)
		if err != nil {
			return nil, fmt.Errorf("error building ssh client, invalid key file: %s, err: %+v", b.keyPath, err)
		}

		b.key = key
	}

	var signer ssh.Signer
	var err error
	if len(b.password) == 0 {
		signer, err = ssh.ParsePrivateKey(b.key)
		if err != nil {
			return nil, fmt.Errorf("error pasing private key: %+v", err)
		}
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(b.key, b.password)
		if err != nil {
			return nil, fmt.Errorf("error pasing private key with pass: %+v", err)
		}
	}

	if b.user == "" {
		b.user = "root"
	}

	config := &ssh.ClientConfig{
		User:            b.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	if b.port == 0 {
		b.port = 22
	}

	addr := fmt.Sprintf("%s:%d", b.host, b.port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("error creating ssh client: %+v", err)
	}

	sshClient := &SSHClient{
		client: client,
	}

	return sshClient, nil
}

// func main() {
// 	// user := "root"
// 	host := "192.168.15.12"
// 	// port := 22
// 	privKeyFile := "/Users/wenlinwu/.ssh/id_rsa"
// 	password := ""
// 	builder := NewSSHClientBuilder()
// 	client, err := builder.WithPrivateKeyPath(privKeyFile).WithHost(host).WithKeyPass(password).Build()
// 	if err != nil {
// 		panic(err)
// 	}
// 	ret, err := client.Run("ls -l /root")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("stdout: %+v\n", ret.Stdout)
// 	fmt.Printf("stderr: %+v\n", ret.Stderr)
// 	fmt.Printf("err: %+v\n", ret.Err)
// }
