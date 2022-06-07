package WINRM

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func Process(data map[string]interface{}) {

	defer func() {

		if r := recover(); r != nil {

			res := make(map[string]interface{})

			res["error"] = r

			errorDisplay(res)

		}
	}()

	host := data["ip"].(string)

	port := int((data["port"]).(float64))

	name := (data["username"]).(string)

	password := (data["password"]).(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)

	var errorList []string

	client, err := winrm.NewClient(endpoint, name, password)

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	_, err = client.CreateShell()

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	if len(errorList) == 0 {

		commandForProcess := "get-process"

		a := "aa"

		process, _, _, err := client.RunPSWithString(commandForProcess, a)

		var processList []map[string]string

		processStringArray := strings.Split(process, "\n")

		flag := 1

		for _, v := range processStringArray {

			if flag <= 3 {

				flag++

				continue

			}

			eachWord := strings.SplitN(standardizeSpaces(v), " ", 8)

			if len(eachWord) <= 7 {

				break

			}

			temp := map[string]string{

				"process.name": eachWord[7],

				"process.id": eachWord[6],

				"process.cpu": eachWord[5],

				"process.virtualMemory": eachWord[4],

				"process.pageableMemory": eachWord[2],

				"process.handles": eachWord[0],
			}

			processList = append(processList, temp)

		}

		var processMap = make(map[string]interface{})

		processMap["process"] = processList

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

func errorDisplay(res map[string]interface{}) {

	bytes, _ := json.Marshal(res)

	fmt.Println(string(bytes))

}
