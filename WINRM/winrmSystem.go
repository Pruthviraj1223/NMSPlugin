package WINRM

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func System(data map[string]interface{}) {

	defer func() {

		if r := recover(); r != nil {

			res := make(map[string]interface{})

			res["error"] = r

			errorDisplay(res)

		}

	}()

	host := data["ip"].(string)

	port := int((data["port"]).(float64))

	name := (data["name"]).(string)

	password := (data["password"]).(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, name, password)

	systemMap := make(map[string]interface{})

	commandForName := "(Get-WmiObject win32_operatingsystem).name"

	a := "aa"

	sysName, _, _, err := client.RunPSWithString(commandForName, a)

	updateName := strings.Replace(sysName, "\r\n", " ", -1)

	systemMap["system.name"] = strings.Replace(updateName, "\\", "", -1)

	commandForVersion := "(Get-WMIObject win32_operatingsystem).version"

	sysVersion, _, _, err := client.RunPSWithString(commandForVersion, a)

	systemMap["system.version"] = strings.Replace(sysVersion, "\r\n", " ", -1)

	username := "whoami"

	uname, _, _, err := client.RunPSWithString(username, a)

	systemMap["uname"] = strings.Replace(strings.Replace(uname, "\r\n", " ", -1), "\\", "", -1)

	sysUpTime := "(Get-WMIObject win32_operatingsystem).LastBootUpTime;"

	sysTime, _, _, err := client.RunPSWithString(sysUpTime, a)

	systemMap["system.up.time"] = strings.Replace(sysTime, "\r\n", " ", -1)

	bytes, err := json.Marshal(systemMap)

	if err != nil {

		response := make(map[string]interface{})

		response["error"] = err.Error()

		errorDisplay(response)

	} else {

		fmt.Println(string(bytes))

	}

}
