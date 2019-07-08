package main

import (
	"net/http"
	"videoProject/api/defs"
	"videoProject/api/session"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_UNAME = "X-User-Name"

func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}
	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}
	//如果没有过期就把用户名加入到Head_Field_Uname中
	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

//简单的用户校验
func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_UNAME)
	if len(uname) == 0 {
		SendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}
	return true
}

//IAM
//SSO
//Rbac, role based access control
