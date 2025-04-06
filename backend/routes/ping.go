package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-parquet-read/internal/service/ping"
	"go-parquet-read/model"
)

func setPingRoutes(router *httprouter.Router, recoverHandler alice.Chain) {
	router.GET("/ping", wrapHandler(recoverHandler.ThenFunc(Ping)))
}

var res model.ResStruct

// Ping godoc
// @Summary ping api
// @Description do ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} dtos.ResStruct
// @Failure 500 {object} dtos.Res500Struct
// @Router /ping [get]

func Ping(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	p := ping.New(rd.l, rd.parquetDirectory)
	err := p.Ping()
	if err != nil {
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusInternalServerError, rd)
	} else {
		writeJSONMessage("pong", MSG, http.StatusOK, rd)
	}
}
