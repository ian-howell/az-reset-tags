package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	helpFlag         bool
	groupFlag        string
	subscriptionFlag string
)

const subscription = "Edge_AzureLocal_EphemeralEngineeringDevEnvironments"

func initFlags() {
	flag.BoolVar(&helpFlag, "help", false, "Show help")
	flag.StringVar(&groupFlag, "group", "", "Azure resource group name")
	flag.StringVar(&subscriptionFlag, "subscription", "", "Azure subscription name")
	flag.Parse()
}

func main() {
	initFlags()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if err := run(); err != nil {
		log.Println(redF("%v", err))
		os.Exit(1)
	}
}

func run() error {
	groupName := getGroupName()
	if groupName == "" {
		return fmt.Errorf("group name not set")
	}

	now := time.Now()
	ts := now.Format("2006-01-02T15:04:05Z")
	log.Printf("Resetting tags for group %s", groupName)
	if err := resetTags(groupName, ts); err != nil {
		return fmt.Errorf("failed to reset tags: %w", err)
	}

	log.Println(green("😊  Finished resetting tags!"))
	oneWeekFromNow := now.AddDate(0, 0, 7).Format("2006-01-02")
	log.Println(greenF("📅  Don't forget to run again prior to %s", oneWeekFromNow))
	log.Println(green("💡  Here's a handy line of cron to just run this command daily at 2PM (local time)"))

	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	fmt.Printf("    0 14 * * * %s -group %v\n\n", executable, groupName)
	return nil
}

func resetTags(groupName, timestamp string) error {
	cmd := exec.Command(
		"az", "group", "update",
		"--name", groupName,
		"--subscription", subscription,
		"--tags",
		"CleanupFrequency=Weekly",
		"Created="+timestamp,
		"--only-show-errors",
		"--output", "none",
	)

	log.Printf("Running command: %v\n", cmd.Args)

	done := make(chan struct{})
	go tick(done)
	defer func() {
		close(done)
		fmt.Println()
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		trimmed := strings.Trim(string(output), " \n")
		indented := indentString(trimmed, 4)
		return fmt.Errorf("failed to reset tags: %v\n\n%s", err, indented)
	}

	return nil
}

func getGroupName() string {
	// The flag gets the highest priority
	if groupFlag != "" {
		return groupFlag
	}

	// Next highest priority is the environment variable
	const groupEnvVar = "DB_AZ_RESOURCE_GROUP"
	if groupName := os.Getenv(groupEnvVar); groupName != "" {
		return groupName
	}

	// Otherwise, we read the config file
	config := Config{}
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "nc-devbox", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	return config.Azure.ResourceGroup
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

func greenF(s string, a ...interface{}) string {
	return green(fmt.Sprintf(s, a...))
}

func red(s string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", s)
}

func redF(s string, a ...interface{}) string {
	return red(fmt.Sprintf(s, a...))
}

func indentString(s string, n int) string {
	return strings.Repeat(" ", n) + strings.ReplaceAll(s, "\n", "\n"+strings.Repeat(" ", n))
}
