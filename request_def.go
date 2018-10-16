package main

type RequestHead struct {
	RequestId string `json:"requestId"`
	UserId    string `json:"userId"`
	Cmd       int32  `json:"cmd"`
}

type ResponseHead struct {
	RequestId string `json:"requestId"`
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Cmd       int32  `json:"cmd"`
}

type WxpluginProgramManagerData struct {
	ProgId      string `json:"progId"`
	ProgName    string `json:"progName"`
	WsId        string `json:"wsId"`
	AppId       string `json:"appId"`
	AppSecrete  string `json:"appSecrete"`
	ProgramType int32  `json:"programType"`
}

type WxpluginProgramManagerReq struct {
	RequestHead
	Data WxpluginProgramManagerData `json:"data"`
}

type WxpluginProgramManagerResp struct {
	ResponseHead
	Data WxpluginProgramManagerData `json:"data"`
}

type GoodsVarietyManagerData struct {
	VarietyId   string `json:"varietyId"`
	VarietyName string `json:"varietyName"`
}

type GoodsVarietyManagerReq struct {
	RequestHead
	Data GoodsVarietyManagerData `json:"data"`
}

type GoodsVarietyManagerResp struct {
	ResponseHead
	Data GoodsVarietyManagerData `json:"data"`
}

type GoodsManagerData struct {
	GoodsId     string                    `json:"goodsId"`
	GoogdsName  string                    `json:"googdsName"`
	VarietyId   string                    `json:"varietyId"`
	GoodsBrief  string                    `json:"goodsBrief"`
	GoodsPrice  int32                     `json:"goodsPrice"`
	ChargeUnit  int32                     `json:"chargeUnit"`
	GoodsWeight int32                     `json:"goodsWeight"`
	WeightUnit  int32                     `json:"weightUnit"`
	GoodsCode   string                    `json:"goodsCode"`
	FactoryId   string                    `json:"factoryId"`
	pictureList []GoodsManagerDataPicture `json:"pictureList"`
}

type GoodsManagerDataPicture struct {
	PictureId    string `json:"pictureId"`
	PicturePath  string `json:"pictureId"`
	PictureOrder int32  `json:"pictureOrder"`
	PictureName  string `json:"pictureName"`
	PictureDesc  string `json:"pictureDesc"`
}

type GoodsManagerReq struct {
	RequestHead
	Data GoodsManagerData `json:"data"`
}

type GoodsManagerResp struct {
	ResponseHead
	Data GoodsManagerData `json:"data"`
}
