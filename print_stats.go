package main

import (
	"fmt"
	"sort"
	"time"

	//"gopkg.in/src-d/go-git.v4"
	//"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const outOfRange = 99999
const daysInLastSixMonths = 183
const weeksInLastSixMonths = 26

type column []int

// printCell given a cell value prints it with a different format
// based on the value amount, and on the `today` flag.
func printCell(val int, today bool) {
    escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

// func printDayCol given the day number (0 is Sunday) prints the day name,
// alternating the rows (prints just 2,4,6)
func printDayCol(day int) {
	out := "     "
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}

	fmt.Printf(out)
}

// func printMonths prints the month names in the first line,
// determining when the month changed between switching weeks
func printMonths() {
	week := getBeginningOfDay(time.Now().Add(-(daysInLastSixMonths*time.Hour*24)))
	month := week.Month()
	fmt.Printf("         ")

	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7*time.Hour*24)
		if week.After(time.Now()) {
			break
		}
	}

	fmt.Printf("\n")
}

// func printCells prints the cells of the graph
func printCells(cols map[int]column) {
	printMonths()
	for j := 6; j >= 0; j-- {	// iterate each day of week
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {	// iterate each week
			if i == weeksInLastSixMonths+1 {
				printDayCol(j)
			}
			if col, ok := cols[i]; ok {
				//special case today
				if i == 0 && j == calcOffset()-1 {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

// func buildCols generates a map with rows and columns ready to be printed
func buildCols(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := int(k/7)
		day_of_week := k%7

		if day_of_week == 0 {
			col = column{}
		}

		col = append(col, commits[k])

		if day_of_week == 6 {
			cols[week] = col
		}
	}

	return cols
}

// func sortMapIntoSlice sorts the indexes of a map
func sortMapIntoSlice(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	return keys
}

// func printCommitStats prints the commit stats
func printCommitStats(commits map[int]int) {
	keys := sortMapIntoSlice(commits)
	cols := buildCols(keys, commits)
	printCells(cols)
}