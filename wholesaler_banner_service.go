package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type wholesaler_banner_service struct {
	d *dbOperator
}

func (m *wholesaler_banner_service) addWholesalerBanner(req WholesalerBannerManagerData) (string, error) {
	args1 := []interface{}{}
	uid, _ := uuid.NewV4()
	args1 = append(args1, req.SalerId)
	args1 = append(args1, uid.String())
	args1 = append(args1, req.BannerName)
	args1 = append(args1, req.BannerPic)
	args1 = append(args1, req.ShowOrder)
	args1 = append(args1, req.LinkUri)
	queryReq := &SqlExecRequest{
		SQL:  "insert into T_wholesaler_banner(Saler_uuid, Banner_uuid, Banner_name, Banner_pic, Show_order, Banner_status, Open_time, Close_time, Link_uri) values(?,?,?,?,?,1,NULL,NULL,?)",
		Args: args1}
	excuteRep := m.d.dbCli.Query(queryReq)
	if excuteRep.Error() != nil {
		zap.L().Error(fmt.Sprintf("add WholesalerBanner[%s] error:%s", req.BannerName, excuteRep.Error()))
		return "", excuteRep.Error()
	}
	return uid.String(), nil
}

func (m *wholesaler_banner_service) updateWholesalerBanner(req WholesalerBannerManagerData) error {
	args1 := []interface{}{}
	args1 = append(args1, req.SalerId)
	args1 = append(args1, req.BannerName)
	args1 = append(args1, req.BannerPic)
	args1 = append(args1, req.ShowOrder)
	args1 = append(args1, req.LinkUri)
	args1 = append(args1, req.SalerId)
	args1 = append(args1, req.BannerId)
	queryReq := &SqlExecRequest{
		SQL:  "update T_wholesaler_banner  set Saler_uuid=?, Banner_name = ?, Banner_pic = ?, Show_order = ?, Link_uri = ? where  Banner_uuid = ?",
		Args: args1}
	excuteRep := m.d.dbCli.Query(queryReq)
	if excuteRep.Error() != nil {
		zap.L().Error(fmt.Sprintf("update WholesalerBanner[%s] error:%s", req.BannerId, excuteRep.Error()))
		return excuteRep.Error()
	}
	return nil
}

func (m *wholesaler_banner_service) deleteWholesalerBanner(req WholesalerBannerManagerData) error {
	args1 := []interface{}{}
	args1 = append(args1, req.BannerId)
	queryReq := &SqlExecRequest{
		SQL:  "update T_wholesaler_banner  set Banner_status = 0 where Banner_uuid = ?",
		Args: args1}
	excuteRep := m.d.dbCli.Query(queryReq)
	if excuteRep.Error() != nil {
		zap.L().Error(fmt.Sprintf("delete WholesalerBanner[%s] error:%s", req.BannerId, excuteRep.Error()))
		return excuteRep.Error()
	}
	return nil
}

func (m *wholesaler_banner_service) queryWholesalerBannerByExample(req WholesalerBannerManagerData) ([]*TWholesalerBanner, error) {
	args1 := []interface{}{}
	var sql string
	sql = "select Banner_id, Saler_uuid, Banner_uuid, Banner_name, Banner_pic, Show_order, Banner_status, Open_time, Close_time, Link_uri from T_wholesaler_banner where 1=1 "
	if len(req.BannerId) != 0 {
		sql += " and Banner_uuid = ?"
		args1 = append(args1, req.BannerId)
	}
	if len(req.BannerName) != 0 {
		sql += " and variety_name = ?"
		args1 = append(args1, req.BannerName)
	}
	if len(req.SalerId) != 0 {
		sql += " and Saler_uuid = ?"
		args1 = append(args1, req.SalerId)
	}
	tmp := TWholesalerBanner{}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args1,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query T_wholesaler_banner error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	var returnMenus []*TWholesalerBanner = []*TWholesalerBanner{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnMenus = append(returnMenus, queryRep.Rows[i].(*TWholesalerBanner))
	}
	return returnMenus, nil
}
