package config

import (
  "time"
  "github.com/alecthomas/kingpin"
)

type Config struct {
	KubeConfig	string
	Namespace   string
	Interval 	time.Duration
}

var defaultConfig = &Config{
	KubeConfig:	"",
	Namespace:  "",
	Interval: 	time.Minute,
}

func NewConfig() *Config {
	return &Config{}
}

// ParseFlags adds and parses flags from command line
func (cfg *Config) ParseFlags(args []string) error {
	app := kingpin.New("eip-controller", "")
	app.Flag("kubeconfig", "Retrieve target cluster configuration from a Kubernetes configuration file (default: auto-detect)").Default(defaultConfig.KubeConfig).StringVar(&cfg.KubeConfig)
	app.Flag("namespace", "Limit sources of endpoints to a specific namespace (default: all namespaces)").Default(defaultConfig.Namespace).StringVar(&cfg.Namespace)
	app.Flag("interval", "The interval between two consecutive synchronizations in duration format (default: 1m)").Default(defaultConfig.Interval.String()).DurationVar(&cfg.Interval)

	_, err := app.Parse(args)
	if err != nil {
		return err
	}

	return nil
}
