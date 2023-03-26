package main

import (
	"fmt"
	"github.com/benabernathy/roundabout/internal"
	cli "github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "options for configuration",
				Subcommands: []*cli.Command{
					{
						Name:   "init",
						Usage:  "initializes a default configuration",
						Action: initConfig,
					},
					{
						Name:   "validate",
						Usage:  "validates a configuration file",
						Action: validateConfig,
					},
					{
						Name:   "show",
						Usage:  "shows an effective configuration",
						Action: showEffectiveConfig,
					},
				},
			},
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "runs the listen server",
				Action:  serve,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initConfig(cCtx *cli.Context) error {

	if cCtx.Args().Len() != 1 {
		log.Fatal("Expected config destination file path. e.g. ~/config.yaml")
	}

	configDest := cCtx.Args().First()
	defaultConfig := internal.GetDefaultConfig()
	internal.WriteConfig(defaultConfig, configDest)

	fmt.Println("Generated default configuration and saved at:", configDest)

	return nil
}

func validateConfig(cCtx *cli.Context) error {

	if cCtx.Args().Len() != 1 {
		log.Fatal("Expected config file path. e.g. ~/config.yaml")
	}

	configPath := cCtx.Args().First()

	config := internal.ReadConfigFile(configPath)

	_ = config.GetNodes()

	log.Println("Validated configuration.")
	return nil
}

func showEffectiveConfig(cCtx *cli.Context) error {
	log.Println("Effective configuration is...")

	if cCtx.Args().Len() != 1 {
		log.Fatal("Expected config file path. e.g. ~/config.yaml")
	}

	configPath := cCtx.Args().First()

	config := internal.ReadConfigFile(configPath)

	configYaml, err := yaml.Marshal(&config)

	if err != nil {
		log.Fatal("Could not marshal config data.", err)
	}

	log.Println("\n\n", string(configYaml))

	return nil
}

func serve(cCtx *cli.Context) error {

	if cCtx.Args().Len() != 1 {
		log.Fatal("Expected config file path. e.g. ~/config.yaml")
	}

	configPath := cCtx.Args().First()
	config := internal.ReadConfigFile(configPath)

	server := internal.Server{}
	server.Serve(config)

	return nil
}
