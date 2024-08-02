package filehandler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
)

const (
	// have to use these or similar permissions to avoid permission denied errors in some cases
	folderPerms fs.FileMode = 0755
	// numWorkers  int         = 5
)

// Rezip is based on https://stackoverflow.com/a/63233911
// func Rezip(src, dest string) error {
// 	file, err := os.Create(dest)
// 	if err != nil {
// 		return err
// 	}

// 	defer file.Close()

// 	w := zip.NewWriter(file)
// 	defer w.Close()

// 	// var mimetypePath = src + string(os.PathSeparator) + "mimetype"
// 	// err = copyMimetypeToZip(w, src, mimetypePath)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	walker := func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		// skip empty directories
// 		if info.IsDir() {
// 			return nil
// 		}

// 		// if mimetypePath == path {
// 		// 	return nil
// 		// }

// 		err = writeToZip(w, src, path)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}
// 	err = filepath.Walk(src, walker)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func writeToZip(w *zip.Writer, src, path string) error {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// need a zip relative path to avoid creating extra directories inside of the zip
// 	var zipRelativePath = strings.Replace(path, src+string(os.PathSeparator), "", 1)
// 	f, err := w.Create(zipRelativePath)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = io.Copy(f, file)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func copyMimetypeToZip(w *zip.Writer, src, path string) error {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// need a zip relative path to avoid creating extra directories inside of the zip
// 	var zipRelativePath = strings.Replace(path, src+string(os.PathSeparator), "", 1)
// 	f, err := w.CreateHeader(&zip.FileHeader{
// 		Name:   strings.ReplaceAll(zipRelativePath, string(os.PathSeparator), "/"),
// 		Method: zip.Store,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	_, err = io.Copy(f, file)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func GetFilesFromZip(src string) (*zip.ReadCloser, map[string]*zip.File, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, nil, err
	}

	var zipFiles = make(map[string]*zip.File, 0)
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		zipFiles[f.Name] = f
	}

	return r, zipFiles, nil
}

func WriteZipCompressedString(w *zip.Writer, filename, contents string) error {
	var reader = bytes.NewReader([]byte(contents))
	f, err := w.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}

	return nil
}

func WriteZipCompressedFile(w *zip.Writer, zipFile *zip.File) error {
	file, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := w.Create(zipFile.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func WriteZipCompressedBytes(w *zip.Writer, filename string, data []byte) error {
	var reader = bytes.NewReader(data)
	f, err := w.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}

	return nil
}

func WriteZipUncompressedFile(w *zip.Writer, zipFile *zip.File) error {
	file, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := w.CreateHeader(&zip.FileHeader{
		Name:   zipFile.Name,
		Method: zip.Store,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func ReadInZipFileContents(zipFile *zip.File) (string, error) {
	file, err := zipFile.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var fileBytes = &bytes.Buffer{}
	_, err = io.Copy(fileBytes, file)
	if err != nil {
		return "", fmt.Errorf(`could not read in zip file contents for %q: %w`, zipFile.Name, err)
	}

	return fileBytes.String(), nil
}

func ReadInZipFileBytes(zipFile *zip.File) ([]byte, error) {
	file, err := zipFile.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileBytes = &bytes.Buffer{}
	_, err = io.Copy(fileBytes, file)
	if err != nil {
		return nil, fmt.Errorf(`could not read in zip file bytes for %q: %w`, zipFile.Name, err)
	}

	return fileBytes.Bytes(), nil
}
