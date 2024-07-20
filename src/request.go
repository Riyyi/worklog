package src

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Request[T ~map[string]string | ~map[string]interface{}](url string, data T, status int) ([]byte, error) {
    jsonData, err := json.Marshal(data)
    if err != nil { return nil, fmt.Errorf("error marshaling JSON: %s", err) }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil { return nil, fmt.Errorf("error creating request: %s", err) }

    auth := username + ":" + password
    authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
    req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil { return nil, fmt.Errorf("error making request: %s", err) }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil { return nil, fmt.Errorf("error reading response body: %s", err) }

	if resp.StatusCode != status {
		return nil, fmt.Errorf("invalid Jira request:\n%s", string(body))
	}

	return body, nil
}
