package ghp

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Code-Hex/ghp/internal/license"
	"github.com/Code-Hex/ghp/internal/ui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/tcnksm/go-gitconfig"
)

const github = "github.com"

// GHP struct
type GHP struct {
	Username string
	GhqRoot  string
	Project  string
	Dest     string
	Options
}

// LICENSE struct
type LICENSE struct {
	Year       string
	Project    string
	GithubUser string
}

// New returns GHP struct pointer
func New() *GHP {
	return new(GHP)
}

// Run method will create a project and returns exit code
func (g *GHP) Run() int {
	if e := g.run(); e != nil {
		exitCode, err := UnwrapErrors(e)
		if err != nil {
			if g.StackTrace {
				fmt.Fprintf(os.Stderr, "Error: %+v\n", e)
			} else {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			return exitCode
		}
	}
	return 0
}

func (g *GHP) prepare() error {
	args, err := parseOptions(&g.Options, os.Args[1:])
	if err != nil {
		return err
	}
	if len(args) != 1 {
		return errors.Errorf("Please pass me an argument for project name")
	}
	g.Project = args[0]
	return nil
}

func (g *GHP) run() error {
	if err := g.prepare(); err != nil {
		return errors.Wrap(err, "Failed to setup")
	}
	g.checkGhqRoot()
	if err := g.checkGitHubUser(); err != nil {
		return errors.Wrap(err, "Failed to check github user")
	}
	if err := g.getDestination(); err != nil {
		return errors.Wrap(err, "Failed to get destination path")
	}
	if err := os.Mkdir(g.Dest, 0755); err != nil {
		return errors.New(err.Error())
	}
	if err := runGitInit(g.Dest); err != nil {
		return errors.Wrap(err, "Failed to run git init")
	}
	if g.Options.WithLicense {
		choose := ui.Choose("Which license do you want to use? (number)", []string{
			"MIT License",
			"The Unlicense",
			"Apache License 2.0",
			"Mozilla Public License 2.0",
			"GNU General Public License v3.0",
			"GNU Affero General Public License v3.0",
			"GNU Lesser General Public License v3.0",
		}, "MIT License")
		t := template.New("init license")
		if err := g.generateLICENSE(t, choose); err != nil {
			return errors.Wrapf(err, "Failed to create %s", choose)
		}
	}
	fmt.Println(g.Dest)
	return nil
}

func (g *GHP) getDestination() error {
	path := filepath.Join(g.GhqRoot, github, g.Username, g.Project)
	if path[:2] == "~/" {
		homedir, err := homedir.Dir()
		if err != nil {
			return err
		}
		path = filepath.Join(homedir, path[2:])
	}
	g.Dest = path
	return nil
}

func (g *GHP) generateLICENSE(t *template.Template, kind string) error {
	if err := os.Chdir(g.Dest); err != nil {
		return err
	}
	f, err := os.Create("LICENSE")
	if err != nil {
		return err
	}
	defer f.Close()

	var choose string
	switch kind {
	case "MIT License":
		choose = license.MIT
	case "The Unlicense":
		choose = license.Unlicense
	case "Apache License 2.0":
		choose = license.Apache
	case "Mozilla Public License 2.0":
		choose = license.MPL2
	case "GNU General Public License v3.0":
		choose = license.GPLv3
	case "GNU Affero General Public License v3.0":
		choose = license.AGPLv3
	case "GNU Lesser General Public License v3.0":
		choose = license.LGPLv3
	}

	var licenseFile LICENSE
	licenseFile.Project = g.Project
	licenseFile.Year = fmt.Sprintf("%d", time.Now().Year())
	licenseFile.GithubUser, err = gitconfig.GithubUser()
	if err != nil {
		return err
	}

	ui.Printf("Writing %s\n", kind)

	tmpl, err := t.Parse(choose)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, licenseFile)
}

func (g *GHP) checkGhqRoot() {
	ghqRoot := ghqRoot()
	if ghqRoot == "" {
		g.GhqRoot = "~/.ghq"
	}
	g.GhqRoot = ghqRoot
}

func (g *GHP) checkGitHubUser() error {
	username, _ := gitconfig.GithubUser()
	if username != "" {
		g.Username = username
		return nil
	}
	err := g.readUsername()
	if err != nil {
		return err
	}
	return g.registerGitConfig()
}

// readUsername read username from prompt input.
func (g *GHP) readUsername() (err error) {
	fmt.Printf("Username for https://%s: ", github)
	g.Username, err = readline()
	return
}

// allow input only by one line
func readline() (line string, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}
	err = scanner.Err()
	return
}

func ghqRoot() string {
	p, _ := gitconfig.Entire("ghq.root")
	if p == "" {
		return os.Getenv("GHQ_ROOT")
	}
	return p
}

func (g *GHP) registerGitConfig() error {
	cmd := exec.Command("git", "config", "--global", "github.user", g.Username)
	return cmd.Run()
}

func runGitInit(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	return cmd.Run()
}

func parseOptions(opts *Options, argv []string) ([]string, error) {
	o, err := opts.parse(argv)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse arguments")
	}
	if opts.Help {
		return nil, makeUsageError(errors.New(opts.usage()))
	}
	return o, nil
}
