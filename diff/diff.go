package diff

import (
	"encoding/json"
	"os"
)

type User map[string]any

type Change struct {
	Type string `json:"type"`

	Key string `json:"key"`

	Old User `json:"old,omitempty"`

	New User `json:"new,omitempty"`
}

func Generate(
	oldFile string,
	newFile string,
	output string,
) error {

	oldUsers, err := loadJSON(oldFile)

	if err != nil {
		return err
	}

	newUsers, err := loadJSON(newFile)

	if err != nil {
		return err
	}

	changes := compare(
		oldUsers,
		newUsers,
	)

	out, err := os.Create(output)

	if err != nil {
		return err
	}

	defer out.Close()

	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")

	return encoder.Encode(changes)
}

func loadJSON(
	filename string,
) (map[string]User, error) {

	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var users []User

	if err := json.Unmarshal(
		data,
		&users,
	); err != nil {
		return nil, err
	}

	result := make(map[string]User)

	for _, user := range users {

		username, ok := user["sAMAccountName"].(string)

		if !ok || username == "" {
			continue
		}

		result[username] = user
	}

	return result, nil
}

func compare(
	old map[string]User,
	current map[string]User,
) []Change {

	changes := make([]Change, 0)

	// Added or modified users

	for username, newUser := range current {

		oldUser, exists := old[username]

		if !exists {

			changes = append(
				changes,
				Change{
					Type: "added",
					Key:  username,
					New:  newUser,
				},
			)

			continue
		}

		if !equal(
			oldUser,
			newUser,
		) {

			changes = append(
				changes,
				Change{
					Type: "modified",
					Key:  username,
					Old:  oldUser,
					New:  newUser,
				},
			)
		}
	}

	// Removed users

	for username, oldUser := range old {

		if _, exists := current[username]; !exists {

			changes = append(
				changes,
				Change{
					Type: "removed",
					Key:  username,
					Old:  oldUser,
				},
			)
		}
	}

	return changes
}

func equal(
	a User,
	b User,
) bool {

	aJSON, _ := json.Marshal(a)

	bJSON, _ := json.Marshal(b)

	return string(aJSON) == string(bJSON)
}
