package main

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
)

// Initializer ...
type Initializer struct {
	ctx             *InitContext
	client          gitcmd.Client
	fs              FileSystem
	logger          *Logger
	questioner      Questioner
	configBuilder   ConfigBuilder
	templateBuilder TemplateBuilder
}

// NewInitializer ...
func NewInitializer(
	ctx *InitContext,
	fs FileSystem,
	questioner Questioner,
	configBuilder ConfigBuilder,
	templateBuilder TemplateBuilder,
) *Initializer {
	return &Initializer{
		ctx:             ctx,
		fs:              fs,
		logger:          NewLogger(colorable.NewColorableStdout(), colorable.NewColorableStderr(), false, false),
		questioner:      questioner,
		configBuilder:   configBuilder,
		templateBuilder: templateBuilder,
	}
}

// Run ...
func (init *Initializer) Run() int {
	ans, err := init.questioner.Ask()
	if err != nil {
		init.logger.Error(err.Error())
		return ExitCodeError
	}

	if err = init.fs.MkdirP(filepath.Join(init.ctx.WorkingDir, ans.ConfigDir)); err != nil {
		init.logger.Error(err.Error())
		return ExitCodeError
	}

	if err = init.generateConfig(ans); err != nil {
		init.logger.Error(err.Error())
		return ExitCodeError
	}

	if err = init.generateTemplate(ans); err != nil {
		init.logger.Error(err.Error())
		return ExitCodeError
	}

	success := color.CyanString("✔")
	init.logger.Log(fmt.Sprintf(`
:sparkles: %s
  %s %s
  %s %s
`,
		color.GreenString("Configuration file and template generation completed!"),
		success,
		filepath.Join(ans.ConfigDir, defaultConfigFilename),
		success,
		filepath.Join(ans.ConfigDir, defaultTemplateFilename),
	))

	return ExitCodeOK
}

func (init *Initializer) generateConfig(ans *Answer) error {
	s, err := init.configBuilder.Build(ans)
	if err != nil {
		return err
	}

	return init.fs.WriteFile(filepath.Join(init.ctx.WorkingDir, ans.ConfigDir, defaultConfigFilename), []byte(s))
}

func (init *Initializer) generateTemplate(ans *Answer) error {
	s, err := init.templateBuilder.Build(ans)
	if err != nil {
		return err
	}

	return init.fs.WriteFile(filepath.Join(init.ctx.WorkingDir, ans.ConfigDir, defaultTemplateFilename), []byte(s))
}
