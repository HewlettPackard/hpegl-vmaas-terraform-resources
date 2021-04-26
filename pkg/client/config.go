// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const GltformExtension = ".gltform"

type Gljwt struct {
	SpaceName string `yaml:"space_name,omitempty"`
	ProjectID string `yaml:"project_id"`
	RestURL   string `yaml:"rest_url"`
	Token     string `yaml:"access_token"`
}

func getGLConfig() (gljwt *Gljwt, err error) {
	homeDir, _ := os.UserHomeDir()
	workingDir, _ := os.Getwd()
	for _, p := range []string{homeDir, workingDir} {
		gljwt, err = loadGLConfig(p)
		if err == nil {
			break
		}
	}

	return gljwt, err
}

func loadGLConfig(dir string) (*Gljwt, error) {
	f, err := os.Open(filepath.Clean(filepath.Join(dir, GltformExtension)))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseGLStream(f)
}

func parseGLStream(s io.Reader) (*Gljwt, error) {
	contents, err := ioutil.ReadAll(s)
	if err != nil {
		return nil, err
	}

	q := &Gljwt{}
	err = yaml.Unmarshal(contents, q)

	return q, err
}
