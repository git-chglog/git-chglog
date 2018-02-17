package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func main() {
	ttl := color.New(color.FgYellow).SprintFunc()

	cli.AppHelpTemplate = fmt.Sprintf(`
%s
   {{.Name}} [options] [tag revision]

%s
   {{.Version}}{{if .Author}}

%s
  {{.Author}}{{end}}

%s
   {{range .Flags}}{{.}}
   {{end}}
%s
   {{.Name}} todo1
   {{.Name}} todo2
   {{.Name}} todo3
   {{.Name}} todo4
`,
		ttl("USAGE:"),
		ttl("VERSION:"),
		ttl("AUTHOR:"),
		ttl("OPTIONS:"),
		ttl("EXAMPLE:"),
	)

	app := cli.NewApp()
	app.Name = "git-chglog"
	app.Usage = "todo usage for git-chglog"
	app.Version = Version

	app.Flags = []cli.Flag{
		// init
		cli.BoolFlag{
			Name:  "init",
			Usage: "generate the git-chglog configuration file in interactive",
		},

		// config
		cli.StringFlag{
			Name:  "config, c",
			Usage: "specifies a different configuration file to pick up",
			Value: ".chglog/config.yml",
		},

		// output
		cli.StringFlag{
			Name:  "output, o",
			Usage: "output path and filename for the changelogs (default: output to stdout)",
		},

		// silent
		cli.BoolFlag{
			Name:  "silent",
			Usage: "disable stdout output",
		},

		// no-color
		cli.BoolFlag{
			Name:   "no-color",
			Usage:  "disable color output",
			EnvVar: "NO_COLOR",
		},

		// no-emoji
		cli.BoolFlag{
			Name:   "no-emoji",
			Usage:  "disable emoji output",
			EnvVar: "NO_EMOJI",
		},

		// help & version
		cli.HelpFlag,
		cli.VersionFlag,
	}

	app.Action = func(c *cli.Context) error {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to get working directory", err)
			os.Exit(1)
		}

		// initializer
		if c.Bool("init") {
			os.Exit(NewInitializer().Run())
		}

		// chglog
		ctx := &Context{
			WorkingDir: wd,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			ConfigPath: c.String("config"),
			OutputPath: c.String("output"),
			Silent:     c.Bool("silent"),
			NoColor:    c.Bool("no-color"),
			NoEmoji:    c.Bool("no-emoji"),
			Query:      c.Args().First(),
		}

		chglogCLI := NewCLI(
			ctx,
			fs,
			NewConfigLoader(),
			NewGenerator(),
		)

		os.Exit(chglogCLI.Run())

		return nil
	}

	app.Run(os.Args)
}
