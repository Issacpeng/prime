package handler

import (
	"encoding/json"
	_"fmt"
	"net/http"
	_"regexp"

	"github.com/astaxie/beego/logs"
	"gopkg.in/macaron.v1"

	_"github.com/prime/models"
	_"github.com/wrench/db"
	_"github.com/wrench/setting"
	_"github.com/wrench/utils"
)

func ListPullRequestHandler(ctx *macaron.Context, log *logs.BeeLogger) (int, []byte) {


	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

func GetSinglePullRequestHandler(ctx *macaron.Context, log *logs.BeeLogger) (int, []byte) {


	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

func CreateNewPullHandler(ctx *macaron.Context, log *logs.BeeLogger) (int, []byte) {


	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

func UpdateSinglePullHandler(ctx *macaron.Context, log *logs.BeeLogger) (int, []byte) {


	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

func ListCommitsOnPullHandler(ctx *macaron.Context, log *logs.BeeLogger) (int, []byte) {


	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

func ListFilesOnPullHandler(ctx *macaron.Context, log *logs.BeeLogger) (int, []byte) {


	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

func IfPullMergedHandler(ctx *macaron.Context, log *logs.BeeLogger) (int, []byte) {


	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

func MergeNewPullHandler(ctx *macaron.Context, log *logs.BeeLogger) (int, []byte) {


	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}