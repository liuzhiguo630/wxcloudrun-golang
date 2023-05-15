package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)
	http.HandleFunc("/wx/notify", service.WxNotifyHandler)

	http.HandleFunc("/getQrCode/freesite", func(w http.ResponseWriter, r *http.Request) {
		// request to https://api.weixin.qq.com/cgi-bin/qrcode/create
		// with body {"action_name":"QR_LIMIT_STR_SCENE","action_info":{"scene":{"from":"freesite","randomId":"uiasdufuadsf"}}}
		// get Id from request params
		queryParams, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		rid := queryParams.Get("rid")

		var requestBody = []byte(`{"expire_seconds": 21600, "action_name": "QR_SCENE", "action_info": {"scene": {"from":"freesite","randomId":"` + rid + `"}}}`)
		resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/qrcode/create", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(all)
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
