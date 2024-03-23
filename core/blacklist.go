package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"slices"
)

type Blacklist struct {
	Blacklist []BlacklistGroup `yaml:"blacklist"`
}
type BlacklistGroup struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Addresses   []string `yaml:"addresses"`
}

func loadBlacklist(filename string) (Blacklist, error) {
	var blacklist Blacklist
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return blacklist, err
	}
	err = yaml.Unmarshal(yamlFile, &blacklist)
	return blacklist, err
}

func CheckIPisBlacklisted(blacklist Blacklist, packet PacketDataStruct) (bool, string) {
	for _, group := range blacklist.Blacklist {
		for range group.Addresses {
			if slices.Contains(group.Addresses, packet.Destination) {
				return true, fmt.Sprintf("Machine with IP %s attempted to connect to a blacklisted IP %s", packet.Source, packet.Destination)
			}
		}
	}
	return false, ""
}
