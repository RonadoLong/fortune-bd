use `wq_fotune`;
drop table if exists `wq_exchange`;
create table wq_exchange
(
    id         int(11)     not null auto_increment comment '主键id',
    exchange   varchar(20) not null comment '交易所代号 okex  huobi bitmex',
    status     tinyint(1)  not null comment '0 close  1 open',
    created_at datetime    not null comment '创建时间',
    updated_at datetime    not null comment '更新时间',
    primary key (id),
    unique key (exchange)
);

drop table if exists `wq_exchange_api`;
create table wq_exchange_api
(
    id            bigint       not null auto_increment comment '主键id',
    user_id       varchar(255) not null comment '',
    exchange_id   int(11)      not null comment '交易所id',
    exchange_name varchar(50)  not null comment '',
    api_key       varchar(100) not null comment '',
    secret        varchar(255) not null comment '',
    passphrase    varchar(255) not null default '' comment 'okex only',
    created_at    datetime     not null comment '创建时间',
    updated_at    datetime     not null comment '更新时间',
    primary key (id),
    unique index `user_id_api_key` (user_id, api_key)
);


drop table if exists `wq_strategy`;

create table wq_user_strategy
(
    id          int(11)        not null auto_increment comment '主键id',
    user_id     varchar(64)    not null comment '',
    strategy_id varchar(64)    not null comment '策略id',
    platform    varchar(11)    not null comment 'exchange_name',
    api_key     varchar(64)    not null comment '',
    balance     decimal(20, 8) not null comment '初始资金',
    status      tinyint(1)     not null comment '0 close  1 running  2 suspend',
    start_time  datetime       not null comment '',
    end_time    datetime       not null comment '',
    validity    int(11)        null default null comment '',
    created_at  datetime       not null comment '创建时间',
    updated_at  datetime       not null comment '更新时间',
    primary key (id),
    unique key (strategy_id),
    index `strategy_id_status` (strategy_id, status)
);


drop table if exists `wq_trade`;

create table wq_trade
(
    id          bigint      not null auto_increment comment '主键id',
    user_id     varchar(64) not null comment '',
    trade_id    varchar(64) not null comment '交易id',
    api_key     varchar(64) not null comment '',
    strategy_id varchar(64) not null comment '',
    symbol      varchar(25) not null comment '',
    open_price  varchar(12) not null comment '开仓价格',
    close_price varchar(12) not null comment '平仓价格',
    volume      int(50)     not null comment '手数',
    profit      varchar(12) not null comment '盈亏',
    pos         varchar(12) not null comment 'pos',
    created_at  datetime    not null comment '创建时间',
    updated_at  datetime    not null comment '更新时间',
    primary key (id),
    index `strategy_id_user_id_api_key` (strategy_id, user_id, api_key)
);

#
# 已实现盈亏
# 未实现盈亏

drop table if exists `wq_profit`;
create table wq_profit
(
    id                bigint       not null auto_increment comment '主键id',
    user_id           varchar(64)  not null comment '',
    api_key           varchar(64)  not null comment '',
    strategy_id       varchar(64)  not null comment '',
    symbol            varchar(25)  not null comment '',
    realize_profit    varchar(12)  not null comment '已实现盈亏',
    un_realize_profit varchar(12)  not null comment '未实现盈亏',
    position          int(100)     not null comment '',
    rate_return       float(10, 2) not null comment '收益率',
    created_at        datetime     not null comment '创建时间',
    updated_at        datetime     not null comment '更新时间',
    primary key (id),
    index `user_id_api_key` (user_id, api_key)
);


drop table if exists wq_strategy;

create table wq_strategy
(
    id               bigint primary key not null auto_increment comment '主键id',
    group_id         varchar(32)        not null comment '策略组合id',
    name             varchar(20)        not null comment '',
    remark           tinytext           not null comment '',
    exchange_name    varchar(12)        not null,
    exchange_id      int(11)            not null,
    symbol           varchar(100)       not null,
    rate_return      varchar(12)        not null,
    rate_return_year varchar(12)        not null,
    price            int                not null default 0,
    is_free          tinyint            not null default 0,
    duration         int                not null,
    status           tinyint            not null default 1,
    rate_win         varchar(12)        not null,
    created_at       datetime           not null comment '创建时间',
    updated_at       datetime           not null comment '更新时间',
    deleted_at       datetime
);


#     ID             int64  `gorm:"column:id" json:"id"`
# 	Symbol         string `gorm:"column:symbol" json:"symbol"`
# 	RateReturnYear string `gorm:"column:rate_return_year" json:"rate_return_year"`
# 	State          int32  `gorm:"column:state" json:"state"`
# 	Url            string `gorm:"column:url" json:"url"`
drop table if exists wq_symbol_recommend;

use `wq_fotune`;
create table wq_symbol_recommend
(
    id               bigint primary key not null auto_increment comment '主键id',
    rate_return_year varchar(12)        not null,
    symbol           varchar(100)       not null,
    url              varchar(100)       null,
    state            tinyint            not null default 1
);
