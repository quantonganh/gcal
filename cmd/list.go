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

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns the calendars on the user's calendar list",
	Run: func(cmd *cobra.Command, args []string) {
		srv, err := newCalendar()
		if err != nil {
			log.Fatalf("Unable to retrieve calendar client: %v", err)
		}

		for {
			calendarList, err := srv.CalendarList.List().Do()
			if err != nil {
				log.Fatalf("Cannot get calendars list: %v", err)
			}
			for _, item := range calendarList.Items {
				fmt.Println(item.Summary)
			}
			pageToken := calendarList.NextPageToken
			if pageToken == "" {
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
