package daos

import (
	"errors"
	"os"

	"github.com/FenixAra/go-util/log"

	"go-parquet-read/model"
)

type LogEntries struct {
	l                *log.Logger
	parquetDirectory string
}

func NewBranchObj(l *log.Logger, parquetDirectory string) *LogEntries {
	return &LogEntries{
		l:                l,
		parquetDirectory: parquetDirectory,
	}
}

type LogEntriesDao interface {
	GetLogEntries(limit, offset int64) (*[]model.LogEntry, error)
	GetFilesFromDirectory() (*[]string, error)
}

func (le *LogEntries) GetLogEntries(limit, offset int64) (*[]model.LogEntry, error) {

	return nil, nil

}

func (le *LogEntries) GetFilesFromDirectory() (*[]string, error) {

	dir, err := os.Open(le.parquetDirectory)
	if err != nil {
		le.l.Error("Error opening directory: %v\n", err)
		return nil, err
	}
	defer dir.Close()

	// Open the directory
	// Read the contents of the directory
	files, err := dir.Readdir(-1) // -1 means read all files
	if err != nil {
		le.l.Error("Error reading directory contents: %v\n", err)
		return nil, err
	}

	// Check if the directory contains any files
	if len(files) == 0 {
		le.l.Error("No files found in the directory.")
		return nil, errors.New("no files found in the directory")
	}
	var filesArray []string
	for _, file := range files {
		if !file.IsDir() && file.Name()[0] != '.' { // Ensure it's a file, not a subdirectory & no hidden
			filesArray = append(filesArray, file.Name())
		}
	}
	return &filesArray, nil
}
