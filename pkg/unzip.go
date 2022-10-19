package pkg

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// UnzipSource принимает путь до архива и путь куда архив будет разарх.
// Ищет в архиве директорию tdata и разархивирует её, отбрасывая все остальные файлы.
// Если в архиве отсутствует директория tdata, вернет ошибку
// В случае успеха, вернет путь до директории tdata
func UnzipSource(source, destination string) (string, error) {
	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return "", err
	}
	var filePath, tdataPath string
	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		// fmt.Println("file name: ", f.Name, "is dir: ", f.FileInfo().IsDir())

		// 4. Check if file paths are not vulnerable to Zip Slip
		filePath = filepath.Join(destination, f.Name)
		// fmt.Println("filePath:  ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return "", fmt.Errorf("invalid file path: %s", filePath)
		}

		if strings.Contains(filePath, "tdata") {
			tdataPath = filePath
			err := unzipFile(f, filePath)
			if err != nil {
				return "", err
			}
		}
	}

	r := strings.Split(tdataPath, "tdata/")
	if len(r) == 1 {
		return "", errors.New("tdata not found in archive: " + source)
	}
	res := r[0] + "tdata/"

	return res, nil
}

func unzipFile(f *zip.File, filePath string) error {

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}

	return nil
}

// tdata, err := UnzipSource("10046.zip", "")
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Println("RES:  ", tdata)
