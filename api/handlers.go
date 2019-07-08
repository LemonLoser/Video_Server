package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"videoProject/api/dbops"
	"videoProject/api/defs"
	"videoProject/api/session"
	"videoProject/api/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//从request中读出create user相关的信息
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		SendErrorResponse(w, defs.ErrorRequsetBodyParaseFailed)
		return
	}
	if err := dbops.AddUserCredential(ubody.UserName, ubody.Password); err != nil {
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(ubody.UserName)
	su := &defs.SignedUp{Success: true, SessionId: id}

	if resp, err := json.Marshal(su); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		SendNormalResponse(w, string(resp), 201)
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s", res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		//io.WriteString(w,"wrong")
		SendErrorResponse(w, defs.ErrorRequsetBodyParaseFailed)
		return
	}
	//validate the request body
	uname := p.ByName("username")
	log.Printf("login url name:%s", uname)
	log.Printf("login body name:%s", ubody.UserName)
	if uname != ubody.UserName {
		SendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	log.Printf("%s", ubody.UserName)
	//从数据库中查询ubody.UserName所对应的pwd和ubody.Password做对比判断密码是否正确
	pwd, err := dbops.GetUserCredential(ubody.UserName)
	log.Printf("login pwd:%s", pwd)
	log.Printf("login body pwd:%s", ubody.Password)
	if pwd != ubody.Password || len(pwd) == 0 || err != nil {
		SendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	id := session.GenerateNewSessionId(ubody.UserName)
	si := &defs.SignedIn{Success: true, SessionId: id}
	if resp, err := json.Marshal(si); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 200)
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Printf("Unauthorized user\n")
		return
	}

	uname := p.ByName("username")
	u, err := dbops.GetUser(uname)
	if err != nil {
		log.Printf("error in GetUsrInfo:%s", err)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}
	ui := &defs.UserInfo{Id: u.Id}
	if resp, err := json.Marshal(ui); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 200)
	}
}

func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Printf("unauthorized usr\n")
		return
	}
	res, _ := ioutil.ReadAll(r.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal(res, nvbody); err != nil {
		log.Printf("%s", err)
		SendErrorResponse(w, defs.ErrorRequsetBodyParaseFailed)
		return
	}
	vi, err := dbops.AddNewVideo(nvbody.AuthorId, nvbody.Name)
	if err != nil {
		log.Printf("Author id:%v,name:%v", nvbody.AuthorId, nvbody.Name)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}
	if resp, err := json.Marshal(vi); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 201)
	}
}

func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}
	uname := p.ByName("username")
	vs, err := dbops.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("error in listAllVideo:%s", err)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}
	vsi := &defs.VideosInfo{Videos: vs}
	if resp, err := json.Marshal(vsi); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 200)
	}
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}
	vid := p.ByName("vid-id")
	if err := dbops.DeleteVideoInfo(vid); err != nil {
		log.Printf("error in DeleteVideo:%s", err)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}
	go utils.SendDeleteVideoRequest(vid)
	SendNormalResponse(w, "", 204)
}

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}
	reqBody, _ := ioutil.ReadAll(r.Body)

	cbody := &defs.NewComment{}
	if err := json.Unmarshal(reqBody, cbody); err != nil {
		log.Printf("%s", err)
		SendErrorResponse(w, defs.ErrorRequsetBodyParaseFailed)
		return
	}
	vid := p.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cbody.AuthorId, cbody.Content); err != nil {
		log.Printf("error in PostComment:%s", err)
		SendErrorResponse(w, defs.ErrorDBError)
	} else {
		SendNormalResponse(w, "ok", 201)
	}
}

func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}
	vid := p.ByName("vid-id")
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("error in ShowComments:%s", err)
		SendErrorResponse(w, defs.ErrorDBError)
		return
	}
	cms := &defs.Comments{Comments: cm}
	if resp, err := json.Marshal(cms); err != nil {
		SendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		SendNormalResponse(w, string(resp), 200)
	}
}

//pagination
