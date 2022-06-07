package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"time"
)

func Disk(data map[string]interface{}) {

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

		diskMap := make(map[string]interface{})

		var diskList []map[string]interface{}

		diskData, _ := session.Output("df")

		diskUtilizationString := string(diskData)

		diskStringArray := strings.Split(diskUtilizationString, "\n")

		count := 1

		for _, v := range diskStringArray {

			if count == 1 {

				count++

				continue
			}

			eachWord := strings.Split(standardizeSpaces(v), " ")

			if len(eachWord) <= 5 {

				continue
			}

			usePercentString := strings.Trim(eachWord[4], "-%")

			usePercent, _ := strconv.Atoi(usePercentString)

			freePercent := 100 - usePercent

			usedBytes, _ := strconv.Atoi(eachWord[2])

			available, _ := strconv.Atoi(eachWord[3])

			temp := make(map[string]interface{})

			temp["disk.name"] = eachWord[0]

			temp["disk.total.bytes"] = eachWord[1]

			temp["disk.used.bytes"] = usedBytes

			temp["disk.available.bytes"] = available

			temp["disk.used.percent"] = usePercent

			temp["disk.free.percent"] = freePercent

			diskList = append(diskList, temp)

		}

		diskMap["disk"] = diskList

		bytes, err := json.Marshal(diskMap)

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
