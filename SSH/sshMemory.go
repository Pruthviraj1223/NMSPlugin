package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Memory(data map[string]interface{}) {

	defer func() {

		if r := recover(); r != nil {

			res := make(map[string]interface{})

			res["error"] = r

			errorDisplay(res)

		}
	}()

	sshUser := (data["name"]).(string)

	sshPassword := (data["password"]).(string)

	sshHost := (data["ip"]).(string)

	sshPort := int((data["port"]).(float64))

	config := &ssh.ClientConfig{

		Timeout: 10 * time.Second,

		User: sshUser,

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Config: ssh.Config{Ciphers: []string{

			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}

	address := fmt.Sprintf("%s:%d", sshHost, sshPort)

	sshClient, err := ssh.Dial("tcp", address, config)

	var errorList []string

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()

	available, err := session.Output("free -b | grep Mem | awk '{ printf(\"%.4f\\n\", $7) }'") // available

	session.Close()

	session, err = sshClient.NewSession()

	freeMemory, err := session.Output("free -b | grep Mem | awk '{ printf(\"%.4f\\n\", $4) }'") //free memory
	session.Close()

	session, err = sshClient.NewSession()

	usedMemory, err := session.Output("free -b | grep Mem | awk '{ printf(\"%.4f \\n\", $3) }'") // used memory
	session.Close()

	session, err = sshClient.NewSession()

	totalMemory, err := session.Output("free -b | grep Mem | awk '{ printf(\"%i\\n\", $2) }'") // used memory
	session.Close()

	session, err = sshClient.NewSession()

	swapUsed, err := session.Output("free -b | grep Swap | awk '{ printf(\"%i\\n\", $3) }'")
	session.Close()

	session, err = sshClient.NewSession()

	swapFree, err := session.Output("free -b | grep Swap | awk '{ printf(\"%i\\n\", $3) }'")

	session, err = sshClient.NewSession()

	swapTotal, err := session.Output("free -b | grep Swap | awk '{ printf(\"%i\\n\", $3) }'")
	session.Close()

	session, err = sshClient.NewSession() // changes

	result := map[string]interface{}{

		"memory.free.bytes": strings.Trim(string(freeMemory), "\n"),

		"memory.used.bytes": strings.Trim(string(usedMemory), "\n"),

		"memory.total.bytes": strings.Trim(string(totalMemory), "\n"),

		"memory.available.bytes": strings.Trim(string(available), "\n"),

		"memory.swap.total.bytes": strings.Trim(string(swapTotal), "\n"),

		"memory.swap.used.bytes": strings.Trim(string(swapUsed), "\n"),

		"memory.swap.free.bytes": strings.Trim(string(swapFree), "\n"),
	}

	bytes, err := json.Marshal(result)

	if err != nil {

		response := make(map[string]interface{})

		response["error"] = err.Error()

		errorDisplay(response)

	} else {

		fmt.Println(string(bytes))

	}

}
