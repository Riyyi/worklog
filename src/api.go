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
	if itemID == "break" || itemID == "lunch" || itemID == "pauze" || itemID == "T1-break" || itemID == "T1-lunch" || itemID == "T1-pauze" { return nil }

	if date == "" || fromTime == "" || toTime == "" || itemID == "" {
		return fmt.Errorf("incomplete log entry: %s, %s-%s, %s, %s", date, fromTime, toTime, itemID, description)
	}

	const layout = "2006-01-02 15:04:05"
	timestamp1 := date + " " + fromTime + ":00"
	timestamp2 := date + " " + toTime + ":00"

	location, err := time.LoadLocation("Local")
	if err != nil {
		return fmt.Errorf("error loading location: %s", err)
	}

	time1, err := time.ParseInLocation(layout, timestamp1, location)
	if err != nil { return fmt.Errorf("error parsing from_time: %s", err) }

	time2, err := time.ParseInLocation(layout, timestamp2, location)
	if err != nil { return fmt.Errorf("error parsing to_time: %s", err) }

	// Convert local timezone to UTC time
	time1UTC := time1.UTC()

	duration := time2.Sub(time1)
	seconds := int(duration.Seconds())
	if seconds < 0 { return fmt.Errorf("from_time is later than to_time: %s > %s", fromTime, toTime) }

	var url string = baseUrl + "/rest/api/2/issue/" + itemID + "/worklog"

    data := map[string]string{
		"comment": description,
		"started": fmt.Sprintf("%s:00.000+0000",
			time1UTC.Format("2006-01-02T15:04")), // "2021-01-17T12:34:00.000+0000",
		"timeSpentSeconds": strconv.Itoa(seconds),
    }

	_, err = Request(url, data, 201) // "Created"

	return err
}
