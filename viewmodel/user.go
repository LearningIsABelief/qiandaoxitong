package viewmodel

// 用户模块的 request 和 response
// 用户模块的 request 和 response
// 用户模块的 request 和 response

// RegisterRequest 用户注册接口 接收参数
type RegisterRequest struct {
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Role     int    `json:"role" form:"role" binding:"required"`
	ClassId  string `json:"class_id" form:"class_id" binding:"required"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Token     string `json:"token"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
	ClassId   string `json:"class_id"`
	ClassName string `json:"class_name"`
	RealName  string `json:"real_name"`
}

// UpdateUserInfoRequest 修改用户信息请求结构体
type UpdateUserInfoRequest struct {
	UserId   string `json:"user_id" form:"user_id"`
	Email    string `json:"email" form:"email"`
	RealName string `json:"real_name" form:"real_name"`
	Hobby    string `json:"hobby" form:"hobby"`
	Address  string `json:"address" form:"address"`
	Sex      int    `json:"sex" form:"sex"`
	Age      int    `json:"age" form:"age"`
}

// UpdateEmailRequest 修改邮箱请求结构体
type UpdateEmailRequest struct {
	UserId string `json:"user_id" form:"user_id"`
	Email  string `json:"email" form:"email"`
}

// UpdateNickNameRequest 修改昵称请求结构体
type UpdateNickNameRequest struct {
	UserId   string `json:"user_id" form:"user_id"`
	NickName string `json:"nick_name" form:"nick_name"`
}
