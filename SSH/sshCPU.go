package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Cpu(data map[string]interface{}) {

	defer func() {

		if r := recover(); r != nil {

			res := make(map[string]interface{})

			res["error"] = r

			errorDisplay(res)

		}

	}()

	time.Sleep(60 * time.Second)

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
	var errorList []string

	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}

	address := fmt.Sprintf("%s:%d", sshHost, sshPort)

	sshClient, err := ssh.Dial("tcp", address, config)

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()

	if err != nil {
		errorList = append(errorList, err.Error())
	}

	var CPU = make(map[string]interface{})

	if len(errorList) == 0 {

		cpuUtilization, _ := session.CombinedOutput("mpstat -P ALL")

		var cpuList []map[string]string

		cpuUtilizationString := string(cpuUtilization)

		cpuStringArray := strings.Split(cpuUtilizationString, "\n")

		flag1 := 1

		for _, v := range cpuStringArray {

			if flag1 <= 3 {

				flag1++

				continue

			}

			eachWord := strings.Split(standardizeSpaces(v), " ")

			if len(eachWord) <= 13 {

				continue

			}

			if eachWord[3] == "all" {

				CPU["cpu.name"] = eachWord[3]

				CPU["cpu.user.percent"] = eachWord[4]

				CPU["cpu.sys.percent"] = eachWord[6]

				CPU["cpu.idle.percent"] = eachWord[13]

			} else {

				temp1 := map[string]string{

					"cpu.name": eachWord[3],

					"cpu.user.percent": eachWord[4],

					"cpu.sys.percent": eachWord[6],

					"cpu.idle.percent": eachWord[13],
				}

				cpuList = append(cpuList, temp1)

			}
		}

		CPU["CPU"] = cpuList

		bytes, err := json.Marshal(CPU)

		if err != nil {

			response := make(map[string]interface{})

			response["error"] = err.Error()

			errorDisplay(response)

		} else {

			fmt.Println(string(bytes))

		}

	} else {

		response := make(map[string]interface{})

		response["error"] = errorList

		errorDisplay(response)

	}

}

func errorDisplay(res map[string]interface{}) {

	bytes, _ := json.Marshal(res)

	fmt.Println(string(bytes))

}
