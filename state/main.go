package state

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// Create a new state
func New() *State {
	return &State{
		LastPollTimestamp: time.Now().Add(-1 * time.Hour * 24 * 1).Format(time.RFC3339),
	}
}

// Check if state exists
func Exists(statePath string) bool {
	info, err := os.Stat(statePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Save state
func Save(currentState *State, statePath string) {
	// Marshal to JSON
	file, _ := json.MarshalIndent(&currentState, "", " ")

	// Write to file
	_ = ioutil.WriteFile(statePath, file, 0644)
}

// Restore state
func Restore(statePath string) (*State, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(statePath)

	// if os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read the opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Initialize our state struct
	var state State

	// unmarshal our byteArray which contains our
	// jsonFile's content into 'state' which we defined above
	err = json.Unmarshal(byteValue, &state)

	// if json.Unmarshal returns an error then handle it
	if err != nil {
		return nil, err
	}

	return &state, nil
}
