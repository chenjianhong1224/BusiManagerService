package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type goods_service struct {
	d *dbOperator
}

func (m *goods_service) addGoods(req GoodsManagerData, opId string) (string, error) {
	args1 := []interface{}{}
	uid, _ := uuid.NewV4()
	args1 = append(args1, uid.String())
	args1 = append(args1, req.GoogdsName)
	args1 = append(args1, req.GoodsCode)
	args1 = append(args1, req.FactoryId)
	args1 = append(args1, req.GoodsPrice)
	args1 = append(args1, req.ChargeUnit)
	args1 = append(args1, req.GoodsWeight)
	args1 = append(args1, req.WeightUnit)
	args1 = append(args1, req.VarietyId)
	args1 = append(args1, req.GoodsBrief)
	args1 = append(args1, opId)
	args1 = append(args1, opId)
	execReq1 := SqlExecRequest{
		SQL:  "insert into T_Goods(goods_uuid, goods_name, goods_bar_code, factory_uuid, goods_price, charge_unit, goods_weight, weight_unit, goods_desc, goods_status, whole_pack, pack_unit, variety_uuid, create_time, create_user, update_time, update_user) values(?,?,?,?,?,?,?,?,?,1,NULL,NULL,?,now(),?,now(),?)",
		Args: args1,
	}
	var execReqList = []SqlExecRequest{execReq1}

	for i := 0; len(req.pictureList) < i; i++ {
		args2 := []interface{}{}
		args2 = append(args2, uid.String())
		args2 = append(args2, req.pictureList[i].PictureOrder)
		args2 = append(args2, req.pictureList[i].PicturePath)
		args2 = append(args2, req.pictureList[i].PictureName)
		args2 = append(args2, req.pictureList[i].PictureDesc)
		args2 = append(args2, req.pictureList[i].PictureId)
		execReq2 := SqlExecRequest{
			SQL:  "update T_Goods_picture set goods_uuid = ?, show_order = ?, picture_path = ?, picture_name = ?, picture_desc = ? where picture_uuid = ?",
			Args: args2,
		}
		execReqList = append(execReqList, execReq2)
		args3 := []interface{}{}
		args3 = append(args3, uid.String())
		args3 = append(args3, req.pictureList[i].PicturePath)
		execReq3 := SqlExecRequest{
			SQL:  "update T_Goods set Goods_picture = ? where goods_uuid = ?",
			Args: args3,
		}
		execReqList = append(execReqList, execReq2)
		execReqList = append(execReqList, execReq3)
	}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return uid.String(), nil
	}
	zap.L().Error(fmt.Sprintf("add goods[%s] error:%s", req.GoogdsName, err.Error()))
	return "", err
}

func (m *goods_service) updateGoods(req GoodsManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, req.GoodsBrief)
	args1 = append(args1, req.GoogdsName)
	args1 = append(args1, req.GoodsCode)
	args1 = append(args1, req.FactoryId)
	args1 = append(args1, req.GoodsPrice)
	args1 = append(args1, req.ChargeUnit)
	args1 = append(args1, req.GoodsWeight)
	args1 = append(args1, req.WeightUnit)
	args1 = append(args1, req.VarietyId)
	args1 = append(args1, opId)
	args1 = append(args1, req.GoodsId)
	execReq1 := SqlExecRequest{
		SQL:  "update T_Goods set goods_desc = ?, goods_name = ?, goods_bar_code = ?, factory_uuid = ?, goods_price = ?, charge_unit = ?, goods_weight = ?, weight_unit = ?, variety_uuid = ?, update_time = now(), update_user = ? where goods_uuid = ?",
		Args: args1,
	}
	args2 := []interface{}{}
	args2 = append(args2, req.GoodsId)
	execReq2 := SqlExecRequest{
		SQL:  "update T_Goods_picture set picture_status = 0 where goods_uuid = ?",
		Args: args2,
	}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	for i := 0; i < len(req.pictureList); i++ {
		args3 := []interface{}{}
		args3 = append(args3, req.GoodsId)
		args3 = append(args3, req.pictureList[i].PictureOrder)
		args3 = append(args3, req.pictureList[i].PicturePath)
		args3 = append(args3, req.pictureList[i].PictureName)
		args3 = append(args3, req.pictureList[i].PictureDesc)
		args3 = append(args3, req.pictureList[i].PictureId)
		execReq3 := SqlExecRequest{
			SQL:  "update T_Goods_picture set goods_uuid = ?, show_order = ?, picture_path = ?, picture_name = ?, picture_desc = ? where picture_uuid = ?",
			Args: args3,
		}
		execReqList = append(execReqList, execReq3)
	}
	err := m.d.dbCli.TransationExcute(execReqList)
	return err
}

