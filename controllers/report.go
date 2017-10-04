package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/iotdog/daily-report/models"
	"github.com/iotdog/daily-report/utils"
	"github.com/leesper/holmes"
)

func SubmitReport(w http.ResponseWriter, r *http.Request) {
	input := models.SubmitReportParams{}
	holmes.Infoln("IP address: ", r.RemoteAddr)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		holmes.Errorln(err)
		utils.Jsonify(w, models.CommonResponse{
			Code: 1,
			Msg:  "读取上传数据失败",
		})
	} else {
		rsp := models.SubmitReport(input, r.RemoteAddr)
		utils.Jsonify(w, rsp)
	}
}

func GetDailyReport(w http.ResponseWriter, r *http.Request) {
	rsp := models.GetDailyReport()
	utils.Jsonify(w, rsp)
}

func SendDailyReport4Date(w http.ResponseWriter, r *http.Request) {
	input := models.SendDailyReportParams{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		holmes.Errorln(err)
		utils.Jsonify(w, models.CommonResponse{
			Code: 1,
			Msg:  "读取上传数据失败",
		})
	} else {
		rsp := models.SendReport4Date(input)
		utils.Jsonify(w, rsp)
	}
}
