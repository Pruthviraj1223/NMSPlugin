package WINRM

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func System(data map[string]interface{}) {

	host := data["ip"].(string)

	port := int((data["port"]).(float64))

	name := (data["name"]).(string)

	password := (data["password"]).(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, name, password)

	_, err = client.CreateShell()

	if err != nil {

		err.Error()

	}

	systemMap := make(map[string]interface{})

	commandForName := "(Get-WmiObject win32_operatingsystem).name"

	a := "aa"

	sysName, _, _, err := client.RunPSWithString(commandForName, a)

	systemMap["SystemName"] = strings.Replace(sysName, "\r\n", " ", 3)

	commandForVersion := "(Get-WMIObject win32_operatingsystem).version"

	sysVersion, _, _, err := client.RunPSWithString(commandForVersion, a)

	systemMap["systemVersion"] = strings.Replace(sysVersion, "\r\n", " ", 3)

	name1 := "whoami"

	uname, _, _, err := client.RunPSWithString(name1, a)

	systemMap["uname"] = strings.Replace(uname, "\r\n", " ", 3)

	sysUpTime := "(Get-WMIObject win32_operatingsystem).LastBootUpTime;"

	sysTime, _, _, err := client.RunPSWithString(sysUpTime, a)

	systemMap["systemUpTime"] = strings.Replace(sysTime, "\r\n", " ", 3)

	bytes, _ := json.Marshal(systemMap)

	fmt.Println(string(bytes))

}
