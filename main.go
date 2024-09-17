package main

import (
	"fmt"
	"io"
)

const (
	ElvUI = "ElvUI"
	TukUI = "tukui-org"
)

func main() {
	done := make(chan bool)
	message := make(chan string)

	go StartSpinner(done, message)

	// get latest tagged version
	message <- "Getting latest tagged ElvUI version"
	latest, err := GetLatestTaggedVersion()
	if err != nil {
		fmt.Printf("could not get latest tagged ElvUI version: %v", err)
		return
	}

	// create output file
	message <- "Creating temporary output file"
	output, err := CreateOutputFile(latest)
	if err != nil {
		fmt.Printf("could not create output file: %v", err)
		return
	}
	defer output.Close()

	// download zip from github
	message <- "downloading addon from github"
	resp, err := DownloadZip(latest)
	if err != nil {
		fmt.Printf("could now download zip file: %v", err)
		return
	}
	defer resp.Body.Close()

	// copy zip contents to temp file
	message <- "copying addon contents to temp file"
	_, err = io.Copy(output, resp.Body)
	if err != nil {
		fmt.Printf("could not copy zipball contents to %s: %v", output.Name(), err)
		return
	}

	// unzip elvui directory
	message <- "unzipping ElvUI files"
	topLevelDir, err := Unzip(output.Name(), "./")
	if err != nil {
		fmt.Printf("could not unzip elvui %s", err)
		return
	}

	// copy to addons folder
	message <- "moving files to addon folder"
	if err = CopyDirectory(topLevelDir, "./dest"); err != nil {
		fmt.Printf("could not move unzipped files to destination: %v", err)
		return
	}

	// cleanup artifacts
	message <- "cleaning up artifacts"
	artifacts := RunArtifacts{Tempfile: output.Name(), TopLevelUnzippedDirectory: topLevelDir}
	if err = CleanupRunArtifacts(artifacts); err != nil {
		fmt.Printf("could not clean up run artifacts: %v", err)
		return
	}

	done <- true

	fmt.Printf("successfully downloaded ElvUI version %s\n", *latest.Name)
}
