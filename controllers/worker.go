package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/iotdog/daily-report/models"
	"github.com/iotdog/daily-report/utils"
)

func AddWorker(w http.ResponseWriter, r *http.Request) {
	input := models.AddWorkerParams{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.Jsonify(w, models.CommonResponse{
			Code: 1,
			Msg:  "读取上传数据失败",
		})
	} else {
		rsp := models.AddWorker(input)
		utils.Jsonify(w, rsp)
	}
}
