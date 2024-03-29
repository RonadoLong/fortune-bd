package model

import (
	"time"
)


type WqWallet struct {
	ID            int       `gorm:"column:id" json:"id"`
	UserID        string    `gorm:"user_id" json:"user_id"`
	ApiKey        string    `gorm:"api_key" json:"api_key"`
	Secret        string    `gorm:"secret" json:"secret"`
	SubAccountID  string    `gorm:"sub_account_id" json:"sub_account_id"`
	DepositAddr   string    `gorm:"deposit_addr" json:"deposit_addr"`
	WqCoinBalance string    `gorm:"wq_coin_balance" json:"wq_coin_balance"`
	CreatedAt     time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"updated_at" json:"updated_at"`
}

func NewWqWalletModel(userId, apikey, secret, subAccountId, DepositAddr string) *WqWallet {
	return &WqWallet{
		UserID:        userId,
		ApiKey:        apikey,
		Secret:        secret,
		SubAccountID:  subAccountId,
		DepositAddr:   DepositAddr,
		WqCoinBalance: "0",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

type WqCoinInfo struct {
	Coin      string    `gorm:"coin" json:"coin"`
	Price     string    `gorm:"price" json:"price"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at"`
}


type WqTransferRecord struct {
	ID                 int       `gorm:"column:id" json:"id"`
	UserID             string    `gorm:"user_id" json:"user_id"`
	Coin               string    `gorm:"coin" json:"coin"`
	Amount             string    `gorm:"amount" json:"amount"`
	AmountBefore       string    `gorm:"amount_before" json:"amount_before"`
	AmountAfter        string    `gorm:"amount_after" json:"amount_after"`
	ToCoin             string    `gorm:"to_coin" json:"to_coin"`
	ToCoinAmount       string    `gorm:"to_coin_amount" json:"to_coin_amount"`
	ToCoinAmountBefore string    `gorm:"to_coin_amount_before" json:"to_coin_amount_before"`
	ToCoinAmountAfter  string    `gorm:"to_coin_amount_after" json:"to_coin_amount_after"`
	TxID               string    `gorm:"tx_id" json:"tx_id"`
	State              int       `gorm:"state" json:"state"`
	CreatedAt          time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"updated_at" json:"updated_at"`
}

func NewWqTransferRecord(userId, coin, amount, amountBefore, amountAfter, toCoin, toCoinAmount, toCoinAmountBefore, toCoinAmountAfter, txID string) *WqTransferRecord {
	return &WqTransferRecord{
		UserID:             userId,
		Coin:               coin,
		Amount:             amount,
		AmountBefore:       amountBefore,
		AmountAfter:        amountAfter,
		ToCoin:             toCoin,
		ToCoinAmount:       toCoinAmount,
		ToCoinAmountBefore: toCoinAmountBefore,
		ToCoinAmountAfter:  toCoinAmountAfter,
		TxID:               txID,
		State:              1,
		CreatedAt:          time.Time{},
		UpdatedAt:          time.Time{},
	}
}

type WqWithdrawal struct {
	ID        int       `gorm:"column:id" json:"id"`
	UserID    string    `gorm:"user_id" json:"user_id"`
	Coin      string    `gorm:"coin" json:"coin"`
	CashAddr  string    `gorm:"cash_addr" json:"cash_addr"`
	Amount    string    `gorm:"amount" json:"amount"`
	State     int       `gorm:"state" json:"state"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at"`
}

type WqIfcGiftRecord struct {
	ID        int       `gorm:"column:id" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	InUserID  string    `gorm:"column:in_user_id" json:"in_user_id"`
	Volume    string    `gorm:"column:volume" json:"volume"`
	Type      string    `gorm:"column:type" json:"type"`
	Exchange  string    `gorm:"column:exchange" json:"exchange"`
	Before    string    `gorm:"column:before" json:"before"`
	After     string    `gorm:"column:after" json:"after"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
