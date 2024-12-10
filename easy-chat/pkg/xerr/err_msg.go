package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "server exception, please try again later",
	REQUEST_PARAM_ERROR: "request parameter error",
	DB_ERROR:            "database error",
}

func ErrMsg(errcode int) string {
	if msg, ok := codeText[errcode]; ok {
		return msg
	}
	return codeText[SERVER_COMMON_ERROR]
}
