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

	_, err := ssh.Dial("tcp", addr, config)

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	//session, err := sshClient.NewSession()
	//
	//if err != nil {
	//	fmt.Println("err ", err.Error())
	//
	//}
	//
	//defer func(session *ssh.Session) {
	//	err := session.Close()
	//	if err != nil {
	//
	//	}
	//}(session)
	//
	//defer func(Conn ssh.Conn) {
	//	err := Conn.Close()
	//	if err != nil {
	//
	//	}
	//}(sshClient.Conn)
	//
	//_, err = session.Output("uname") // available
	//
	//if err != nil {
	//	errorList = append(errorList, err.Error())
	//}

	if len(errorList) == 0 {

		result["status"] = "success"

	} else {

		result["status"] = "fail"

		result["error"] = errorList

	}

	return result
}

func standardizeSpaces(s string) string {

	return strings.Join(strings.Fields(s), " ")

}
