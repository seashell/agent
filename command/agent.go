package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"github.com/dimiro1/banner"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/joho/godotenv"
	agent "github.com/seashell/agent/agent"
	cli "github.com/seashell/agent/pkg/cli"
	log "github.com/seashell/agent/pkg/log"
	"github.com/seashell/agent/pkg/log/simple"
)

// AgentCommand :
type AgentCommand struct {
	UI cli.UI
}

// Name :
func (c *AgentCommand) Name() string {
	return "agent"
}

// Synopsis :
func (c *AgentCommand) Synopsis() string {
	return "Runs the seashell agent"
}

// Run :
func (c *AgentCommand) Run(ctx context.Context, args []string) int {

	displayBanner()

	config := c.parseConfig(args)

	// logger, err := logrus.NewLoggerAdapter(logrus.Config{
	// 	LoggerOptions: log.LoggerOptions{
	// 		Level:  config.LogLevel,
	// 		Prefix: "agent: ",
	// 	},
	// })

	// logger, err := zap.NewLoggerAdapter(zap.Config{
	// 	LoggerOptions: log.LoggerOptions{
	// 		Level:  config.LogLevel,
	// 		Prefix: "agent: ",
	// 	},
	// })

	logger, err := simple.NewLoggerAdapter(simple.Config{
		LoggerOptions: log.LoggerOptions{
			Level:  config.LogLevel,
			Prefix: "agent: ",
		},
	})

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output("==> Starting Seashell agent...")

	// Create DataDir and other subdirectories if they do not exist
	if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
		os.Mkdir(config.DataDir, 0755)
	}

	c.printConfig(config)

	agent, err := agent.New(config, logger)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error starting agent: %s\n", err.Error()))
		return 1
	}

	<-ctx.Done()

	agent.Shutdown()

	return 0
}

func (c *AgentCommand) parseConfig(args []string) *agent.Config {

	flags := FlagSet(c.Name())

	configFromFlags := c.parseFlags(flags, args)
	configFromFile := c.parseConfigFiles(flags.configPaths...)
	configFromEnv := c.parseEnv(flags.envPaths...)

	config := agent.DefaultConfig()

	config = config.Merge(configFromFile)
	config = config.Merge(configFromEnv)
	config = config.Merge(configFromFlags)

	if err := config.Validate(); err != nil {
		c.UI.Error(fmt.Sprintf("Invalid input: %s", err.Error()))
		os.Exit(1)
	}

	return config
}

func (c *AgentCommand) parseFlags(flags *RootFlagSet, args []string) *agent.Config {

	flags.Usage = func() {
		c.UI.Output("\n" + c.Help() + "\n")
	}

	config := agent.EmptyConfig()

	var devMode bool

	// Agent mode
	flags.BoolVar(&devMode, "dev", false, "")

	// General options (available in both client and server modes)
	flags.StringVar(&config.DataDir, "data-dir", "", "")
	flags.StringVar(&config.LogLevel, "log-level", "", "")

	// Client-only options
	flags.StringVar(&config.Client.OutputDir, "output-dir", "", "")
	flags.StringVar(&config.Client.DeviceID, "device-id", "", "")
	flags.StringVar(&config.Client.SecretID, "secret-id", "", "")

	if err := flags.Parse(args); err != nil {
		c.UI.Error("==> Error: " + err.Error() + "\n")
		os.Exit(1)
	}

	return config
}

func (c *AgentCommand) parseConfigFiles(paths ...string) *agent.Config {

	config := agent.EmptyConfig()

	if len(paths) > 0 {
		c.UI.Info(fmt.Sprintf("==> Loading configurations from: %v", paths))
		for _, s := range paths {
			err := hclsimple.DecodeFile(s, nil, config)
			if err != nil {
				c.UI.Error("Failed to load configuration: " + err.Error())
				os.Exit(0)
			}
		}
	} else {
		c.UI.Output("==> No configuration files loaded")
	}

	return config
}

func (c *AgentCommand) parseEnv(paths ...string) *agent.Config {

	config := agent.EmptyConfig()

	if len(paths) > 0 {

		c.UI.Info(fmt.Sprintf("==> Loading environment variables from: %v", paths))
		c.UI.Warn(fmt.Sprintf("  - This will not override already existing variables!"))

		err := godotenv.Load(paths...)

		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing env files: %s", err.Error()))
			os.Exit(1)
		}
	}

	env.Parse(config)

	return config
}

func (c *AgentCommand) printConfig(config *agent.Config) {

	info := map[string]string{
		"data dir":  config.DataDir,
		"device id": config.Client.DeviceID,
		"log level": config.LogLevel,
		"version":   config.Version.VersionNumber(),
	}

	padding := 18
	c.UI.Output("==> Seashell agent configuration:\n")
	for k := range info {
		c.UI.Info(fmt.Sprintf(
			"%s%s: %v",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}

	c.UI.Output("")
}

// Help :
func (c *AgentCommand) Help() string {
	h := `
Usage: seashell agent [options]
	
  Starts the Seashell agent and runs until an interrupt is received.
  The agent runs in client mode, and interacts with the Seashell Cloud.
  
  The Seashell agent's configuration primarily comes from the config
  files used, but a subset of the options may also be passed directly
  as CLI arguments.

General Options:
` + GlobalOptions() + `

Agent Options:

  --data-dir=<path>
    The data directory where all state will be persisted. On Seashell 
    clients this is used to store local network configurations, whereas
    on server nodes, the data dir is also used to keep the desired state
	for all the managed networks. Overrides the DRAGO_DATA_DIRenvironment
	variable if set.

  --log-level=<level>
    The logging level Seashell should log at. Valid values are INFO, WARN, DEBUG, ERROR, FATAL.
    Overrides the DRAGO_LOG_LEVEL environment variable if set.
	
`
	return strings.TrimSpace(h)
}

// Prints an ASCII banner to the standard output
func displayBanner() {
	banner.Init(os.Stdout, true, true, strings.NewReader(agent.Banner))
}
