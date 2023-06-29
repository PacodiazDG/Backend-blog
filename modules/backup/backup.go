package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/PacodiazDG/Backend-blog/modules/logs"
)

func GenerateBackup(folder string) error {
	t := time.Now()
	fb, err := os.Create("Backup" + t.Format("20060102150405") + ".zip")
	if err != nil {
		logs.WriteLogs(err, logs.CriticalError)
		return err
	}
	defer fb.Close()
	w := zip.NewWriter(fb)
	defer w.Close()
	WalkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		f, err := w.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		return file.Sync()
	}
	err = filepath.Walk(folder, WalkFunc)
	if err != nil {
		logs.WriteLogs(err, logs.CriticalError)
		return err
	}
	return fb.Close()
}
