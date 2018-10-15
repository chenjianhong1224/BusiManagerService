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
	variety_id     int64
	variety_uuid   string
	variety_name   string
	variety_ico    sql.NullString
	variety_status int32
	create_time    mysql.NullTime
	create_user    sql.NullString
	update_time    mysql.NullTime
	update_user    sql.NullString
	remark         sql.NullString
}

type TGoods struct {
	goods_id       int64
	goods_uuid     string
	goods_name     string
	goods_bar_code sql.NullString
	factory_uuid   sql.NullString
	goods_price    int32
	charge_unit    int32
	goods_weight   int32
	weight_unit    int32
	goods_desc     sql.NullString
	goods_status   int32
	whole_pack     int32
	pack_unit      int32
	variety_uuid   sql.NullString
	create_time    mysql.NullTime
	create_user    sql.NullString
	update_time    mysql.NullTime
	update_user    sql.NullString
}

type TGoodsPicture struct {
	picture_id     int64
	picture_uuid   string
	picture_name   sql.NullString
	picture_path   string
	goods_uuid     sql.NullString
	show_order     int32
	picture_desc   sql.NullString
	picture_status int32
	size_variety   int32
}

type TFactory struct {
	facotry_id      int64
	factory_uuid    string
	factory_name    string
	link_person     sql.NullString
	link_phone      sql.NullString
	factory_desc    sql.NullString
	factory_address sql.NullString
	factory_status  int32
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
	banner_id     int64
	saler_uuid    string
	banner_uuid   string
	banner_name   sql.NullString
	banner_pic    string
	show_order    int32
	banner_status int32
	open_time     mysql.NullTime
	close_time    mysql.NullTime
	link_uri      sql.NullString
}

type TSalseman struct {
	salesman_id     int64
	salesman_uuid   string
	saler_uuid      string
	salesman_name   string
	salesman_phone  string
	entry_time      string
	departure_time  mysql.NullTime
	salesman_status int32
	remark          sql.NullString
}
