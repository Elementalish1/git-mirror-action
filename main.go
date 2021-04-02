package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	ORIGINAL = ""
	MIRROR   = ""
	PAT      = ""

	UpstreamRemote = "upstream"
	MirrorRemote   = "mirror"
	MasterBranch   = "master"
	NewBranch = "mirror"

	TempDir = "tmp"

	ErrUrlNotHttps = "url is not https"
)

type Config struct {
	originalURL string
	mirrorURL   string
	pat         string
}

func main() {

	config := Config{
		originalURL: ORIGINAL,
		mirrorURL:   MIRROR,
		pat:         PAT,
	}

	err := config.usePAT()
	if err != nil {
		log.Println(err)
		return
	}

	err = gitInit()
	if err != nil {
		log.Println(err)
		return
	}

	err = addRemote(UpstreamRemote, config.originalURL)
	if err != nil {
		log.Println(err)
		return
	}

	err = fetch(UpstreamRemote)
	if err != nil {
		log.Println(err)
		return
	}

	err = addRemote(MirrorRemote, config.mirrorURL)
	if err != nil {
		log.Println(err)
		return
	}

	err = fetch(MirrorRemote)
	if err != nil {
		log.Println(err)
		return
	}

	err = checkout(UpstreamRemote, MasterBranch)
	if err != nil {
		log.Println(err)
		return
	}

	err = pull(UpstreamRemote, MasterBranch)
	if err != nil {
		log.Println(err)
		return
	}

	err = branch(NewBranch)
	if err != nil {
		log.Println(err)
		return
	}
	err = push(MirrorRemote, NewBranch)
	if err != nil {
		log.Println(err)
	}
}

func gitInit() error {
	log.Printf("initializing git")
	return command("init")
}

func addRemote(name string, repo string) error {
	log.Printf("adding remote: %v\n", name)
	return command("remote", "add", name, repo)
}

func fetch(remote string) error {
	log.Printf("fetching remote: %v\n", remote)
	return command("fetch", remote)
}

func checkout(remote string, branch string) error {
	log.Printf("checking out: %v/%v\n", remote, branch)
	return command("checkout", fmt.Sprintf("%v/%v", remote, branch))
}

func pull(remote string, branch string) error {
	log.Printf("pulling: %v/%v", remote, branch)
	return command("pull", remote, branch)
}

func push(remote string, branch string) error {
	log.Printf("pushing: %v/%v\n", remote, branch)
	return command("push","--set-upstream", remote, branch)
}

func branch(name string) error {
	log.Printf("creating branch: %v", name)
	return command("branch", name)
}

func command(args ...string) error {
	cwd, _ := os.Getwd()
	pathArgs := []string{"-C", fmt.Sprintf("%v/%v", cwd, TempDir)}
	args = append(pathArgs, args...)

	cmd := exec.Command("git", args...)

	if _, err := os.Stat(fmt.Sprintf("%v/%v", cwd, TempDir)); os.IsNotExist(err) {
		err := os.Mkdir(TempDir, 0777)
		if err != nil {
			return err
		}
	}

	err := cmd.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	err = cmd.Wait()
	return err
}

func (c *Config) usePAT() (err error) {
	c.originalURL, err = convertURL(c.originalURL, c.pat)
	if err != nil {
		return err
	}

	c.mirrorURL, err = convertURL(c.mirrorURL, c.pat)
	return err
}

func convertURL(url string, pat string) (s string, err error) {
	if url[:8] != "https://" {
		return "", errors.New(ErrUrlNotHttps)
	}

	return fmt.Sprintf("https://git-mirror-action:%v@%v", pat, url[8:]), err
}