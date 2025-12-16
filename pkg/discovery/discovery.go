package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Sync sends the given Entry to the discovery server and retrieves a list of matching Entries in response.
func Sync(serverURL string, masterKey string, entry Entry) (*[]Entry, error) {
	data, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", serverURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Master-Key", masterKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var entries []Entry
	if err := json.Unmarshal(respBody, &entries); err != nil {
		return nil, err
	}

	return &entries, nil
}
