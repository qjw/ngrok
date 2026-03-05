// interactive web user interface
package web

import (
	"io/fs"
	"net/http"

	"github.com/gorilla/websocket"
	assets "github.com/qjw/ngrok/assets/client"
	"github.com/qjw/ngrok/ngrok/client/mvc"
	"github.com/qjw/ngrok/ngrok/log"
	"github.com/qjw/ngrok/ngrok/proto"
	"github.com/qjw/ngrok/ngrok/util"
)

type WebView struct {
	log.Logger

	ctl mvc.Controller

	// messages sent over this broadcast are sent to all websocket connections
	wsMessages *util.Broadcast
}

func NewWebView(ctl mvc.Controller, addr string) *WebView {
	wv := &WebView{
		Logger:     log.NewPrefixLogger("view", "web"),
		wsMessages: util.NewBroadcast(),
		ctl:        ctl,
	}

	// for now, always redirect to the http view
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/http/in", 302)
	})

	// handle web socket connections
	http.HandleFunc("/_ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)

		if err != nil {
			http.Error(w, "Failed websocket upgrade", 400)
			wv.Warn("Failed websocket upgrade: %v", err)
			return
		}

		msgs := wv.wsMessages.Reg()
		defer wv.wsMessages.UnReg(msgs)
		for m := range msgs {
			err := conn.WriteMessage(websocket.TextMessage, m.([]byte))
			if err != nil {
				// connection is closed
				break
			}
		}
	})

	// 从 embed.FS 中创建子文件系统
	// staticFiles 中包含 "static/" 前缀，我们需要去掉这个前缀
	staticFS, err := fs.Sub(assets.StaticFiles, "static")
	if err != nil {
		panic(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	wv.Info("Serving web interface on %s", addr)
	wv.ctl.Go(func() { http.ListenAndServe(addr, nil) })
	return wv
}

func (wv *WebView) NewHttpView(proto *proto.Http) *WebHttpView {
	return newWebHttpView(wv.ctl, wv, proto)
}

func (wv *WebView) Shutdown() {
}
