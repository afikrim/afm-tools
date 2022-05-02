package config

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/afikrim/afm-tools/helper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PostmanApiKey              string `yaml:"postman_api_key"`
	PostmanWorkspaceID         string `yaml:"postman_workspace_id"`
	PostmanPersonalApiKey      string `yaml:"postman_personal_api_key"`
	PostmanPersonalWorkspaceID string `yaml:"postman_personal_workspace_id"`
}

func LoadConfig(init bool) (*Config, error) {
	config := &Config{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dirPath := homeDir + "/.afm-tools"
	if _, err := os.Stat(dirPath); err != nil && errors.Is(err, os.ErrNotExist) {
		if !init {
			return nil, fmt.Errorf("config directory not found")
		}
		if err := os.Mkdir(dirPath, 0777); err != nil {
			return nil, err
		}
	}
	if _, err := os.Stat(dirPath + "/config.yaml"); err != nil && errors.Is(err, os.ErrNotExist) {
		if !init {
			return nil, fmt.Errorf("config file not found")
		}
		if _, err := os.Create(dirPath + "/config.yaml"); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	configRaw, err := ioutil.ReadFile(dirPath + "/config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configRaw, config)
	if err != nil {
		return nil, err
	}

	postman := helper.NewPostman(http.DefaultClient)
	if init {
		askForPostmanApiKey(config)
		askForPostmanWorkspace(config, postman)
		askForPersonalPostmanApiKey(config)
		askForPersonalPostmanWorkspace(config, postman)
	}

	updateConfigFile(dirPath, config)

	return config, nil
}

func createConfigFile(path string) error {

	return nil
}

func askForPostmanApiKey(config *Config) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Main API Key [%s]: ", config.PostmanApiKey)
	postmanApiKey, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	postmanApiKey = strings.TrimSuffix(postmanApiKey, "\n")
	if postmanApiKey == "" {
		postmanApiKey = config.PostmanApiKey
	}

	config.PostmanApiKey = postmanApiKey
	return nil
}

func askForPostmanWorkspace(config *Config, postman helper.Postman) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Main Workspace [%s]: ", "My Workspace")
	postmanWorkspace, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	postmanWorkspace = strings.TrimSuffix(postmanWorkspace, "\n")
	if postmanWorkspace == "" {
		postmanWorkspace = "My Workspace"
	}

	workspaceID, err := postman.GetPostmanWorkspaceID(postmanWorkspace, &config.PostmanApiKey)
	if err != nil {
		return err
	}
	if workspaceID == nil {
		workspaceID, err = postman.CreatePostmanWorkspaceReturnID(postmanWorkspace, &config.PostmanApiKey)
		if err != nil {
			return err
		}
	}

	config.PostmanWorkspaceID = *workspaceID
	return nil
}

func askForPersonalPostmanApiKey(config *Config) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Personal API Key [%s]: ", config.PostmanPersonalApiKey)
	personalPostmanApiKey, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	personalPostmanApiKey = strings.TrimSuffix(personalPostmanApiKey, "\n")
	if personalPostmanApiKey == "" {
		personalPostmanApiKey = config.PostmanPersonalApiKey
	}

	config.PostmanPersonalApiKey = personalPostmanApiKey
	return nil
}

func askForPersonalPostmanWorkspace(config *Config, postman helper.Postman) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Personal Workspace [%s]: ", "My Workspace")
	personalPostmanWorkspace, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	personalPostmanWorkspace = strings.TrimSuffix(personalPostmanWorkspace, "\n")
	if personalPostmanWorkspace == "" {
		personalPostmanWorkspace = "My Workspace"
	}

	personalWorkspaceId, err := postman.GetPostmanWorkspaceID(personalPostmanWorkspace, &config.PostmanPersonalApiKey)
	if err != nil {
		return err
	}
	if personalWorkspaceId == nil {
		personalWorkspaceId, err = postman.CreatePostmanWorkspaceReturnID(personalPostmanWorkspace, &config.PostmanPersonalApiKey)
		if err != nil {
			return err
		}
	}

	config.PostmanPersonalWorkspaceID = *personalWorkspaceId
	return nil
}

func updateConfigFile(path string, config *Config) error {
	configRaw, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path+"/config.yaml", configRaw, 0644)
	if err != nil {
		return err
	}

	return nil
}
