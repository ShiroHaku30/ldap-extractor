package extract

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func FindPrevious(
	dir string,
	suffix string,
	current string,
) (string, error) {

	files, err := filepath.Glob(
		filepath.Join(
			dir,
			"*-"+suffix,
		),
	)

	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", nil
	}

	sort.Strings(files)

	// Remove current file from list
	filtered := make([]string, 0)

	for _, file := range files {

		if file != current {
			filtered = append(
				filtered,
				file,
			)
		}
	}

	if len(filtered) == 0 {
		return "", nil
	}

	// Since filenames start with timestamp,
	// the last sorted file is newest
	return filtered[len(filtered)-1], nil
}

func Cleanup(
	dir string,
	suffix string,
	retention int,
) error {

	if retention <= 0 {
		return nil
	}

	files, err := filepath.Glob(
		filepath.Join(
			dir,
			"*-"+suffix,
		),
	)

	if err != nil {
		return err
	}

	sort.Strings(files)

	if len(files) <= retention {
		return nil
	}

	removeCount := len(files) - retention

	for i := 0; i < removeCount; i++ {

		fmt.Println(
			"Removing old file:",
			files[i],
		)

		if err := os.Remove(
			files[i],
		); err != nil {

			return fmt.Errorf(
				"remove %s: %w",
				files[i],
				err,
			)
		}
	}

	return nil
}
