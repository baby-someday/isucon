package build

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func Compress(src, dst string, ignore []string) error {
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	zipWriter := zip.NewWriter(dstFile)
	defer zipWriter.Close()

	err = archive(src, src, ignore, zipWriter)
	if err != nil {
		return err
	}

	return nil
}

func archive(dir, p string, ignore []string, zipWriter *zip.Writer) error {
	// TODO: 変数名ちゃんとつける
	fileInfos, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		fp := path.Join(p, fileInfo.Name())
		name := strings.TrimPrefix(fp, dir)

		var shouldFileBeIgnored = false
		for _, i := range ignore {
			if strings.TrimPrefix(name, "/") == i {
				shouldFileBeIgnored = true
				break
			}
		}
		if shouldFileBeIgnored {
			continue
		}

		if fileInfo.IsDir() {
			err = archive(dir, fp, ignore, zipWriter)
			if err != nil {
				return err
			}
			continue
		}

		if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			header, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				return err
			}

			header.Method = zip.Deflate
			srcFileWriter, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			_, err = srcFileWriter.Write([]byte(fp))
			if err != nil {
				return err
			}
		} else {
			srcFile, err := os.Open(fp)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			srcFileWriter, err := zipWriter.Create(name)
			if err != nil {
				return err
			}

			_, err = io.Copy(srcFileWriter, srcFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
