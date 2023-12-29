package bumpver

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Owner string
	Repo  string
	Sha   string
	token string
}

func NewConfig() (*Config, error) {
	c := &Config{}

	if err := c.getRepository(); err != nil {
		return nil, err
	}

	if err := c.getSha(); err != nil {
		return nil, err
	}

	if err := c.getGitHubToken(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) getRepository() error {
	repoFullName, found := os.LookupEnv("GITHUB_REPOSITORY")
	if !found {
		return fmt.Errorf("GITHUB_REPOSITORY not found")
	}

	split := strings.Split(repoFullName, "/")
	if len(split) != 2 {
		return fmt.Errorf("invalid GITHUB_REPOSITORY format")
	}

	c.Owner = split[0]
	c.Repo = split[1]

	return nil
}

func (c *Config) getSha() error {
	sha, found := os.LookupEnv("GITHUB_SHA")
	if !found {
		return fmt.Errorf("GITHUB_SHA not found")
	}

	c.Sha = sha

	return nil
}

func (c *Config) getGitHubToken() error {
	token, found := os.LookupEnv("GITHUB_TOKEN")
	if found {
		c.token = token
		return nil
	}

	cmd := exec.Command("gh", "auth", "token")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get GITHUB_TOKEN: %w", err)
	}

	c.token = strings.TrimSpace(string(out))

	return nil
}
