package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// LoadJSONConfig loads the given JSON config into the given config struct.
// config must be a pointer to a struct.
func LoadJSONConfig(config interface{}, configPath string) {
	ext := path.Ext(configPath)
	if ext != ".json" {
		panic(fmt.Sprintf("invalid file extension: %s (expected .json)", ext))
	}
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	// this will panic if the given config is not a pointer
	if err = json.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}
}

// LoadFromPathOrEnvIfSet loads the given JSON from the given path or env variable (if set).
// Optionally copies the config file from the given path to the env path.
func LoadFromPathOrEnvIfSet(config interface{}, configPath string, envPath string, copySample bool) {
	configEnvPath := os.Getenv(envPath)
	if len(configEnvPath) == 0 {
		LoadJSONConfig(config, configPath)
	} else {
		if copySample {
			// only copy sample config if it doesn't exist in the dest
			if _, err := os.Stat(configEnvPath); err != nil {
				if os.IsNotExist(err) {
					// check if the sample config exists
					if _, err2 := os.Stat(configPath); err2 != nil {
						panic(err2)
					}
					if mvErr := MoveSampleConfig(configPath, configEnvPath); mvErr != nil {
						panic(mvErr)
					}
				} else {
					panic(err)
				}
			}
		}
		LoadJSONConfig(config, configEnvPath)
	}
}

// MoveSampleConfig copies the given sample config to the given destination path.
func MoveSampleConfig(samplePath string, destPath string) error {
	source, err := os.Open(samplePath)
	if err != nil {
		return err
	}
	defer source.Close()
	dest, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer dest.Close()
	if _, err := io.Copy(dest, source); err != nil {
		return err
	}
	return nil
}
