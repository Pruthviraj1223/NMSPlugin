package WINRM

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func Cpu(data map[string]interface{}) {

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

	commandForCpu := "Get-WmiObject win32_Processor | select DeviceID, SystemName, LoadPercentage | Foreach-Object {$_.DeviceId,$_.SystemName,$_.LoadPercentage -join \" \"}"

	var cpu string

	var cpuList []map[string]string

	a := "aa"

	cpu, _, _, err = client.RunPSWithString(commandForCpu, a)

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

	bytes, _ := json.MarshalIndent(cpuMap, " ", " ")

	fmt.Println(string(bytes))

}
