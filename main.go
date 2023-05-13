package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	ReleaseURL     string `yaml:"releaseUrl"`
	TargetFolder   string `yaml:"targetFolder"`
	TargetFileName string `yaml:"targetFileName"`
}

func main() {
	config := Config{}
	configFile, err := ioutil.ReadFile("pineapple-update.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(config.ReleaseURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		panic(err)
	}

	if len(release.Assets) == 0 {
		panic("no assets")
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

	_, err = ioutil.ReadFile(asset.Name)
	if err == nil {
		fmt.Println("already up to date")
		return
	}

	resp, err = http.Get(asset.BrowserDownloadURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(asset.Name, body, 0755)
	if err != nil {
		panic(err)
	}

	fmt.Println("downloaded " + asset.Name)

	if config.TargetFileName != "" {
		if _, err := os.Lstat(config.TargetFileName); err == nil {
			os.Remove(config.TargetFileName)
		}

		err = os.Symlink(asset.Name, config.TargetFileName)
		if err != nil {
			panic(err)
		}
	}
}