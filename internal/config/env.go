package config

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnv(path string) error {
	open, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			panic(err)
		}
	}(open)

	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}
