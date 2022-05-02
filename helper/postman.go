package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	postmanApiCollectionsUrl = "https://api.getpostman.com/collections"
	postmanApiWorkspaceUrl   = "https://api.getpostman.com/workspaces"
)

type Postman interface {
	GetPostmanWorkspaces(apiKey *string) ([]interface{}, error)
	CreatePostmanWorkspace(workspaceName string, apiKey *string) (interface{}, error)
	GetPostmanCollections(workspaceId string, apiKey *string) ([]interface{}, error)
	CreatePostmanCollection(workspaceId string, payload map[string]interface{}, apiKey *string) (interface{}, error)
	GetPostmanWorkspaceID(workspaceName string, apiKey *string) (*string, error)
	CreatePostmanWorkspaceReturnID(workspaceName string, apiKey *string) (*string, error)
}

type postman struct {
	httpClient *http.Client
}

func NewPostman(httpClient *http.Client) Postman {
	return &postman{httpClient: httpClient}
}

func (p *postman) GetPostmanWorkspaces(apiKey *string) ([]interface{}, error) {
	var workspaces map[string]interface{}

	url := postmanApiWorkspaceUrl
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if apiKey == nil || *apiKey == "" {
		return nil, fmt.Errorf("api key not provided")
	}
	req.Header.Add("X-Api-Key", *apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get postman workspace: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces["workspaces"].([]interface{}), nil
}

func (p *postman) CreatePostmanWorkspace(workspaceName string, apiKey *string) (interface{}, error) {
	var workspace map[string]interface{}

	url := postmanApiWorkspaceUrl
	workspacePayload := map[string]interface{}{"name": workspaceName, "type": "personal"}
	body, err := json.Marshal(map[string]interface{}{"workspace": workspacePayload})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if apiKey == nil || *apiKey == "" {
		return nil, fmt.Errorf("api key not provided")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", *apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to create postman workspace: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}

	return workspace["workspace"], nil
}

func (p *postman) GetPostmanCollections(workspaceId string, apiKey *string) ([]interface{}, error) {
	var collections map[string]interface{}

	url := fmt.Sprintf("%s?workspace=%s", postmanApiCollectionsUrl, workspaceId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if apiKey == nil || *apiKey == "" {
		return nil, fmt.Errorf("api key not provided")
	}
	req.Header.Add("X-Api-Key", *apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get postman collections: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &collections)
	if err != nil {
		return nil, err
	}

	return collections["collections"].([]interface{}), nil
}

func (p *postman) CreatePostmanCollection(workspaceId string, payload map[string]interface{}, apiKey *string) (interface{}, error) {
	var collection map[string]interface{}

	url := fmt.Sprintf("%s?workspace=%s", postmanApiCollectionsUrl, workspaceId)
	body, err := json.Marshal(map[string]interface{}{"collection": payload})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if apiKey == nil || *apiKey == "" {
		return nil, fmt.Errorf("api key not provided")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", *apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to create postman collection: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &collection)
	if err != nil {
		return nil, err
	}

	return collection["collection"], nil
}

func (p *postman) GetPostmanWorkspaceID(workspaceName string, apiKey *string) (*string, error) {
	var workspaceID *string

	workspaces, err := p.GetPostmanWorkspaces(apiKey)
	if err != nil {
		return nil, err
	}

	for _, w := range workspaces {
		if w.(map[string]interface{})["name"] == workspaceName {
			tempWorkspaceID := w.(map[string]interface{})["id"].(string)
			workspaceID = &tempWorkspaceID
		}
	}

	return workspaceID, nil
}

func (p *postman) CreatePostmanWorkspaceReturnID(workspaceName string, apiKey *string) (*string, error) {
	workspace, err := p.CreatePostmanWorkspace(workspaceName, apiKey)
	if err != nil {
		return nil, err
	}

	workspaceID := workspace.(map[string]interface{})["id"].(string)
	return &workspaceID, nil
}
