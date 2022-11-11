package archive

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// RecursiveZip | Will recursivly zip the folder to a destination path
func RecursiveZip(pathToZip string, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	defer destinationFile.Close()
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		fsFile.Close()
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

// Unzip | will unzip the current vault
func Unzip(src, dst string) error {
	archive, err := zip.OpenReader(src)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return errors.New("invalid file path whilst extracting data")
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}

// SingleToZip will take a single file and put it in to a zip file
func SingleToZip(filePath string, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	defer destinationFile.Close()
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	zipFile, err := myZip.Create(filePath)
	if err != nil {
		return err
	}
	fsFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(zipFile, fsFile)
	if err != nil {
		return err
	}
	fsFile.Close()
	myZip.Close()

	return nil
}
