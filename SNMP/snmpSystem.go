package SNMP

import (
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strings"
	"time"
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

	community := (data["community"]).(string)

	params := &g.GoSNMP{

		Target: host,

		Port: uint16(port),

		Community: community,

		Version: g.Version2c,

		Timeout: time.Duration(2) * time.Second,
	}

	err := params.Connect()

	var errorList []string

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	defer params.Conn.Close()

	var sysName string

	var sysLocation string

	var sysDescription string

	var sysOID string

	var sysUpTime string

	var systemMap = make(map[string]interface{})

	oids := []string{"1.3.6.1.2.1.1.5.0", "1.3.6.1.2.1.1.6.0", "1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.2.0", "1.3.6.1.2.1.1.3.0", "1.3.6.1.2.1.1.3.0"} // sysName , sysLocation, sysDiscription,sysOID,

	result, err := params.Get(oids)

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	if len(errorList) == 0 {

		for _, variable := range result.Variables {

			switch variable.Name {

			case ".1.3.6.1.2.1.1.5.0":

				sysName = string(variable.Value.([]byte))

				systemMap["system.name"] = sysName

			case ".1.3.6.1.2.1.1.6.0":

				sysLocation = string(variable.Value.([]byte))

				systemMap["system.location"] = sysLocation

			case ".1.3.6.1.2.1.1.1.0":

				sysDescription = string(variable.Value.([]byte))

				systemMap["system.description"] = strings.Replace(sysDescription, "\r\n", " ", 3)

			case ".1.3.6.1.2.1.1.2.0":

				sysOID = fmt.Sprintf("%v", variable.Value)

				systemMap["system.oid"] = sysOID

			case ".1.3.6.1.2.1.1.3.0":

				sysUpTime = fmt.Sprintf("%v", variable.Value)

				systemMap["system.uptime"] = sysUpTime

			}
		}

		bytes, err := json.Marshal(systemMap)

		if err != nil {

			response := make(map[string]interface{})

			response["error"] = err.Error()

			errorDisplay(response)

		} else {

			fmt.Println(string(bytes))

		}

	} else {

		response := make(map[string]interface{})

		response["error"] = errorList

		errorDisplay(response)

	}

}
