package core

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/olekukonko/tablewriter"
)

var ProcessPath string

func init() {
	if RootCheck(false) {
		fmt.Println("Do not run Apx as root!")
		os.Exit(1)
	}
}

func RootCheck(display bool) bool {
	if os.Geteuid() != 0 {
		if display {
			fmt.Println("You must be root to run this command")
		}
		return false
	}
	return true
}

func AskConfirmation(s string) bool {
	var response string
	fmt.Print(s + " [y/N]: ")
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {
		return true
	}
	return false
}

func CopyToUserTemp(path string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(user.HomeDir, ".cache", "apx")
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return "", err
		}
	}

	fileName := filepath.Base(path)
	newPath := filepath.Join(cacheDir, fileName)

	pathContents, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer pathContents.Close()

	newPathContents, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer newPathContents.Close()

	_, err = newPathContents.ReadFrom(pathContents)
	if err != nil {
		return "", err
	}

	return newPath, nil
}

// getPrettifiedDate returns a human readable date from a timestamp
func getPrettifiedDate(date string) string {
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", date)
	if err != nil {
		return date
	}

	// If the date is less than 24h ago, return the time since
	if t.After(time.Now().Add(-24 * time.Hour)) {
		duration := time.Since(t).Round(time.Hour)
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh ago", hours)
	}

	return t.Format("02 Jan 2006 15:04:05")
}

func CreateApxTable(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetColumnSeparator("┊")
	table.SetCenterSeparator("┼")
	table.SetRowSeparator("┄")
	table.SetHeaderLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)

	return table
}

func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}
