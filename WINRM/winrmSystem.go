package WINRM

import (
	"fmt"
	"github.com/masterzen/winrm"
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

	systemMap["SystemName"] = sysName

	commandForVersion := "(Get-WMIObject win32_operatingsystem).version"

	sysVersion, _, _, err := client.RunPSWithString(commandForVersion, a)

	systemMap["systemVersion"] = sysVersion

	name1 := "whoami"

	uname, _, _, err := client.RunPSWithString(name1, a)

	systemMap["uname"] = uname

	sysUpTime := "(Get-WMIObject win32_operatingsystem).LastBootUpTime;"

	sysTime, _, _, err := client.RunPSWithString(sysUpTime, a)

	systemMap["systemUpTime"] = sysTime

	fmt.Println(systemMap)

}
