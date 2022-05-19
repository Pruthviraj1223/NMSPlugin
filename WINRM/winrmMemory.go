package WINRM

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func Memory(data map[string]interface{}) {

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

	commandForMemory := "Get-WmiObject win32_OperatingSystem |%{\"{0} {1} {2} {3}\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory} "

	a := "aa"

	memory, _, _, err := client.RunPSWithString(commandForMemory, a)

	memoryStringArray := strings.Split(standardizeSpaces(memory), " ")

	result := map[string]interface{}{
		"free.memory":          memoryStringArray[1],
		"free.virtual.memory":  memoryStringArray[3],
		"total.memory":         memoryStringArray[0],
		"total.virtual.memory": memoryStringArray[2],
	}

	bytes, _ := json.MarshalIndent(result, " ", " ")

	fmt.Println(string(bytes))

}
