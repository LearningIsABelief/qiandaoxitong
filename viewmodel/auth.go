package viewmodel

// CodeInfoResponse 图片验证码响应体
type CodeInfoResponse struct {
	Uuid   string `json:"uuid"`
	Base64 string `json:"base64"`
}

// OnlineUserInfo 在线用户信息
type OnlineUserInfo struct {
	// 通过认证的账号的 id
	Id string `json:"id"`
	// 通过认证的帐号的账号
	UserAccount string `json:"user_account"`
	// 通过认证的账号的 token
	Token string `json:"token"`
	// 用户登录设备的 ip 地址
	Ip string `json:"ip"`
	// 用户登录所在的地址
	Address string `json:"address"`
	// 用户登录的时间 (时间戳形式)
	LoginTime int64 `json:"login_time"`
}
