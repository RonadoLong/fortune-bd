-- we don't know how to generate schema shop (class Schema) :(
create table goods
(
	goods_id bigint auto_increment
		primary key,
	merchant_id bigint null comment '厂商ID',
	code varchar(16) null comment '商品编号',
	category_id int null comment '商品类别ID',
	title varchar(128) not null comment '标题',
	en_title varchar(128) null comment '英文标题',
	sell_point varchar(255) null comment '买点',
	en_sell_point varchar(255) null comment '英文： 买点',
	tag_id int null comment '标签',
	goods_images varchar(255) not null comment '商品列表小图',
	has_activity int(1) default '0' null comment '是否有活动',
	goods_type int(3) default '1' not null comment '商品类型 1: 单品 2：
	一种规格',
	status int(1) default '1' null comment '0下架 1.上架 2.卖完',
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create table goods_class
(
	id int auto_increment
		primary key,
	code varchar(32) null comment '分类编码',
	parentId int null,
	name varchar(8) null comment '标题',
	enName varchar(8) null comment '英文标题',
	imgUrl varchar(100) null,
	sort int(2) null comment '排序',
	isParentId int(1) default '0' null,
	status int(1) default '1' null comment '是否显示',
	updateTime timestamp default CURRENT_TIMESTAMP not null,
	createTime timestamp default CURRENT_TIMESTAMP not null
)
comment '商品分类'
;

create table goods_detail
(
	detail_id bigint auto_increment
		primary key,
	goods_id bigint not null,
	goods_banners mediumtext not null comment '详情轮播图',
	goods_detail text not null comment '详情图',
	goods_desc text null comment '商品描述',
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null
)
;

create index goods_detail_INDEX
	on goods_detail (goods_id)
;

create table goods_nav
(
	class_id int auto_increment
		primary key,
	title varchar(10) not null,
	en_title varchar(20) not null,
	sort int(1) null,
	status tinyint(1) default '1' not null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null
)
;

create table goods_sku
(
	sku_id bigint auto_increment
		primary key,
	goods_id bigint not null comment '商品ID',
	member_price int not null comment '会员价',
	price int not null comment '原价',
	low_price int null comment '进货价格',
	activity_price int null comment '活动价格',
	discount_price int null comment '折扣价格',
	is_active int(1) default '0' null comment '默认选中',
	sold_count int(8) default '0' null comment '销量',
	stock int not null comment '库存',
	lock_stock int default '0' null comment '锁定库存',
	commission int null comment '佣金',
	integral int null comment '积分',
	sort int(1) null comment '属性排序',
	sku_type int(1) default '0' null comment '0 为单品 1 为一种规格  2 为两种规格',
	status int(1) default '1' null comment '1.上架，0.下架',
	update_time timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	buy_num_max int(3) null comment '最大购买数',
	sku_pic varchar(100) null
)
;

create table goods_sku_property
(
	id bigint auto_increment
		primary key,
	sku_id bigint not null,
	sku_value varchar(30) null,
	en_sku_value varchar(30) null,
	sku_name varchar(30) null,
	en_sku_name varchar(30) null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null
)
;

create index idx_goods_sku_property
	on goods_sku_property (sku_id)
;

create table goods_stock_flow
(
	id bigint auto_increment
		primary key,
	order_id varchar(32) null,
	sku_id bigint not null,
	stock_before int null,
	stock_after int null,
	stock_change int null,
	lock_stock_before int null,
	lock_stock_after int null,
	lock_stock_change int null,
	check_status int(1) default '0' null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create index idx_stock_flow_orderId
	on goods_stock_flow (order_id, check_status)
;

create index idx_stock_flow_skuId
	on goods_stock_flow (sku_id, check_status)
;

create table goods_tag
(
	id int auto_increment
		primary key,
	name varchar(10) null,
	status int(1) default '1' null,
	createTime datetime null,
	updateTime datetime null
)
;

create table home_carousel
(
	id int auto_increment
		primary key,
	title varchar(100) null,
	url text null,
	img_url varchar(255) null,
	sort int(1) null,
	status int(1) default '1' null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null
)
comment '首页轮播图'
;

create table home_goods
(
	id bigint not null
		primary key,
	type int(1) default '1' null comment '类型 1 普通商品 2 秒杀',
	status tinyint default '1' not null
)
comment '首页商品表'
;

create table home_nav
(
	id int auto_increment
		primary key,
	title varchar(10) not null,
	en_title varchar(10) null,
	img_url varchar(255) not null,
	type int(1) null,
	jump_url varchar(100) null,
	sort int(1) null,
	status int(1) default '1' null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null
)
comment '首页导航'
;

create table home_news
(
	id bigint auto_increment comment '自增主键'
		primary key,
	title varchar(100) default '' null comment '标题',
	thumb_url varchar(200) null comment '文章缩略图',
	author varchar(100) null comment '作者名称',
	avatar varchar(200) null comment '作者头像url',
	source_type_name varchar(100) null comment '文章来源类型名称',
	read_count int default '0' null comment '阅读量',
	comment_count int default '0' null comment '评论数',
	like_count int default '0' null comment '点赞数',
	category varchar(10) null,
	view_type int(1) null comment '显示类型 1单图 2多图 34',
	is_recommend int(1) default '0' null comment '是否推荐',
	content longtext null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null,
	status int(1) default '1' null comment '是否可用'
)
comment '资讯文章'
;

create table news
(
	id bigint auto_increment comment '自增主键'
		primary key,
	title varchar(100) default '' null comment '标题',
	thumb_url text not null comment '文章缩略图',
	author varchar(100) null comment '作者名称',
	avatar varchar(200) null comment '作者头像url',
	source_type_name varchar(100) null comment '文章来源类型名称',
	read_count int default '0' null comment '阅读量',
	comment_count int default '0' null comment '评论数',
	like_count int default '0' null comment '点赞数',
	category varchar(10) null,
	view_type int(1) null comment '显示类型 1单图 2多图 34',
	is_recommend int(1) default '0' null comment '是否推荐',
	content longtext null,
	create_time timestamp default CURRENT_TIMESTAMP null,
	update_time timestamp default CURRENT_TIMESTAMP null,
	status int(1) default '1' null comment '是否可用'
)
comment '资讯文章'
;

create table news_category
(
	id int auto_increment
		primary key,
	title varchar(10) not null,
	sort int default '0' not null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null,
	status int(1) null
)
;

create index idx_news_category_title
	on news_category (title, status)
;

create table order_after_sale
(
	id bigint auto_increment
		primary key,
	orderId varchar(32) null comment '订单id',
	userId varchar(32) null comment '申请人编号',
	dealUserId varchar(64) null comment '处理人用户ID',
	dealUserName varchar(16) null comment '处理人用户名称',
	dealStatus int(1) null comment '0待审核 1审核通过、2审核不通过',
	createTime timestamp default CURRENT_TIMESTAMP null,
	updateTime timestamp default CURRENT_TIMESTAMP null
)
comment '订单售后'
;

create table order_deal_spoor
(
	id bigint auto_increment
		primary key,
	orderId varchar(32) null comment '订单id',
	dealUserId varchar(64) null comment '处理人用户ID',
	dealUserName varchar(16) null comment '处理人用户名称',
	createTime timestamp default CURRENT_TIMESTAMP null,
	updateTime timestamp default CURRENT_TIMESTAMP null
)
comment '订单处理记录'
;

create table order_goods
(
	id bigint auto_increment
		primary key,
	order_id varchar(32) null comment '订单id',
	product_id bigint null comment '商品id',
	goods_title varchar(128) null comment '商品标题',
	goods_price int null comment '商品价格',
	total_price int not null comment '商品总价',
	goods_image varchar(255) null comment '商品小图',
	goods_count int null comment '商品数量',
	goods_number varchar(32) null comment '商品编号',
	sku_id bigint null comment '属性id',
	sku_values varchar(100) null comment '属性名字',
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null
)
comment '订单商品表，存放一个订单所包含的所有商品'
;

create index idx_goods_orderId
	on order_goods (order_id, product_id)
;

create table order_info
(
	order_id varchar(32) not null comment '订单ID'
		primary key,
	user_id varchar(64) not null comment '收货人id',
	username varchar(20) not null comment '收货人',
	merchant_id varchar(32) null comment '商家id',
	total_integral varchar(20) not null,
	order_address varchar(255) not null comment '订单地址',
	order_type int(1) default '0' null comment '订单类型: 0 普通',
	is_post_fee int(1) default '0' null comment '是否包邮',
	post_fee int null comment '邮费。精确到2位小数;单位:元。如:200.07，表示:200元7分',
	coupon_id varchar(32) null comment '优惠券',
	coupon_paid int null comment '优惠券金额',
	goods_count int not null comment '商品数量',
	total_amount int not null comment '总金额',
	really_amount int not null comment '实际支付金额',
	order_identifier varchar(32) null comment '订单编码',
	shipping_name varchar(20) null comment '物流名称',
	shipping_code varchar(20) null comment '物流单号',
	buyer_msg varchar(100) null comment '买家留言',
	payment_time timestamp null comment '付款时间',
	total_settlement_price int null comment '订单结算总价',
	buyer_rate int(1) default '0' null comment '买家是否已经评价',
	pay_type varchar(20) not null comment '支付类型 1 pay pal  2 信用卡',
	order_status int(1) not null comment '状态：1未确认 2已确认 3退款 4交易成功(已收货) 5交易关闭 6无效',
	shipping_status int(1) not null comment '发货状态 1未发货 2已发货 3已收货',
	pay_status int(1) not null comment '支付状态 1未支付 2支付中 3已支付',
	create_time timestamp default CURRENT_TIMESTAMP not null comment '订单创建时间',
	update_time timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '订单更新时间',
	consign_time timestamp null
)
comment '订单主表'
;

create index idx_buyer_userId
	on order_info (user_id)
;

create index idx_createTime
	on order_info (create_time)
;

create index idx_order_payType
	on order_info (pay_type)
;

create index idx_order_status
	on order_info (order_status)
;

create table order_pay
(
	pay_id varchar(32) not null comment '支付编号'
		primary key,
	order_id varchar(32) null comment '订单编号',
	pay_amount int null comment '支付金额',
	is_paid char null comment '是否已支付 0否 1是'
)
;

create table product
(
	product_id bigint auto_increment
		primary key,
	merchant_id bigint null comment '厂商ID',
	code varchar(16) null comment '商品编号',
	category_id int null comment '商品类别ID',
	title varchar(128) not null comment '标题',
	en_title varchar(128) null comment '英文标题',
	sell_point varchar(255) null comment '买点',
	en_sell_point varchar(255) null comment '英文： 买点',
	postage int null comment '邮费',
	goods_images varchar(255) not null comment '商品列表小图',
	member_price int not null,
	has_activity int(1) default '0' null comment '是否有活动',
	price int not null,
	stock int not null,
	low_price int null,
	commission int null,
	activity_price int null,
	integral int null,
	status int(1) default '1' null comment '0下架 1.上架 2.卖完',
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
	sold_count int null,
	buy_max int null
)
;

create table service_category
(
	id int auto_increment
		primary key,
	name varchar(20) null,
	en_name varchar(20) null,
	imageurl varchar(200) null,
	settings varchar(100) null,
	status tinyint default '1' null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create table service_order
(
	order_id varchar(32) not null,
	user_id varchar(32) not null,
	service_id varchar(32) not null,
	auto_pay tinyint(1) default '0' null comment '是否自动续费',
	is_renew tinyint(1) default '0' null comment '是否续费',
	pay_price int not null,
	pay_status tinyint(1) not null comment '1 待支付 2 已支付 3 到期 4 续期',
	status tinyint(1) default '1' null,
	expire_time timestamp null,
	start_time timestamp null,
	create_at timestamp default CURRENT_TIMESTAMP not null,
	update_at timestamp default CURRENT_TIMESTAMP not null,
	ex_time datetime null,
	constraint order_id
		unique (order_id)
)
;

create index idx_service_order
	on service_order (service_id, user_id)
;

alter table service_order
	add primary key (order_id)
;

create table service_payment_setting
(
	id int auto_increment
		primary key,
	name varchar(50) not null,
	en_name varchar(50) not null,
	price int not null comment '单位价格',
	time int not null comment '时长 单位天',
	status int(1) default '1' null,
	create_at timestamp default CURRENT_TIMESTAMP not null,
	update_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create table service_setting
(
	id int auto_increment
		primary key,
	name varchar(10) not null,
	status int(1) default '1' null,
	create_at timestamp default CURRENT_TIMESTAMP not null,
	update_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create table shopping_cart
(
	id bigint auto_increment
		primary key,
	user_id varchar(32) not null,
	product_id bigint not null,
	check_status tinyint(1) default '1' null,
	sku_values varchar(20) not null,
	goods_count int not null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null,
	status tinyint(1) unsigned default '1' not null
)
;

create index idx_cart_userId
	on shopping_cart (user_id)
;

create table sys_user
(
	id int auto_increment
		primary key,
	username varchar(20) null,
	password varchar(100) null,
	avatar varchar(100) null,
	nickname varchar(10) null,
	status int(1) default '1' null,
	lastLoginTime timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null,
	create_time timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create table toh_areas_usa
(
	area_id int unsigned auto_increment
		primary key,
	parent_id int unsigned not null comment '上一级的id值',
	area_name varchar(50) not null comment '地区名称',
	sort int unsigned default '99' not null comment '排序'
)
comment '地区信息' charset=utf8
;

create index parent_id
	on toh_areas_usa (parent_id)
;

create index sort
	on toh_areas_usa (sort)
;

create table user
(
	user_id varchar(32) not null comment 'UUID',
	nickname varchar(100) not null comment '昵称',
	real_name varchar(100) default '' null comment '真实姓名',
	sex char default '0' not null comment '性别 1男 2女  0保密',
	avatar varchar(200) default '' not null comment '头像原图',
	hometown varchar(200) default '' null comment '家乡',
	remark varchar(200) default 'welcome to tohnet' null comment '个性签名',
	create_time timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
	update_time timestamp default CURRENT_TIMESTAMP not null comment '更新时间',
	login_Time timestamp default CURRENT_TIMESTAMP not null,
	is_recommend char default '0' not null comment '是否推荐用户',
	recommend_code varchar(20) null,
	duration int default '0' not null comment '使用时长',
	status char default '1' not null comment '是否可用：1表示可以',
	integral int default '0' null comment '积分',
	commission int default '0' null comment '佣金',
	birth timestamp default CURRENT_TIMESTAMP not null,
	constraint user_id_uindex
		unique (user_id),
	constraint user_recommendCode_uindex
		unique (recommend_code)
)
;

alter table user
	add primary key (user_id)
;

create table user_address
(
	id bigint auto_increment
		primary key,
	user_id varchar(36) not null,
	contact_name varchar(50) not null,
	mobile varchar(20) not null,
	address varchar(100) not null,
	state varchar(20) not null,
	postal_code varchar(20) not null,
	status int(1) default '1' null,
	create_at timestamp default CURRENT_TIMESTAMP not null,
	update_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create table user_agreement
(
	id int auto_increment
		primary key,
	user_id varchar(36) not null,
	status int(1) default '1' null,
	create_at timestamp default CURRENT_TIMESTAMP not null,
	update_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create table user_auth
(
	id varchar(36) not null comment '主键:Id,类型为UUID',
	user_id varchar(36) not null comment '相关的用户ID',
	identify_type varchar(50) not null comment '登录方式',
	identify varchar(100) not null comment '登录账号',
	credential varchar(100) null comment '登录密码/Token',
	trade_id varchar(100) null comment '交易Id，用于支付场景的Idn	对应微信的OpenId',
	update_time timestamp default CURRENT_TIMESTAMP not null,
	status char default '1' not null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	constraint user_auth_id_uindex
		unique (id)
)
comment '用户授权表'
;

create index user_auth__fk_user_id
	on user_auth (user_id)
;

alter table user_auth
	add primary key (id)
;

create table user_integral_flow
(
	id int auto_increment
		primary key,
	user_id varchar(50) not null,
	r_user_id varchar(50) null,
	order_id varchar(50) null,
	integral int not null,
	create_at timestamp default CURRENT_TIMESTAMP not null,
	update_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create table user_payment
(
	id bigint auto_increment
		primary key,
	user_id varchar(36) not null,
	card_number varchar(50) not null,
	code varchar(20) not null,
	mm varchar(4) not null,
	yy varchar(4) not null,
	fisrt_name varchar(20) not null,
	last_name varchar(20) not null,
	status int(1) default '1' null,
	create_at timestamp default CURRENT_TIMESTAMP not null,
	update_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
;

create table user_recommend
(
	id bigint auto_increment
		primary key,
	recommend_user_id varchar(32) not null,
	user_id varchar(32) not null,
	recommend_code varchar(20) null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null,
	status char default '1' null,
	constraint idx_user_recommend_uniq
		unique (recommend_user_id, user_id, recommend_code)
)
;

create table video
(
	id bigint auto_increment comment '自增主键'
		primary key,
	title varchar(100) null comment '标题',
	thumb_url varchar(200) not null comment '文章缩略图',
	author varchar(100) null comment '作者名称',
	duration varchar(10) not null,
	tags varchar(200) null,
	read_count int default '0' null comment '阅读量',
	comment_count int default '0' null comment '评论数',
	like_count int default '0' null comment '点赞数',
	category varchar(10) null,
	is_recommend int(1) default '0' null comment '是否推荐',
	content text null,
	create_time timestamp default CURRENT_TIMESTAMP null,
	update_time timestamp default CURRENT_TIMESTAMP null,
	status int(1) default '1' null comment '是否可用',
	video_desc tinytext null,
	pusher_info tinytext null
)
comment '资讯文章'
;

create index idx_video_title
	on video (title)
;

create table video_category
(
	id int auto_increment
		primary key,
	title varchar(10) not null,
	sort int default '0' not null,
	create_time timestamp default CURRENT_TIMESTAMP not null,
	update_time timestamp default CURRENT_TIMESTAMP not null,
	status int(1) null
)
;

create index idx_video_category_title
	on video_category (title, status)
;

