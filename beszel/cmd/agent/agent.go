package main

import (
	"beszel"
	"beszel/internal/agent"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Define flags for key and port
	keyFlag := flag.String("key", "", "Public key")
	portFlag := flag.String("port", "45876", "Port number")

	// Parse the flags
	flag.Parse()

	// handle flags / subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println(beszel.AppName+"-agent", beszel.Version)
			os.Exit(0)
		case "help":
			fmt.Printf("Usage: %s [options] [subcommand]\n", os.Args[0])
			fmt.Println("\nOptions:")
			flag.PrintDefaults()
			fmt.Println("\nSubcommands:")
			fmt.Println("  version      Display the version")
			fmt.Println("  help         Display this help message")
			fmt.Println("  update       Update the agent to the latest version")
			os.Exit(0)
		case "update":
			agent.Update()
			os.Exit(0)
		}
	}

	var pubKey []byte
	// Override the key if the -key flag is provided
	if *keyFlag != "" {
		pubKey = []byte(*keyFlag)
	} else {
		// Try to get the key from the KEY environment variable.
		pubKey = []byte(os.Getenv("KEY"))
	}

	// If KEY is not set, try to read the key from the file specified by KEY_FILE.
	if len(pubKey) == 0 {
		keyFile, exists := os.LookupEnv("KEY_FILE")
		if !exists {
			log.Fatal("Must set KEY or KEY_FILE environment variable")
		}
		var err error
		pubKey, err = os.ReadFile(keyFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	addr := ":" + *portFlag
	if portEnvVar, exists := os.LookupEnv("PORT"); exists {
		// allow passing an address in the form of "127.0.0.1:45876"
		if !strings.Contains(portEnvVar, ":") {
			portEnvVar = ":" + portEnvVar
		}
		addr = portEnvVar
	}

	// Override the port if the -port flag is provided
	if *portFlag != "45876" {
		addr = ":" + *portFlag
	}

	agent.NewAgent().Run(pubKey, addr)
}
