package WINRM

import (
	"github.com/masterzen/winrm"
	"strings"
)

func Discovery(data map[string]interface{}) map[string]interface{} {

	host := data["ip"].(string)

	port := int((data["port"]).(float64))

	name := (data["username"]).(string)

	password := (data["password"]).(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)

	var result = make(map[string]interface{})

	client, err := winrm.NewClient(endpoint, name, password)

	if err != nil {

		result["status"] = "fail"

		result["error"] = err.Error()

		return result
	}

	_, err = client.CreateShell()

	if err != nil {

		result["status"] = "fail"

		result["error"] = err.Error()

		return result
	}

	a := "aa"

	hostname, _, _, err := client.RunPSWithString("hostname", a)

	if err != nil {

		result["status"] = "fail"

		result["error"] = err.Error()

		return result
	}

	result["status"] = "success"

	result["hostname"] = strings.Trim(hostname, "\r\n")

	return result
}
