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

type FactoryManagerData struct {
	FactoryId      string `json:"factoryId"`
	FactoryName    string `json:"factoryName"`
	LinkPerson     string `json:"linkPerson"`
	LinkPhone      string `json:"linkPhone"`
	FactoryAddress string `json:"factoryAddress"`
}

type FactoryManagerReq struct {
	RequestHead
	Data FactoryManagerData `json:"data"`
}

type FactoryManagerResp struct {
	ResponseHead
	Data FactoryManagerData `json:"data"`
}

type WholesalerManagerData struct {
	WholesalerId   string `json:"wholesalerId"`
	WholesalerName string `json:"wholesalerName"`
	LinkPhone      string `json:"linkPhone"`
	Company        string `json:"company"`
}

type WholesalerManagerReq struct {
	RequestHead
	Data WholesalerManagerData `json:"data"`
}

type WholesalerManagerResp struct {
	ResponseHead
	Data WholesalerManagerData `json:"data"`
}

type WholesalerBannerManagerData struct {
	BannerId   string `json:"bannerId"`
	SalerId    string `json:"salerId"`
	BannerName string `json:"bannerName"`
	BannerPic  string `json:"bannerPic"`
	ShowOrder  int32
	LinkUri    string `json:"linkUri"`
}

type WholesalerBannerManagerReq struct {
	RequestHead
	Data WholesalerBannerManagerData `json:"data"`
}

type WholesalerBannerManagerResp struct {
	ResponseHead
	Data WholesalerBannerManagerData `json:"data"`
}

type SalsemanManagerData struct {
	SalesmanId   string `json:"salesmanId"`
	SalesmanName string `json:"salesmanName"`
	LinkPhone    string `json:"linkPhone"`
	WholesalerId string `json:"wholesalerId"`
}

type SalsemanManagerReq struct {
	RequestHead
	Data SalsemanManagerData `json:"data"`
}

type SalsemanManagerResp struct {
	ResponseHead
	Data SalsemanManagerData `json:"data"`
}

type WholesalerMemberManagerData struct {
	MemberId     string `json:"memberId"`
	MemberName   string `json:"memberName"`
	LinkPhone    string `json:"linkPhone"`
	WholesalerId string `json:"wholesalerId"`
	SalesmanId   string `json:"salesmanId"`
}

type WholesalerMemberManagerReq struct {
	RequestHead
	Data WholesalerMemberManagerData `json:"data"`
}

type WholesalerMemberManagerResp struct {
	ResponseHead
	Data WholesalerMemberManagerData `json:"data"`
}
