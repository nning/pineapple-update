package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Release struct {
	Assets []Asset `json:"assets"`
}

type Config struct {
	TargetFolder      string `yaml:"targetFolder"`
	Symlink           string `yaml:"symlink"`
	SymlinkFileName   string `yaml:"symlinkFileName"`
	RemoveOldVersions string `yaml:"removeOldVersions"`
}

const releaseUrl = "https://api.github.com/repos/pineappleEA/pineapple-src/releases/latest"

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	err = os.Chdir(filepath.Dir(ex))
	if err != nil {
		panic(err)
	}

	configs := []string{
		".pineapple-update.yml",
		"pineapple-update.yml",
	}

	// check if config file exists
	var configPath string
	for _, c := range configs {
		if _, err := os.Lstat(c); err == nil {
			configPath = c
			break
		}
	}

	config := Config{}
	if configPath != "" {
		configFile, err := os.ReadFile(configPath)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			panic(err)
		}
	}

	resp, err := http.Get(releaseUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		panic(err)
	}

	if len(release.Assets) == 0 {
		panic("release does not contain any assets")
	}

	var asset Asset
	for _, a := range release.Assets {
		if strings.HasSuffix(a.Name, ".AppImage") {
			asset = a
			break
		}
	}

	if config.TargetFolder != "" {
		err = os.Chdir(config.TargetFolder)
		if err != nil {
			panic(err)
		}
	}

	_, err = os.ReadFile(asset.Name)
	if err == nil {
		fmt.Println("already up to date")
		if config.RemoveOldVersions != "false" {
			goto cleanup
		}
		return
	}

	resp, err = http.Get(asset.BrowserDownloadURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(asset.Name, body, 0755)
	if err != nil {
		panic(err)
	}

	fmt.Println("downloaded " + asset.Name)

cleanup:
	if config.Symlink != "false" {
		symlinkTarget := config.SymlinkFileName
		if symlinkTarget == "" {
			symlinkTarget = "yuzu-ea.AppImage"
		}

		if _, err := os.Lstat(symlinkTarget); err == nil {
			os.Remove(symlinkTarget)
		}

		err = os.Symlink(asset.Name, symlinkTarget)
		if err != nil {
			panic(err)
		}
	}

	if config.RemoveOldVersions != "false" {
		files, err := os.ReadDir(".")
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			if strings.HasPrefix(f.Name(), "Linux-Yuzu-EA-") && strings.HasSuffix(f.Name(), ".AppImage") && f.Name() != asset.Name {
				err = os.Remove(f.Name())
				if err != nil {
					panic(err)
				}

				fmt.Println("removed " + f.Name())
			}
		}
	}
}
