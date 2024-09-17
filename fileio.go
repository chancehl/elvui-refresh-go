package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v64/github"
)

// createOutputFile creates a temporary output file for storing the downloaded zip file.
// The file name is based on the provided GitHub repository tag.
func CreateOutputFile(tag *github.RepositoryTag) (*os.File, error) {
	output, err := os.Create(fmt.Sprintf("./%s.zip", *tag.Name))
	if err != nil {
		return nil, fmt.Errorf("could not create outfile at %s: %v", output.Name(), err)
	}
	return output, nil
}

// downloadZip fetches the zip file corresponding to the latest tagged version from the GitHub repository.
func DownloadZip(tag *github.RepositoryTag) (*http.Response, error) {
	resp, err := http.Get(*tag.ZipballURL)
	if err != nil {
		return nil, fmt.Errorf("could not download zipball from %s: %v", *tag.ZipballURL, err)
	}
	return resp, nil
}

// unzip extracts the contents of the zip file from the source path to the destination directory.
// It returns the top-level directory extracted.
func Unzip(src string, dest string) (string, error) {
	// Open the zip file
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var topLevelDirectory string

	// Iterate through the files in the zip archive
	for index, f := range r.File {
		// Construct the full path for the extracted file
		fpath := filepath.Join(dest, f.Name)

		if index == 0 {
			topLevelDirectory = fpath
		}

		// Check if it's a directory
		if f.FileInfo().IsDir() {
			// Make the directory
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return "", err
			}
			continue
		}

		// Make the parent directory for the file (in case it doesn't exist)
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}

		// Open the file inside the zip archive
		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		// Create the destination file
		outFile, err := os.Create(fpath)
		if err != nil {
			rc.Close() // Close the file from the archive
			return "", err
		}

		// Copy the content from the archive to the destination file
		_, err = io.Copy(outFile, rc)

		// Close both files
		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err
		}
	}

	return topLevelDirectory, nil
}

// copyFile copies the contents of a source file to a destination file.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the file content
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Flush to disk
	err = destinationFile.Sync()
	return err
}

// copyDirectory recursively copies the contents of a source directory to a destination directory.
func CopyDirectory(srcDir, dstDir string) error {
	err := filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create the destination path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		// If it's a directory, create the destination directory
		if info.IsDir() {
			if err = os.MkdirAll(dstPath, info.Mode()); err != nil {
				return err
			}
		} else {
			// Otherwise, copy the file
			if err = CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
