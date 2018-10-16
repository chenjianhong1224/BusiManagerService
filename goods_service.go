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
	args1 = append(args1, opId)
	args1 = append(args1, opId)
	execReq1 := SqlExecRequest{
		SQL:  "insert into T_Goods(goods_uuid, goods_name, goods_bar_code, factory_uuid, goods_price, charge_unit, goods_weight, weight_unit, goods_desc, goods_status, whole_pack, pack_unit, variety_uuid, create_time, create_user, update_time, update_user) values(?,?,?,?,?,?,?,?,'',1,NULL,NULL,?,now(),?,now(),?)",
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
	}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return uid.String(), nil
	}
	zap.L().Error(fmt.Sprintf("add goods[%s] error:%s", req.GoogdsName, err.Error()))
	return "", err
}

func (m *goods_service) updateGoods(req GoodsManagerData, opId string) error {
	return nil
}
