package main

import (
	"crypto/md5"
	"fmt"
	//	"math/rand"
	//	"time"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type wholesaler_member_service struct {
	d *dbOperator
}

func (m *wholesaler_member_service) addWholesalerMember(req WholesalerMemberManagerData, opId string) (ud string, e error) {
	passwd := ""
	args1 := []interface{}{}
	uid, _ := uuid.NewV4()
	args1 = append(args1, uid.String())
	//手机号为用户名
	args1 = append(args1, req.MemberName)
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
	args2 = append(args2, req.SalesmanId)
	args2 = append(args2, req.LinkPhone)
	args2 = append(args2, req.MemberName)
	args2 = append(args2, req.LinkPhone)
	args2 = append(args2, opId)
	args2 = append(args2, opId)
	execReq2 := SqlExecRequest{
		SQL:  "insert into t_wholesaler_member(member_uuid, saler_uuid, salesman_uuid, member_name, mobile, member_status, open_id, other_from, member_bonus, create_time, create_user, update_time, update_user, remark) values(?,?,?,?,?,1,NULL,NULL,NULL,now(),?,now(),?,NULL)",
		Args: args2,
	}
	var execReqList = []SqlExecRequest{execReq2, execReq1}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return uid.String(), nil
	}
	zap.L().Error(fmt.Sprintf("add WholesalerMember[%s] error:%s", req.MemberName, err.Error()))
	return "", err
}

func (m *wholesaler_member_service) updateWholesalerMember(req WholesalerMemberManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, req.WholesalerId)
	args1 = append(args1, req.SalesmanId)
	args1 = append(args1, req.MemberName)
	args1 = append(args1, req.LinkPhone)
	args1 = append(args1, opId)
	args1 = append(args1, req.MemberId)
	execReq1 := SqlExecRequest{
		SQL:  "update t_wholesaler_member  set saler_uuid=?, salesman_uuid = ?, member_name=?, Mobile =?, Update_time=now(), update_user = ? where member_uuid = ?",
		Args: args1}
	args2 := []interface{}{}
	args2 = append(args2, req.MemberName)
	args2 = append(args2, opId)
	args2 = append(args2, req.MemberId)
	args2 = append(args2, req.MemberId)
	execReq2 := SqlExecRequest{
		SQL:  "update t_user  set User_name=?, Update_time=now(), update_user = ? where User_uuid = ? and Agent_uuid = ?",
		Args: args2}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err != nil {
		zap.L().Error(fmt.Sprintf("update WholesalerMemberMember[%s] error:%s", req.MemberId, err.Error()))
		return err
	}
	return nil
}

func (m *wholesaler_member_service) deleteWholesalerMember(req WholesalerMemberManagerData, opId string) error {
	args1 := []interface{}{}
	args1 = append(args1, opId)
	args1 = append(args1, req.MemberId)
	execReq1 := SqlExecRequest{
		SQL:  "update t_wholesaler_member  set member_status=0, Update_time=now(), update_user = ? where Saler_uuid = ?",
		Args: args1}
	args2 := []interface{}{}
	args2 = append(args2, opId)
	args2 = append(args2, req.MemberId)
	args2 = append(args2, req.MemberId)
	execReq2 := SqlExecRequest{
		SQL:  "update t_user  set User_status=0, Update_time=now(), update_user = ? where User_uuid = ? and Agent_uuid = ?",
		Args: args2}
	var execReqList = []SqlExecRequest{execReq1, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err != nil {
		zap.L().Error(fmt.Sprintf("delete WholesalerMember[%s] error:%s", req.MemberId, err.Error()))
		return err
	}
	return nil
}

func (m *wholesaler_member_service) queryWholesalerMemberByExample(req WholesalerMemberManagerData) ([]*TWholeSalerMember, error) {
	args := []interface{}{}
	var sql string
	sql = "select member_id, member_uuid, saler_uuid, salesman_uuid, member_name, mobile, member_status, open_id, other_from, member_bonus, create_time, create_user, update_time, update_user, remark from t_wholesaler_member where 1=1 "
	if len(req.MemberId) > 0 {
		args = append(args, req.MemberId)
		sql += " and member_uuid = ? "
	}
	if len(req.MemberName) > 0 {
		args = append(args, req.MemberName)
		sql += " and member_name = ? "
	}
	if len(req.SalesmanId) > 0 {
		args = append(args, req.SalesmanId)
		sql += " and salesman_uuid = ? "
	}
	if len(req.WholesalerId) > 0 {
		args = append(args, req.WholesalerId)
		sql += " and saler_uuid = ? "
	}
	tmp := TWholeSalerMember{}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query WholesalerMember error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	var returnMembers []*TWholeSalerMember = []*TWholeSalerMember{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnMembers = append(returnMembers, queryRep.Rows[i].(*TWholeSalerMember))
	}
	return returnMembers, nil
}
