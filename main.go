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

// config struct which holds all information required
type config struct {
	originalURL    string
	originalBranch string
	mirrorURL      string
	mirrorBranch   string
	pat            string
	useForce       bool
	useVerbose     bool
	useTags        bool
}

type byteSlice []byte

func main() {

	// get originalURL input (required)
	originalURL := githubactions.GetInput(OriginalURLInputField)
	if originalURL == "" {
		githubactions.Fatalf(ErrNoOriginalURL)
		return
	}

	// get originalBranch input (optional)
	originalBranch := githubactions.GetInput(OriginalBranchInputField)
	if originalBranch == "" {
		log.Printf(InfoNoOriginalBranch)
		originalBranch = OriginalDefaultBranch
	}

	// get mirrorURL input (required)
	mirrorURL := githubactions.GetInput(MirrorURLInputField)
	if mirrorURL == "" {
		githubactions.Fatalf(ErrNoMirrorURL)
		return
	}

	// get mirrorBranch input (optional)
	mirrorBranch := githubactions.GetInput(MirrorBranchInputField)
	if mirrorBranch == "" {
		log.Printf(InfoNoMirrorBranch)
		mirrorBranch = MirrorDefaultBranch
	}

	// get Personal Access Token encoded in base64 input (required)
	patEncoded := githubactions.GetInput(PATInputField)
	if patEncoded == "" {
		githubactions.Fatalf(ErrNoPAT)
		return
	}

	// add encoded PAT to mask to make sure it doesn't get logged to console
	githubactions.AddMask(patEncoded)

	// base64 decode PAT
	patBytes, err := base64.StdEncoding.DecodeString(patEncoded)
	if err != nil {
		githubactions.Fatalf(ErrFailedToBase64DecodePAT)
		return
	}

	pat := strings.TrimSuffix(string(patBytes), "\n")

	// add true PAT to mask
	githubactions.AddMask(pat)

	// get force input to see if push can use the argument `--force`
	var useForce = false
	useForceInput := githubactions.GetInput(UseForceInputField)
	if useForceInput == UseForceTrue {
		log.Printf(InfoUsingForce)
		useForce = true
	}

	// get verbose input to check whether or not to use verbose mode
	var useVerbose = false
	useVerboseInput := githubactions.GetInput(UseVerboseInputField)
	if useVerboseInput == UseVerboseTrue {
		log.Println(InfoUsingVerbose)
		useVerbose = true
	}

	// get tag input to check whether or not to use tags
	var useTags = false
	useTagsInput := githubactions.GetInput(UseTagsInputField)
	if useTagsInput == UseTagsTrue {
		log.Println(InfoUsingTags)
		useTags = true
	}

	// make config
	config := config{
		originalURL:    originalURL,
		originalBranch: originalBranch,
		mirrorURL:      mirrorURL,
		mirrorBranch:   mirrorBranch,
		pat:            pat,
		useForce:       useForce,
		useVerbose:     useVerbose,
		useTags:        useTags,
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
		log.Printf("Output: %v\n", out)
		githubactions.Fatalf(err.Error())
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
	log.Printf("initializing git")

	defaultCommand := []string{"init"}

	return c.command(defaultCommand...)

}

// adds new remote to local git repository
// also fetches branches from the repository (--fetch flag at the end)
func (c *config) addRemote(name string, repo string) (output string, err error) {
	log.Printf("adding remote: %v\n", name)

	defaultCommand := []string{"remote", "add", name, repo, "--fetch"}

	return c.command(defaultCommand...)
}

// checks out specific branch from specific remote in local git repository
func (c *config) checkout(remote string, branch string) (output string, err error) {
	log.Printf("checking out: %v/%v\n", remote, branch)

	defaultCommand := []string{"checkout", fmt.Sprintf("%v/%v", remote, branch)}

	return c.command(defaultCommand...)
}

// pulls specific branch from specific remote
func (c *config) pull(remote string, branch string) (output string, err error) {
	log.Printf("pulling: %v/%v", remote, branch)

	defaultCommand := []string{"pull", remote, branch}

	return c.command(defaultCommand...)
}

// pushes to specific branch on remote
func (c *config) push(remote string, branch string) (output string, err error) {
	log.Printf("pushing: %v/%v\n", remote, branch)

	defaultCommand := []string{"push", "--set-upstream", remote, branch}

	if c.useForce {
		defaultCommand = append(defaultCommand, "--force")
	}

	if c.useTags {
		defaultCommand = append(defaultCommand, "--tags")
	}

	return c.command(defaultCommand...)
}

// creates a new branch with specific name
func (c *config) branch(name string) (output string, err error) {
	log.Printf("creating branch: %v", name)

	defaultCommand := []string{"branch", name}

	return c.command(defaultCommand...)
}

// executes git commands in the TempDir folder
func (c *config) command(args ...string) (out string, err error) {

	cwd, _ := os.Getwd()
	pathArgs := []string{"-C", fmt.Sprintf("%v/%v", cwd, TempDir)}
	args = append(pathArgs, args...)

	cmd := exec.Command("git", args...)

	// verbose mode
	if c.useVerbose {
		log.Println(cmd.Args)
	}

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

	if url[len(url)-4:] != ".git" {
		url = fmt.Sprintf("%v.git", url)
		fmt.Println(errors.New(ErrURLEndNotDotGit))
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
