package WINRM

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func Memory(data map[string]interface{}) {

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

	if len(errorList) == 0 {

		commandForMemory := "Get-WmiObject win32_OperatingSystem |%{\"{0} {1} {2} {3}\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory} "

		memory, _, _, err := client.RunPSWithString(commandForMemory, "")

		memoryStringArray := strings.Split(standardizeSpaces(memory), " ")

		result := map[string]interface{}{

			"free.memory.bytes": memoryStringArray[1],

			"free.virtual.memory.bytes": memoryStringArray[3],

			"total.memory.bytes": memoryStringArray[0],

			"total.virtual.memory.bytes": memoryStringArray[2],
		}

		bytes, err := json.Marshal(result)

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
