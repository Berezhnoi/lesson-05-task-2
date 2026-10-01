package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config %q: %w", path, err)
	}

	var values map[string]any
	if err := json.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("decode config %q: %w", path, err)
	}
	return nil
}
