package main

import (
	"encoding/json"
	"fmt"
	"log"
	"openrc-on-the-go/config"
	"openrc-on-the-go/structures"
	"openrc-on-the-go/wrapper"
	"os"
	"strings"
	"time"
)

func handleErr(err error, message string) {
	if err != nil {
		if message != "" {
			log.Println(message)
		}
		log.Fatal(err)
	}
}

func printHelp() {
	log.Println("Usage:")
	log.Println("--[[type]] --ini=[[path to file]]")
	log.Println("Or:")
	log.Println("--cfg=[[path to file]]")
	log.Println("")
	log.Println("Types:")
	log.Println("--simple - simple output, human-readable")
	log.Println("--mon - continious output, human-readable")
	log.Println("--json - json output")
	log.Fatal("--cfg=[[path]] - generate config file and exit")
}

func printMotd() {
	log.Println("/\\/\\/\\/\\/\\/\\/\\/\\/\\/\\/\\/\\/\\")
	log.Println("/\\/\\/\\/\\🐧OpenRC🐧/\\/\\/\\/\\")
	log.Println("/\\/\\/\\/\\/\\/\\/\\/\\/\\/\\/\\/\\/\\")
	log.Println("")
}

func loopForeverMonitoring(services map[string]string, monName string, monInterval int) {
	aliasMaxLen, serviceMaxLen := len("Alias"), len("Name")
	for k, v := range services {
		if len(k) > aliasMaxLen {
			aliasMaxLen = len(k)
		}
		if len(v) > serviceMaxLen {
			serviceMaxLen = len(v)
		}
	}
	stoppedCounter := 0
	for {
		stoppedCounter = 0
		wrapper.ClearTerminal()
		printMotd()
		if monName == "//hostname//" {
			displayMonName, err := os.Hostname()
			handleErr(err, "Cannot determine hostname.")
			monName = displayMonName
		}
		if monName != "nouse" {
			log.Println("")
			log.Printf("\\\\\\\\\\\\/%s/\\\\\\\\\\\\", strings.Repeat("/", len(monName)))
			log.Printf("\\\\\\\\\\\\ %s \\\\\\\\\\\\", monName)
			log.Printf("\\\\\\\\\\\\/%s/\\\\\\\\\\\\", strings.Repeat("/", len(monName)))
			log.Println("")
		}
		log.Printf("%s%s :: %s%s :: %s",
			"Alias", strings.Repeat(" ", aliasMaxLen-len("Alias")), "Name",
			strings.Repeat(" ", serviceMaxLen-len("Name")), "Status")
		log.Printf("%s%s%s%s", strings.Repeat("-", aliasMaxLen),
			strings.Repeat("-", serviceMaxLen), strings.Repeat("-", 4),
			strings.Repeat("-", len("Status")))
		for k, v := range services {
			if wrapper.IsServiceStarted(v) {
				log.Printf("%s%s :: %s%s :: %s",
					k, strings.Repeat(" ", aliasMaxLen-len(k)), v,
					strings.Repeat(" ", serviceMaxLen-len(v)), "🟢")
			} else {
				stoppedCounter++
				log.Printf("%s%s :: %s%s :: %s",
					k, strings.Repeat(" ", aliasMaxLen-len(k)), v,
					strings.Repeat(" ", serviceMaxLen-len(v)), "🔴")
			}
		}
		if stoppedCounter > 0 {
			log.Println("")
			log.Println("")
			log.Println("")
			log.Println("🔴 NOT ALL SERVICES ARE OPERATING! 🔴")
		}
		time.Sleep(time.Second * time.Duration(monInterval))
	}
}

func generateCfgIfWanted(cmdArgs []string) {
	if len(cmdArgs[0]) > 5 {
		if cmdArgs[0][:6] == "--cfg=" {
			configFileName := cmdArgs[0][6:]
			err := config.GenerateExampleConfig(configFileName)
			handleErr(err, fmt.Sprintf("Error while generating %s", configFileName))
			log.Printf("%s created", configFileName)
			log.Println("Exiting...")
			os.Exit(0)
		}
	}
}

func checkMode(cmdArgs []string) (bool, bool, bool, structures.JsonOutput) {
	simpleMode, jsonMode, monMode := false, false, false
	if cmdArgs[0] == "--simple" {
		jsonMode = false
		simpleMode = true
		monMode = false
	}
	if cmdArgs[0] == "--json" {
		simpleMode = false
		jsonMode = true
		monMode = false
	}
	if cmdArgs[0] == "--mon" {
		simpleMode = false
		jsonMode = false
		monMode = true
	}
	outputJson := structures.JsonOutput{Started: 0, Stopped: 0, Services: []structures.Service{}}
	return simpleMode, jsonMode, monMode, outputJson
}

func basicCheckForArgs(cmdArgs []string) {
	if len(cmdArgs) == 0 {
		printHelp()
	}
	if cmdArgs[0] != "--simple" && cmdArgs[0] != "--json" && cmdArgs[0] != "--mon" && cmdArgs[0][:6] != "--cfg=" {
		printHelp()
	}
}

func checkForIniArg(cmdArgs []string) {
	if cmdArgs[1][:6] != "--ini=" {
		log.Println("Provide path to .ini config as second argument --ini=")
		printHelp()
	}
}

func listServicesToCheck(simpleMode bool, services map[string]string) {
	if simpleMode {
		log.Println("Services to check:")
		for _, v := range services {
			log.Println(v)
		}
		log.Println("")
		log.Println("Checking services...")
	}
}

func checkForAliveServices(services map[string]string, simpleMode bool,
	jsonMode bool, outputJson structures.JsonOutput) structures.JsonOutput {
	for _, v := range services {
		if simpleMode {
			output := wrapper.CheckServiceStatus(v)
			log.Printf("%s: %s", v, wrapper.OutTrimmed(output))
		}
		if jsonMode {
			if wrapper.IsServiceStarted(v) {
				outputJson.Started++
				outputJson.Services = append(outputJson.Services,
					structures.Service{Name: v, Started: true})
			} else {
				outputJson.Stopped++
				outputJson.Services = append(outputJson.Services,
					structures.Service{Name: v, Started: false})
			}
		}
	}
	return outputJson
}

func main() {
	wrapper.CheckOS()

	cmdArgs := os.Args[1:]
	basicCheckForArgs(cmdArgs)
	generateCfgIfWanted(cmdArgs)

	simpleMode, jsonMode, monMode, outputJson := checkMode(cmdArgs)
	checkForIniArg(cmdArgs)
	if simpleMode {
		printMotd()
	}

	openrcBinary := wrapper.IsOpenrcExecutable()
	if !openrcBinary {
		log.Fatal("No OpenRC binary found")
	} else {
		if simpleMode {
			log.Println("OpenRC executable found in PATH.")
			log.Println("")
		}
	}

	services, monInterval, monName, err := config.LoadConfigCli(cmdArgs[1][6:])
	handleErr(err, "Cannot load config file.")
	if monMode {
		loopForeverMonitoring(services, monName, monInterval)
	}

	listServicesToCheck(simpleMode, services)
	outputJson = checkForAliveServices(services, simpleMode, jsonMode, outputJson)

	if simpleMode {
		log.Println("")
		log.Println("Check complete.")
		log.Println("Exiting...")
	}
	if jsonMode {
		output, err := json.Marshal(outputJson)
		handleErr(err, "Error building JSON output.")
		fmt.Print(string(output[:]))
	}
}
