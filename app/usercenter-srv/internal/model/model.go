package model

import (
    "time"
    "wq-fotune-backend/libs/ucode"
    "wq-fotune-backend/pkg/bcrypt2"
    "wq-fotune-backend/pkg/snowflake"
)

type WqUserBase struct {
    ID             int       `gorm:"column:id"`
    UserID         string    `gorm:"column:user_id"`
    Name           string    `gorm:"column:name"`
    Phone          string    `gorm:"column:phone"`
    Password       string    `gorm:"column:password"`
    Avatar         string    `gorm:"column:avatar"`
    Status         int       `gorm:"column:status"`          // 0 禁用  1可用
    InvitationCode string    `gorm:"column:invitation_code"` //邀请码
    LastLogin      time.Time `gorm:"column:last_login"`
    LoginCount     int       `gorm:"column:login_count"`
    CreatedAt      time.Time `gorm:"column:created_at"`
    UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func NewWqUserBase(phone, password string) *WqUserBase {
    return &WqUserBase{
        UserID:         snowflake.SNode.Generate().String(),
        Name:           ucode.RandStringRunesNoNum(20),
        Phone:          phone,
        Password:       bcrypt2.CryptPassword(password),
        InvitationCode: ucode.GetRandomString(8),
        LoginCount:     0,
        Status:         1,
        LastLogin:      time.Now(),
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }
}

type WqUserInvite struct {
    ID            int32  `gorm:"column:id"`
    UserID        string `gorm:"column:user_id"`
    InvitedUserID string `gorm:"column:in_user_id"`
}

type WqMember struct {
    ID        int32     `gorm:"column:id"`
    Name      string    `gorm:"column:name"`
    Remark    string    `gorm:"column:remark"`
    Price     int32     `gorm:"column:price"`
    OldPrice  int32     `gorm:"column:old_price"`
    Duration  int32     `gorm:"column:duration"`
    State     int32     `gorm:"column:state"`
    CreatedAt time.Time `gorm:"column:created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
}

type WqPayment struct {
    ID        int32     `gorm:"column:id"`
    Name      string    `gorm:"column:name"`
    Remark    string    `gorm:"column:remark"`
    BitAddr   string    `gorm:"column:bit_addr"`
    BitCode   string    `gorm:"column:bit_code"`
    State     int32     `gorm:"column:state"`
    CreatedAt time.Time `gorm:"column:created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
}

//id         int                    not null auto_increment primary key,
//user_id    varchar(32)            not null,
//pay_id     int                    not null comment '',
//member_id  int                    not null comment '',
//price      int                    not null comment '',
//state      tinyint(1)             not null comment '0取消 1成功 2支付中',
//from_addr  varchar(255)           null comment '数字货币支付地址',
//created_at datetime               not null comment '创建时间',
//updated_at datetime default now() not null comment '更新时间'

type WqUserPay struct {
    ID        int32     `gorm:"column:id"`
    UserID    string    `gorm:"column:user_id"`
    PayID     int32     `gorm:"column:pay_id"`
    MemberID  int32     `gorm:"column:member_id"`
    Price     int32     `gorm:"column:price"`
    State     int32     `gorm:"column:state"`
    FromAddr  string    `gorm:"column:from_addr"`
    CreatedAt time.Time `gorm:"column:created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
}
