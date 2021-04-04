package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

const (

	// name of remote urls
	UpstreamRemote = "upstream"
	MirrorRemote   = "mirror"

	// default branch names
	OriginalDefaultBranch = "master"
	MirrorDefaultBranch   = "mirror"

	// name used as temporary directory
	TempDir = "tmp"

	// error messages
	ErrURLNotHTTPS             = "url is not https"
	ErrNoOriginalURL           = "no original repository url provided"
	ErrNoMirrorURL             = "no mirror repository url provided"
	ErrNoPAT                   = "no personal access token provided"
	ErrFailedToBase64DecodePAT = "failed to decode PAT from b64"

	// info messages
	InfoNoOriginalBranch = "no original branch provided, using 'master'"
	InfoNoMirrorBranch   = "no mirror branch provided, using 'mirror'"
	InfoUsingForce       = "git will now use --force to push"

	// input constants
	OriginalURLInput    = "originalURL"
	OriginalBranchInput = "originalBranch"
	MirrorURLInput      = "mirrorURL"
	MirrorBranchInput   = "mirrorBranch"
	PATInput            = "pat"
	UseForceInput       = "useForce"

	// input responses
	UseForceYes = "yes"
)

// config struct which holds all information required
type config struct {
	originalURL    string
	originalBranch string
	mirrorURL      string
	mirrorBranch   string
	pat            string
	useForce       bool
}

type byteSlice []byte

func main() {

	// get originalURL input (required)
	originalURL := githubactions.GetInput(OriginalURLInput)
	if originalURL == "" {
		githubactions.Fatalf(ErrNoOriginalURL)
		return
	}

	// get originalBranch input (optional)
	originalBranch := githubactions.GetInput(OriginalBranchInput)
	if originalBranch == "" {
		log.Printf(InfoNoOriginalBranch)
		originalBranch = OriginalDefaultBranch
	}

	// get mirrorURL input (required)
	mirrorURL := githubactions.GetInput(MirrorURLInput)
	if mirrorURL == "" {
		githubactions.Fatalf(ErrNoMirrorURL)
		return
	}

	// get mirrorBranch input (optional)
	mirrorBranch := githubactions.GetInput(MirrorBranchInput)
	if mirrorBranch == "" {
		log.Printf(InfoNoMirrorBranch)
		mirrorBranch = MirrorDefaultBranch
	}

	// get Personal Access Token encoded in base64 input (required)
	patEncoded := githubactions.GetInput(PATInput)
	if patEncoded == "" {
		githubactions.Fatalf(ErrNoPAT)
		return
	}

	// add encoded PAT to mask to make sure it doesn't get logged to console
	githubactions.AddMask(patEncoded)

	// base64 decode PAT
	pat, err := base64.StdEncoding.DecodeString(patEncoded)
	if err != nil {
		githubactions.Fatalf(ErrFailedToBase64DecodePAT)
		return
	}

	// add true PAT to mask
	githubactions.AddMask(string(pat))

	// get useForce input to see if push can use the argument `--force`
	var useForce = false
	useForceInput := githubactions.GetInput(UseForceInput)
	if useForceInput == UseForceYes {
		log.Printf(InfoUsingForce)
		useForce = true
	}

	// make config
	config := config{
		originalURL:    originalURL,
		originalBranch: originalBranch,
		mirrorURL:      mirrorURL,
		mirrorBranch:   mirrorBranch,
		pat:            string(pat),
		useForce:       useForce,
	}

	// convert URLs to use PAT
	err = config.usePAT()
	if err != nil {
		githubactions.Fatalf(err.Error())
		return
	}

	// init git repository
	out, err := config.gitInit()
	if err != nil {
		log.Printf("Output: %v\n", out)
		githubactions.Fatalf(err.Error())

		return
	}

	// add upstream remote (url of repo we want to clone)
	out, err = config.addRemote(UpstreamRemote, config.originalURL)
	if err != nil {
		log.Printf("Output: %v\n", out)
		githubactions.Fatalf(err.Error())

		return
	}

	// add mirror remote (url of repo we want to mirror to)
	out, err = config.addRemote(MirrorRemote, config.mirrorURL)
	if err != nil {
		log.Printf("Output: %v\n", out)
		githubactions.Fatalf(err.Error())

		return
	}

	// checks out branch on upstream remote
	out, err = config.checkout(UpstreamRemote, config.originalBranch)
	if err != nil {
		log.Printf("Output: %v\n", out)
		githubactions.Fatalf(err.Error())

		return
	}

	// pulls branch on upstream remote
	out, err = config.pull(UpstreamRemote, config.originalBranch)
	if err != nil {
		githubactions.Fatalf(err.Error())
		log.Printf("Output: %v\n", out)
		return
	}

	// makes new branch
	out, err = config.branch(config.mirrorBranch)
	if err != nil {
		log.Printf("Output: %v\n", out)
		githubactions.Fatalf(err.Error())

		return
	}

	// pushes new branch to mirror remote
	out, err = config.push(MirrorRemote, config.mirrorBranch)
	if err != nil {
		log.Printf("Output: %v\n", out)
		githubactions.Fatalf(err.Error())

		return
	}
}

// initialize git repository
func (c *config) gitInit() (output string, err error) {
	githubactions.Debugf("initializing git")
	return command("init")
}

// adds new remote to local git repository
// also fetches branches from the repository (--fetch flag at the end)
func (c *config) addRemote(name string, repo string) (output string, err error) {
	githubactions.Debugf("adding remote: %v\n", name)
	return command("remote", "add", name, repo, "--fetch")
}

// checks out specific branch from specific remote in local git repository
func (c *config) checkout(remote string, branch string) (output string, err error) {
	githubactions.Debugf("checking out: %v/%v\n", remote, branch)
	return command("checkout", fmt.Sprintf("%v/%v", remote, branch))
}

// pulls specific branch from specific remote
func (c *config) pull(remote string, branch string) (output string, err error) {
	githubactions.Debugf("pulling: %v/%v", remote, branch)
	return command("pull", remote, branch)
}

// pushes to specific branch on remote
func (c *config) push(remote string, branch string) (output string, err error) {
	githubactions.Debugf("pushing: %v/%v\n", remote, branch)
	if c.useForce {
		return command("push", "--set-upstream", remote, branch, "--force")
	}

	return command("push", "--set-upstream", remote, branch)
}

// creates a new branch with specific name
func (c *config) branch(name string) (output string, err error) {
	githubactions.Debugf("creating branch: %v", name)
	return command("branch", name)
}

// executes git commands in the TempDir folder
func command(args ...string) (out string, err error) {

	cwd, _ := os.Getwd()
	pathArgs := []string{"-C", fmt.Sprintf("%v/%v", cwd, TempDir)}
	args = append(pathArgs, args...)

	cmd := exec.Command("git", args...)

	log.Println(cmd.Args)

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
		return "", err
	}

	out = output.ToString()
	return out, err
}

// converts url in Config type to use Personal Access Token
func (c *config) usePAT() (err error) {
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
		return "", errors.New(ErrURLNotHTTPS)
	}

	return fmt.Sprintf("https://%v@%v", pat, url[8:]), err
}

// convert a byteSlice to a string
func (b byteSlice) ToString() string {
	var a []string
	for _, by := range b {
		a = append(a, string(by))
	}

	return strings.Join(a, " ")
}
