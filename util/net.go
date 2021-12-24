package util

import (
	"net/http"
	"strconv"
)

func GetUsername(re *http.Request) string {
	username := re.FormValue("id")
	if len(username) == 0 {
		username = re.FormValue("username")
	}
	return username
}
func GetStart(re *http.Request) int {
	value := re.FormValue("start")
	start, err := strconv.Atoi(value)
	if err == nil {
		return start
	}
	return 0
}
func GetSize(re *http.Request) int {
	value := re.FormValue("size")
	start, err := strconv.Atoi(value)
	if err == nil {
		return start
	}
	return 10
}
func GetMessage(re *http.Request) string {
	msg := re.FormValue("msg")
	if len(msg) == 0 {
		msg = re.FormValue("message")
	}
	return msg
}

func HttpCross(w http.ResponseWriter) {
	h := w.Header()
	h.Add("Access-Control-Allow-Origin", "*")
	h.Add("Content-Type", "text/html; charset=utf-8")
}

func HttpCrossChunked(w http.ResponseWriter) {
	h := w.Header()
	HttpCross(w)
	h.Add("Transfer-Encoding", "chunked")
}
