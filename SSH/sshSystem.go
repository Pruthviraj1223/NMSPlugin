package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func System(data map[string]interface{}) {

	defer func() {

		if r := recover(); r != nil {

			res := make(map[string]interface{})

			res["error"] = r

			errorDisplay(res)

		}
	}()

	sshUser := (data["username"]).(string)

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

	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)

	sshClient, err := ssh.Dial("tcp", addr, config)

	var errorList []string

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	defer session.Close()

	var systemMap = make(map[string]interface{})

	if len(errorList) == 0 {

		res, _ := session.Output("uname -n && vmstat")

		splittedString := strings.Split(string(res), "\n")

		flag := 1

		for _, v := range splittedString {

			if flag == 1 {

				systemMap["system.user.name"] = v

				flag++

			}

			if flag == 2 || flag == 3 || v == "" {

				flag++

				continue

			}

			output := strings.SplitN(standardizeSpaces(v), " ", 17)

			systemMap["system.running.process"] = output[0]

			systemMap["system.blocking.process"] = output[1]

			systemMap["system.context.switches"] = output[11]

		}

		bytes, _ := json.Marshal(systemMap)

		fmt.Println(string(bytes))

	} else {

		var response = make(map[string]interface{})

		response["error"] = errorList

		errorDisplay(response)

	}

}
