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
		} else if req.Cmd == 3000 || req.Cmd == 3002 || req.Cmd == 3004 || req.Cmd == 3006 {
			m.processFactory(body, w)
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

//func (m *httpHandler) addSysUser(body []byte, w http.ResponseWriter) {
//	var req SystemManagerUserReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	sysUserUUid, err := m.systemUserSv.addSysUser(req.Data, req.UserId)
//	var resp SystemManagerUserResp
//	if err == nil {
//		resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1001, ErrorCode: 0}, SystemManagerUserRespData{SysUserId: sysUserUUid, UserName: req.Data.UserName, LoginName: req.Data.LoginName, UserMobile: req.Data.UserMobile, UserEMail: req.Data.UserEMail}}
//	} else {
//		resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1001, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerUserRespData{}}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) updateSysUser(body []byte, w http.ResponseWriter) {
//	var req SystemManagerUserReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerUserResp
//	if len(req.Data.SysUserId) == 0 {
//		resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1003, ErrorCode: 9999, ErrorMsg: "sysUserId不能为空"}, SystemManagerUserRespData{}}
//	} else {
//		errMsg, err := m.systemUserSv.updateSysUser(req.Data, req.UserId)
//		if err == nil {
//			if len(errMsg) == 0 {
//				resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1003, ErrorCode: 0}, SystemManagerUserRespData{SysUserId: req.Data.SysUserId, UserName: req.Data.UserName, LoginName: req.Data.LoginName, UserMobile: req.Data.UserMobile, UserEMail: req.Data.UserEMail}}
//			} else {
//				resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1003, ErrorCode: 9999, ErrorMsg: errMsg}, SystemManagerUserRespData{}}
//			}
//		} else {
//			resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1003, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerUserRespData{}}
//		}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) deleteSysUser(body []byte, w http.ResponseWriter) {
//	var req SystemManagerUserReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerUserResp
//	if len(req.Data.SysUserId) == 0 {
//		resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1005, ErrorCode: 9999, ErrorMsg: "sysUserId不能为空"}, SystemManagerUserRespData{}}

//	} else {
//		err := m.systemUserSv.deleteSysUser(req.Data.SysUserId, req.UserId)
//		if err == nil {
//			resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1005, ErrorCode: 0}, SystemManagerUserRespData{SysUserId: req.Data.SysUserId}}
//		} else {
//			resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1005, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerUserRespData{}}
//		}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) querySysUser(body []byte, w http.ResponseWriter) {
//	var req SystemManagerUserReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerUserResp
//	var tSysUsers []*TSysUser
//	tSysUsers, err = m.systemUserSv.queryAvailableSysUserByExample(req.Data)
//	if err == nil {
//		if tSysUsers == nil || len(tSysUsers) == 0 {
//			resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1007, ErrorCode: 9999, ErrorMsg: "未查到对应的系统用户"}, SystemManagerUserRespData{}}
//		} else {
//			resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1007, ErrorCode: 0}, SystemManagerUserRespData{SysUserId: tSysUsers[0].User_uuid, UserName: tSysUsers[0].User_name.String, LoginName: tSysUsers[0].Login_name.String, UserMobile: tSysUsers[0].User_phone.String, UserEMail: tSysUsers[0].User_email.String}}
//		}
//	} else {
//		resp = SystemManagerUserResp{ResponseHead{RequestId: req.RequestId, Cmd: 1007, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerUserRespData{}}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) addSysRole(body []byte, w http.ResponseWriter) {
//	var req SystemManagerRoleReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerRoleResp
//	roleUUid, err := m.systemRoleSv.addSysRole(req.Data, req.UserId)
//	if err == nil {
//		resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1021, ErrorCode: 0}, SystemManagerRoleRespData{RoleId: roleUUid, RoleName: req.Data.RoleName, RoleLevel: req.Data.RoleLevel}}
//	} else {
//		resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1021, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerRoleRespData{}}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) updateSysRole(body []byte, w http.ResponseWriter) {
//	var req SystemManagerRoleReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerRoleResp
//	if len(req.Data.RoleId) == 0 {
//		resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1023, ErrorCode: 9999, ErrorMsg: "roleId不能为空"}, SystemManagerRoleRespData{}}
//	} else {
//		err := m.systemRoleSv.updateSysRole(req.Data, req.UserId)
//		if err == nil {
//			resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1023, ErrorCode: 0}, SystemManagerRoleRespData{RoleId: req.Data.RoleId, RoleName: req.Data.RoleName, RoleLevel: req.Data.RoleLevel}}
//		} else {
//			resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1023, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerRoleRespData{}}
//		}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) deleteSysRole(body []byte, w http.ResponseWriter) {
//	var req SystemManagerRoleReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerRoleResp
//	if len(req.Data.RoleId) == 0 {
//		resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1025, ErrorCode: 9999, ErrorMsg: "roleId不能为空"}, SystemManagerRoleRespData{}}
//	} else {
//		err := m.systemRoleSv.deleteSysRole(req.Data, req.UserId)
//		if err == nil {
//			resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1025, ErrorCode: 0}, SystemManagerRoleRespData{RoleId: req.Data.RoleId}}
//		} else {
//			resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1025, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerRoleRespData{}}
//		}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) querySysRole(body []byte, w http.ResponseWriter) {
//	var req SystemManagerRoleReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerRoleResp
//	var tSysRoles []*TSysRole
//	tSysRoles, err = m.systemRoleSv.queryAvailableSysRole(req.Data)
//	if err == nil {
//		if tSysRoles == nil || len(tSysRoles) == 0 {
//			resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1027, ErrorCode: 9999, ErrorMsg: "未查到对应的角色"}, SystemManagerRoleRespData{}}

//		} else {
//			resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1027, ErrorCode: 0}, SystemManagerRoleRespData{RoleId: tSysRoles[0].Role_uuid, RoleName: tSysRoles[0].Role_name.String, RoleLevel: tSysRoles[0].Role_level}}
//		}
//	} else {
//		resp = SystemManagerRoleResp{ResponseHead{RequestId: req.RequestId, Cmd: 1027, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerRoleRespData{}}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) addSysMenu(body []byte, w http.ResponseWriter) {
//	var req SystemManagerMenuReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	sysMenuUUid, err := m.systemMenuSv.addSysMenu(req.Data, req.UserId)
//	var resp SystemManagerMenuResp
//	if err == nil {
//		resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1041, ErrorCode: 0}, SystemManagerMenuRespData{MenuId: sysMenuUUid, MenuName: req.Data.MenuName, MenuLevel: req.Data.MenuLevel}}
//	} else {
//		resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1041, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerMenuRespData{}}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) updateSysMenu(body []byte, w http.ResponseWriter) {
//	var req SystemManagerMenuReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerMenuResp
//	if len(req.Data.MenuId) == 0 {
//		resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1043, ErrorCode: 9999, ErrorMsg: "sysUserId不能为空"}, SystemManagerMenuRespData{}}
//	} else {
//		errMsg, err := m.systemMenuSv.updateSysMenu(req.Data, req.UserId)
//		if err == nil {
//			if len(errMsg) == 0 {
//				resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1043, ErrorCode: 0}, SystemManagerMenuRespData{MenuId: req.Data.MenuId, MenuName: req.Data.MenuName, MenuLevel: req.Data.MenuLevel}}
//			} else {
//				resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1043, ErrorCode: 9999, ErrorMsg: errMsg}, SystemManagerMenuRespData{}}
//			}
//		} else {
//			resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1043, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerMenuRespData{}}
//		}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) deleteSysMenu(body []byte, w http.ResponseWriter) {
//	var req SystemManagerMenuReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerMenuResp
//	if len(req.Data.MenuId) == 0 {
//		resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1045, ErrorCode: 9999, ErrorMsg: "MenuId不能为空"}, SystemManagerMenuRespData{}}
//	} else {
//		err := m.systemMenuSv.deleteSysMenu(req.Data.MenuId, req.UserId)
//		if err == nil {
//			resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1045, ErrorCode: 0}, SystemManagerMenuRespData{MenuId: req.Data.MenuId}}
//		} else {
//			resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1045, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerMenuRespData{}}
//		}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) querySysMenu(body []byte, w http.ResponseWriter) {
//	var req SystemManagerMenuReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerMenuResp
//	var tSysMenus []*TSysMenu
//	tSysMenus, err = m.systemMenuSv.queryAvailableSysMeunByExample(req.Data)
//	if err == nil {
//		if tSysMenus == nil || len(tSysMenus) == 0 {
//			resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1047, ErrorCode: 9999, ErrorMsg: "未查到对应的菜单"}, SystemManagerMenuRespData{}}

//		} else {
//			resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1047, ErrorCode: 0}, SystemManagerMenuRespData{MenuId: tSysMenus[0].Menu_uuid, MenuName: tSysMenus[0].Menu_name, MenuLevel: tSysMenus[0].Menu_level}}
//		}
//	} else {
//		resp = SystemManagerMenuResp{ResponseHead{RequestId: req.RequestId, Cmd: 1047, ErrorCode: 9999, ErrorMsg: err.Error()}, SystemManagerMenuRespData{}}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) addSysPrivilege(body []byte, w http.ResponseWriter) {
//	var req SystemManagerPrivilegeReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var systemMangerPrivilegeRespDatas []SystemManagerPrivilegeRespData
//	systemMangerPrivilegeRespDatas, err = m.systemPrivilegeSv.addSysPrivilege(req.Data, req.UserId)
//	var resp SystemManagerPrivilegeResp
//	if err == nil {
//		resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1061, ErrorCode: 0}, []SystemManagerPrivilegeRespData{}}
//		for i := 0; i < len(systemMangerPrivilegeRespDatas); i++ {
//			resp.Data = append(resp.Data, systemMangerPrivilegeRespDatas[i])
//		}
//	} else {
//		resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1061, ErrorCode: 9999, ErrorMsg: err.Error()}, []SystemManagerPrivilegeRespData{}}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) updateSysPrivilege(body []byte, w http.ResponseWriter) {
//	var req SystemManagerPrivilegeReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerPrivilegeResp
//	if len(req.Data.PowerId) == 0 {
//		resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1063, ErrorCode: 9999, ErrorMsg: "powerId不能为空"}, []SystemManagerPrivilegeRespData{}}
//	} else {
//		err := m.systemPrivilegeSv.updateSysPrivilege(req.Data, req.UserId)
//		if err == nil {
//			resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1063, ErrorCode: 0}, []SystemManagerPrivilegeRespData{}}
//		} else {
//			resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1063, ErrorCode: 9999, ErrorMsg: err.Error()}, []SystemManagerPrivilegeRespData{}}
//		}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) deleteSysPrivilege(body []byte, w http.ResponseWriter) {
//	var req SystemManagerPrivilegeReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerPrivilegeResp
//	if len(req.Data.PowerId) == 0 {
//		resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1065, ErrorCode: 9999, ErrorMsg: "powerId不能为空"}, []SystemManagerPrivilegeRespData{}}
//	} else {
//		err := m.systemPrivilegeSv.deleteSysPrivilege(req.Data.PowerId, req.UserId)
//		if err == nil {
//			resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1065, ErrorCode: 0}, []SystemManagerPrivilegeRespData{}}
//		} else {
//			resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1065, ErrorCode: 9999, ErrorMsg: err.Error()}, []SystemManagerPrivilegeRespData{}}
//		}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) querySysPrivilege(body []byte, w http.ResponseWriter) {
//	var req SystemManagerPrivilegeReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp SystemManagerPrivilegeResp
//	var tSysRoleMenus []*TSysRoleMenu
//	tSysRoleMenus, err = m.systemPrivilegeSv.queryAvailableSysPrivilegeByExample(req.Data)
//	if err == nil {
//		if tSysRoleMenus == nil || len(tSysRoleMenus) == 0 {
//			resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1067, ErrorCode: 9999, ErrorMsg: "未查到对应的权限记录"}, []SystemManagerPrivilegeRespData{}}
//		} else {
//			resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1067, ErrorCode: 0}, []SystemManagerPrivilegeRespData{}}
//			for i := 0; i < len(tSysRoleMenus); i++ {
//				resp.Data = append(resp.Data, SystemManagerPrivilegeRespData{PowerId: tSysRoleMenus[i].Power_uuid, RoleId: tSysRoleMenus[i].Role_uuid, MenuId: tSysRoleMenus[i].Menu_uuid})
//			}
//		}
//	} else {
//		resp = SystemManagerPrivilegeResp{ResponseHead{RequestId: req.RequestId, Cmd: 1067, ErrorCode: 9999, ErrorMsg: err.Error()}, []SystemManagerPrivilegeRespData{}}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}
