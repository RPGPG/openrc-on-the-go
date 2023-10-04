package main

import (
	"fmt"
	"log"
	"net/http"
	"openrc-on-the-go/config"
	"openrc-on-the-go/wrapper"
	"os"
	"strings"
	"time"
)

func main() {
	wrapper.CheckOS()
	if !wrapper.IsOpenrcExecutable() {
		log.Fatal("No OpenRC executable found.")
	}

	cmdArgs := os.Args[1:]
	configPath, err := config.GetConfigPath(cmdArgs)
	if err != nil {
		log.Fatal(err)
	}
	services, topic, interval, err := config.LoadConfigNtfy(configPath)
	if err != nil {
		log.Fatal(err)
	}

	statuses := map[string]bool{}
	for _, v := range services {
		statuses[v] = true
	}
	upDown := map[bool]string{true: "UP 🟢", false: "DOWN 🔴"}

	for {
		for k, v := range statuses {
			currentStatus := wrapper.IsServiceStarted(k)
			if currentStatus != v {
				statuses[k] = currentStatus
				message := fmt.Sprintf(
					"service %s changed status to %s",
					k,
					upDown[currentStatus],
				)
				log.Println(message)
				http.Post(fmt.Sprintf("https://ntfy.sh/%s", topic),
					"text/plain",
					strings.NewReader(
						message,
					),
				)
			} else {
				log.Printf("no change for %s", k)
			}
		}
		time.Sleep(time.Second * time.Duration(interval))
	}
}
