package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type goodsVariety_service struct {
	d *dbOperator
}

func (m *goodsVariety_service) addGoodsVariety(req GoodsVarietyManagerData, opId string) (string, error) {
	args1 := []interface{}{}
	uid, _ := uuid.NewV4()
	args1 = append(args1, uid.String())
	args1 = append(args1, req.VarietyName)
	args1 = append(args1, opId)
	args1 = append(args1, opId)

	queryReq := &SqlExecRequest{
		SQL:  "insert into T_Goods_Variety(variety_uuid, variety_name, variety_status, create_time, create_user, update_time, update_user) values(?,?,1,now(),?,now(),?)",
		Args: args1}
	excuteRep := m.d.dbCli.Query(queryReq)
	if excuteRep.Error() != nil {
		zap.L().Error(fmt.Sprintf("add goodsVariety[%s] error:%s", req.VarietyName, excuteRep.Error()))
		return "", excuteRep.Error()
	}
	return uid.String(), nil
}

func (m *goodsVariety_service) updateGoodsVariety(req GoodsVarietyManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, req.VarietyName)
	args1 = append(args1, opId)
	args1 = append(args1, req.VarietyId)
	queryReq := &SqlExecRequest{
		SQL:  "update T_Goods_Variety  set variety_name = ?, update_time = now(), update_user = ? where variety_uuid = ?",
		Args: args1}
	excuteRep := m.d.dbCli.Query(queryReq)
	if excuteRep.Error() != nil {
		zap.L().Error(fmt.Sprintf("update goodsVariety[%s] error:%s", req.VarietyId, excuteRep.Error()))
		return excuteRep.Error()
	}
	return nil
}

func (m *goodsVariety_service) deleteGoodsVariety(req GoodsVarietyManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, opId)
	args1 = append(args1, req.VarietyId)
	queryReq := &SqlExecRequest{
		SQL:  "update T_Goods_Variety  set variety_status = 0, update_time = now(), update_user = ? where variety_uuid = ?",
		Args: args1}
	excuteRep := m.d.dbCli.Query(queryReq)
	if excuteRep.Error() != nil {
		zap.L().Error(fmt.Sprintf("delete goodsVariety[%s] error:%s", req.VarietyId, excuteRep.Error()))
		return excuteRep.Error()
	}
	return nil
}

func (m *goodsVariety_service) queryGoodsVarietyByExample(req GoodsVarietyManagerData) ([]*TGoodsVariety, error) {
	args1 := []interface{}{}
	var sql string
	sql = "select variety_id, variety_uuid, variety_name, variety_status, create_time, create_user, update_time, update_user from T_Goods_Variety where 1=1 "
	if len(req.VarietyId) != 0 {
		sql += " and variety_uuid = ?"
		args1 = append(args1, req.VarietyId)
	}
	if len(req.VarietyName) != 0 {
		sql += " and variety_name = ?"
		args1 = append(args1, req.VarietyName)
	}
	tmp := TGoodsVariety{}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args1,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query T_Goods_Variety error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	var returnMenus []*TGoodsVariety = []*TGoodsVariety{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnMenus = append(returnMenus, queryRep.Rows[i].(*TGoodsVariety))
	}
	return returnMenus, nil
}
