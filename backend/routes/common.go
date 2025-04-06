package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/FenixAra/go-util/log"
)

const (
	ERR_MSG = "ERROR_MESSAGE"
	MSG     = "MESSAGE"
)

type ResStruct struct {
	Status   string `json:"status" example:"SUCCESS" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"200" example:"500"`
	Message  string `json:"message" example:"pong" example:"could not connect to db"`
}

type Res500Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"500"`
	Message  string `json:"message" example:"could not connect to db"`
}

type Res400Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"400"`
	Message  string `json:"message" example:"Invalid param"`
}

type RequestData struct {
	l                *log.Logger
	Start            time.Time
	w                http.ResponseWriter
	r                *http.Request
	parquetDirectory string
}

type RenderData struct {
	Data  interface{}
	Paths []string
}

type TemplateData struct {
	Data interface{}
}

func (t *TemplateData) SetConstants() {

}

func logAndGetContext(w http.ResponseWriter, r *http.Request) *RequestData {
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("X-Frame-Options", "DENY")

	//url := strings.TrimSpace("http://0.0.0.0:9008/v1/logs")

	//Set config according to the use case
	cfg := log.NewConfig("parquet")
	//cfg.SetRemoteConfig(url, "", "admin")
	cfg.SetLevelStr("Debug")
	cfg.SetFilePathSizeStr("Full")
	cfg.SetReference(r.Header.Get("ReferenceID"))

	l := log.New(cfg)
	parquetDirectory := os.Getenv("PARQUET_FILE_DIRECTORT")

	start := time.Now()
	l.LogAPIInfo(r, 0, 0)

	return &RequestData{
		l:                l,
		Start:            start,
		r:                r,
		w:                w,
		parquetDirectory: parquetDirectory,
	}
}

func jsonifyMessage(msg string, msgType string, httpCode int) ([]byte, int) {
	var data []byte
	var Obj struct {
		Status   string `json:"status"`
		HTTPCode int    `json:"code"`
		Message  string `json:"message"`
		Err      error  `json:"error"`
	}
	Obj.Message = msg
	Obj.HTTPCode = httpCode
	switch msgType {
	case ERR_MSG:
		Obj.Status = "FAILED"

	case MSG:
		Obj.Status = "SUCCESS"
	}
	data, _ = json.Marshal(Obj)
	return data, httpCode
}

func writeJSONMessage(msg string, msgType string, httpCode int, rd *RequestData) {
	d, code := jsonifyMessage(msg, msgType, httpCode)
	writeJSONResponse(d, code, rd)
}

func writeJSONStruct(v interface{}, code int, rd *RequestData) {
	d, err := json.Marshal(v)
	if err != nil {
		writeJSONMessage("Unable to marshal data. Err: "+err.Error(), ERR_MSG, http.StatusInternalServerError, rd)
		return
	}
	writeJSONResponse(d, code, rd)
}

func writeJSONResponse(d []byte, code int, rd *RequestData) {
	rd.l.LogAPIInfo(rd.r, time.Since(rd.Start).Seconds(), code)
	if code == http.StatusInternalServerError {
		rd.l.Info(rd.r.URL, "Status Code:", code, ", Response time:", time.Since(rd.Start), rd.r.URL, " Response:", string(d))
	} else {
		rd.l.Info(rd.r.URL, "Status Code:", code, ", Response time:", time.Since(rd.Start))
	}
	rd.w.Header().Set("Access-Control-Allow-Origin", "*")
	rd.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rd.w.WriteHeader(code)
	rd.w.Write(d)
}
