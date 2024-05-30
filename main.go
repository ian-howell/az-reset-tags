package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const groupNameEnvVar = "AZ_GROUP"

func main() {
	groupName, ok := os.LookupEnv(groupNameEnvVar)
	if !ok {
		log.Fatalf("Environment variable %s not set\n", groupNameEnvVar)
	}

	log.Printf("Resetting Janitor for group %s\n", groupName)

	ts := time.Now().Format("2006-01-02T15:04:05Z")
	_ = ts

	//"az group update --name %s --tags 'CleanupFrequency=Weekly' 'Created=%s' 'Owner=ianhowell@microsoft.com' --only-show-errors --output none", groupName, ts)
	cmd := exec.Command(
		"az", "group", "update",
		"--name", groupName,
		"--tags",
		"CleanupFrequency=Weekly",
		"Created="+ts,
		"--only-show-errors",
		"--output", "none",
	)

	log.Printf("Running command: %v\n", cmd.Args)

	done := make(chan struct{})
	go tick(done)

	output, err := cmd.CombinedOutput()
	if err != nil {
		close(done)
		fmt.Println()
		trimmed := strings.Trim(string(output), " \n")
		log.Fatalf("Received the following output (%v):\n%s", err, indented(red(trimmed), 4))
	}

	close(done)
	fmt.Println()
	log.Println(green("😊  Finished resetting Janitor!"))
}

func tick(done chan struct{}) {
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fmt.Print(".")
		}
	}
}

func green(s string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", s)
}

func red(s string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", s)
}

func indented(s string, n int) string {
	return strings.Repeat(" ", n) + strings.ReplaceAll(s, "\n", "\n"+strings.Repeat(" ", n))
}
