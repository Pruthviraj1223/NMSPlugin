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

	host := data["ip.address"].(string)

	port := int((data["port"]).(float64))

	community := (data["community"]).(string)

	params := &g.GoSNMP{

		Target: host,

		Port: uint16(port),

		Community: community,

		Version: g.Version2c,

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

	var interfaceMap = make(map[string]interface{})

	var listOfMap []map[string]interface{}

	for i := 0; i < len(list); i++ {

		var newList []string

		for j := 0; j < len(oidList); j++ {
			newList = append(newList, oidList[j]+strconv.Itoa(list[i]))
		}

		ans, _ := params.Get(newList)

		for _, result := range ans.Variables {

			VariableName := strings.SplitAfter(result.Name, ".1.3.6.1.2.1.2.2.1.")

			strArr := strings.Split(VariableName[1], ".")

			ch, _ := strconv.Atoi(strArr[0])

			var typevalue string

			switch ch {

			case 2:

				interfaceMap["interface.Description"] = string(result.Value.([]byte))

			case 3:

				switch (result.Value).(int) {

				case 6:
					typevalue = "ethernetCsmacd"
				case 1:
					typevalue = "other"
				case 135:
					typevalue = "l2vlan"
				case 53:
					typevalue = "propVirtual"
				case 24:
					typevalue = "softwareLoopback"
				case 131:
					typevalue = "tunnel"
				}

				interfaceMap["interface.Type"] = typevalue
				// 5 6 10 16

			case 5:

			case 6:

			case 7:
				var Adminvalue string

				if result.Value.(int) == 1 {

					Adminvalue = "Up"

				}

				if result.Value.(int) == 2 {

					Adminvalue = "Down"

				}

				interfaceMap["interface.admin.status"] = Adminvalue

			case 8:

				var operatingvalue string

				if result.Value.(int) == 1 {

					operatingvalue = "Up"

				}

				if result.Value.(int) == 2 {

					operatingvalue = "Down"

				}

				interfaceMap["interface.operating.status"] = operatingvalue

			case 10:

			case 14:

				if result.Value == nil {

					interfaceMap["interface.InError"] = ""

				} else {

					interfaceMap["interface.InError"] = result.Value

				}

			case 16:

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

	bytes, _ := json.MarshalIndent(dataMap, " ", " ")

	fmt.Println(string(bytes))

}

func walkFunc(pdu g.SnmpPDU) error {

	list = append(list, pdu.Value.(int))

	return nil

}

//index,desc,alias,op-status,name