func (m *goods_service) deleteGoods(req GoodsManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, opId)
	args1 = append(args1, req.GoodsId)
	execReq1 := SqlExecRequest{
		SQL:  "update T_Goods goods_status = 0, update_time = now(), update_user = ? where goods_uuid = ?",
		Args: args1,
	}
	args2 := []interface{}{}
	args2 = append(args2, req.GoodsId)
	execReq2 := SqlExecRequest{
		SQL:  "update T_Goods_picture set picture_status = 0 where goods_uuid = ?",
		Args: args2,
	}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	return err
}

func (m *goods_service) queryGoodsByExample(req GoodsManagerData) ([]*GoodsManagerData, error) {
	args := []interface{}{}
	var sql string
	sql = "select goods_id, goods_uuid, goods_name, goods_bar_code, factory_uuid, goods_price, charge_unit, goods_weight, weight_unit, goods_desc, goods_status, whole_pack, pack_unit, variety_uuid, create_time, create_user, update_time, update_user from T_Goods where 1=1 "
	if len(req.GoodsId) > 0 {
		sql += " and goods_uuid = ?"
		args = append(args, req.GoodsId)
	}
	if len(req.GoogdsName) > 0 {
		sql += " and goods_name =?"
		args = append(args, req.GoogdsName)
	}
	if len(req.FactoryId) > 0 {
		sql += " and factory_uuid"
		args = append(args, req.FactoryId)
	}
	if len(req.VarietyId) > 0 {
		sql += "and variety_id = ?"
		args = append(args, req.VarietyId)
	}
	tmp := TGoods{}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query goods error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	var returngoodsManagerData []*GoodsManagerData = []*GoodsManagerData{}
	for i := 0; i < len(queryRep.Rows); i++ {
		goods := queryRep.Rows[i].(*TGoods)
		tGoods := &GoodsManagerData{
			GoodsId:     queryRep.Rows[i].(*TGoods).Goods_uuid,
			GoogdsName:  queryRep.Rows[i].(*TGoods).Goods_name,
			VarietyId:   queryRep.Rows[i].(*TGoods).Variety_uuid.String,
			GoodsBrief:  queryRep.Rows[i].(*TGoods).Goods_desc.String,
			GoodsPrice:  queryRep.Rows[i].(*TGoods).Goods_price,
			ChargeUnit:  queryRep.Rows[i].(*TGoods).Charge_unit,
			GoodsWeight: queryRep.Rows[i].(*TGoods).Goods_weight,
			WeightUnit:  queryRep.Rows[i].(*TGoods).Weight_unit,
			GoodsCode:   queryRep.Rows[i].(*TGoods).Goods_bar_code.String,
			FactoryId:   queryRep.Rows[i].(*TGoods).Factory_uuid.String,
			pictureList: []GoodsManagerDataPicture{},
		}
		args2 := []interface{}{}
		args2 = append(args2, goods.Goods_uuid)
		tmp2 := TGoodsPicture{}
		queryReq2 := &SqlQueryRequest{
			SQL:         "select picture_uuid, picture_path, show_order, picture_name, picture_desc from t_goods_picture where goods_uuid = ?",
			Args:        args2,
			RowTemplate: tmp2}
		reply2 := m.d.dbCli.Query(queryReq2)
		queryRep2, _ := reply2.(*SqlQueryReply)
		if queryRep2.Err != nil {
			zap.L().Error(fmt.Sprintf("query goods error:%s", queryRep2.Err.Error()))
			return nil, queryRep2.Err
		}
		for j := 0; j < len(queryRep2.Rows); j++ {
			goodsPicture := queryRep.Rows[i].(*TGoodsPicture)
			goodsManagerDataPicture := GoodsManagerDataPicture{
				PictureId:    goodsPicture.Picture_uuid,
				PicturePath:  goodsPicture.Picture_path,
				PictureOrder: goodsPicture.Show_order,
				PictureName:  goodsPicture.Picture_name.String,
				PictureDesc:  goodsPicture.Picture_desc.String,
			}
			tGoods.pictureList = append(tGoods.pictureList, goodsManagerDataPicture)
		}
		returngoodsManagerData = append(returngoodsManagerData, tGoods)
	}
	return returngoodsManagerData, nil
}
