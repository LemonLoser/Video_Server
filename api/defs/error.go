package defs

var (
	ErrorRequsetBodyParaseFailed = ErrorResponse{400, Err{Error: "request body is not correct",
		ErrorCode: "001"}}
	ErrorNotAuthUser     = ErrorResponse{401, Err{Error: "User authentication failed", ErrorCode: "002"}}
	ErrorDBError         = ErrorResponse{HttpSC: 500, Error: Err{Error: "DB ops failed", ErrorCode: "003"}}
	ErrorInternalFaults  = ErrorResponse{HttpSC: 500, Error: Err{Error: "Internal service error", ErrorCode: "004"}}
	ErrorUserLoginFailed = ErrorResponse{HttpSC: 404, Error: Err{Error: "Userlogin error", ErrorCode: "005"}}
)

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErrorResponse struct {
	HttpSC int
	Error  Err
}
