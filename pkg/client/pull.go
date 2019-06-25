package client

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/syncromatics/idl-repository/pkg/config"

	"github.com/pkg/errors"
)

type PullOptions struct {
	Configuration *config.Configuration
}

func Pull(options PullOptions) error {
	if len(options.Configuration.Dependencies) < 1 {
		return errors.New("nothing to pull")
	}

	for _, dependency := range options.Configuration.Dependencies {
		path := fmt.Sprintf("%s/v1/projects/%s/types/%s/versions/%s/data.tar.gz",
			options.Configuration.ResolveRepository(dependency),
			dependency.Name,
			dependency.Type,
			dependency.Version)

		resp, err := http.Get(path)
		if err != nil {
			return errors.Wrap(err, "failed getting dependency")
		}

		if resp.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("status code %d is not OK", resp.StatusCode))
		}

		err = unPackDependency(options.Configuration, dependency, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func unPackDependency(configuration *config.Configuration, dependency config.Dependency, file io.ReadCloser) error {
	defer file.Close()

	pth := path.Join(configuration.IdlDirectory, dependency.Name, dependency.Type, dependency.Version)

	err := os.MkdirAll(pth, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create directories")
	}

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(pth, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
