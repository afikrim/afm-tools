package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/afikrim/afm-tools/config"
)

const (
	postmanApiCollectionsUrl = "https://api.getpostman.com/collections"
	postmanApiWorkspaceUrl   = "https://api.getpostman.com/workspaces"
)

type SyncPostman interface {
	SyncPostmanCollection(collectionId string) error
}

type syncPostman struct {
	config     *config.Config
	httpClient *http.Client
}

func NewSyncPostman(config *config.Config, httpClient *http.Client) SyncPostman {
	return &syncPostman{config: config, httpClient: httpClient}
}

func (s *syncPostman) SyncPostmanCollection(collectionName string) error {
	collectionId, err := s.getPostmanCollectionId(collectionName)
	if err != nil {
		return err
	}

	collection, err := s.getPostmanCollection(*collectionId)
	if err != nil {
		return err
	}

	err = s.createPostmanCollection(*collection)
	if err != nil {
		return err
	}

	return nil
}

func (s *syncPostman) getPostmanCollectionId(collectionName string) (*string, error) {
	url := postmanApiCollectionsUrl
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Api-Key", s.config.PostmanApiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var collections map[string]interface{}
	err = json.Unmarshal(body, &collections)
	if err != nil {
		return nil, err
	}

	var collectionId string
	for _, c := range collections["collections"].([]interface{}) {
		if c.(map[string]interface{})["name"] == collectionName {
			collectionId = c.(map[string]interface{})["id"].(string)
			break
		}
	}
	if collectionId == "" {
		return nil, fmt.Errorf("collection %s not found", collectionName)
	}

	return &collectionId, nil
}

func (s *syncPostman) getPostmanCollection(collectionId string) (*map[string]interface{}, error) {
	url := postmanApiCollectionsUrl + "/" + collectionId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Api-Key", s.config.PostmanApiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get postman collection: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var collection map[string]interface{}
	err = json.Unmarshal(body, &collection)

	return &collection, nil
}

func (s *syncPostman) createPostmanCollection(collection map[string]interface{}) error {
	url := postmanApiCollectionsUrl
	if s.config.PostmanPersonalWorkspaceID != "" {
		url = url + "?workspace=" + s.config.PostmanPersonalWorkspaceID
	}

	body, err := json.Marshal(collection)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", s.config.PostmanPersonalApiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to create postman collection: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)

	return nil
}
