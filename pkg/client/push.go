package client

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/syncromatics/idl-repository/pkg/config"

	"github.com/coreos/go-semver/semver"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type PushOptions struct {
	Configuration *config.Configuration
	Version       *semver.Version
}

func Push(options PushOptions) error {
	if len(options.Configuration.Provides) < 1 {
		return errors.New("nothing to push")
	}

	for _, provider := range options.Configuration.Provides {
		excludes := getExcludes(provider.IdlIgnore)
		tempFile, err := gzipRoot(provider.Root, excludes)
		if err != nil {
			return err
		}

		f, err := os.Open(tempFile)
		if err != nil {
			return errors.Wrap(err, "failed opening source.zip")
		}
		defer f.Close()

		url := fmt.Sprintf("%s/v1/projects/%s/types/%s/versions/%s",
			options.Configuration.Repository,
			options.Configuration.Name,
			provider.Type,
			options.Version.String())

		resp, err := http.Post(url, "", f)
		if err != nil {
			return errors.Wrap(err, "failed posting module to registry")
		}

		if resp.StatusCode != http.StatusCreated {
			return errors.New("upload failed")
		}
	}
	return nil
}

func getExcludes(idlIgnore string) []string {
	if idlIgnore == "" {
		excludes, err := readIdlIgnoreFile(".idlignore")
		if err != nil {
			return nil
		}
		return excludes
	}

	split := strings.Split(strings.Replace(strings.TrimSpace(idlIgnore), "\r\n", "\n", -1), "\n")
	if len(split) > 1 {
		return split
	}

	excludes, err := readIdlIgnoreFile(idlIgnore)
	if err != nil {
		return nil
	}
	return excludes

}

func readIdlIgnoreFile(path string) ([]string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func gzipRoot(root string, excludes []string) (string, error) {
	name := xid.New()
	tempLocation := fmt.Sprintf("%s/%s", os.TempDir(), name)

	_, err := os.Stat(root)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read provider root '%s'", root)
	}

	f, err := os.Create(tempLocation)
	if err != nil {
		return "", errors.Wrap(err, "failed creating temp file for gzip")
	}

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	cleanRoot := strings.TrimPrefix(root, "./")

	pm, err := fileutils.NewPatternMatcher(excludes)
	if err != nil {
		return "", err
	}

	err = filepath.Walk(root, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		relFile, err := filepath.Rel(cleanRoot, file)
		if err != nil {
			// Error getting relative path OR we are looking
			// at the source directory path. Skip in both situations.
			return nil
		}

		skip, err := pm.Matches(relFile)
		if err != nil {
			return err
		}

		if skip {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		if file == root {
			return nil
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(file, cleanRoot)

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})

	if err != nil {
		return "", errors.Wrap(err, "failed writing files to tar")
	}

	return tempLocation, nil
}
