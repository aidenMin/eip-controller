package config

import (
	"github.com/alecthomas/kingpin"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Interval 	time.Duration
	KubeConfig	string
	LogLevel 	string
	NodeLabel	string
}

var defaultConfig = &Config{
	Interval: 	time.Minute,
	KubeConfig:	"",
	LogLevel:  	logrus.InfoLevel.String(),
	NodeLabel:	"upbit.com/eip-group",
}

func NewConfig() *Config {
	return &Config{}
}

// ParseFlags adds and parses flags from command line
func (cfg *Config) ParseFlags(args []string) error {
	app := kingpin.New("eip-controller", "")
	app.Flag("kubeconfig", "Retrieve target cluster configuration from a Kubernetes configuration file (default: auto-detect)").Default(defaultConfig.KubeConfig).StringVar(&cfg.KubeConfig)
	app.Flag("interval", "The interval between two consecutive synchronizations in duration format (default: 1m)").Default(defaultConfig.Interval.String()).DurationVar(&cfg.Interval)
	app.Flag("log-level", "Set the level of logging. (default: info, options: panic, debug, info, warning, error, fatal").Default(defaultConfig.LogLevel).EnumVar(&cfg.LogLevel, allLogLevelsAsStrings()...)
	app.Flag("node-label", "Set the label of node. (default: upbit.com/eip-group").Default(defaultConfig.NodeLabel).StringVar(&cfg.NodeLabel)

	_, err := app.Parse(args)
	if err != nil {
		return err
	}

	return nil
}

// allLogLevelsAsStrings returns all logrus levels as a list of strings
func allLogLevelsAsStrings() []string {
	var levels []string
	for _, level := range logrus.AllLevels {
		levels = append(levels, level.String())
	}
	return levels
}