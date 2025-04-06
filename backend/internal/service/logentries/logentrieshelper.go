package logentries

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/FenixAra/go-util/log"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"

	"go-parquet-read/model"
)

func (le *LogEntries) ReadParquetWithPagination(parquetDirectory string, parquetFiles *[]string, limit, offset int, search string) ([]model.LogEntry, error) {
	var result []model.LogEntry
	totalRead := 0
	currentOffset := offset

	for _, fileName := range *parquetFiles {
		parquetFile := filepath.Join(parquetDirectory, fileName)
		le.l.Info("parquetFile : %v", parquetFile)
		fr, err := local.NewLocalFileReader(parquetFile)
		if err != nil {
			le.l.Errorf("ERROR: NewLocalFileReader: %v", err)
			return result, fmt.Errorf("NewLocalFileReader error: %w", err)
		}
		defer fr.Close()

		pr, err := reader.NewParquetReader(fr, new(model.LogEntry), 5)
		if err != nil {
			le.l.Errorf("ERROR: NewParquetReader: %v", err)
			return result, fmt.Errorf("NewParquetReader error: %w", err)
		}
		defer pr.ReadStop()

		numRows := int(pr.GetNumRows())

		if currentOffset >= numRows {
			le.l.Errorf("Offset %d exceeds numRows %d in file %s, skipping", currentOffset, numRows, fileName)
			currentOffset -= numRows
			continue
		}

		le.l.Info("currentOffset ", currentOffset)
		if err := pr.SkipRows(int64(currentOffset)); err != nil {
			le.l.Errorf("ERROR: SkipRows: %v", err)
			return result, fmt.Errorf("SkipRows error: %w", err)
		}

		rowsToRead := limit - totalRead
		if rowsToRead > (numRows - currentOffset) {
			rowsToRead = numRows - currentOffset
		}

		le.l.Info("Println ", rowsToRead)
		logEntrys := make([]model.LogEntry, rowsToRead)
		if err := pr.Read(&logEntrys); err != nil {
			le.l.Errorf("ERROR: Read: %v", err)
			return result, fmt.Errorf("read error: %v", err)
		}

		le.l.Info("Read %d rows from file %s", len(logEntrys), fileName)
		logEntrysFilter, err1 := le.AddFilter(search, logEntrys)
		if err1 != nil {
			le.l.Errorf("ERROR: addFilter: %v", err1)
			return result, fmt.Errorf("addFilter error: %v", err)
		}

		result = append(result, logEntrysFilter...)
		totalRead += len(logEntrys)
		currentOffset = 0 // Reset offset for following the files

		if totalRead >= limit {
			break
		}
	}

	le.l.Info("Total rows read: %d", totalRead)
	return result, nil
}

func (le *LogEntries) AddFilter(search string, logEntrys []model.LogEntry) ([]model.LogEntry, error) {

	if search == "" {
		return logEntrys, nil
	}

	var filteredEntries []model.LogEntry
	for _, logEntry := range logEntrys {
		if le.ContainsSearchString(logEntry, search) {
			filteredEntries = append(filteredEntries, logEntry)
		}
	}

	return filteredEntries, nil
}

func (le *LogEntries) ContainsSearchString(logEntry model.LogEntry, search string) bool {
	searchLower := strings.ToLower(search)

	if containsCaseInsensitive(logEntry.MsgId, searchLower) ||
		containsCaseInsensitive(logEntry.Timestamp, searchLower) ||
		containsCaseInsensitive(logEntry.Hostname, searchLower) ||
		containsCaseInsensitive(logEntry.AppName, searchLower) ||
		containsCaseInsensitive(logEntry.ProcId, searchLower) ||
		containsCaseInsensitive(logEntry.Message, searchLower) ||
		containsCaseInsensitive(logEntry.Tag, searchLower) ||
		containsCaseInsensitive(logEntry.Sender, searchLower) ||
		containsCaseInsensitive(logEntry.Groupings, searchLower) ||
		containsCaseInsensitive(logEntry.Event, searchLower) ||
		containsCaseInsensitive(logEntry.EventId, searchLower) ||
		containsCaseInsensitive(logEntry.EventId, searchLower) ||
		containsCaseInsensitive(logEntry.NanoTimeStamp, searchLower) ||
		containsCaseInsensitive(logEntry.SeverityString, searchLower) ||
		containsCaseInsensitive(logEntry.FacilityString, searchLower) ||
		containsCaseInsensitive(logEntry.MessageRaw, searchLower) ||
		containsCaseInsensitive(logEntry.StructuredData, searchLower) ||
		containsCaseInsensitive(logEntry.Namespace, searchLower) {
		return true
	}

	// Check if search is numeric and compare with int64 fields
	if isNumeric(searchLower) {
		searchInt64, err := strconv.ParseInt(searchLower, 10, 64)
		if err != nil {
			// Not a valid integer, skip
			return false
		}
		if logEntry.PartitionId == searchInt64 ||
			int64(logEntry.Priority) == searchInt64 ||
			int64(logEntry.Facility) == searchInt64 ||
			int64(logEntry.Severity) == searchInt64 {
			return true
		}
	}

	return false
}

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func containsCaseInsensitive(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}

func (le *LogEntries) GetParquetFlies() (*[]string, error) {

	res, errA := le.logEntriesDao.GetFilesFromDirectory()
	if errA != nil {
		le.l.Error("ERROR: getParquetFlies", errA)
		return nil, errA
	}
	return res, nil
}

func (le *LogEntries) CheckSizeImage(file *multipart.FileHeader, limit int64, ul *log.Logger) bool {
	size := file.Size / 1024
	ul.Info("file size kb - ", size)
	return size <= limit
}

func (le *LogEntries) CheckDirectory(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
