package main

import (
	"crypto/md5"
	"fmt"
	//	"math/rand"
	//	"time"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type salseman_service struct {
	d *dbOperator
}

func (m *salseman_service) addSalseman(req SalsemanManagerData, opId string) (ud string, e error) {
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
	args2 = append(args2, req.WholesalerId)
	args2 = append(args2, req.SalesmanName)
	args2 = append(args2, req.LinkPhone)
	execReq2 := SqlExecRequest{
		SQL:  "insert into t_Salseman(Salesman_uuid, Saler_uuid, Salesman_name, Salesman_phone, Entry_time, Departure_time, Salesman_status, Remark) values(?,?,?,?,NULL,NULL,1,NULL)",
		Args: args2,
	}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return uid.String(), nil
	}
	zap.L().Error(fmt.Sprintf("add Salseman[%s,%s] error:%s", req.SalesmanName, req.LinkPhone, err.Error()))
	return "", err
}

func (m *salseman_service) updateSalseman(req SalsemanManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, req.SalesmanName)
	args1 = append(args1, req.WholesalerId)
	args1 = append(args1, req.LinkPhone)
	args1 = append(args1, req.SalesmanId)
	execReq1 := SqlExecRequest{
		SQL:  "update t_Salseman  set Salesman_name=?, Saler_uuid = ?, Salesman_phone =? where Salesman_uuid = ?",
		Args: args1}
	args2 := []interface{}{}
	args2 = append(args2, req.LinkPhone)
	args2 = append(args2, opId)
	args2 = append(args2, req.SalesmanId)
	args2 = append(args2, req.SalesmanId)
	execReq2 := SqlExecRequest{
		SQL:  "update t_user  set User_name=?, Update_time=now(), update_user = ? where User_uuid = ? and Agent_uuid = ?",
		Args: args2}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err != nil {
		zap.L().Error(fmt.Sprintf("update Salseman[%s] error:%s", req.SalesmanId, err.Error()))
		return err
	}
	return nil
}

func (m *salseman_service) deleteSalseman(req SalsemanManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, req.SalesmanId)
	execReq1 := SqlExecRequest{
		SQL:  "update t_Salseman  set Salesman_status=0 where Salesman_uuid = ?",
		Args: args1}
	args2 := []interface{}{}
	args2 = append(args2, opId)
	args2 = append(args2, req.SalesmanId)
	args2 = append(args2, req.SalesmanId)
	execReq2 := SqlExecRequest{
		SQL:  "update t_user  set User_status=0, Update_time=now(), update_user = ? where User_uuid = ? and Agent_uuid = ?",
		Args: args2}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err != nil {
		zap.L().Error(fmt.Sprintf("delete Salseman[%s] error:%s", req.SalesmanId, err.Error()))
		return err
	}
	return nil
}

func (m *salseman_service) querySalsemanByExample(req SalsemanManagerData) ([]*TSalseman, error) {
	args := []interface{}{}
	var sql string
	sql = "select Salesman_id, Salesman_uuid, Saler_uuid, Salesman_name, Salesman_phone, Entry_time, Departure_time, Salesman_status, Remark from t_Salseman where 1=1 "
	if len(req.SalesmanName) > 0 {
		args = append(args, req.SalesmanName)
		sql += " and Salesman_name = ? "
	}
	if len(req.LinkPhone) > 0 {
		args = append(args, req.LinkPhone)
		sql += " and Salesman_phone = ? "
	}
	if len(req.WholesalerId) > 0 {
		args = append(args, req.WholesalerId)
		sql += " and Saler_uuid = ? "
	}
	if len(req.SalesmanId) > 0 {
		args = append(args, req.SalesmanId)
		sql += " and Salesman_uuid = ? "
	}
	tmp := TSalseman{}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query Salseman error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	var returnMenus []*TSalseman = []*TSalseman{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnMenus = append(returnMenus, queryRep.Rows[i].(*TSalseman))
	}
	return returnMenus, nil
}
