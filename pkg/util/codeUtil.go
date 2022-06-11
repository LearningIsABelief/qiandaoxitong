package util

import (
	"github.com/mojocn/base64Captcha"
	"github.com/spf13/viper"
	"image/color"
)

var Result = base64Captcha.DefaultMemStore

// CreateCode
// @Description: 生成图片验证码的 base64编码和ID
// @Author YangXuZheng 2022-06-11 13:21
// @Result id 验证码id
// @Result base64 图片base64编码
// @Result error 错误
func CreateCode() (string, string, error) {
	var driver base64Captcha.Driver
	switch viper.GetString("code.captcha_type") {
	case "audio":
		driver = autoConfig()
	case "string":
		driver = stringConfig()
	case "math":
		driver = mathConfig()
	case "chinese":
		driver = chineseConfig()
	case "digit":
		driver = digitConfig()
	}
	if driver == nil {
		panic("生成验证码的类型没有配置，请在yaml文件中配置完再次重试启动项目")
	}
	c := base64Captcha.NewCaptcha(driver, Result)
	id, b64s, err := c.Generate()
	return id, b64s, err
}

// VerifyCaptcha
// @Description: 校验验证码
// @Author YangXuZheng 2022-06-11 13:21
// @Pram id 验证码id
// @Pram VerifyValue 答案
// @Result true：正确，false：失败
func VerifyCaptcha(id, VerifyValue string) bool {
	return Result.Verify(id, VerifyValue, true)
}

// GetCodeAnswer
// @Description: 获取验证码答案
// @Author YangXuZheng 2022-06-11 13:31
// @Pram codeId 验证码id
// @Result 验证码答案
func GetCodeAnswer(codeId string) string {
	return Result.Get(codeId, false)
}

// mathConfig 生成图形化算术验证码配置
func mathConfig() *base64Captcha.DriverMath {
	mathType := &base64Captcha.DriverMath{
		Height:          viper.GetInt("code.height"),
		Width:           viper.GetInt("code.width"),
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine,
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return mathType
}

// digitConfig 生成图形化数字验证码配置
func digitConfig() *base64Captcha.DriverDigit {
	digitType := &base64Captcha.DriverDigit{
		Height:   viper.GetInt("code.height"),
		Width:    viper.GetInt("code.width"),
		Length:   5,
		MaxSkew:  0.45,
		DotCount: 80,
	}
	return digitType
}

// stringConfig 生成图形化字符串验证码配置
func stringConfig() *base64Captcha.DriverString {
	stringType := &base64Captcha.DriverString{
		Height:          viper.GetInt("code.height"),
		Width:           viper.GetInt("code.width"),
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:          5,
		Source:          viper.GetString("code.string.source"),
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return stringType
}

// chineseConfig 生成图形化汉字验证码配置
func chineseConfig() *base64Captcha.DriverChinese {
	chineseType := &base64Captcha.DriverChinese{
		Height:          viper.GetInt("code.height"),
		Width:           viper.GetInt("code.width"),
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowSlimeLine,
		Length:          2,
		Source:          viper.GetString("code.chinese.source"),
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return chineseType
}

// autoConfig 生成图形化数字音频验证码配置
func autoConfig() *base64Captcha.DriverAudio {
	chineseType := &base64Captcha.DriverAudio{
		Length:   viper.GetInt("code.auto.autoLength"),
		Language: viper.GetString("code.auto.language"),
	}
	return chineseType
}
