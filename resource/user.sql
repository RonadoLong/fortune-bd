use `wq_fotune`;
drop table  if exists `wq_user_base`;
create table wq_user_base
(
    id              int(11)      not null auto_increment comment '主键id',
    user_id         varchar(32)  not null comment '用户id',
    name            varchar(25)  not null comment '昵称',
    phone           varchar(15)  not null comment '手机号',
    password        varchar(255) not null comment '密码',
    avatar          varchar(255) default '' not null comment '头像地址',
    status          tinyint(1) default 1 not null comment '0不可用  1账号可用',
    invitation_code varchar(10)  not null comment '邀请码',
    last_login      datetime     null default null comment '最后登陆时间',
    login_count     int(11)      not null default 0 comment '登陆次数',
    created_at      datetime     not null comment '创建时间',
    updated_at      datetime     not null comment '更新时间',
    primary key (id),
    unique index `user_id` (user_id),
    unique index `phone` (phone),
    unique index `invitation_code` (invitation_code)
);

drop table  if exists `wq_user_invite`;
create table wq_user_invite
(
    id         int(11)     not null auto_increment comment '主键id',
    user_id    varchar(32) not null comment '用户id',
    in_user_id varchar(32) not null comment '被邀请用户id',
    primary key (id),
    unique index `user_in_user`(user_id, in_user_id)
);
drop table  if exists `wq_user_platform`;
create table wq_user_platform
(
    id         int(11)              not null auto_increment comment '主键id',
    user_id    varchar(32)          not null comment '用户id',
    login_type varchar(20)          not null comment '登陆类型',
    identifier varchar(64)          not null comment '微博id 微信id',
    credential varchar(255)         not null comment '第三方token',
    status     tinyint(1) default 1 not null comment '是否已经验证 0未绑定, 1已绑定',
    primary key (id),
    unique index `user_id_identifier` (user_id, identifier)
);


create table wq_member
(
    id         int         not null auto_increment primary key,
    name       varchar(50) not null default '',
    remark     text(500)   null  ,
    price      int         not null default 0,
    old_price  int         not null default 0,
    duration   int         not null default 30 comment '有效期 30 90 365',
    state      tinyint(1)  not null default 0 comment '0 close 1 open',
    created_at datetime    not null comment '创建时间',
    updated_at datetime             default now() not null comment '更新时间'
);

create table wq_pay
(
    id          int          not null auto_increment primary key,
    name        varchar(50)  not null default '',
    remark      text(500)    null ,
    bit_addr varchar(255) null     default '',
    bit_code    text         null      ,
    state       tinyint(1)   not null default 0 comment '0 close 1 open',
    created_at  datetime     not null comment '创建时间',
    updated_at  datetime              default now() not null comment '更新时间'
);

create table wq_user_pay
(
    id         int                    not null auto_increment primary key,
    user_id    varchar(32)            not null,
    pay_id     int                    not null comment '',
    member_id  int                    not null comment '',
    price      int                    not null comment '',
    state      tinyint(1)             not null comment '0取消 1成功 2支付中',
    from_addr  varchar(255)           null comment '数字货币支付地址',
    created_at datetime               not null comment '创建时间',
    updated_at datetime default now() not null comment '更新时间'
);

create table wq_ifc_gift_record
(
    id int not null  auto_increment primary key,
    user_id varchar(32) not null comment '被充值用户',
    in_user_id varchar(32) not null comment '被邀请用户',
    volume varchar(50) not null comment '',
    type varchar(10) not null comment 'register 代表注册送, bind_api 代表绑定api送',
    exchange varchar(10) not null comment '如果是绑定api 填写这个字段 否则填'' ',
    created_at datetime               not null comment '创建时间',
    updated_at datetime default now() not null comment '更新时间'
)