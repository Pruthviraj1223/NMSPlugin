package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"time"
)

func Memory(data map[string]interface{}) {

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

	available, err := session.Output("free -b | grep Mem | awk '{ printf(\"%.4f\\n\", $7) }'") // available

	session.Close()

	session, err = sshClient.NewSession()

	freeMemory, err := session.Output("free -b | grep Mem | awk '{ printf(\"%.4f\\n\", $4) }'") //free memory
	session.Close()

	session, err = sshClient.NewSession()

	usedMemory, err := session.Output("free -b | grep Mem | awk '{ printf(\"%.4f \\n\", $3) }'") // used memory
	session.Close()

	session, err = sshClient.NewSession()

	totalMemory, err := session.Output("free -b | grep Mem | awk '{ printf(\"%i\\n\", $2) }'") // used memory
	session.Close()

	session, err = sshClient.NewSession()

	swapUsed, err := session.Output("free -b | grep Swap | awk '{ printf(\"%i\\n\", $3) }'")
	session.Close()

	session, err = sshClient.NewSession()

	swapFree, err := session.Output("free -b | grep Swap | awk '{ printf(\"%i\\n\", $3) }'")
	session.Close()

	session, err = sshClient.NewSession()

	swapTotal, err := session.Output("free -b | grep Swap | awk '{ printf(\"%i\\n\", $3) }'")
	session.Close()

	session, err = sshClient.NewSession()

	fmt.Println("ready")

	fmt.Println(string(available), " ", string(freeMemory), " ", string(usedMemory), " ", string(totalMemory), " ", string(swapTotal), " ", string(swapUsed), " ", string(swapFree))

	usedMemoryInt, _ := strconv.ParseFloat(string(usedMemory), 64)

	totalInt, _ := strconv.Atoi(string(totalMemory))

	fmt.Println(usedMemoryInt, totalInt)

	result := map[string]interface{}{
		"Device":          "linux",
		"freeMemory":      string(freeMemory),
		"usedMemory":      string(usedMemory),
		"totalMemory":     string(totalMemory),
		"availableMemory": string(available),
		"swapTotal":       string(swapTotal),
		"swapUsed":        string(swapUsed),
		"swapFree":        string(swapFree),
	}

	bytes, _ := json.MarshalIndent(result, " ", " ")

	fmt.Println(string(bytes))

}
