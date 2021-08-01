use `wq_fotune`;
drop table if exists `wq_common_carousel`;

create table wq_common_carousel
(
    id         int(11)      not null auto_increment comment '主键id',
    image      varchar(500) not null comment '轮播图地址',
    click_url  varchar(500) not null comment '点击地址',
    created_at datetime     not null comment '创建时间',
    updated_at datetime     not null comment '更新时间',
    primary key (id)
);

drop table if exists `wq_common_contact`;
create table wq_common_contact
(
    id         int(11)      not null auto_increment comment '主键id',
    image      varchar(500) not null comment '二维码地址',
    content    text(1000)   not null comment '',
    contact    varchar(255) not null comment '',
    created_at datetime     not null comment '创建时间',
    updated_at datetime     not null comment '更新时间',
    primary key (id)
);

create table wq_app_version
(
    id int(11) primary key not null auto_increment ,
    app_version varchar(10) not null comment '版本号',
    download_code varchar(500) not null comment '下载二维码',
    created_at datetime     not null comment '创建时间',
    updated_at datetime     not null comment '更新时间'
)