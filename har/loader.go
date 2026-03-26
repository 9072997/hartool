package har

import (
	"encoding/json"
	"os"
)

func Load(path string) (*HAR, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var h HAR
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	return &h, nil
}
