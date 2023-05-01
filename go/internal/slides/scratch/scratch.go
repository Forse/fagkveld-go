package main

import "os"

func SetPermissions() error {
	f, err := os.Open("file.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	f.Chmod(0644)
	f.Chown(0, 0)
	return nil
}
