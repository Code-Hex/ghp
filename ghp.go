package ghp

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
	Options
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
		return err
	}
	path := filepath.Join(g.GhqRoot, github, g.Username, g.Project)
	if path[:2] == "~/" {
		homedir, err := homedir.Dir()
		if err != nil {
			return err
		}
		path = filepath.Join(homedir, path[2:])
	}
	if err := os.Mkdir(path, 0755); err != nil {
		return errors.New(err.Error())
	}
	if err := runGitInit(path); err != nil {
		return err
	}
	fmt.Println(path)
	return nil
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
