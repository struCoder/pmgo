package utils

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

const (
	MINUTE = 60
	HOUR   = MINUTE * 60
	DAY    = HOUR * 24
	MONTH  = DAY * 30
	YEAR   = MONTH * 12
	BYTE   = 1
	KB     = 1024 * BYTE
	MB     = 1024 * KB
	GB     = 1024 * MB
)

// PadString will add totalSize spaces evenly to the right and left side of str.
// Returns str after applying the pad.
func PadString(str string, totalSize int) string {
	turn := 0
	for {
		if len(str) >= totalSize {
			break
		}
		if turn == 0 {
			str = " " + str
			turn ^= 1
		} else {
			str = str + " "
			turn ^= 1
		}
	}
	return str
}

// FormatUptime will figure out current proc uptime
func FormatUptime(startTime, currentTime int64) string {
	val := currentTime - startTime
	if val < MINUTE {
		return strconv.Itoa(int(val)) + "s"
	} else if val >= MINUTE && val < HOUR {
		return strconv.Itoa(int(val/MINUTE)) + "m"
	} else if val >= HOUR && val < DAY {
		return strconv.Itoa(int(val/HOUR)) + "h"
	} else if val >= DAY && val < MONTH {
		return strconv.Itoa(int(val/DAY)) + "d"
	} else if val >= MONTH && val < YEAR {
		return strconv.Itoa(int(val/MONTH)) + "M"
	}
	return strconv.Itoa(int(val/YEAR)) + "y"
}

// FormatMemory will format memory val
func FormatMemory(input int) string {
	if input < KB {
		return strconv.Itoa(input) + "KB"
	} else if input >= KB && input < MB {
		return strconv.Itoa(input/KB) + "KB"
	} else if input >= MB && input < GB {
		return strconv.Itoa(input/MB) + "MB"
	}
	return strconv.Itoa(input/GB) + "GB"
}

// GetTableWriter will return instance of tablewriter
func GetTableWriter() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator("─")
	table.SetRowSeparator("─")
	table.SetColumnSeparator("│")
	return table
}

// CheckSourceFolderExit is deprecated - no longer used in modern Go modules
// func CheckSourceFolderExit(sourceFolder string) (bool, error) {
// 	// This function was used for GOPATH structure and is no longer needed
// 	return true, nil
// }
