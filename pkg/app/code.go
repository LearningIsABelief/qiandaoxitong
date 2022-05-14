package app

var (
	// 服务器状态码
	OK                  = &Errno{Code: 200, Message: "OK"}
	InternalServerError = &Errno{Code: 500, Message: "服务器异常."}
	ErrBind             = &Errno{Code: 10002, Message: "将请求正文绑定到结构时发生错误."}
	ErrParamNull        = &Errno{Code: 11001, Message: "参数结果为空"}

	// 课程错误
	ErrLessonExist    = &Errno{Code: 11004, Message: "课程已存在"}
	ErrRecordNotExist = &Errno{Code: 11005, Message: "查询记录不存在"}
	ErrUpdated        = &Errno{Code: 11006, Message: "更新失败"}
	ErrInserted       = &Errno{Code: 11007, Message: "插入失败"}
	ErrDeleted        = &Errno{Code: 11008, Message: "删除失败"}

	// 用户错误
	ErrEncrypt             = &Errno{Code: 20101, Message: "加密用户密码时出错"}
	ErrPhoneExist          = &Errno{Code: 20102, Message: "手机号已被注册"}
	ErrEmailExist          = &Errno{Code: 20103, Message: "该邮箱已被绑定"}
	ErrNickNameExist       = &Errno{Code: 20104, Message: "该昵称已被占用"}
	ErrPassword            = &Errno{Code: 20105, Message: "原密码不正确"}
	ErrOldNewInconsistent  = &Errno{Code: 20106, Message: "请确保两次输入的密码一样"}
	ErrPhoneDoesNotExist   = &Errno{Code: 20107, Message: "手机号不存在"}
	ErrPhoneBinEmail       = &Errno{Code: 20108, Message: "请输入手机号绑定的正确邮箱"}
	ErrAccountDoesNotExist = &Errno{Code: 20109, Message: "账号不存在"}
	ErrLoginPassword       = &Errno{Code: 20110, Message: "密码错误"}
	ErrUserNotExist        = &Errno{Code: 11002, Message: "用户不存在"}

	// token
	ErrTokenInvalid = &Errno{Code: 20103, Message: "无效token"}

	// 班级错误
	ErrClassExist    = &Errno{Code: 20201, Message: "班级已存在，不能重复创建"}
	ErrClassNotExist = &Errno{Code: 11003, Message: "班级不存在"}
)
