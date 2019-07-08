package streamserver

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./videos/upload.html")
	t.Execute(w, nil)
}
func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid //链接
	video, err := os.Open(vl)
	if err != nil {
		log.Printf("error when try to open the file;%v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "internet error")
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "file is too big")
		return
	}
	//读取文件
	file, _, err := r.FormFile("file") //<form name = "file"拿到file句柄和file文件
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}
	data, err := ioutil.ReadAll(file) //二进制流
	if err != nil {
		log.Printf("read file error:%v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "internal error")
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+fn, data, 0666)
	if err != nil {
		log.Printf("write file error:%v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "write file error")
		return
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "uploaded video successfully")
}
