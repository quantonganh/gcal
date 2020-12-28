/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/api/calendar/v3"
)

const layoutISO = "2006-01-02"

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Creates an event",
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := newCalendar()
		if err != nil {
			log.Fatalf("Unable to retrieve calendar client: %v", err)
		}

		calendarSummary, err := cmd.Flags().GetString("calendar")
		if err != nil {
			log.Fatalf("Unable to get value of calendar flag: %v", err)
		}

		calendarID, err := getCalendarID(srv, calendarSummary)
		if err != nil {
			log.Fatalf("Unable to get calendar ID of %s: %v", calendarSummary, err)
		}

		eventTitle, err := cmd.Flags().GetString("title")
		if err != nil {
			log.Fatalf("Unable to get value of title flag: %v", err)
		}

		now := time.Now().Format(layoutISO)

		startDate, err := cmd.Flags().GetString("start-date")
		if err != nil {
			log.Fatalf("Unable to get value of start-date flag: %v", err)
		}
		if startDate == "" {
			startDate = now
		}

		endDate, err := cmd.Flags().GetString("end-date")
		if err != nil {
			log.Fatalf("Unable to get value of end-date flag: %v", err)
		}
		if endDate == "" {
			endDate = now
		}

		event := &calendar.Event{
			Summary: eventTitle,
			Start: &calendar.EventDateTime{
				Date: startDate,
			},
			End: &calendar.EventDateTime{
				Date: endDate,
			},
		}
		event, err = srv.Events.Insert(calendarID, event).Do()
		if err != nil {
			log.Fatalf("Unable to create event: %v\n", err)
		}
		fmt.Printf("Event created: %+v\n", event)
	},
}

func init() {
	now := time.Now().Format(layoutISO)

	eventCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringP("calendar", "c", "", "Calendar summary")
	insertCmd.Flags().StringP("title", "t", "", "Event title")
	insertCmd.MarkFlagRequired("title")
	insertCmd.Flags().StringP("start-date", "s", now, "Start date (inclusive)")
	insertCmd.Flags().StringP("end-date", "e", now, "End date (exclusive)")
}
