package SNMP

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
	"time"
)

func Discovery(data map[string]interface{}) map[string]interface{} {

	host := data["ip"].(string)

	port := (data["port"]).(float64)

	community := (data["community"]).(string)

	version := g.Version2c

	switch data["version"] {

	case "v1":
		version = g.Version1
		break

	case "v2":
		version = g.Version2c
		break

	}
	params := &g.GoSNMP{

		Target: host,

		Port: uint16(port),

		Community: community,

		Version: version,

		Timeout: time.Duration(2) * time.Second,
	}

	err := params.Connect()

	defer params.Conn.Close()

	var result = make(map[string]interface{})

	var errorList []string

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	res, err := params.Get([]string{".1.3.6.1.2.1.1.5.0"})

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	for _, outcome := range res.Variables {

		result["systemName"] = string(outcome.Value.([]byte))

	}

	var oidList []string

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.1.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.2.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.8.")

	oidList = append(oidList, ".1.3.6.1.2.1.31.1.1.1.1.")

	oidList = append(oidList, ".1.3.6.1.2.1.31.1.1.1.18.")

	err = params.Walk(".1.3.6.1.2.1.2.2.1.1", walkFunc)

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	var listOfMap []map[string]interface{}

	for i := 0; i < len(list); i++ {

		var newList []string

		for j := 0; j < len(oidList); j++ {

			newList = append(newList, oidList[j]+strconv.Itoa(list[i]))

		}

		ans, _ := params.Get(newList)

		var interfaceMap = make(map[string]interface{})

		for _, outcome := range ans.Variables {

			// copy plugin to that

			if strings.Contains(outcome.Name, "1.31.1") {

				VariableName := strings.SplitAfter(outcome.Name, ".1.3.6.1.2.1.31.1.1.1.")

				strArr := strings.Split(VariableName[1], ".")

				ch, _ := strconv.Atoi(strArr[0])

				fmt.Println("new one ", ch)

				switch ch {

				case 1:
					interfaceMap["interfaceName"] = string(outcome.Value.([]byte))

				case 18:
					interfaceMap["alias"] = string(outcome.Value.([]byte))

				}

			} else {
				VariableName := strings.SplitAfter(outcome.Name, ".1.3.6.1.2.1.2.2.1.")

				strArr := strings.Split(VariableName[1], ".")

				ch, _ := strconv.Atoi(strArr[0])

				fmt.Println("name ", outcome.Name, " val = ", outcome.Value, " strArr = ", ch)

				switch ch {

				case 1:
					interfaceMap["index"] = outcome.Value

				case 2:

					fmt.Println(outcome.Value.([]byte))

					interfaceMap["interface.Description"] = string(outcome.Value.([]byte))

				case 8:

					var operationalStatus string

					if outcome.Value.(int) == 1 {

						operationalStatus = "Up"

					}

					if outcome.Value.(int) == 2 {

						operationalStatus = "Down"

					}

					interfaceMap["interface.operational.status"] = operationalStatus

				}
			}

		}

		listOfMap = append(listOfMap, interfaceMap)

	}

	var dataMap = make(map[string]interface{})

	dataMap["interface"] = listOfMap

	if len(errorList) == 0 {

		result["status"] = "success"

	} else {

		result["status"] = "fail"

		result["error"] = errorList

	}

	return result
}
