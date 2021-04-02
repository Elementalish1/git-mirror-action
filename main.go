package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	ORIGINAL = ""
	MIRROR   = ""
	PAT      = ""

	UpstreamRemote = "upstream"
	MirrorRemote   = "mirror"
	OriginalBranch = "master"
	MirrorBranch   = "mirror"

	TempDir = "tmp"

	ErrUrlNotHttps = "url is not https"
)

type Config struct {
	originalURL    string
	originalBranch string
	mirrorURL      string
	mirrorBranch   string
	pat            string
}

type byteSlice []byte

func main() {

	config := Config{
		originalURL:    ORIGINAL,
		originalBranch: OriginalBranch,
		mirrorURL:      MIRROR,
		mirrorBranch:   MirrorBranch,
		pat:            PAT,
	}


	err := config.usePAT()
	if err != nil {
		log.Println(err)
		return
	}

	// init
	out, err := gitInit()
	if err != nil {
		log.Println(err)
		log.Printf("Output: %v\n", out)
		return
	}

	// add upstream remote (url of repo we want to clone)
	out, err = addRemote(UpstreamRemote, config.originalURL)
	if err != nil {
		log.Println(err)
		log.Printf("Output: %v\n", out)
		return
	}

	// add mirror remote (url of repo we want to mirror to)
	out, err = addRemote(MirrorRemote, config.mirrorURL)
	if err != nil {
		log.Println(err)
		log.Printf("Output: %v\n", out)
		return
	}

	// checks out branch on upstream remote
	out, err = checkout(UpstreamRemote, config.originalBranch)
	if err != nil {
		log.Println(err)
		log.Printf("Output: %v\n", out)
		return
	}

	// pulls branch on upstream remote
	out, err = pull(UpstreamRemote, config.originalBranch)
	if err != nil {
		log.Println(err)
		log.Printf("Output: %v\n", out)
		return
	}

	// makes new branch
	out, err = branch(config.mirrorBranch)
	if err != nil {
		log.Println(err)
		log.Printf("Output: %v\n", out)
		return
	}

	// pushes new branch to mirror remote
	out, err = push(MirrorRemote, config.mirrorBranch)
	if err != nil {
		log.Println(err)
		log.Printf("Output: %v\n", out)
	}
}

// initialize git repository
func gitInit() (output string, err error) {
	log.Printf("initializing git")
	return command("init")
}

// adds new remote to local git repository
// also fetches branches from the repository (--fetch flag at the end)
func addRemote(name string, repo string) (output string, err error) {
	log.Printf("adding remote: %v\n", name)
	return command("remote", "add", name, repo, "--fetch")
}

// checks out specific branch from specific remote in local git repository
func checkout(remote string, branch string) (output string, err error) {
	log.Printf("checking out: %v/%v\n", remote, branch)
	return command("checkout", fmt.Sprintf("%v/%v", remote, branch))
}

// pulls specific branch from specific remote
func pull(remote string, branch string) (output string, err error) {
	log.Printf("pulling: %v/%v", remote, branch)
	return command("pull", remote, branch)
}

// pushes to specific branch on remote
func push(remote string, branch string) (output string, err error) {
	log.Printf("pushing: %v/%v\n", remote, branch)
	return command("push", "--set-upstream", remote, branch, "--force")
}

// creates a new branch with specific name
func branch(name string) (output string, err error) {
	log.Printf("creating branch: %v", name)
	return command("branch", name)
}

// executes git commands in the TempDir folder
func command(args ...string) (out string, err error) {
	cwd, _ := os.Getwd()
	pathArgs := []string{"-C", fmt.Sprintf("%v/%v", cwd, TempDir)}
	args = append(pathArgs, args...)

	cmd := exec.Command("git", args...)

	// makes sure the TempDir exist
	if _, err := os.Stat(fmt.Sprintf("%v/%v", cwd, TempDir)); os.IsNotExist(err) {
		err := os.Mkdir(TempDir, 0777)
		if err != nil {
			return "", err
		}
	}

	// starts command
	var output byteSlice
	output, err = cmd.Output()
	if err != nil {
		log.Println(err)
		return "", err
	}

	out = output.byteSliceToString()
	return out, err
}

// converts url in Config type to use Personal Access Token
func (c *Config) usePAT() (err error) {
	c.originalURL, err = convertURL(c.originalURL, c.pat)
	if err != nil {
		return err
	}

	c.mirrorURL, err = convertURL(c.mirrorURL, c.pat)
	return err
}

// adds pat to https url
func convertURL(url string, pat string) (s string, err error) {
	if url[:8] != "https://" {
		return "", errors.New(ErrUrlNotHttps)
	}

	return fmt.Sprintf("https://%v@%v", pat, url[8:]), err
}

// convert a byteSlice to a string
func (b byteSlice) byteSliceToString() string {
	var a []string
	for _, byte := range b {
		a = append(a, string(byte))
	}

	return strings.Join(a, " ")
}
