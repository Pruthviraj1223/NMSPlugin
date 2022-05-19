package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Cpu(data map[string]interface{}) {

	sshUser := (data["name"]).(string)

	sshPassword := (data["password"]).(string)

	sshHost := (data["ip.address"]).(string)

	sshPort := int((data["port"]).(float64))

	config := &ssh.ClientConfig{

		Timeout: 10 * time.Second,

		User: sshUser,

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Config: ssh.Config{Ciphers: []string{

			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}

	address := fmt.Sprintf("%s:%d", sshHost, sshPort)

	sshClient, err := ssh.Dial("tcp", address, config)

	if err != nil {

	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()

	if err != nil {
		return
	}

	var CPU = make(map[string]interface{})

	cpuUtilization, err := session.CombinedOutput("mpstat -P ALL")

	var cpuList []map[string]string

	cpuUtilizationString := string(cpuUtilization)

	cpuStringArray := strings.Split(cpuUtilizationString, "\n")

	flag1 := 1
	for _, v := range cpuStringArray {

		if flag1 <= 3 {
			flag1++
			continue
		}
		eachWord := strings.Split(standardizeSpaces(v), " ")
		if len(eachWord) <= 13 {
			continue
		}

		fmt.Println("values = ", v, " each = > ", eachWord[3])

		if eachWord[3] == "all" {
			CPU["cpu.name"] = eachWord[3]
			CPU["cpu.user.percent"] = eachWord[4]
			CPU["cpu.sys.percent"] = eachWord[6]
			CPU["cpu.idle.percent"] = eachWord[13]
		} else {
			temp1 := map[string]string{
				"cpu.name":         eachWord[3],
				"cpu.user.percent": eachWord[4],
				"cpu.sys.percent":  eachWord[6],
				"cpu.idle.percent": eachWord[13],
			}
			cpuList = append(cpuList, temp1)
		}
	}

	CPU["CPU"] = cpuList

	bytes, err := json.MarshalIndent(CPU, " ", " ")

	if err != nil {
		return
	}

	fmt.Println(string(bytes))

}
