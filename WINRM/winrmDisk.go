package WINRM

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strings"
)

func standardizeSpaces(s string) string {

	return strings.Join(strings.Fields(s), " ")

}

func Disk(data map[string]interface{}) {

	defer func() {

		if r := recover(); r != nil {

			res := make(map[string]interface{})

			res["error"] = r

			errorDisplay(res)

		}

	}()

	var errorList []string

	//host := data["ip"].(string)

	port := int((data["port"]).(float64))
	//
	//name := (data["name"]).(string)
	//
	//password := (data["password"]).(string)

	endpoint := winrm.NewEndpoint(data["ip"].(string), port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, data["name"].(string), data["password"].(string))

	if err != nil {

		errorList = append(errorList, err.Error())

	}

	if len(errorList) == 0 {

		var disk string

		commandForDisk := "Get-WmiObject win32_logicaldisk | Foreach-Object {$_.DeviceId,$_.Freespace,$_.Size -join \" \"}"

		disk, aa, bb, err := client.RunPSWithString(commandForDisk, "")

		fmt.Println(aa)

		fmt.Println(bb)

		var diskList []map[string]string

		diskStringArray := strings.Split(disk, "\n")

		for _, v := range diskStringArray {

			eachWord := strings.Split(standardizeSpaces(v), " ")

			if len(eachWord) == 0 {

				break

			}

			if len(eachWord) == 3 {

				temp := map[string]string{

					"disk.name.bytes": eachWord[0],

					"disk.free.bytes": eachWord[1],

					"disk.size.bytes": eachWord[2],
				}

				diskList = append(diskList, temp)

			}
			if len(eachWord) == 1 {

				temp := map[string]string{

					"disk.name.bytes": eachWord[0],

					"disk.free.bytes": "0",

					"disk.size.bytes": "0",
				}

				diskList = append(diskList, temp)

			}

		}

		var diskMap = make(map[string]interface{})

		diskMap["disk"] = diskList

		bytes, err := json.Marshal(diskMap)

		if err != nil {

			response := make(map[string]interface{})

			response["error"] = err.Error()

			errorDisplay(response)

		} else {

			fmt.Println(string(bytes))

		}

	} else {

		res := make(map[string]interface{})

		res["error"] = errorList

		errorDisplay(res)

	}

}
