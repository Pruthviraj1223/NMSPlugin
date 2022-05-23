package WINRM

import (
	"github.com/masterzen/winrm"
	"strings"
)

func Discovery(data map[string]interface{}) map[string]interface{} {

	host := data["ip"].(string)

	port := int((data["port"]).(float64))

	name := (data["name"]).(string)

	password := (data["password"]).(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, name, password)

	_, err = client.CreateShell()

	a := "aa"

	hostname, _, _, err := client.RunPSWithString("hostname", a)

	var errorList []string

	if err != nil {
		errorList = append(errorList, err.Error())
	}

	var result = make(map[string]interface{})

	if len(errorList) == 0 {

		result["status"] = "success"

		result["hostname"] = strings.Trim(hostname, "\r\n")

	} else {

		result["status"] = "fail"

		result["error"] = errorList

	}

	return result
}
