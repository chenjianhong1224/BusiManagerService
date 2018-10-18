package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type factory_service struct {
	d *dbOperator
}

func (m *factory_service) addFactory(req FactoryManagerData) (string, error) {
	args1 := []interface{}{}
	uid, _ := uuid.NewV4()
	args1 = append(args1, uid.String())
	args1 = append(args1, req.FactoryName)
	args1 = append(args1, req.LinkPerson)
	args1 = append(args1, req.LinkPhone)
	args1 = append(args1, req.FactoryAddress)
	execReq1 := SqlExecRequest{
		SQL:  "insert into T_Factory(Factory_uuid,Factory_name,Link_person,Link_phone,Factory_desc,Factory_address,Factory_status) values(?,?,?,?,'',?,1)",
		Args: args1,
	}
	var execReqList = []SqlExecRequest{execReq1}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return uid.String(), nil
	}
	zap.L().Error(fmt.Sprintf("add factory[%s] error:%s", req.FactoryName, err.Error()))
	return "", err
}

func (m *factory_service) updateFactory(req FactoryManagerData) error {
	args1 := []interface{}{}
	args1 = append(args1, req.FactoryName)
	args1 = append(args1, req.LinkPerson)
	args1 = append(args1, req.LinkPhone)
	args1 = append(args1, req.FactoryAddress)
	args1 = append(args1, req.FactoryId)
	execReq1 := SqlExecRequest{
		SQL:  "update T_Factory set Factory_name = ?,Link_person = ?,Link_phone = ?, Factory_address = ? where factory_uuid = ?",
		Args: args1,
	}
	var execReqList = []SqlExecRequest{execReq1}
	err := m.d.dbCli.TransationExcute(execReqList)
	return err
}

func (m *factory_service) deleteFactory(req FactoryManagerData) error {
	args1 := []interface{}{}
	args1 = append(args1, req.FactoryId)
	execReq1 := SqlExecRequest{
		SQL:  "update T_Factory factory_status = 0 where factory_uuid = ?",
		Args: args1,
	}
	var execReqList = []SqlExecRequest{execReq1}
	err := m.d.dbCli.TransationExcute(execReqList)
	return err
}

func (m *factory_service) queryFactoryByExample(req FactoryManagerData) ([]*TFactory, error) {
	args := []interface{}{}
	var sql string
	sql = "select factory_id, Factory_uuid,Factory_name,Link_person,Link_phone,Factory_desc,Factory_address,Factory_status from T_Factory where 1=1 "
	if len(req.FactoryId) > 0 {
		sql += " and factory_uuid = ?"
		args = append(args, req.FactoryId)
	}
	if len(req.FactoryName) > 0 {
		sql += " and factory_name =?"
		args = append(args, req.FactoryName)
	}
	tmp := TFactory{}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query factory error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	var returnTFactory []*TFactory = []*TFactory{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnTFactory = append(returnTFactory, queryRep.Rows[i].(*TFactory))
	}
	return returnTFactory, nil
}
