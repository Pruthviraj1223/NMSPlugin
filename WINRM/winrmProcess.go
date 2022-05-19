package WINRM

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func Process(data map[string]interface{}) {

	host := data["ip.address"].(string)

	port := int((data["port"]).(float64))

	name := (data["name"]).(string)

	password := (data["password"]).(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, name, password)

	_, err = client.CreateShell()

	if err != nil {

		err.Error()

	}

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

	bytes, _ := json.MarshalIndent(processMap, " ", " ")

	fmt.Println(string(bytes))
}