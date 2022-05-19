package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func Process(data map[string]interface{}) {

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

	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)

	sshClient, err := ssh.Dial("tcp", addr, config)

	if err != nil {
		panic(err.Error())
	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()

	if err != nil {
		panic(err)
	}

	psaux, err := session.Output("ps aux")

	if err != nil {
		panic(err)
	}

	psauxString := string(psaux)

	myStringArray := strings.Split(psauxString, "\n")

	var processList []map[string]string
	flag := 1
	for _, v := range myStringArray {
		if flag == 1 {
			flag = 0
			continue
		}
		eachWorld := strings.SplitN(standardizeSpaces(v), " ", 11)
		if len(eachWorld) <= 10 {
			break
		}

		temp1 := map[string]string{
			"process.user":   eachWorld[0],
			"process.id":     eachWorld[1],
			"process.cpu":    eachWorld[2],
			"process.memory": eachWorld[3],
			"Process.commad": eachWorld[10],
		}

		processList = append(processList, temp1)

	}

	processMap := make(map[string]interface{})

	processMap["Process"] = processList

	bytes, _ := json.MarshalIndent(processMap, " ", " ")

	fmt.Println(string(bytes))

}
