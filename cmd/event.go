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
)

const defaultCalendar = "primary"

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Returns events on the specified calendar",
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
			log.Fatalf("Unable to get calendar ID of %s: %s", calendarSummary, err)
		}

		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next 10 of the user's events: %v", err)
		}
		fmt.Println("Upcoming events:")
		if len(events.Items) == 0 {
			fmt.Println("No upcoming events found")
		} else {
			for _, item := range events.Items {
				date := item.Start.DateTime
				if date == "" { // "All day" event
					date = item.Start.Date
				}
				fmt.Printf("%v (%v)\n", item.Summary, date)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(eventCmd)
	eventCmd.Flags().StringP("calendar", "c", "", "Calendar summary")
}
