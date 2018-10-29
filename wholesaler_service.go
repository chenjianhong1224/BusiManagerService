package main

import (
	"crypto/md5"
	"fmt"
	//	"math/rand"
	//	"time"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type wholesaler_service struct {
	d *dbOperator
}

func (m *wholesaler_service) addWholesaler(req WholesalerManagerData, opId string) (ud string, e error) {
	passwd := req.LinkPhone[len(req.LinkPhone)-7 : len(req.LinkPhone)]
	args1 := []interface{}{}
	uid, _ := uuid.NewV4()
	args1 = append(args1, uid.String())
	//手机号为用户名
	args1 = append(args1, req.LinkPhone)
	data := []byte(passwd)
	has := md5.Sum(data)
	args1 = append(args1, fmt.Sprintf("%x", has))
	args1 = append(args1, uid.String())
	args1 = append(args1, opId)
	args1 = append(args1, opId)
	execReq1 := SqlExecRequest{
		SQL:  "insert into t_user(User_uuid, User_name, Passwd, Open_id, Other_from, Nickname, Head_portrait, Agent_uuid, User_type, User_status, User_token, Expiry_time, Create_time, Create_user, Update_time, Update_user, Remark) values(?,?,?,NULL,NULL,NULL,NULL,?,0,1,NULL,NULL,now(),?,now(),?,NULL)",
		Args: args1,
	}
	args2 := []interface{}{}
	args2 = append(args2, uid.String())
	args2 = append(args2, req.WholesalerName)
	args2 = append(args2, req.Company)
	args2 = append(args2, req.LinkPhone)
	args2 = append(args2, opId)
	args2 = append(args2, opId)
	execReq2 := SqlExecRequest{
		SQL:  "insert into t_wholesaler(Saler_uuid, Saler_name, Company, Mobile, Saler_status, Create_time, Create_user, Update_time, Update_user, Remark, Salesperson) values(?,?,?,?,1,now(),?,now(),?,NULL,NULL)",
		Args: args2,
	}
	var execReqList = []SqlExecRequest{execReq2, execReq1}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return uid.String(), nil
	}
	zap.L().Error(fmt.Sprintf("add wholesaler[%s,%s] error:%s", req.WholesalerName, req.LinkPhone, err.Error()))
	return "", err
}

func (m *wholesaler_service) updateWholesaler(req WholesalerManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, req.WholesalerName)
	args1 = append(args1, req.Company)
	args1 = append(args1, req.LinkPhone)
	args1 = append(args1, opId)
	args1 = append(args1, req.WholesalerId)
	execReq1 := SqlExecRequest{
		SQL:  "update t_wholesaler  set Saler_name=?, Company = ?, Mobile =?, Update_time=now(), update_user = ? where Saler_uuid = ?",
		Args: args1}
	args2 := []interface{}{}
	args2 = append(args2, req.LinkPhone)
	args2 = append(args2, opId)
	args2 = append(args2, req.WholesalerId)
	args2 = append(args2, req.WholesalerId)
	execReq2 := SqlExecRequest{
		SQL:  "update t_user  set User_name=?, Update_time=now(), update_user = ? where User_uuid = ? and Agent_uuid = ?",
		Args: args2}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err != nil {
		zap.L().Error(fmt.Sprintf("update Wholesaler[%s] error:%s", req.WholesalerId, err.Error()))
		return err
	}
	return nil
}

func (m *wholesaler_service) deleteWholesaler(req WholesalerManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, opId)
	args1 = append(args1, req.WholesalerId)
	execReq1 := SqlExecRequest{
		SQL:  "update t_wholesaler  set Saler_status=0, Update_time=now(), update_user = ? where Saler_uuid = ?",
		Args: args1}
	args2 := []interface{}{}
	args2 = append(args2, opId)
	args2 = append(args2, req.WholesalerId)
	args2 = append(args2, req.WholesalerId)
	execReq2 := SqlExecRequest{
		SQL:  "update t_user  set User_status=0, Update_time=now(), update_user = ? where User_uuid = ? and Agent_uuid = ?",
		Args: args2}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err != nil {
		zap.L().Error(fmt.Sprintf("delete Wholesaler[%s] error:%s", req.WholesalerId, err.Error()))
		return err
	}
	return nil
}

func (m *wholesaler_service) queryWholesalerByExample(req WholesalerManagerData) ([]*TWholeSaler, error) {
	args := []interface{}{}
	var sql string
	sql = "select Saler_id, Saler_uuid, Saler_name, Company, Mobile, Saler_status, Create_time, Create_user, Update_time, Update_user, Remark, Salesperson from t_wholesaler where 1=1 "
	if len(req.Company) > 0 {
		args = append(args, req.Company)
		sql += " and Company = ? "
	}
	if len(req.LinkPhone) > 0 {
		args = append(args, req.LinkPhone)
		sql += " and Mobile = ? "
	}
	if len(req.WholesalerId) > 0 {
		args = append(args, req.WholesalerId)
		sql += " and Saler_uuid = ? "
	}
	if len(req.WholesalerName) > 0 {
		args = append(args, req.WholesalerName)
		sql += " and Saler_name = ? "
	}
	tmp := TWholeSaler{}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query wholesaler error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	var returnMenus []*TWholeSaler = []*TWholeSaler{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnMenus = append(returnMenus, queryRep.Rows[i].(*TWholeSaler))
	}
	return returnMenus, nil
}
