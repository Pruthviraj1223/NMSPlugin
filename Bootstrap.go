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

		panic(err)

	}

	data := make(map[string]interface{})

	err = json.Unmarshal(jsonStr, &data)

	if err != nil {

		fmt.Println("error ", err.Error())

	}

	if data["type"] == "linux" {

		if data["category"] == "discovery" {

			var ans = SSH.Discovery(data)

			bytes, _ := json.Marshal(ans)

			fmt.Println(string(bytes))

		} else if data["category"] == "polling" {

			if data["counter"] == "disk" {

				SSH.Disk(data)

			} else if data["counter"] == "CPU" {

				SSH.Cpu(data)

			} else if data["counter"] == "Memory" {

				SSH.Memory(data)

			} else if data["counter"] == "Process" {

				SSH.Process(data)

			} else if data["counter"] == "SystemInfo" {

				SSH.System(data)
			}

		}

	} else if data["type"] == "windows" {

		if data["category"] == "discovery" {

			var ans = WINRM.Discovery(data)

			bytes, _ := json.Marshal(ans)

			fmt.Println(string(bytes))

		} else if data["category"] == "polling" {

			if data["counter"] == "disk" {

				WINRM.Disk(data)

			} else if data["counter"] == "CPU" {

				WINRM.Cpu(data)

			} else if data["counter"] == "process" {

				WINRM.Process(data)

			} else if data["counter"] == "memory" {

				WINRM.Memory(data)

			}

		}
	} else if data["type"] == "networking" {

		if data["category"] == "discovery" {

			var ans = SNMP.Discovery(data)

			bytes, _ := json.Marshal(ans)

			fmt.Println(string(bytes))

		} else if data["category"] == "polling" {

			if data["counter"] == "systemInfo" {

				SNMP.System(data)

			} else if data["counter"] == "interface" {

				SNMP.Interface(data)

			}
		}
	}

}
