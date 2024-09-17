package main

import (
	"fmt"
	"os"
)

type RunArtifacts struct {
	Tempfile                  string
	TopLevelUnzippedDirectory string
}

// CleanupRunArtifacts deletes all the unused files post-run
func CleanupRunArtifacts(artifacts RunArtifacts) error {
	if err := os.Remove(artifacts.Tempfile); err != nil {
		return fmt.Errorf("could not remove tempfile post run: %v", err)
	}

	if err := os.RemoveAll(artifacts.TopLevelUnzippedDirectory); err != nil {
		return fmt.Errorf("could not remove top level unzipped directory post run: %v", err)
	}

	return nil
}
