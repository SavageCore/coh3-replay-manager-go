package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func ExtractZip(zipFile string, destPath string) error {
	// Open the ZIP file for reading
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create the destination directory if it doesn't exist
	err = os.MkdirAll(destPath, 0755)
	if err != nil {
		return err
	}

	// Extract each file from the ZIP archive
	for _, f := range r.File {
		// Open the file in the ZIP archive
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create the destination file
		destPath := filepath.Join(destPath, f.Name)
		if f.FileInfo().IsDir() {
			// If the file entry in the ZIP archive is a directory, create the directory
			err = os.MkdirAll(destPath, f.Mode())
			if err != nil {
				return err
			}
		} else {
			// If the file entry in the ZIP archive is a file, create the file and write its content
			destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
