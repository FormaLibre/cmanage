package utils

import (
  "fmt"
  "os"
  "io"
  "compress/gzip"
  "path/filepath"
  "io/ioutil"

  jww "github.com/spf13/jwalterweatherman"
)

// DiskUsage Caluculated used disk space for a given folder
func DiskUsage(currentPath string, info os.FileInfo) int64 {
  size := info.Size()
  if !info.IsDir() {
    return size
  }
  dir, err := os.Open(currentPath)
  if err != nil {
    fmt.Println(err)
    return size
  }
  defer dir.Close()
  files, err := dir.Readdir(-1)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  for _, file := range files {
    if file.Name() == "." || file.Name() == ".." {
            continue
    }
    size += DiskUsage(currentPath+"/"+file.Name(), file)
  }
  return size
}

// Exists function checks if a path exists on the filesystem
func Exists(path string) (bool, error) {
  _, err := os.Stat(path)
  if err == nil { return true, nil }
  if os.IsNotExist(err) { return false, nil }
  return true, err
}

// NotExists function checks if a path is non existant on the filesystem
func NotExists(path string) (bool, error) {
  _, err := os.Stat(path)
  if err == nil { return false, nil }
  if os.IsNotExist(err) { return true, nil }
  return false, err
}

// Check checks for errors
func Check(e error) {
  if e != nil {
    jww.ERROR.Println(e)
  }
}

// Ungzip unzips target archive
func Ungzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()
	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()
	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = io.Copy(writer, archive)
	return err
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
