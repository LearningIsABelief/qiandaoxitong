package viewmodel

// AddressResponse 高德地图接口响应结果集
type AddressResponse struct {
	Status    string    `json:"status"`
	RegeoCode Regeocode `json:"regeocode"`
	Info      string    `json:"info"`
	InfoCode  string    `json:"infocode"`
}

// Regeocode 逆地理编码列表
type Regeocode struct {
	AddressComponent AddressComponent `json:"addressComponent"`
	Pois             []Pois           `json:"pois"`
}

// AddressComponent 地址元素列表
type AddressComponent struct {
	// 国家
	Country string `json:"country"`
	// 省份
	Province string `json:"province"`
	// 城市
	City string `json:"city"`
	// 区/县
	District string `json:"district"`
	// 乡镇/社区街道
	Township string `json:"township"`
	// 乡镇/社区街道编码
	Towncode string `json:"towncode"`
	// 行政编码
	AdCode string `json:"adcode"`
	// 城市编码
	CityCode string `json:"citycode"`
}

type Pois struct {
	// poi地址信息
	Address string `json:"address"`
	// poi点名称
	Name string `json:"name"`
	// 经纬度
	Location string `json:"location"`
	// 该POI中心点到请求坐标的距离
	Distance string `json:"distance"`
	// POI类型
	PoiType string `json:"type"`
}
