package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type clientInfo struct {
	ipStr string
	ipNum int32
	port  int32
}

type httpHandler struct {
	cfg                *Config
	wxpluginProgramSv  *wxpluginProgram_service
	factorySv          *factory_service
	goodsSv            *goods_service
	goodsVarietySv     *goodsVariety_service
	wholesalerBannerSv *wholesaler_banner_service
	salermanSv         *salseman_service
	wholesalerSv       *wholesaler_service
}

func (ci *clientInfo) inetAton() {
	ip := net.ParseIP(ci.ipStr)
	ci.ipNum = int32(binary.BigEndian.Uint32(ip.To4()))
}

func (m *httpHandler) start() error {
	//start http server
	s := &http.Server{
		Addr:           m.cfg.Server.Endpoint,
		Handler:        nil,
		ReadTimeout:    m.cfg.Server.HttpReadTimeout,
		WriteTimeout:   m.cfg.Server.HttpWriteTimeout,
		MaxHeaderBytes: int(m.cfg.Server.MaxHeadSize),
	}
	http.HandleFunc("/api", m.process)
	go s.ListenAndServe()

	return nil
}

func (m *httpHandler) ivalidResp(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (m *httpHandler) getClientInfo(r *http.Request) *clientInfo {
	cliIp, cliPort, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("userip: %q is not IP:port", r.RemoteAddr))
		return &clientInfo{ipNum: 0, port: 0}
	} else {
		zap.L().Debug(fmt.Sprintf("package from %s:%s", cliIp, cliPort))
		p, e := strconv.Atoi(cliPort)
		if e != nil {
			zap.L().Error(fmt.Sprintf("strconv Atoi port fail"))
			p = 0
		}

		ci := &clientInfo{
			ipStr: cliIp,
			port:  int32(p),
			ipNum: 0,
		}

		ci.inetAton()
		return ci
	}
}

func (m *httpHandler) process(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		zap.L().Info(fmt.Sprintf("get method not support, method:%s", r.Method))
		statObj.statHandler.StatCount(StatInvalidMethodReq)
		m.ivalidResp(w)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		statObj.statHandler.StatCount(StatReadBody)
		m.ivalidResp(w)
		return
	} else {
		zap.L().Debug(fmt.Sprintf("recv body len:%d content:%s", len(body), body))
		var req RequestHead
		err := json.Unmarshal(body, &req)
		if err != nil {
			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
			m.ivalidResp(w)
			return
		}
		if req.Cmd == 2000 || req.Cmd == 2002 || req.Cmd == 2004 || req.Cmd == 2006 {
			m.processWxpluginProgram(body, w)
		} else if req.Cmd == 2020 || req.Cmd == 2022 || req.Cmd == 2024 || req.Cmd == 2026 {
			m.processGoodsVariety(body, w)
		} else if req.Cmd == 2040 || req.Cmd == 2042 || req.Cmd == 2044 || req.Cmd == 2046 {
			m.processGoods(body, w)
		} else if req.Cmd == 2060 || req.Cmd == 2062 || req.Cmd == 2064 || req.Cmd == 2066 {
			m.processFactory(body, w)
		} else if req.Cmd == 2080 || req.Cmd == 2082 || req.Cmd == 2084 || req.Cmd == 2086 {
			m.processWholesaler(body, w)
		} else if req.Cmd == 3000 || req.Cmd == 3002 || req.Cmd == 3004 || req.Cmd == 3006 {
			m.processFactory(body, w)
		} else if req.Cmd == 3020 || req.Cmd == 3022 || req.Cmd == 3024 || req.Cmd == 3026 {
			m.processSalerman(body, w)
		} else {
			var respHead ResponseHead
			respHead = ResponseHead{RequestId: req.RequestId, ErrorCode: 9999, Cmd: req.Cmd, ErrorMsg: "cmd不合法"}
			jsonData, err := json.Marshal(respHead)
			if err != nil {
				zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
				m.ivalidResp(w)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(jsonData))
			return
		}
	}
}

