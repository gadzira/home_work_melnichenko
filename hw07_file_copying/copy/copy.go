package copy

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func checkRegular(fromPath string) error {
	rf, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("failed to get file stat: %w", err)
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
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Step 1.1: define file size
	fs, err := file.Stat()
	fmt.Println("File size:", fs.Size())
	fmt.Println("Copy limit:", limit)
	if err != nil {
		return fmt.Errorf("failed to define size: %w", err)
	}

	// If offset more than file size return error
	if offset > fs.Size() {
		return ErrOffsetExceedsFileSize
	}

	// If limit more than file size, set limit as equal fs
	if limit > fs.Size()-offset || limit == 0 {
		limit = fs.Size() - offset
	}

	// Offset: zero by default, but we always make offset
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to make offset: %w", err)
	}

	// Step 2: Open or create file
	out, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Step 3: Copy source file to destination file
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(file)
	_, err = io.CopyN(out, barReader, limit)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	defer bar.Finish()

	return nil
}
