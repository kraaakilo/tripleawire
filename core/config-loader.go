package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Config struct {
	protocols []string
}

func LoadConfigMap(path string) map[string]string {
	content, err := os.ReadFile(path)
	config := map[string]string{}

	if err != nil {
		fmt.Println(fmt.Sprintf("The path that you specified is not found. More: %v", err.Error()))
		return nil
	}

	for _, line := range strings.Split(string(content), "\n") {
		settings := strings.Split(line, "=")
		name := settings[0]
		value := settings[1]
		config[name] = value
	}
	return config
}

func LoadConfig(configMap map[string]string) Config {
	keys := []string{"PROTOCOLS", "PORTS"}

	config := Config{}

	for key, value := range configMap {
		if slices.Contains(keys, key) {
			switch key {
			case "PROTOCOLS":
				for _, p := range strings.Split(value, ",") {
					if slices.Contains(Protocols, p) {
						config.protocols = append(config.protocols, p)
					}
				}
			}
		}
	}
	return config
}
