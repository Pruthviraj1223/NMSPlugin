package main

import (
	"NMSPlugin/SNMP"
	"NMSPlugin/SSH"
	"NMSPlugin/WINRM"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	encoded := os.Args[1]

	jsonStr, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {

		fmt.Println(err.Error())

	}

	data := make(map[string]interface{})

	err = json.Unmarshal(jsonStr, &data)

	if err != nil {

		//

	}

	if data["type"] == "linux" {

		if data["category"] == "discovery" {

			var ans = SSH.Discovery(data)

			bytes, _ := json.Marshal(ans)

			fmt.Println(string(bytes))

		} else if data["category"] == "polling" {

			if data["metricGroup"] == "disk" {

				SSH.Disk(data)

			} else if data["metricGroup"] == "cpu" {

				SSH.Cpu(data)

			} else if data["metricGroup"] == "memory" {

				SSH.Memory(data)

			} else if data["metricGroup"] == "process" {

				SSH.Process(data)

			} else if data["metricGroup"] == "SystemInfo" {

				SSH.System(data)

			} else {

				res := make(map[string]interface{})

				res["error"] = "Invalid metricGroup"

				bytes, _ := json.Marshal(res)

				fmt.Println(string(bytes))

			}

		}

	} else if data["type"] == "windows" {

		if data["category"] == "discovery" {

			var ans = WINRM.Discovery(data)

			bytes, _ := json.Marshal(ans)

			fmt.Println(string(bytes))

		} else if data["category"] == "polling" {

			if data["metricGroup"] == "disk" {

				WINRM.Disk(data)

			} else if data["metricGroup"] == "cpu" {

				WINRM.Cpu(data)

			} else if data["metricGroup"] == "process" {

				WINRM.Process(data)

			} else if data["metricGroup"] == "memory" {

				WINRM.Memory(data)

			} else if data["metricGroup"] == "SystemInfo" {

				WINRM.System(data)

			} else {

				res := make(map[string]interface{})

				res["error"] = "Invalid metricGroup"

				bytes, _ := json.Marshal(res)

				fmt.Println(string(bytes))

			}

		}
	} else if data["type"] == "networking" {

		if data["category"] == "discovery" {

			var ans = SNMP.Discovery(data)

			bytes, _ := json.Marshal(ans)

			fmt.Println(string(bytes))

		} else if data["category"] == "polling" {

			if data["metricGroup"] == "SystemInfo" {

				SNMP.System(data)

			} else if data["metricGroup"] == "interface" {

				SNMP.Interface(data)

			} else {

				res := make(map[string]interface{})

				res["error"] = "Invalid metricGroup...."

				bytes, _ := json.Marshal(res)

				fmt.Println(string(bytes))

			}
		}
	} else {

		res := make(map[string]interface{})

		res["error"] = "Invalid type"

		bytes, _ := json.Marshal(res)

		fmt.Println(string(bytes))

	}

}
