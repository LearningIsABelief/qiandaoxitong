package util

import (
	"encoding/json"
	"fmt"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"qiandao/viewmodel"
)

// GetAddressByLngAndLat
// @Description: 根据经纬度获取地理位置
// @Author YangXuZheng
// @Date: 2022-06-14 8:56
func GetAddressByLngAndLat(Longitude, latitude string) string {
	// 经纬度
	lonAndLat := Longitude + "," + latitude
	addressResponse := viewmodel.AddressResponse{}
	// 签名
	sig := fmt.Sprintf("extensions=%s&"+
		"key=%s&"+
		"location=%s&"+
		"output=%s&"+
		"radius=%s&"+
		"roadlevel=%s"+
		"%s",
		viper.GetString("gaode.extensions"),
		viper.GetString("gaode.key"),
		lonAndLat,
		viper.GetString("gaode.output"),
		viper.GetString("gaode.radius"),
		viper.GetString("gaode.roadlevel"),
		viper.GetString("gaode.private_key"))

	// url
	getAddressURL := fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?"+
		"key=%s&"+
		"location=%s&"+
		"sig=%s&"+
		"radius=%s&"+
		"extensions=%s&"+
		"roadlevel=%s&"+
		"output=%s&",
		viper.GetString("gaode.key"),
		viper.GetString("gaode.location"),
		MD5(sig),
		viper.GetString("gaode.radius"),
		viper.GetString("gaode.extensions"),
		viper.GetString("gaode.roadlevel"),
		viper.GetString("gaode.output"))
	client := &http.Client{}
	// 调用接口
	reqest, err := http.NewRequest("GET", getAddressURL, nil)
	if err != nil {
		log.Errorf(err, "高德地图API接口调用失败")
		return ""
	}
	reqest.Header.Add("Content--Type", "application/json:charset=UTF-8")
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		log.Errorf(err, "处理返回结果失败")
		return ""
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf(err, "读取响应数据失败")
		return ""
	}
	fmt.Println(ByteSliceToString(body))
	err = json.Unmarshal(body, &addressResponse)
	if err != nil {
		log.Errorf(err, "json反序列化失败")
		return ""
	}
	// 返回 国家/省份/城市/区县/街道名称/街道号/具体位置
	// 例如：中国河南新乡红旗区新二街1625号大数据产业园
	if addressResponse.RegeoCode.AddressComponent.City == "" {
		return addressResponse.RegeoCode.AddressComponent.Country +
			addressResponse.RegeoCode.AddressComponent.Province +
			addressResponse.RegeoCode.AddressComponent.Township +
			addressResponse.RegeoCode.Pois[0].Address +
			addressResponse.RegeoCode.Pois[0].Name
	}
	return addressResponse.RegeoCode.AddressComponent.Country +
		addressResponse.RegeoCode.AddressComponent.Province +
		addressResponse.RegeoCode.AddressComponent.City +
		addressResponse.RegeoCode.AddressComponent.Township +
		addressResponse.RegeoCode.Pois[0].Address +
		addressResponse.RegeoCode.Pois[0].Name
}
