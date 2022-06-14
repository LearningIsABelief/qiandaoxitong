package util

/**
对在线用户管理的工具类
对在线用户管理的工具类
对在线用户管理的工具类
*/
import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"qiandao/store"
	"qiandao/viewmodel"
	"strings"
	"sync"
)

// GetRequestIP
// @Description: 获取用户ip地址
// @Author YangXuZheng
// @Date: 2022-06-11 22:06
func GetRequestIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	//c.Request.Header.Get("X-Forward-For")
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

// GetAllOnlineUser
// @Description: 获取所有在线的用户
// @Author YangXuZheng
// @Date: 2022-06-12 21:03
func GetAllOnlineUser(userId string) ([]viewmodel.OnlineUserInfo, error) {
	var allUser []viewmodel.OnlineUserInfo
	// 获取到匹配的所有在线用户的key
	keys, err := RedisScanKeys("online-token*")
	if err != nil {
		return []viewmodel.OnlineUserInfo{}, nil
	}
	var dateMap viewmodel.OnlineUserInfo
	// 然后将这些key放入
	for _, v := range keys {
		userInfo, _ := RedisGet(v)
		bytes := StringToByteSlice(userInfo)
		err := json.Unmarshal(bytes, &dateMap)
		if err != nil {
			log.Errorf(err, "redis中字符串序列化为结构体失败")
			return nil, errors.New("redis中字符串序列化为结构体失败")
		}
		allUser = append(allUser, dateMap)
	}
	return allUser, nil
}

// CheckLoginOnUser
// @Description: 检测用户是否在之前已经登录，已经登录踢下线
// @Author YangXuZheng
// @@Date: 2022-06-12 21:03
func CheckLoginOnUser(userId, ignoreToken string) error {
	// 查找到所有在线的用户，根据当前通过所有验证的账号的id去遍历查找所有在线的用户中跟这个ID相等的 并且 token解密出来不能和当前账号的token相等的
	allUser, err := GetAllOnlineUser(userId)
	if err != nil {
		return err
	}
	if len(allUser) == 0 {
		return nil
	}
	// 使用go并发踢出当前已经在线的用户
	wg := sync.WaitGroup{}
	for _, v := range allUser {
		if strings.Compare(v.Id, userId) == 0 {
			wg.Add(1)
			go func(userInfo *viewmodel.OnlineUserInfo) {
				defer wg.Done()
				// 将之前加密F过的token转换成解密所需要的类型
				decodeString, _ := base64.StdEncoding.DecodeString(v.Token)
				// 将加密过的token解密后转换为string类型
				decryptToken := ByteSliceToString(DesDecrypt(decodeString))
				if strings.Compare(decryptToken, ignoreToken) != 0 {
					_ = kickOut(decryptToken)
				}
			}(&v)
		} else {
			continue
		}
	}
	go func() {
		wg.Wait()
	}()
	return nil
}

// KickOut
// @Description: 踢出用户
// @Author YangXuZheng
// @Date: 2022-06-12 21:03
func kickOut(token string) error {
	key := viper.GetString("jwt.online_key") + token
	if err := store.RedisDB.Self.Del(key).Err(); err != nil {
		return err
	}
	return nil
}
