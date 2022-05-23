package SSH

import (
	_ "encoding/json"
	"fmt"
	_ "fmt"
	"golang.org/x/crypto/ssh"
	_ "golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Discovery(data map[string]interface{}) map[string]interface{} {

	sshPort := int((data["port"]).(float64))

	sshHost := (data["ip"]).(string)

	sshPassword := (data["password"]).(string)

	sshUser := (data["name"]).(string)

	config := &ssh.ClientConfig{

		Timeout: 10 * time.Second,

		User: sshUser,

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Config: ssh.Config{Ciphers: []string{

			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}

	var result = make(map[string]interface{})

	var errorList []string

	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}

	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)

	fmt.Println(addr)

	sshClient, err := ssh.Dial("tcp", addr, config)

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	session, err := sshClient.NewSession()

	res, err := session.Output("uname -n") // available

	if err != nil {
		errorList = append(errorList, err.Error())
	}

	ans := string(res)

	if len(errorList) == 0 {

		result["status"] = "success"

		result["hostname"] = strings.Trim(ans, "\n")

	} else {

		result["status"] = "fail"

		result["error"] = errorList

	}

	return result
}

func standardizeSpaces(s string) string {

	return strings.Join(strings.Fields(s), " ")

}
