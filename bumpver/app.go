package bumpver

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v57/github"
)

type App struct {
	client *github.Client
	cfg    *Config
}

type TagVersion struct {
	Tag    *github.RepositoryTag
	SemVer *semver.Version
}

func IncrementVersion(tag *semver.Version, bumpType string) (string, error) {
	var newTag semver.Version

	switch bumpType {
	case "major":
		newTag = tag.IncMajor()
	case "minor":
		newTag = tag.IncMinor()
	case "patch":
		newTag = tag.IncPatch()
	default:
		return "", fmt.Errorf("invalid bump type: %s", bumpType)
	}

	// TODO make the leading "v" configurable
	return "v" + newTag.String(), nil
}

func New(cfg *Config) *App {
	app := &App{
		cfg: cfg,
	}

	app.initializingGitHubClient(cfg)

	return app
}

func (a *App) initializingGitHubClient(cfg *Config) {
	client := github.NewClient(nil).WithAuthToken(cfg.token)

	a.client = client

	return
}

func (a *App) GetLatestTag(ctx context.Context) (*TagVersion, error) {
	owner := a.cfg.Owner
	repo := a.cfg.Repo

	opts := &github.ListOptions{PerPage: 100}
	var allTags []*TagVersion

	for {
		tags, resp, err := a.client.Repositories.ListTags(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		for _, tag := range tags {
			v, err := semver.NewVersion(*tag.Name)
			if err != nil {
				return nil, fmt.Errorf("invalid semver tag: %s", *tag.Name)
			}
			allTags = append(allTags, &TagVersion{Tag: tag, SemVer: v})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// If no tags are found, return a default TagVersion of v0.0.0
    if len(allTags) == 0 {
        defaultVersion, _ := semver.NewVersion("v0.0.0")
        return &TagVersion{SemVer: defaultVersion}, nil
    }

	sort.Slice(allTags, func(i, j int) bool {
		return allTags[i].SemVer.GreaterThan(allTags[j].SemVer)
	})

	return allTags[0], nil
}

func (a *App) CreateNewTag(ctx context.Context, newTag string, dryRun bool) {
	owner := a.cfg.Owner
	repo := a.cfg.Repo
	sha := a.cfg.Sha

	ref := &github.Reference{
		Ref: github.String("refs/tags/" + newTag),
		Object: &github.GitObject{
			SHA: &sha,
		},
	}

	if dryRun {
		fmt.Printf("Dry run: would have created ref with owner: %s, repo: %s, ref: %v\n", owner, repo, ref)
	} else {
		_, _, err := a.client.Git.CreateRef(ctx, owner, repo, ref)
		if err != nil {
			panic(err)
		}
	}
}

func (a *App) GetBumpTypeFromCommitMessage(ctx context.Context) string {
	owner := a.cfg.Owner
	repo := a.cfg.Repo

	commits, _, err := a.client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		panic(err)
	}

	commitMessage := *commits[0].Commit.Message
	commitMessage = strings.ToLower(commitMessage)

	if strings.Contains(commitMessage, "#major") {
		return "major"
	} else if strings.Contains(commitMessage, "#minor") {
		return "minor"
	} else {
		return "patch"
	}
}
