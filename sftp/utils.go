package sftp

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

func getHostKey(host string) ssh.PublicKey {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read known_hosts file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing %q: %v\n", fields[2], err)
				os.Exit(1)
			}
			break
		}
	}

	if hostKey == nil {
		fmt.Fprintf(os.Stderr, "No hostkey found for %s", host)
		os.Exit(1)
	}

	return hostKey
}

func humanReadableSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
