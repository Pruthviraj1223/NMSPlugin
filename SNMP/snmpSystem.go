package SNMP

import (
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"log"
	"strings"
	"time"
)

func System(data map[string]interface{}) {

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

	if err != nil {
		err.Error()
	}

	defer params.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.5.0", "1.3.6.1.2.1.1.6.0", "1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.2.0", "1.3.6.1.2.1.1.3.0", "1.3.6.1.2.1.1.3.0"} // sysName , sysLocation, sysDiscription,sysOID,

	result, err2 := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	var sysName string
	var sysLocation string
	var sysDescription string
	var sysOID string
	var sysUpTime string

	var systemMap = make(map[string]interface{})

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
			systemMap["system.description"] = strings.Trim(sysDescription, "\r\n")

		case ".1.3.6.1.2.1.1.2.0":
			sysOID = fmt.Sprintf("%v", variable.Value)
			systemMap["system.OID"] = sysOID

		case ".1.3.6.1.2.1.1.3.0":
			sysUpTime = fmt.Sprintf("%v", variable.Value)
			systemMap["system.uptime"] = sysUpTime

		}
	}

	bytes, _ := json.Marshal(systemMap)

	fmt.Println(string(bytes))

}
