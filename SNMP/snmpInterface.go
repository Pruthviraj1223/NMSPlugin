package SNMP

import (
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
	"time"
)

var list []int

func Interface(data map[string]interface{}) {

	host := data["ip"].(string)

	port := int((data["port"]).(float64))

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

	_, err = params.Get([]string{"1.3.6.1.2.1.1.5.0"})

	if err != nil {

		err.Error()

	}

	defer params.Conn.Close()

	var oidList []string

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.2.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.3.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.5.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.6.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.7.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.8.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.10.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.14.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.16.")

	oidList = append(oidList, ".1.3.6.1.2.1.2.2.1.20.")

	err = params.Walk(".1.3.6.1.2.1.2.2.1.1", walkFunc)

	if err != nil {
		return
	}

	var listOfMap []map[string]interface{}

	for i := 0; i < len(list); i++ {

		var newList []string

		for j := 0; j < len(oidList); j++ {

			newList = append(newList, oidList[j]+strconv.Itoa(list[i]))

		}

		ans, _ := params.Get(newList)

		var interfaceMap = make(map[string]interface{})

		fmt.Println("new ", newList)

		for _, result := range ans.Variables {

			fmt.Println("name ", result.Name, " ", result.Value)

			VariableName := strings.SplitAfter(result.Name, ".1.3.6.1.2.1.2.2.1.")

			strArr := strings.Split(VariableName[1], ".")

			ch, _ := strconv.Atoi(strArr[0])

			switch ch {

			case 2:

				interfaceMap["interface.Description"] = string(result.Value.([]byte))

			case 3:

				switch (result.Value).(int) {

				case 6:
					interfaceMap["interface.Type"] = "ethernetCsmacd"
				case 1:
					interfaceMap["interface.Type"] = "other"
				case 135:
					interfaceMap["interface.Type"] = "l2vlan"
				case 53:
					interfaceMap["interface.Type"] = "propVirtual"
				case 24:
					interfaceMap["interface.Type"] = "softwareLoopback"
				case 131:
					interfaceMap["interface.Type"] = "tunnel"
				}

				// 5 6 10 16

			case 5:

				interfaceMap["interface.ifSpeed"] = result.Value

			case 6:

				interfaceMap["interface.ifPhysAddress"] = fmt.Sprintf("%x", result.Value)

			case 7:

				if result.Value.(int) == 1 {

					interfaceMap["interface.admin.status"] = "Up"

				}

				if result.Value.(int) == 2 {

					interfaceMap["interface.admin.status"] = "Down"

				}

			case 8:

				if result.Value.(int) == 1 {

					interfaceMap["interface.operating.status"] = "Up"

				}

				if result.Value.(int) == 2 {

					interfaceMap["interface.operating.status"] = "Down"

				}

			case 10:

				if result.Value != nil {

					interfaceMap["interface.ifInOctets"] = result.Value.(uint)

				} else {

					interfaceMap["interface.ifInOctets"] = ""

				}

			case 14:

				if result.Value == nil {

					interfaceMap["interface.InError"] = ""

				} else {

					interfaceMap["interface.InError"] = result.Value

				}

			case 16:

				if result.Value != nil {

					interfaceMap["interface.ifOutOctets"] = result.Value.(uint)

				} else {

					interfaceMap["interface.ifOutOctets"] = ""

				}

			case 20:

				if (result.Value) == nil {

					interfaceMap["interface.OutError"] = ""

				} else {

					interfaceMap["interface.OutError"] = result.Value

				}
			}
		}
		listOfMap = append(listOfMap, interfaceMap)
	}

	var dataMap = make(map[string]interface{})

	dataMap["interface"] = listOfMap

	bytes, _ := json.Marshal(dataMap)

	fmt.Println(string(bytes))

}

func walkFunc(pdu g.SnmpPDU) error {

	list = append(list, pdu.Value.(int))

	return nil

}

//index,desc,alias,op-status,name
