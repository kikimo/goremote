# goremote

golang implementation of remote execution base on ssh.

example:

```go
package main

import (
	"fmt"

	"github.com/kikimo/goremote"
)

func main() {
	// user := "root"
	// port := 22
	host := "192.168.15.12"
	privKeyFile := "/Users/wenlinwu/.ssh/id_rsa"
	password := ""
	builder := goremote.NewSSHClientBuilder()
	client, err := builder.WithPrivateKeyPath(privKeyFile).WithHost(host).WithKeyPass(password).Build()
	if err != nil {
		panic(err)
	}
	ret, err := client.Run("ls -l /root")
	if err != nil {
		panic(err)
	}

	fmt.Printf("stdout: %+v\n", ret.Stdout)
	fmt.Printf("stderr: %+v\n", ret.Stderr)
	fmt.Printf("err: %+v\n", ret.Err)
}
```


