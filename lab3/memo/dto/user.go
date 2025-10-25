package dto

/* 注册请求体 */
type SignupReq struct {
	Username string `form:"username,required"`
	Password string `form:"password,required"`
}

/* 登录请求体 */
type LoginReq struct {
	Username string `form:"username,required"`
	Password string `form:"password,required"`
}

/* 登录返回的token数据 */
type TokenData struct {
	Token string
}
