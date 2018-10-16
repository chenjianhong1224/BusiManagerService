package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type wxpluginProgram_service struct {
	d *dbOperator
}

func (m *wxpluginProgram_service) addWxpluginProgram(req WxpluginProgramManagerData) (string, error) {
	args1 := []interface{}{}
	uid, _ := uuid.NewV4()
	args1 = append(args1, uid.String())
	args1 = append(args1, req.ProgName)
	args1 = append(args1, req.AppId)
	args1 = append(args1, req.AppSecrete)
	args1 = append(args1, req.WsId)
	args1 = append(args1, req.ProgramType)

	queryReq := &SqlExecRequest{
		SQL:  "insert into t_Wxplugin_program(Program_uuid, Program_name, Appid, Appsecrete, Program_status, Saler_uuid, Program_type) values(?,?,?,?,1,?,?)",
		Args: args1}
	excuteRep := m.d.dbCli.Query(queryReq)
	if excuteRep.Error() != nil {
		zap.L().Error(fmt.Sprintf("add wxpluginProgram[%s,%s] error:%s", req.ProgName, req.AppId, excuteRep.Error()))
		return "", excuteRep.Error()
	}
	return uid.String(), nil
}

func (m *wxpluginProgram_service) updateWxpluginProgram(req WxpluginProgramManagerData) error {
	args1 := []interface{}{}
	args1 = append(args1, req.ProgName)
	args1 = append(args1, req.AppId)
	args1 = append(args1, req.AppSecrete)
	args1 = append(args1, req.WsId)
	args1 = append(args1, req.ProgramType)
	args1 = append(args1, req.ProgId)
	queryReq := &SqlExecRequest{
		SQL:  "update t_Wxplugin_program  set Program_name = ?, Appid = ?, Appsecrete = ?, Saler_uuid = ?, Program_type = ? where Program_uuid = ?",
		Args: args1}
	excuteRep := m.d.dbCli.Query(queryReq)
	if excuteRep.Error() != nil {
		zap.L().Error(fmt.Sprintf("update wxpluginProgram[%s ] error:%s", req.ProgId, excuteRep.Error()))
		return excuteRep.Error()
	}
	return nil
}

func (m *wxpluginProgram_service) deleteWxpluginProgram(req WxpluginProgramManagerData) error {
	args1 := []interface{}{}
	args1 = append(args1, req.ProgId)
	queryReq := &SqlExecRequest{
		SQL:  "update t_Wxplugin_program  set Program_status = 0 where Program_uuid = ?",
		Args: args1}
	excuteRep := m.d.dbCli.Query(queryReq)
	if excuteRep.Error() != nil {
		zap.L().Error(fmt.Sprintf("delete wxpluginProgram[%s ] error:%s", req.ProgId, excuteRep.Error()))
		return excuteRep.Error()
	}
	return nil
}

func (m *wxpluginProgram_service) queryWxpluginProgramByExample(req WxpluginProgramManagerData) ([]*TWxpluginProgram, error) {
	args1 := []interface{}{}
	var sql string
	sql = "select Program_id, Program_uuid, Program_name, Appid, Appsecrete, Program_status, Saler_uuid, Program_type from t_Wxplugin_program where 1=1 "
	if len(req.AppId) != 0 {
		sql += " and Appid = ?"
		args1 = append(args1, req.AppId)
	}
	if len(req.AppSecrete) != 0 {
		sql += " and Appsecrete = ?"
		args1 = append(args1, req.AppSecrete)
	}
	if len(req.ProgId) != 0 {
		sql += " and Program_uuid = ?"
		args1 = append(args1, req.ProgId)
	}
	if len(req.ProgName) != 0 {
		sql += " and Program_name = ?"
		args1 = append(args1, req.ProgName)
	}
	sql += " and Program_type = ?"
	args1 = append(args1, req.ProgramType)
	if len(req.WsId) != 0 {
		sql += " and Saler_uuid = ?"
		args1 = append(args1, req.WsId)
	}
	tmp := TWxpluginProgram{}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args1,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query t_Wxplugin_program error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	var returnMenus []*TWxpluginProgram = []*TWxpluginProgram{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnMenus = append(returnMenus, queryRep.Rows[i].(*TWxpluginProgram))
	}
	return returnMenus, nil
}
