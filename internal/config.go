package internal

import (
	_ "embed"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	//go:embed resources/default_config.yml
	defaultConfigResource string
)

type Config struct {
	Service struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"service"`
	Limits struct {
		IdleTimeout time.Duration `yaml:"idle"`
		MaxTime     time.Duration `yaml:"maxTime"`
	} `yaml:"limits"`
	MotdPath string   `yaml:"motdPath"`
	NodesDef []string `yaml:"nodes"`
}

func GetDefaultConfig() Config {
	var config Config

	err := yaml.Unmarshal([]byte(defaultConfigResource), &config)
	if err != nil {
		log.Fatal("Error while getting default configuration", err)
	}
	return config
}

func ReadConfigFile(path string) Config {
	var config Config

	configYaml, err := os.ReadFile(path)

	if err != nil {
		log.Fatal("Could not load yaml config file: ", path, err)
	}

	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		log.Fatal("Error while getting parsing yaml configuration", err)
	}
	return config
}

func WriteConfig(config Config, path string) {
	configYaml, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal("Could not marshal config data.", err)
	}

	err = os.WriteFile(path, configYaml, 0644)

	if err != nil {
		log.Fatal("Could not write config yaml to file:", path, err)
	}
}

func portRangeSplit(r rune) bool {
	return r == '-' || r == ','
}

func (cfg *Config) GetNodes() []string {

	var ports []string

	for _, nodeDef := range cfg.NodesDef {
		parts := strings.Split(nodeDef, ":")

		if len(parts) != 2 {
			log.Fatal("Node definition expected format: host:port_range")
		}

		host := parts[0]

		possiblePorts := parts[1]

		possiblePortRanges := strings.Split(possiblePorts, ",")

		for _, possiblePortRange := range possiblePortRanges {

			if strings.Contains(possiblePortRange, "-") {

				portRange := strings.Split(possiblePortRange, "-")

				if len(portRange) != 2 {
					log.Fatal("Invalid port range: ", nodeDef)
				}

				min, err := strconv.Atoi(portRange[0])

				if err != nil {
					log.Fatal("Invalid port number: ", portRange[0], err)
				}

				max, err := strconv.Atoi(portRange[1])

				if err != nil {
					log.Fatal("Invalid port number: ", portRange[1], err)
				}

				if min >= max {
					log.Fatal("Invalid port range: ", min, "-", max)
				}

				for p := min; p < max; p = p + 1 {
					ports = append(ports, host+":"+strconv.Itoa(p))
				}
			} else {
				port, err := strconv.Atoi(possiblePortRange)

				if err != nil {
					log.Fatal("Invalid port number: ", port, err)
				}

				ports = append(ports, host+":"+strconv.Itoa(port))
			}
		}

	}

	return ports
}
