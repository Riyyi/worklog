/*
 * Copyright (C) 2024 Riyyi
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package src

import (
	"fmt"
	"strconv"
	"time"
)

var Api api

type api struct {}

func (api) CallApi(date string, fromTime string, toTime string, itemID string, description string) error {
	if itemID == "break" || itemID == "lunch" || itemID == "pauze" { return nil }

	if date == "" || fromTime == "" || toTime == "" || itemID == "" {
		return fmt.Errorf("incomplete log entry: %s, %s-%s, %s, %s", date, fromTime, toTime, itemID, description)
	}

	time1, err := time.Parse("15:04", fromTime)
	if err != nil { return fmt.Errorf("error parsing from_time: %s", err) }

	time2, err := time.Parse("15:04", toTime)
	if err != nil { return fmt.Errorf("error parsing to_time: %s", err) }

	// Convert local timezone to UTC time
	_, offset := time.Now().Zone()
	time1.Add(-time.Duration(offset) * time.Second);
	time2.Add(-time.Duration(offset) * time.Second);

	duration := time2.Sub(time1)
	seconds := int(duration.Seconds())
	if seconds < 0 { return fmt.Errorf("from_time is later than to_time: %s > %s", fromTime, toTime) }

	var url string = baseUrl + "/rest/api/2/issue/" + itemID + "/worklog"

    data := map[string]string{
		"comment": description,
		"started": fmt.Sprintf("%sT%s:00.000+0000", date, fromTime), // "2021-01-17T12:34:00.000+0000",
		"timeSpentSeconds": strconv.Itoa(seconds),
    }

	_, err = Request(url, data, 201) // "Created"

	return err
}