func (m *httpHandler) processWxpluginProgram(body []byte, w http.ResponseWriter) {
	var req WxpluginProgramManagerReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	var resp WxpluginProgramManagerResp
	resp = WxpluginProgramManagerResp{
		ResponseHead{
			RequestId: req.RequestId,
			ErrorCode: 0,
			Cmd:       req.Cmd + 1,
		},
		WxpluginProgramManagerData{
			ProgName:    req.Data.ProgName,
			WsId:        req.Data.WsId,
			AppId:       req.Data.AppId,
			AppSecrete:  req.Data.AppSecrete,
			ProgramType: req.Data.ProgramType,
		},
	}
	var progId string
	if req.Cmd == 2000 {
		progId, err = m.wxpluginProgramSv.addWxpluginProgram(req.Data)
		resp.Data.ProgId = progId
	} else if req.Cmd == 2002 {
		err = m.wxpluginProgramSv.updateWxpluginProgram(req.Data)
	} else if req.Cmd == 2004 {
		err = m.wxpluginProgramSv.deleteWxpluginProgram(req.Data)
	} else if req.Cmd == 2006 {
		var tWxpluginPrograms []*TWxpluginProgram
		tWxpluginPrograms, err = m.wxpluginProgramSv.queryWxpluginProgramByExample(req.Data)
		if len(tWxpluginPrograms) == 0 {
			err = errors.New("查询不到对应的数据")
		} else {
			resp.Data.AppId = tWxpluginPrograms[0].Appid
			resp.Data.AppSecrete = tWxpluginPrograms[0].Appsecrete
			resp.Data.ProgId = tWxpluginPrograms[0].Program_uuid
			resp.Data.ProgName = tWxpluginPrograms[0].Program_name
			resp.Data.ProgramType = tWxpluginPrograms[0].Program_type
			resp.Data.WsId = tWxpluginPrograms[0].Saler_uuid.String
		}
	}
	if err != nil {
		resp.ErrorCode = 9999
		resp.ErrorMsg = err.Error()
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}

func (m *httpHandler) processGoodsVariety(body []byte, w http.ResponseWriter) {
	var req GoodsVarietyManagerReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	var resp GoodsVarietyManagerResp
	resp = GoodsVarietyManagerResp{
		ResponseHead{
			RequestId: req.RequestId,
			ErrorCode: 0,
			Cmd:       req.Cmd + 1,
		},
		GoodsVarietyManagerData{
			VarietyId:   req.Data.VarietyId,
			VarietyName: req.Data.VarietyName,
		},
	}
	var varietyId string
	if req.Cmd == 2020 {
		varietyId, err = m.goodsVarietySv.addGoodsVariety(req.Data, req.UserId)
		resp.Data.VarietyId = varietyId
	} else if req.Cmd == 2022 {
		err = m.goodsVarietySv.updateGoodsVariety(req.Data, req.UserId)
	} else if req.Cmd == 2024 {
		err = m.goodsVarietySv.deleteGoodsVariety(req.Data, req.UserId)
	} else if req.Cmd == 2026 {
		var tGoodsVariety []*TGoodsVariety
		tGoodsVariety, err = m.goodsVarietySv.queryGoodsVarietyByExample(req.Data)
		if len(tGoodsVariety) == 0 {
			err = errors.New("查询不到对应的数据")
		} else {
			resp.Data.VarietyId = tGoodsVariety[0].Variety_uuid
			resp.Data.VarietyName = tGoodsVariety[0].Variety_name
		}
	}
	if err != nil {
		resp.ErrorCode = 9999
		resp.ErrorMsg = err.Error()
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}

func (m *httpHandler) processWholesaler(body []byte, w http.ResponseWriter) {
	var req WholesalerManagerReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	var resp WholesalerManagerResp
	resp = WholesalerManagerResp{
		ResponseHead{
			RequestId: req.RequestId,
			ErrorCode: 0,
			Cmd:       req.Cmd + 1,
		},
		WholesalerManagerData{
			WholesalerId:   req.Data.WholesalerId,
			WholesalerName: req.Data.WholesalerName,
			LinkPhone:      req.Data.LinkPhone,
			Company:        req.Data.Company,
		},
	}
	var wholesaler_uuid string
	if req.Cmd == 2080 {
		wholesaler_uuid, err = m.wholesalerSv.addWholesaler(req.Data, req.UserId)
		resp.Data.WholesalerId = wholesaler_uuid
	} else if req.Cmd == 2082 {
		err = m.wholesalerSv.updateWholesaler(req.Data, req.UserId)
	} else if req.Cmd == 2084 {
		err = m.wholesalerSv.deleteWholesaler(req.Data, req.UserId)
	} else if req.Cmd == 2086 {
		var tWholeSaler []*TWholeSaler
		tWholeSaler, err = m.wholesalerSv.queryWholesalerByExample(req.Data)
		if len(tWholeSaler) == 0 {
			err = errors.New("查询不到对应的数据")
		} else {
			resp.Data.WholesalerId = tWholeSaler[0].Saler_uuid
			resp.Data.WholesalerName = tWholeSaler[0].Saler_name.String
			resp.Data.Company = tWholeSaler[0].Company
			resp.Data.LinkPhone = tWholeSaler[0].Mobile
		}
	}
	if err != nil {
		resp.ErrorCode = 9999
		resp.ErrorMsg = err.Error()
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}

func (m *httpHandler) processGoods(body []byte, w http.ResponseWriter) {
	var req GoodsManagerReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	var resp GoodsManagerResp
	resp = GoodsManagerResp{
		ResponseHead{
			RequestId: req.RequestId,
			ErrorCode: 0,
			Cmd:       req.Cmd + 1,
		},
		GoodsManagerData{
			GoodsId:     req.Data.GoodsId,
			GoogdsName:  req.Data.GoogdsName,
			VarietyId:   req.Data.VarietyId,
			GoodsBrief:  req.Data.GoodsBrief,
			GoodsPrice:  req.Data.GoodsPrice,
			ChargeUnit:  req.Data.ChargeUnit,
			GoodsWeight: req.Data.GoodsWeight,
			WeightUnit:  req.Data.WeightUnit,
			GoodsCode:   req.Data.GoodsCode,
			FactoryId:   req.Data.FactoryId,
			pictureList: req.Data.pictureList,
		},
	}
	var goodsUUid string
	if req.Cmd == 2040 {
		goodsUUid, err = m.goodsSv.addGoods(req.Data, req.UserId)
		resp.Data.GoodsId = goodsUUid
	} else if req.Cmd == 2042 {
		err = m.goodsSv.updateGoods(req.Data, req.UserId)
	} else if req.Cmd == 2044 {
		err = m.goodsSv.deleteGoods(req.Data, req.UserId)
	} else if req.Cmd == 2026 {
		var goodsManagerDatas []*GoodsManagerData
		goodsManagerDatas, err = m.goodsSv.queryGoodsByExample(req.Data)
		if len(goodsManagerDatas) == 0 {
			err = errors.New("查询不到对应的数据")
		} else {
			resp.Data = *goodsManagerDatas[0]
		}
	}
	if err != nil {
		resp.ErrorCode = 9999
		resp.ErrorMsg = err.Error()
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}

func (m *httpHandler) processFactory(body []byte, w http.ResponseWriter) {
	var req FactoryManagerReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	var resp FactoryManagerResp
	resp = FactoryManagerResp{
		ResponseHead{
			RequestId: req.RequestId,
			ErrorCode: 0,
			Cmd:       req.Cmd + 1,
		},
		FactoryManagerData{},
	}
	resp.Data = req.Data
	var factoryId string
	if req.Cmd == 2060 {
		factoryId, err = m.factorySv.addFactory(req.Data)
		resp.Data.FactoryId = factoryId
	} else if req.Cmd == 2062 {
		err = m.factorySv.updateFactory(req.Data)
	} else if req.Cmd == 2064 {
		err = m.factorySv.deleteFactory(req.Data)
	} else if req.Cmd == 2066 {
		var tFactory []*TFactory
		tFactory, err = m.factorySv.queryFactoryByExample(req.Data)
		if len(tFactory) == 0 {
			err = errors.New("查询不到对应的数据")
		} else {
			resp.Data.FactoryId = tFactory[0].Factory_uuid
			resp.Data.FactoryName = tFactory[0].Factory_name
			resp.Data.LinkPerson = tFactory[0].Link_person.String
			resp.Data.LinkPhone = tFactory[0].Link_phone.String
			resp.Data.FactoryAddress = tFactory[0].Factory_address.String
		}
	}
	if err != nil {
		resp.ErrorCode = 9999
		resp.ErrorMsg = err.Error()
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}

func (m *httpHandler) processBanner(body []byte, w http.ResponseWriter) {
	var req WholesalerBannerManagerReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	var resp WholesalerBannerManagerResp
	resp = WholesalerBannerManagerResp{
		ResponseHead{
			RequestId: req.RequestId,
			ErrorCode: 0,
			Cmd:       req.Cmd + 1,
		},
		WholesalerBannerManagerData{},
	}
	resp.Data = req.Data
	var bannerId string
	if req.Cmd == 3000 {
		bannerId, err = m.wholesalerBannerSv.addWholesalerBanner(req.Data)
		resp.Data.BannerId = bannerId
	} else if req.Cmd == 3002 {
		err = m.wholesalerBannerSv.updateWholesalerBanner(req.Data)
	} else if req.Cmd == 3004 {
		err = m.wholesalerBannerSv.deleteWholesalerBanner(req.Data)
	} else if req.Cmd == 3006 {
		var tWholesalerBanners []*TWholesalerBanner
		tWholesalerBanners, err = m.wholesalerBannerSv.queryWholesalerBannerByExample(req.Data)
		if len(tWholesalerBanners) == 0 {
			err = errors.New("查询不到对应的数据")
		} else {
			resp.Data.BannerId = tWholesalerBanners[0].Banner_uuid
			resp.Data.BannerName = tWholesalerBanners[0].Banner_name.String
			resp.Data.BannerPic = tWholesalerBanners[0].Banner_pic
			resp.Data.LinkUri = tWholesalerBanners[0].Link_uri.String
			resp.Data.SalerId = tWholesalerBanners[0].Saler_uuid
			resp.Data.ShowOrder = tWholesalerBanners[0].Show_order
		}
	}
	if err != nil {
		resp.ErrorCode = 9999
		resp.ErrorMsg = err.Error()
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}

func (m *httpHandler) processSalerman(body []byte, w http.ResponseWriter) {
	var req SalsemanManagerReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	var resp SalsemanManagerResp
	resp = SalsemanManagerResp{
		ResponseHead{
			RequestId: req.RequestId,
			ErrorCode: 0,
			Cmd:       req.Cmd + 1,
		},
		SalsemanManagerData{
			SalesmanId:   req.Data.SalesmanId,
			SalesmanName: req.Data.SalesmanName,
			LinkPhone:    req.Data.LinkPhone,
			WholesalerId: req.Data.WholesalerId,
		},
	}
	resp.Data = req.Data
	var salesman_uuid string
	if req.Cmd == 3020 {
		salesman_uuid, err = m.salermanSv.addSalseman(req.Data, req.UserId)
		resp.Data.SalesmanId = salesman_uuid
	} else if req.Cmd == 3022 {
		err = m.salermanSv.updateSalseman(req.Data, req.UserId)
	} else if req.Cmd == 3024 {
		err = m.salermanSv.deleteSalseman(req.Data, req.UserId)
	} else if req.Cmd == 3026 {
		var tSalseman []*TSalseman
		tSalseman, err = m.salermanSv.querySalsemanByExample(req.Data)
		if len(tSalseman) == 0 {
			err = errors.New("查询不到对应的数据")
		} else {
			resp.Data.SalesmanId = tSalseman[0].Salesman_uuid
			resp.Data.SalesmanName = tSalseman[0].Salesman_name
			resp.Data.LinkPhone = tSalseman[0].Salesman_phone
			resp.Data.WholesalerId = tSalseman[0].Saler_uuid
		}
	}
	if err != nil {
		resp.ErrorCode = 9999
		resp.ErrorMsg = err.Error()
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
		m.ivalidResp(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
	return
}
