package main

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	chglog "github.com/git-chglog/git-chglog"
)

// CLI ...
type CLI struct {
	ctx              *CLIContext
	fs               FileSystem
	logger           *Logger
	configLoader     ConfigLoader
	generator        Generator
	processorFactory *ProcessorFactory
}

// NewCLI ...
func NewCLI(
	ctx *CLIContext, fs FileSystem,
	configLoader ConfigLoader,
	generator Generator,
) *CLI {
	silent := false
	if ctx.Silent || ctx.OutputPath == "" {
		silent = true
	}

	return &CLI{
		ctx:              ctx,
		fs:               fs,
		logger:           NewLogger(ctx.Stdout, ctx.Stderr, silent, ctx.NoEmoji),
		configLoader:     configLoader,
		generator:        generator,
		processorFactory: NewProcessorFactory(),
	}
}

// Run ...
func (c *CLI) Run() int {
	start := time.Now()

	if c.ctx.NoColor {
		color.NoColor = true
	}

	c.logger.Log(":watch: Generating changelog ...")

	config, err := c.prepareConfig()
	if err != nil {
		c.logger.Error(err.Error())
		return ExitCodeError
	}

	changelogConfig, err := c.createChangelogConfig(config)
	if err != nil {
		c.logger.Error(err.Error())
		return ExitCodeError
	}

	w, err := c.createOutputWriter()
	if err != nil {
		c.logger.Error(err.Error())
		return ExitCodeError
	}

	err = c.generator.Generate(w, c.ctx.Query, changelogConfig)
	if err != nil {
		c.logger.Error(err.Error())
		return ExitCodeError
	}

	c.logger.Log(fmt.Sprintf(":sparkles: Generate of %s is completed! (%s)",
		color.GreenString("\""+c.ctx.OutputPath+"\""),
		color.New(color.Bold).SprintFunc()(time.Since(start).String()),
	))

	return ExitCodeOK
}

func (c *CLI) prepareConfig() (*Config, error) {
	config, err := c.configLoader.Load(c.ctx.ConfigPath)
	if err != nil {
		return nil, err
	}

	err = config.Normalize(c.ctx)
	if err != nil {
		return nil, err
	}

	config.Convert(c.ctx)

	return config, err
}

func (c *CLI) createChangelogConfig(config *Config) (*chglog.Config, error) {
	processor, err := c.processorFactory.Create(config)
	if err != nil {
		return nil, err
	}

	changelogConfig := config.Convert(c.ctx)
	changelogConfig.Options.Processor = processor

	return changelogConfig, nil
}

func (c *CLI) createOutputWriter() (io.Writer, error) {
	if c.ctx.OutputPath == "" {
		return c.ctx.Stdout, nil
	}

	out := c.ctx.OutputPath
	dir := filepath.Dir(out)
	err := c.fs.MkdirP(dir)
	if err != nil {
		return nil, err
	}

	file, err := c.fs.Create(out)
	if err != nil {
		return nil, err
	}

	return file, nil
}
