package logentries

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/FenixAra/go-util/log"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"

	"go-parquet-read/internal/daos"
	"go-parquet-read/model"
)

type LogEntries struct {
	l                *log.Logger
	parquetDirectory string
	logEntriesDao    daos.LogEntriesDao
}

func New(l *log.Logger, parquetDirectory string) *LogEntries {
	return &LogEntries{
		l:                l,
		parquetDirectory: parquetDirectory,
		logEntriesDao:    daos.NewBranchObj(l, parquetDirectory),
	}
}

func (le *LogEntries) LogEntries(limit, offset, search string) (*model.LogEntryRes, error) {

	limitI, _ := strconv.ParseInt(limit, 10, 64)
	offsetI, _ := strconv.ParseInt(offset, 10, 64)
	if limitI == 0 {
		limitI = 20
	}
	if offsetI == 0 {
		offsetI = 0
	}

	parquetFlies, errF := le.GetParquetFlies()
	if errF != nil {
		le.l.Error("ERROR: getParquetFlies", errF)
		return nil, errF
	}

	le.l.Info("parquetFile here ", parquetFlies)
	result, err := le.ReadParquetWithPagination(le.parquetDirectory, parquetFlies, int(limitI), int(offsetI), search)
	if err != nil {
		le.l.Error("ERROR: readParquetWithPagination", err)
	}

	branchEntries := model.LogEntryRes{}
	branchEntries.LogEntries = &result
	branchEntries.Limit = limitI
	branchEntries.OffSet = offsetI
	return &branchEntries, nil
}

func (le *LogEntries) GetTotalRecords(search string) (*model.LogEntryResAll, error) {

	parquetFlies, errF := le.GetParquetFlies()
	if errF != nil {
		le.l.Error("ERROR: getParquetFlies", errF)
		return nil, errF
	}
	response := model.LogEntryResAll{}

	var result []model.LogEntry
	var totalRecods int

	for _, fileName := range *parquetFlies {
		parquetFile := filepath.Join(le.parquetDirectory, fileName)
		le.l.Info("parquetFile :", parquetFile)

		fr, err := local.NewLocalFileReader(parquetFile)
		if err != nil {
			le.l.Errorf("ERROR: NewLocalFileReader: %v", err)
			return &response, fmt.Errorf("NewLocalFileReader error: %w", err)
		}
		defer fr.Close()

		pr, err := reader.NewParquetReader(fr, new(model.LogEntry), 1)
		if err != nil {
			le.l.Errorf("ERROR: NewParquetReader: %v", err)
			return &response, fmt.Errorf("NewParquetReader error: %w", err)
		}
		defer pr.ReadStop()

		le.l.Info("NewParquetReader: ", pr.GetNumRows())

		logEntrys := make([]model.LogEntry, int(pr.GetNumRows()))
		if err := pr.Read(&logEntrys); err != nil {
			le.l.Error("ERROR: Read: ", err)
			return &response, fmt.Errorf("read error: %v", err)
		}

		le.l.Info("Found logEntrys: ", parquetFile, len(logEntrys))

		logEntrysFilter, err1 := le.AddFilter(search, logEntrys)
		if err1 != nil {
			le.l.Errorf("ERROR: addFilter: %v", err1)
			return &response, fmt.Errorf("addFilter error: %v", err)
		}

		result = append(result, logEntrysFilter...)
		totalRecods += len(logEntrysFilter)
	}

	response.LogEntries = &result
	response.Total = totalRecods

	return &response, nil
}

func (le *LogEntries) UplaodPerquetFile(imageFor string, file multipart.File, fileHeader *multipart.FileHeader) (*model.Message, error) {

	bool := le.CheckSizeImage(fileHeader, 30000, le.l) // 30 MB max
	if !bool {
		le.l.Error("image size issue ")
		return nil, errors.New("too large file (30 mb limit)")
	}

	isDirExists := le.CheckDirectory(le.parquetDirectory)
	if !isDirExists {
		err := os.MkdirAll(le.parquetDirectory, os.ModePerm) // os.ModePerm sets permissions to 0777
		if err != nil {
			le.l.Error("ERROR: Creating directory ", le.parquetDirectory, err)
			return nil, err
		}
	}
	imageName := fmt.Sprintf("%v_%v", imageFor, fileHeader.Filename)
	fileFullPath := filepath.Join(le.parquetDirectory, imageName)
	le.l.Info("fileFullPath: ", imageFor, fileFullPath)

	out, err := os.Create(fileFullPath)
	if err != nil {
		defer out.Close()
		le.l.Error("fileFullPath: ", imageFor, fileFullPath, err)
		return nil, err
	}
	defer out.Close()

	_, errU := io.Copy(out, file)
	if errU != nil {
		le.l.Error("file Copy error: ", fileFullPath, errU)
		return nil, errU
	}

	le.l.Info("Image uploaded successfully: ", fileFullPath)
	roleResponse := model.Message{}
	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v", fileFullPath)
	return &roleResponse, nil
}
