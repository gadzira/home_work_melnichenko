package copy

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const (
	rootDir string = "/tmp/"
	whence  int    = 0
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func checkRegular(fromPath string) error {
	rf, err := os.Lstat(fromPath)
	if err != nil {
		log.Println(err)
	}

	if !rf.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Step 0: check is regular file
	err := checkRegular(fromPath)
	if err != nil {
		return err
	}

	// Step 1: open file
	file, err := os.OpenFile(fromPath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Step 1.1: define file size
	fs, err := file.Stat()
	if err != nil {
		log.Println(err)
	}

	// If offset more than file size return error
	if offset > fs.Size() {
		return ErrOffsetExceedsFileSize
	}

	// If limit more than file size, set limit as equal fs
	if limit > fs.Size() || limit == 0 {
		limit = fs.Size()
	}

	// Offset: zero by default, but we always make offset
	_, err = file.Seek(offset, whence)
	if err != nil {
		return err
	}

	buf := make([]byte, fs.Size()-offset)
	if _, err := io.ReadFull(file, buf); err != nil {
		return err
	}

	// Step 2: Open or create file
	out, err := os.Create(rootDir + toPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Step 3: Copy source file to destination file
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(bytes.NewReader(buf))
	_, err = io.CopyN(out, barReader, limit)
	if err != nil {
		return err
	}
	bar.Finish()

	return nil
}
