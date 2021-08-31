use `wq_fotune`;
drop table if exists `wq_wallet`;
create table `wq_wallet`
(
    id              int(11)      not null primary key auto_increment comment '主键id',
    user_id         varchar(32)  not null comment '用户id',
    api_key         varchar(100) not null comment '子账户的apikey',
    secret          varchar(255) not null comment '',
    sub_account_id  varchar(50)  not null comment '子账户的交易所id',
    deposit_addr    varchar(150) not null comment '子账户的充值地址',
    deposit_coin    varchar(10)  not null comment '充值币种',
    exchange        varchar(20)  not null comment '交易所名称',
    wq_coin_balance varchar(50)  not null comment '微宽平台币余额',
    created_at      datetime     not null comment '创建时间',
    updated_at      datetime     not null comment '更新时间'
);

#划转
drop table if exists `wq_transfer`;
create table `wq_transfer`
(
    id                   int(11)      not null primary key auto_increment comment '主键id',
    user_id              varchar(32)  not null comment '用户id',
    from_addr            varchar(150) not null comment '转账来源子账户 母账户给子账户转账这里填空',
    to_addr              varchar(150) not null comment '转出目标子账户 子账户给母账户转账这里填空',
    coin                 varchar(10)  not null comment '划转源币种',
    amount               varchar(50)  not null comment '划转数量',
    amount_before        varchar(50)  not null comment '划转前源币种数量',
    amount_after         varchar(50)  not null comment '划转后源币种数量',
    to_coin              varchar(10)  not null comment '目标币种',
    to_coin_amount       varchar(50)  not null comment '换算目标币种数量',
    to_coin_amount_after varchar(50)  not null comment '换算后目标币种数量',
    tx_id                varchar(150) not null comment '交易hashId',
    state                tinyint      not null comment '0失败， 1划转, 2已退还',
    created_at           datetime     not null comment '创建时间',
    updated_at           datetime     not null comment '更新时间'
);

#提现
drop table if exists `wq_withdrawal`;
create table `wq_withdrawal`
(
    id         int(11)      not null primary key auto_increment comment '主键id',
    user_id    varchar(32)  not null comment '用户id',
    coin       varchar(10)  not null comment '币种',
    cash_addr  varchar(150) not null comment '体现地址',
    amount     varchar(50)  not null comment '数量',
    state      tinyint      not null comment '0拒绝， 1待审核, 2已完成',
    created_at datetime     not null comment '创建时间',
    updated_at datetime     not null comment '更新时间'
);

# use `wq_fotune`;
# drop table if exists `wq_deposit_record`;
# create table `wq_deposit_record`
# (
#     id             int(11)      not null primary key auto_increment comment '主键id',
#     user_id        varchar(32)  not null comment '用户id',
#     sub_account_id varchar(50)  not null comment '子账户的交易所id',
#     deposit_addr   varchar(150) not null comment '子账户的收款地址',
#     amount         varchar(50)  not null comment '数量',
#     coin           varchar(10)  not null comment '币种',
#     from_addr      varchar(150) not null comment '充值来源地址',
#     tx_id          varchar(150) not null comment '交易hashId',
#     state          tinyint      not null comment '0,未支付， 1已支付, 2已退款',
#     created_at     datetime     not null comment '创建时间',
#     updated_at     datetime     not null comment '更新时间'
# );

# use `wq_fotune`;
drop table if exists `wq_coin_info`;
create table `wq_coin_info`
(
    coin       varchar(10) not null comment '币种',
    price      varchar(11) not null comment '美元价格',
    created_at datetime    not null comment '创建时间',
    updated_at datetime    not null comment '更新时间'
)

