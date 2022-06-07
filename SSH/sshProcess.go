package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Process(data map[string]interface{}) {

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

	if len(errorList) == 0 {

		psaux, _ := session.Output("ps aux")

		psauxString := string(psaux)

		myStringArray := strings.Split(psauxString, "\n")

		var processList []map[string]string

		flag := 1

		for _, v := range myStringArray {

			if flag == 1 {

				flag = 0

				continue

			}

			eachWorld := strings.SplitN(standardizeSpaces(v), " ", 11)

			if len(eachWorld) <= 10 {

				break

			}

			temp1 := map[string]string{

				"process.user": eachWorld[0],

				"process.id": eachWorld[1],

				"process.cpu": eachWorld[2],

				"process.memory": eachWorld[3],
			}

			processList = append(processList, temp1)

		}

		processMap := make(map[string]interface{})

		processMap["Process"] = processList

		bytes, err := json.Marshal(processMap)

		if err != nil {

			response := make(map[string]interface{})

			response["error"] = err.Error()

			errorDisplay(response)

		} else {

			fmt.Println(string(bytes))

		}

	} else {

		res := make(map[string]interface{})

		res["error"] = errorList

		errorDisplay(res)

	}

}
