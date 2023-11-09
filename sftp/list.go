package sftp

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"github.com/spf13/viper"
)

func ListFiles(sc *sftp.Client) error {
	remoteDir := viper.GetString("SERVER.PATH")
	fmt.Fprintf(os.Stdout, "Listing [%s] ...\n\n", remoteDir)

	files, err := sc.ReadDir(remoteDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to list remote dir: %v\n", err)
		return err
	}

	for _, f := range files {
		var name, modTime, size string

		name = f.Name()
		modTime = f.ModTime().Format("2006-01-02 15:04:05")

		if f.IsDir() {
			name = name + "/"
			modTime = ""
			size = "DIR"
		} else {
			size = humanReadableSize(f.Size())
		}
		fmt.Fprintf(os.Stdout, "%19s %12s %s\n", modTime, size, name)
	}

	return nil
}

func ListFilesToFile(sc *sftp.Client) error {
	remoteDir := viper.GetString("SFTP.SFTP_PATH")
	fmt.Fprintf(os.Stdout, "Listing [%s] ...\n\n", remoteDir)

	files, err := sc.ReadDir(remoteDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to list remote dir: %v\n", err)
		return err
	}

	file, err := os.Create("file_listing.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create CSV file: %v\n", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Modification Time", "Size", "Name"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write CSV file: %v\n", err)
		return err
	}

	for _, f := range files {
		var name, modTime, size string
		if f.IsDir() || !strings.Contains(f.Name(), viper.GetString("GREP_PATTERN")) {
			continue
		}

		name = f.Name()
		modTime = f.ModTime().Format("2006-01-02 15:04:05")

		if f.IsDir() {
			name = name + "/"
			modTime = ""
			size = "DIR"
		} else {
			size = humanReadableSize(f.Size())
		}

		err = writer.Write([]string{modTime, size, name})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write CSV file: %v\n", err)
			return err
		}
	}

	fmt.Println("File listing exported to file_listing.csv")
	return nil
}
