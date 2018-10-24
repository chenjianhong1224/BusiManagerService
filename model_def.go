package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type TWxpluginProgram struct {
	Program_id     int64
	Program_uuid   string
	Program_name   string
	Appid          string
	Appsecrete     string
	Program_status int32
	Saler_uuid     sql.NullString
	Program_type   int32
}

type TGoodsVariety struct {
	Variety_id     int64
	Variety_uuid   string
	Variety_name   string
	Variety_ico    sql.NullString
	Variety_status int32
	Create_time    mysql.NullTime
	Create_user    sql.NullString
	Update_time    mysql.NullTime
	Update_user    sql.NullString
	Remark         sql.NullString
}

type TGoods struct {
	Goods_id       int64
	Goods_uuid     string
	Goods_name     string
	Goods_bar_code sql.NullString
	Factory_uuid   sql.NullString
	Goods_price    int32
	Charge_unit    int32
	Goods_weight   int32
	Weight_unit    int32
	Goods_desc     sql.NullString
	Goods_status   int32
	Whole_pack     int32
	Pack_unit      int32
	Variety_uuid   sql.NullString
	Create_time    mysql.NullTime
	Create_user    sql.NullString
	Update_time    mysql.NullTime
	Update_user    sql.NullString
	Goods_picture  sql.NullString
}

type TGoodsPicture struct {
	Picture_id     int64
	Picture_uuid   string
	Picture_name   sql.NullString
	Picture_path   string
	Goods_uuid     sql.NullString
	Show_order     int32
	Picture_desc   sql.NullString
	Picture_status int32
	Size_variety   int32
}

type TFactory struct {
	Facotry_id      int64
	Factory_uuid    string
	Factory_name    string
	Link_person     sql.NullString
	Link_phone      sql.NullString
	Factory_desc    sql.NullString
	Factory_address sql.NullString
	Factory_status  int32
}

type TWholeSaler struct {
	Saler_id     int64
	Saler_uuid   string
	Saler_name   sql.NullString
	Company      string
	Mobile       string
	Saler_status int32
	Create_time  mysql.NullTime
	Create_user  sql.NullString
	Update_time  mysql.NullTime
	Update_user  sql.NullString
	Remark       sql.NullString
	Salesperson  sql.NullString
}

type TWholesalerBanner struct {
	Banner_id     int64
	Saler_uuid    string
	Banner_uuid   string
	Banner_name   sql.NullString
	Banner_pic    string
	Show_order    int32
	Banner_status int32
	Open_time     mysql.NullTime
	Close_time    mysql.NullTime
	Link_uri      sql.NullString
}

type TSalseman struct {
	Salesman_id     int64
	Salesman_uuid   string
	Saler_uuid      string
	Salesman_name   string
	Salesman_phone  string
	Entry_time      string
	Departure_time  mysql.NullTime
	Salesman_status int32
	Remark          sql.NullString
}

type TWholeSalerMember struct {
	member_id     int64
	member_uuid   string
	saler_uuid    sql.NullString
	salesman_uuid sql.NullString
	member_name   sql.NullString
	mobile        string
	member_status int32
	open_id       string
	other_from    sql.NullString
	member_bonus  int32
	create_time   mysql.NullTime
	create_user   sql.NullString
	update_time   mysql.NullTime
	update_user   sql.NullString
	remark        sql.NullString
}
