package WINRM

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func Cpu(data map[string]interface{}) {

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

	client, err := winrm.NewClient(endpoint, name, password)

	var errorList []string

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	if len(errorList) == 0 {

		commandForCpu := "Get-WmiObject win32_Processor | select DeviceID, SystemName, LoadPercentage | Foreach-Object {$_.DeviceId,$_.SystemName,$_.LoadPercentage -join \" \"}"

		var cpu string

		var cpuList []map[string]string

		a := "aa"

		cpu, _, _, err := client.RunPSWithString(commandForCpu, a)

		cpuStringArray := strings.Split(cpu, "\n")

		for _, v := range cpuStringArray {

			if len(cpuStringArray) == 0 {

				break

			}

			eachWord := strings.Split(standardizeSpaces(v), " ")

			if len(eachWord) <= 2 {

				break

			}

			temp := map[string]string{

				"cpu.name": eachWord[0],

				"cpu.system.name": eachWord[1],

				"cpu.load.percentage": eachWord[2],
			}

			cpuList = append(cpuList, temp)

		}

		var cpuMap = make(map[string]interface{})

		cpuMap["CPU"] = cpuList

		bytes, err := json.Marshal(cpuMap)

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
