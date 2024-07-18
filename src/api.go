/*
 * Copyright (C) 2024 Riyyi
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package src

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var Api api

type api struct {}

func (api) CallApi(date string, from_time string, to_time string, item_id string, description string) error {
	if item_id == "break" || item_id == "lunch" || item_id == "pauze" { return nil }

	if date == "" || from_time == "" || to_time == "" || item_id == "" {
		return fmt.Errorf("incomplete log entry: %s, %s-%s, %s, %s", date, from_time, to_time, item_id, description)
	}

	time1, err := time.Parse("15:04", from_time)
	if err != nil { return fmt.Errorf("error parsing from_time: %s", err) }

	time2, err := time.Parse("15:04", to_time)
	if err != nil { return fmt.Errorf("error parsing to_time: %s", err) }

	// Convert local timezone to UTC time
	_, offset := time.Now().Zone()
	time1.Add(-time.Duration(offset) * time.Second);
	time2.Add(-time.Duration(offset) * time.Second);

	duration := time2.Sub(time1)
	seconds := int(duration.Seconds())
	if seconds < 0 { return fmt.Errorf("from_time is later than to_time: %s > %s", from_time, to_time) }

	var url string = base_url + "/rest/api/2/issue/" + item_id + "/worklog"

    data := map[string]string{
		"comment": description,
		"started": fmt.Sprintf("%sT%s:00.000+0000", date, from_time), // "2021-01-17T12:34:00.000+0000",
		"timeSpentSeconds": strconv.Itoa(seconds),
    }

    json_data, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("error marshaling JSON: %s", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
    if err != nil {
        return fmt.Errorf("error creating request: %s", err)
    }

    auth := username + ":" + password
    authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
    req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil { return fmt.Errorf("error making request: %s", err) }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil { return fmt.Errorf("error reading response body: %s", err) }

	if resp.Status != "201 Created" {
		return fmt.Errorf("invalid Jira request:\n%s", string(body))
	}

	return nil
}
