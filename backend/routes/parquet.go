package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-parquet-read/internal/service/logentries"
)

func parquetReader(router *httprouter.Router, recoverHandler alice.Chain) {
	router.GET("/v1/list/log/entries", wrapHandler(recoverHandler.ThenFunc(listLogEntries)))
	router.GET("/v1/fetch/all/records", wrapHandler(recoverHandler.ThenFunc(getTotalRecords)))
	router.POST("/v1/upload/perquet", wrapHandler(recoverHandler.ThenFunc(uplaodPerquet)))
}

func listLogEntries(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	search := keys.Get("search")
	rd.l.Info("listLogEntries: ", "limit", limit, "offset", offset)
	loge := logentries.New(rd.l, rd.parquetDirectory)
	res, err := loge.LogEntries(limit, offset, search)
	if err != nil {
		rd.l.Error("listLogEntries error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	duration := time.Since(rd.Start)
	res.Duration = fmt.Sprintf("%v ms", duration.Milliseconds())
	writeJSONStruct(res, http.StatusOK, rd)
}

func getTotalRecords(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	loge := logentries.New(rd.l, rd.parquetDirectory)
	keys := r.URL.Query()
	search := keys.Get("search")
	rd.l.Info("listLogEntries: ", "search", search)

	res, err := loge.GetTotalRecords(search)
	if err != nil {
		rd.l.Error("totalRecords error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	duration := time.Since(rd.Start)
	res.Duration = fmt.Sprintf("%v ms", duration.Milliseconds())
	writeJSONStruct(res, http.StatusOK, rd)
}

func uplaodPerquet(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	loge := logentries.New(rd.l, rd.parquetDirectory)
	keys := r.URL.Query()
	imageFor := keys.Get("imageFor")
	// Get the file from the request
	file, header, err := r.FormFile("image")
	if err != nil {
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	res, err := loge.UplaodPerquetFile(imageFor, file, header)
	if err != nil {
		rd.l.Error("uplaodPerquet error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
