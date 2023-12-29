package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/gabeduke/bumpver-action/bumpver"
)


func main() {
	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", true, "If true, print API calls but do not make them.")
	flag.Parse()

	config, err := bumpver.NewConfig()
	if err != nil {
		fmt.Printf("Error getting config: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	app := bumpver.New(config)

	latestTagVersion, err := app.GetLatestTag(ctx)
	if err != nil {
		fmt.Printf("Error getting latest tag: %v\n", err)
		os.Exit(1)
	}

	bumpType := app.GetBumpTypeFromCommitMessage(ctx)
	newTag, err := bumpver.IncrementVersion(latestTagVersion.SemVer, bumpType)
	if err != nil {
		fmt.Printf("Error incrementing version: %v\n", err)
		os.Exit(1)
	}

	app.CreateNewTag(ctx, newTag, dryRun)
}
