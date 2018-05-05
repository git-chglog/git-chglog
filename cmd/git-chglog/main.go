package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
	"github.com/urfave/cli"
)

func main() {
	ttl := color.New(color.FgYellow).SprintFunc()

	cli.AppHelpTemplate = fmt.Sprintf(`
%s
  {{.Name}} [options] <tag query>

    There are the following specification methods for <tag query>.

    1. <old>..<new> - Commit contained in <old> tags from <new>.
    2. <name>..     - Commit from the <name> to the latest tag.
    3. ..<name>     - Commit from the oldest tag to <name>.
    4. <name>       - Commit contained in <name>.

%s
  {{range .Flags}}{{.}}
  {{end}}
%s

  $ {{.Name}}

    If <tag query> is not specified, it corresponds to all tags.
    This is the simplest example.

  $ {{.Name}} 1.0.0..2.0.0

    The above is a command to generate CHANGELOG including commit of 1.0.0 to 2.0.0.

  $ {{.Name}} 1.0.0

    The above is a command to generate CHANGELOG including commit of only 1.0.0.

  $ {{.Name}} $(git describe --tags $(git rev-list --tags --max-count=1))

    The above is a command to generate CHANGELOG with the commit included in the latest tag.

  $ {{.Name}} --output CHANGELOG.md

    The above is a command to output to CHANGELOG.md instead of standard output.

  $ {{.Name}} --config custom/dir/config.yml

    The above is a command that uses a configuration file placed other than ".chglog/config.yml".
`,
		ttl("USAGE:"),
		ttl("OPTIONS:"),
		ttl("EXAMPLE:"),
	)

	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		cli.HelpPrinterCustom(colorable.NewColorableStdout(), templ, data, nil)
	}

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
			Usage: "output path and filename for the changelogs. If not specified, output to stdout",
		},

		cli.StringFlag{
			Name:  "next-tag",
			Usage: "treat unreleased commits as specified tags (EXPERIMENTAL)",
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
			os.Exit(ExitCodeError)
		}

		// initializer
		if c.Bool("init") {
			initializer := NewInitializer(
				&InitContext{
					WorkingDir: wd,
					Stdout:     colorable.NewColorableStdout(),
					Stderr:     colorable.NewColorableStderr(),
				},
				fs,
				NewQuestioner(
					gitcmd.New(&gitcmd.Config{
						Bin: "git",
					}),
					fs,
				),
				NewConfigBuilder(),
				templateBuilderFactory,
			)

			os.Exit(initializer.Run())
		}

		// chglog
		chglogCLI := NewCLI(
			&CLIContext{
				WorkingDir: wd,
				Stdout:     colorable.NewColorableStdout(),
				Stderr:     colorable.NewColorableStderr(),
				ConfigPath: c.String("config"),
				OutputPath: c.String("output"),
				Silent:     c.Bool("silent"),
				NoColor:    c.Bool("no-color"),
				NoEmoji:    c.Bool("no-emoji"),
				Query:      c.Args().First(),
				NextTag:    c.String("next-tag"),
			},
			fs,
			NewConfigLoader(),
			NewGenerator(),
		)

		os.Exit(chglogCLI.Run())

		return nil
	}

	app.Run(os.Args)
}
