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

type GHP struct {
	Username string
	GhqRoot  string
}

func New() *GHP {
	return new(GHP)
}

func (g *GHP) Run() int {
	if err := g.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return 1
	}
	return 0
}

func (g *GHP) run() error {
	g.checkGhqRoot()
	if err := g.checkGitHubUser(); err != nil {
		return err
	}
	if len(os.Args) != 2 {
		return errors.Errorf("Please pass me an argument for project name")
	}
	newProject := os.Args[1]
	path := filepath.Join(g.GhqRoot, github, g.Username, newProject)
	if path[:2] == "~/" {
		homedir, err := homedir.Dir()
		if err != nil {
			return err
		}
		path = filepath.Join(homedir, path[2:])
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return runGitInit(path)
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
	cmd := exec.Command("git", "config", "github.user", fmt.Sprintf(`"%s"`, g.Username))
	return cmd.Run()
}

func runGitInit(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	return cmd.Run()
}
