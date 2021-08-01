package goex

import (
	"strings"
	"sync"
)

type Currency struct {
	Symbol string
	Desc   string
}

func (c Currency) String() string {
	return c.Symbol
}

func (c Currency) Eq(c2 Currency) bool {
	return c.Symbol == c2.Symbol
}

// A->B(A兑换为B)
type CurrencyPair struct {
	CurrencyA Currency
	CurrencyB Currency
}

var (
	UNKNOWN = Currency{"UNKNOWN", ""}
	CNY     = Currency{"CNY", ""}
	USD     = Currency{"USD", ""}
	USDT    = Currency{"USDT", ""}
	//PAX     = Currency{"PAX", "https://www.paxos.com/"}
	//USDC    = Currency{"USDC", "https://www.centre.io/"}
	//EUR     = Currency{"EUR", ""}
	KRW = Currency{"KRW", ""}
	JPY = Currency{"JPY", ""}
	//BTC     = Currency{"BTC", "https://bitcoin.org/"}
	XBT = Currency{"XBT", ""}
	//BCC     = Currency{"BCC", ""}
	//BCH     = Currency{"BCH", ""}
	BCX = Currency{"BCX", ""}
	//LTC     = Currency{"LTC", ""}
	//ETH     = Currency{"ETH", ""}
	//ETC     = Currency{"ETC", ""}
	//EOS     = Currency{"EOS", ""}
	//BTS     = Currency{"BTS", ""}
	//QTUM    = Currency{"QTUM", ""}
	//SC      = Currency{"SC", ""}
	ANS = Currency{"ANS", ""}
	//ZEC     = Currency{"ZEC", ""}
	//DCR     = Currency{"DCR", ""}
	//XRP     = Currency{"XRP", ""}
	BTG = Currency{"BTG", ""}
	BCD = Currency{"BCD", ""}
	//NEO     = Currency{"NEO", ""}
	HSR = Currency{"HSR", ""}
	BSV = Currency{"BSV", ""}
	OKB = Currency{"OKB", "OKB is a global utility token issued by OK Blockchain Foundation"}
	HT  = Currency{"HT", "HuoBi Token"}
	//BNB     = Currency{"BNB", "BNB, or Binance Coin, is a cryptocurrency created by Binance."}
	//TRX     = Currency{"TRX", ""}

	ADADOWN  = Currency{"ADADOWN", ""}
	ADAUP    = Currency{"ADAUP", ""}
	ADA      = Currency{"ADA", ""}
	AION     = Currency{"AION", ""}
	ALGO     = Currency{"ALGO", ""}
	ANKR     = Currency{"ANKR", ""}
	ARDR     = Currency{"ARDR", ""}
	ARPA     = Currency{"ARPA", ""}
	ATOM     = Currency{"ATOM", ""}
	AUD      = Currency{"AUD", ""}
	BAL      = Currency{"BAL", ""}
	BAND     = Currency{"BAND", ""}
	BAT      = Currency{"BAT", ""}
	BCC      = Currency{"BCC", ""}
	BCHABC   = Currency{"BCHABC", ""}
	BCHSV    = Currency{"BCHSV", ""}
	BCH      = Currency{"BCH", ""}
	BEAM     = Currency{"BEAM", ""}
	BEAR     = Currency{"BEAR", ""}
	BKRW     = Currency{"BKRW", ""}
	BLZ      = Currency{"BLZ", ""}
	BNBBEAR  = Currency{"BNBBEAR", ""}
	BNBBULL  = Currency{"BNBBULL", ""}
	BNBDOWN  = Currency{"BNBDOWN", ""}
	BNBUP    = Currency{"BNBUP", ""}
	BNB      = Currency{"BNB", ""}
	BNT      = Currency{"BNT", ""}
	BTCDOWN  = Currency{"BTCDOWN", ""}
	BTCUP    = Currency{"BTCUP", ""}
	BTC      = Currency{"BTC", ""}
	BTS      = Currency{"BTS", ""}
	BTT      = Currency{"BTT", ""}
	BULL     = Currency{"BULL", ""}
	BUSD     = Currency{"BUSD", ""}
	CELR     = Currency{"CELR", ""}
	CHR      = Currency{"CHR", ""}
	CHZ      = Currency{"CHZ", ""}
	COCOS    = Currency{"COCOS", ""}
	COMP     = Currency{"COMP", ""}
	COS      = Currency{"COS", ""}
	COTI     = Currency{"COTI", ""}
	CTSI     = Currency{"CTSI", ""}
	CTXC     = Currency{"CTXC", ""}
	CVC      = Currency{"CVC", ""}
	DAI      = Currency{"DAI", ""}
	DASH     = Currency{"DASH", ""}
	DATA     = Currency{"DATA", ""}
	DCR      = Currency{"DCR", ""}
	DENT     = Currency{"DENT", ""}
	DGB      = Currency{"DGB", ""}
	DOCK     = Currency{"DOCK", ""}
	DOGE     = Currency{"DOGE", ""}
	DOT      = Currency{"DOT", ""}
	DREP     = Currency{"DREP", ""}
	DUSK     = Currency{"DUSK", ""}
	ENJ      = Currency{"ENJ", ""}
	EOSBEAR  = Currency{"EOSBEAR", ""}
	EOSBULL  = Currency{"EOSBULL", ""}
	EOS      = Currency{"EOS", ""}
	ERD      = Currency{"ERD", ""}
	ETC      = Currency{"ETC", ""}
	ETHBEAR  = Currency{"ETHBEAR", ""}
	ETHBULL  = Currency{"ETHBULL", ""}
	ETHDOWN  = Currency{"ETHDOWN", ""}
	ETHUP    = Currency{"ETHUP", ""}
	ETH      = Currency{"ETH", ""}
	EUR      = Currency{"EUR", ""}
	FET      = Currency{"FET", ""}
	FTM      = Currency{"FTM", ""}
	FTT      = Currency{"FTT", ""}
	FUN      = Currency{"FUN", ""}
	GBP      = Currency{"GBP", ""}
	GTO      = Currency{"GTO", ""}
	GXS      = Currency{"GXS", ""}
	HBAR     = Currency{"HBAR", ""}
	HC       = Currency{"HC", ""}
	HIVE     = Currency{"HIVE", ""}
	HOT      = Currency{"HOT", ""}
	ICX      = Currency{"ICX", ""}
	IOST     = Currency{"IOST", ""}
	IOTA     = Currency{"IOTA", ""}
	IOTX     = Currency{"IOTX", ""}
	IRIS     = Currency{"IRIS", ""}
	JST      = Currency{"JST", ""}
	KAVA     = Currency{"KAVA", ""}
	KEY      = Currency{"KEY", ""}
	KMD      = Currency{"KMD", ""}
	KNC      = Currency{"KNC", ""}
	LEND     = Currency{"LEND", ""}
	LINKDOWN = Currency{"LINKDOWN", ""}
	LINKUP   = Currency{"LINKUP", ""}
	LINK     = Currency{"LINK", ""}
	LRC      = Currency{"LRC", ""}
	LSK      = Currency{"LSK", ""}
	LTC      = Currency{"LTC", ""}
	LTO      = Currency{"LTO", ""}
	MANA     = Currency{"MANA", ""}
	MATIC    = Currency{"MATIC", ""}
	MBL      = Currency{"MBL", ""}
	MCO      = Currency{"MCO", ""}
	MDT      = Currency{"MDT", ""}
	MFT      = Currency{"MFT", ""}
	MITH     = Currency{"MITH", ""}
	MKR      = Currency{"MKR", ""}
	MTL      = Currency{"MTL", ""}
	NANO     = Currency{"NANO", ""}
	NEO      = Currency{"NEO", ""}
	NKN      = Currency{"NKN", ""}
	NPXS     = Currency{"NPXS", ""}
	NULS     = Currency{"NULS", ""}
	OGN      = Currency{"OGN", ""}
	OMG      = Currency{"OMG", ""}
	ONE      = Currency{"ONE", ""}
	ONG      = Currency{"ONG", ""}
	ONT      = Currency{"ONT", ""}
	PAX      = Currency{"PAX", ""}
	PERL     = Currency{"PERL", ""}
	PNT      = Currency{"PNT", ""}
	QTUM     = Currency{"QTUM", ""}
	REN      = Currency{"REN", ""}
	REP      = Currency{"REP", ""}
	RLC      = Currency{"RLC", ""}
	RVN      = Currency{"RVN", ""}
	SC       = Currency{"SC", ""}
	SNX      = Currency{"SNX", ""}
	SOL      = Currency{"SOL", ""}
	SRM      = Currency{"SRM", ""}
	STMX     = Currency{"STMX", ""}
	STORJ    = Currency{"STORJ", ""}
	STORM    = Currency{"STORM", ""}
	STPT     = Currency{"STPT", ""}
	STRAT    = Currency{"STRAT", ""}
	STX      = Currency{"STX", ""}
	SXP      = Currency{"SXP", ""}
	TCT      = Currency{"TCT", ""}
	TFUEL    = Currency{"TFUEL", ""}
	THETA    = Currency{"THETA", ""}
	TOMO     = Currency{"TOMO", ""}
	TROY     = Currency{"TROY", ""}
	TRX      = Currency{"TRX", ""}
	TUSD     = Currency{"TUSD", ""}
	USDC     = Currency{"USDC", ""}
	USDSB    = Currency{"USDSB", ""}
	USDS     = Currency{"USDS", ""}
	VEN      = Currency{"VEN", ""}
	VET      = Currency{"VET", ""}
	VITE     = Currency{"VITE", ""}
	VTHO     = Currency{"VTHO", ""}
	WAN      = Currency{"WAN", ""}
	WAVES    = Currency{"WAVES", ""}
	WIN      = Currency{"WIN", ""}
	WRX      = Currency{"WRX", ""}
	WTC      = Currency{"WTC", ""}
	XLM      = Currency{"XLM", ""}
	XMR      = Currency{"XMR", ""}
	XRPBEAR  = Currency{"XRPBEAR", ""}
	XRPBULL  = Currency{"XRPBULL", ""}
	XRP      = Currency{"XRP", ""}
	XTZDOWN  = Currency{"XTZDOWN", ""}
	XTZUP    = Currency{"XTZUP", ""}
	XTZ      = Currency{"XTZ", ""}
	XZC      = Currency{"XZC", ""}
	YFI      = Currency{"YFI", ""}
	ZEC      = Currency{"ZEC", ""}
	ZEN      = Currency{"ZEN", ""}
	ZIL      = Currency{"ZIL", ""}
	ZRX      = Currency{"ZRX", ""}

	//currency pair

	BTC_CNY  = CurrencyPair{BTC, CNY}
	LTC_CNY  = CurrencyPair{LTC, CNY}
	BCC_CNY  = CurrencyPair{BCC, CNY}
	ETH_CNY  = CurrencyPair{ETH, CNY}
	ETC_CNY  = CurrencyPair{ETC, CNY}
	EOS_CNY  = CurrencyPair{EOS, CNY}
	BTS_CNY  = CurrencyPair{BTS, CNY}
	QTUM_CNY = CurrencyPair{QTUM, CNY}
	SC_CNY   = CurrencyPair{SC, CNY}
	ANS_CNY  = CurrencyPair{ANS, CNY}
	ZEC_CNY  = CurrencyPair{ZEC, CNY}

	BTC_KRW = CurrencyPair{BTC, KRW}
	ETH_KRW = CurrencyPair{ETH, KRW}
	ETC_KRW = CurrencyPair{ETC, KRW}
	LTC_KRW = CurrencyPair{LTC, KRW}
	BCH_KRW = CurrencyPair{BCH, KRW}

	BTC_USD = CurrencyPair{BTC, USD}
	LTC_USD = CurrencyPair{LTC, USD}
	ETH_USD = CurrencyPair{ETH, USD}
	ETC_USD = CurrencyPair{ETC, USD}
	BCH_USD = CurrencyPair{BCH, USD}
	BCC_USD = CurrencyPair{BCC, USD}
	XRP_USD = CurrencyPair{XRP, USD}
	BCD_USD = CurrencyPair{BCD, USD}
	EOS_USD = CurrencyPair{EOS, USD}
	BTG_USD = CurrencyPair{BTG, USD}
	BSV_USD = CurrencyPair{BSV, USD}

	//BTC_USDT = CurrencyPair{BTC, USDT}
	//LTC_USDT = CurrencyPair{LTC, USDT}
	//BCH_USDT = CurrencyPair{BCH, USDT}
	//BCC_USDT = CurrencyPair{BCC, USDT}
	//ETC_USDT = CurrencyPair{ETC, USDT}
	//ETH_USDT = CurrencyPair{ETH, USDT}
	BCD_USDT = CurrencyPair{BCD, USDT}
	//NEO_USDT = CurrencyPair{NEO, USDT}
	//EOS_USDT = CurrencyPair{EOS, USDT}
	//XRP_USDT = CurrencyPair{XRP, USDT}
	HSR_USDT = CurrencyPair{HSR, USDT}
	BSV_USDT = CurrencyPair{BSV, USDT}
	OKB_USDT = CurrencyPair{OKB, USDT}
	HT_USDT  = CurrencyPair{HT, USDT}
	//BNB_USDT = CurrencyPair{BNB, USDT}
	//PAX_USDT = CurrencyPair{PAX, USDT}
	//TRX_USDT = CurrencyPair{TRX, USDT}

	XRP_EUR = CurrencyPair{XRP, EUR}

	BTC_JPY = CurrencyPair{BTC, JPY}
	LTC_JPY = CurrencyPair{LTC, JPY}
	ETH_JPY = CurrencyPair{ETH, JPY}
	ETC_JPY = CurrencyPair{ETC, JPY}
	BCH_JPY = CurrencyPair{BCH, JPY}

	LTC_BTC = CurrencyPair{LTC, BTC}
	ETH_BTC = CurrencyPair{ETH, BTC}
	ETC_BTC = CurrencyPair{ETC, BTC}
	BCC_BTC = CurrencyPair{BCC, BTC}
	BCH_BTC = CurrencyPair{BCH, BTC}
	DCR_BTC = CurrencyPair{DCR, BTC}
	XRP_BTC = CurrencyPair{XRP, BTC}
	BTG_BTC = CurrencyPair{BTG, BTC}
	BCD_BTC = CurrencyPair{BCD, BTC}
	NEO_BTC = CurrencyPair{NEO, BTC}
	EOS_BTC = CurrencyPair{EOS, BTC}
	HSR_BTC = CurrencyPair{HSR, BTC}
	BSV_BTC = CurrencyPair{BSV, BTC}
	OKB_BTC = CurrencyPair{OKB, BTC}
	HT_BTC  = CurrencyPair{HT, BTC}
	BNB_BTC = CurrencyPair{BNB, BTC}
	TRX_BTC = CurrencyPair{TRX, BTC}

	ETC_ETH = CurrencyPair{ETC, ETH}
	EOS_ETH = CurrencyPair{EOS, ETH}
	ZEC_ETH = CurrencyPair{ZEC, ETH}
	NEO_ETH = CurrencyPair{NEO, ETH}
	HSR_ETH = CurrencyPair{HSR, ETH}
	LTC_ETH = CurrencyPair{LTC, ETH}

	ADADOWN_USDT  = CurrencyPair{ADADOWN, USDT}
	ADAUP_USDT    = CurrencyPair{ADAUP, USDT}
	ADA_USDT      = CurrencyPair{ADA, USDT}
	AION_USDT     = CurrencyPair{AION, USDT}
	ALGO_USDT     = CurrencyPair{ALGO, USDT}
	ANKR_USDT     = CurrencyPair{ANKR, USDT}
	ARDR_USDT     = CurrencyPair{ARDR, USDT}
	ARPA_USDT     = CurrencyPair{ARPA, USDT}
	ATOM_USDT     = CurrencyPair{ATOM, USDT}
	AUD_USDT      = CurrencyPair{AUD, USDT}
	BAL_USDT      = CurrencyPair{BAL, USDT}
	BAND_USDT     = CurrencyPair{BAND, USDT}
	BAT_USDT      = CurrencyPair{BAT, USDT}
	BCC_USDT      = CurrencyPair{BCC, USDT}
	BCHABC_USDT   = CurrencyPair{BCHABC, USDT}
	BCHSV_USDT    = CurrencyPair{BCHSV, USDT}
	BCH_USDT      = CurrencyPair{BCH, USDT}
	BEAM_USDT     = CurrencyPair{BEAM, USDT}
	BEAR_USDT     = CurrencyPair{BEAR, USDT}
	BKRW_USDT     = CurrencyPair{BKRW, USDT}
	BLZ_USDT      = CurrencyPair{BLZ, USDT}
	BNBBEAR_USDT  = CurrencyPair{BNBBEAR, USDT}
	BNBBULL_USDT  = CurrencyPair{BNBBULL, USDT}
	BNBDOWN_USDT  = CurrencyPair{BNBDOWN, USDT}
	BNBUP_USDT    = CurrencyPair{BNBUP, USDT}
	BNB_USDT      = CurrencyPair{BNB, USDT}
	BNT_USDT      = CurrencyPair{BNT, USDT}
	BTCDOWN_USDT  = CurrencyPair{BTCDOWN, USDT}
	BTCUP_USDT    = CurrencyPair{BTCUP, USDT}
	BTC_USDT      = CurrencyPair{BTC, USDT}
	BTS_USDT      = CurrencyPair{BTS, USDT}
	BTT_USDT      = CurrencyPair{BTT, USDT}
	BULL_USDT     = CurrencyPair{BULL, USDT}
	BUSD_USDT     = CurrencyPair{BUSD, USDT}
	CELR_USDT     = CurrencyPair{CELR, USDT}
	CHR_USDT      = CurrencyPair{CHR, USDT}
	CHZ_USDT      = CurrencyPair{CHZ, USDT}
	COCOS_USDT    = CurrencyPair{COCOS, USDT}
	COMP_USDT     = CurrencyPair{COMP, USDT}
	COS_USDT      = CurrencyPair{COS, USDT}
	COTI_USDT     = CurrencyPair{COTI, USDT}
	CTSI_USDT     = CurrencyPair{CTSI, USDT}
	CTXC_USDT     = CurrencyPair{CTXC, USDT}
	CVC_USDT      = CurrencyPair{CVC, USDT}
	DAI_USDT      = CurrencyPair{DAI, USDT}
	DASH_USDT     = CurrencyPair{DASH, USDT}
	DATA_USDT     = CurrencyPair{DATA, USDT}
	DCR_USDT      = CurrencyPair{DCR, USDT}
	DENT_USDT     = CurrencyPair{DENT, USDT}
	DGB_USDT      = CurrencyPair{DGB, USDT}
	DOCK_USDT     = CurrencyPair{DOCK, USDT}
	DOGE_USDT     = CurrencyPair{DOGE, USDT}
	DOT_USDT      = CurrencyPair{DOT, USDT}
	DREP_USDT     = CurrencyPair{DREP, USDT}
	DUSK_USDT     = CurrencyPair{DUSK, USDT}
	ENJ_USDT      = CurrencyPair{ENJ, USDT}
	EOSBEAR_USDT  = CurrencyPair{EOSBEAR, USDT}
	EOSBULL_USDT  = CurrencyPair{EOSBULL, USDT}
	EOS_USDT      = CurrencyPair{EOS, USDT}
	ERD_USDT      = CurrencyPair{ERD, USDT}
	ETC_USDT      = CurrencyPair{ETC, USDT}
	ETHBEAR_USDT  = CurrencyPair{ETHBEAR, USDT}
	ETHBULL_USDT  = CurrencyPair{ETHBULL, USDT}
	ETHDOWN_USDT  = CurrencyPair{ETHDOWN, USDT}
	ETHUP_USDT    = CurrencyPair{ETHUP, USDT}
	ETH_USDT      = CurrencyPair{ETH, USDT}
	EUR_USDT      = CurrencyPair{EUR, USDT}
	FET_USDT      = CurrencyPair{FET, USDT}
	FTM_USDT      = CurrencyPair{FTM, USDT}
	FTT_USDT      = CurrencyPair{FTT, USDT}
	FUN_USDT      = CurrencyPair{FUN, USDT}
	GBP_USDT      = CurrencyPair{GBP, USDT}
	GTO_USDT      = CurrencyPair{GTO, USDT}
	GXS_USDT      = CurrencyPair{GXS, USDT}
	HBAR_USDT     = CurrencyPair{HBAR, USDT}
	HC_USDT       = CurrencyPair{HC, USDT}
	HIVE_USDT     = CurrencyPair{HIVE, USDT}
	HOT_USDT      = CurrencyPair{HOT, USDT}
	ICX_USDT      = CurrencyPair{ICX, USDT}
	IOST_USDT     = CurrencyPair{IOST, USDT}
	IOTA_USDT     = CurrencyPair{IOTA, USDT}
	IOTX_USDT     = CurrencyPair{IOTX, USDT}
	IRIS_USDT     = CurrencyPair{IRIS, USDT}
	JST_USDT      = CurrencyPair{JST, USDT}
	KAVA_USDT     = CurrencyPair{KAVA, USDT}
	KEY_USDT      = CurrencyPair{KEY, USDT}
	KMD_USDT      = CurrencyPair{KMD, USDT}
	KNC_USDT      = CurrencyPair{KNC, USDT}
	LEND_USDT     = CurrencyPair{LEND, USDT}
	LINKDOWN_USDT = CurrencyPair{LINKDOWN, USDT}
	LINKUP_USDT   = CurrencyPair{LINKUP, USDT}
	LINK_USDT     = CurrencyPair{LINK, USDT}
	LRC_USDT      = CurrencyPair{LRC, USDT}
	LSK_USDT      = CurrencyPair{LSK, USDT}
	LTC_USDT      = CurrencyPair{LTC, USDT}
	LTO_USDT      = CurrencyPair{LTO, USDT}
	MANA_USDT     = CurrencyPair{MANA, USDT}
	MATIC_USDT    = CurrencyPair{MATIC, USDT}
	MBL_USDT      = CurrencyPair{MBL, USDT}
	MCO_USDT      = CurrencyPair{MCO, USDT}
	MDT_USDT      = CurrencyPair{MDT, USDT}
	MFT_USDT      = CurrencyPair{MFT, USDT}
	MITH_USDT     = CurrencyPair{MITH, USDT}
	MKR_USDT      = CurrencyPair{MKR, USDT}
	MTL_USDT      = CurrencyPair{MTL, USDT}
	NANO_USDT     = CurrencyPair{NANO, USDT}
	NEO_USDT      = CurrencyPair{NEO, USDT}
	NKN_USDT      = CurrencyPair{NKN, USDT}
	NPXS_USDT     = CurrencyPair{NPXS, USDT}
	NULS_USDT     = CurrencyPair{NULS, USDT}
	OGN_USDT      = CurrencyPair{OGN, USDT}
	OMG_USDT      = CurrencyPair{OMG, USDT}
	ONE_USDT      = CurrencyPair{ONE, USDT}
	ONG_USDT      = CurrencyPair{ONG, USDT}
	ONT_USDT      = CurrencyPair{ONT, USDT}
	PAX_USDT      = CurrencyPair{PAX, USDT}
	PERL_USDT     = CurrencyPair{PERL, USDT}
	PNT_USDT      = CurrencyPair{PNT, USDT}
	QTUM_USDT     = CurrencyPair{QTUM, USDT}
	REN_USDT      = CurrencyPair{REN, USDT}
	REP_USDT      = CurrencyPair{REP, USDT}
	RLC_USDT      = CurrencyPair{RLC, USDT}
	RVN_USDT      = CurrencyPair{RVN, USDT}
	SC_USDT       = CurrencyPair{SC, USDT}
	SNX_USDT      = CurrencyPair{SNX, USDT}
	SOL_USDT      = CurrencyPair{SOL, USDT}
	SRM_USDT      = CurrencyPair{SRM, USDT}
	STMX_USDT     = CurrencyPair{STMX, USDT}
	STORJ_USDT    = CurrencyPair{STORJ, USDT}
	STORM_USDT    = CurrencyPair{STORM, USDT}
	STPT_USDT     = CurrencyPair{STPT, USDT}
	STRAT_USDT    = CurrencyPair{STRAT, USDT}
	STX_USDT      = CurrencyPair{STX, USDT}
	SXP_USDT      = CurrencyPair{SXP, USDT}
	TCT_USDT      = CurrencyPair{TCT, USDT}
	TFUEL_USDT    = CurrencyPair{TFUEL, USDT}
	THETA_USDT    = CurrencyPair{THETA, USDT}
	TOMO_USDT     = CurrencyPair{TOMO, USDT}
	TROY_USDT     = CurrencyPair{TROY, USDT}
	TRX_USDT      = CurrencyPair{TRX, USDT}
	TUSD_USDT     = CurrencyPair{TUSD, USDT}
	USDC_USDT     = CurrencyPair{USDC, USDT}
	USDSB_USDT    = CurrencyPair{USDSB, USDT}
	USDS_USDT     = CurrencyPair{USDS, USDT}
	VEN_USDT      = CurrencyPair{VEN, USDT}
	VET_USDT      = CurrencyPair{VET, USDT}
	VITE_USDT     = CurrencyPair{VITE, USDT}
	VTHO_USDT     = CurrencyPair{VTHO, USDT}
	WAN_USDT      = CurrencyPair{WAN, USDT}
	WAVES_USDT    = CurrencyPair{WAVES, USDT}
	WIN_USDT      = CurrencyPair{WIN, USDT}
	WRX_USDT      = CurrencyPair{WRX, USDT}
	WTC_USDT      = CurrencyPair{WTC, USDT}
	XLM_USDT      = CurrencyPair{XLM, USDT}
	XMR_USDT      = CurrencyPair{XMR, USDT}
	XRPBEAR_USDT  = CurrencyPair{XRPBEAR, USDT}
	XRPBULL_USDT  = CurrencyPair{XRPBULL, USDT}
	XRP_USDT      = CurrencyPair{XRP, USDT}
	XTZDOWN_USDT  = CurrencyPair{XTZDOWN, USDT}
	XTZUP_USDT    = CurrencyPair{XTZUP, USDT}
	XTZ_USDT      = CurrencyPair{XTZ, USDT}
	XZC_USDT      = CurrencyPair{XZC, USDT}
	YFI_USDT      = CurrencyPair{YFI, USDT}
	ZEC_USDT      = CurrencyPair{ZEC, USDT}
	ZEN_USDT      = CurrencyPair{ZEN, USDT}
	ZIL_USDT      = CurrencyPair{ZIL, USDT}
	ZRX_USDT      = CurrencyPair{ZRX, USDT}

	UNKNOWN_PAIR = CurrencyPair{UNKNOWN, UNKNOWN}
)

func (c CurrencyPair) String() string {
	return c.ToSymbol("_")
}

func (c CurrencyPair) Eq(c2 CurrencyPair) bool {
	return c.String() == c2.String()
}

func (c Currency) AdaptBchToBcc() Currency {
	if c.Symbol == "BCH" || c.Symbol == "bch" {
		return BCC
	}
	return c
}

func (c Currency) AdaptBccToBch() Currency {
	if c.Symbol == "BCC" || c.Symbol == "bcc" {
		return BCH
	}
	return c
}

func NewCurrency(symbol, desc string) Currency {
	switch symbol {
	case "cny", "CNY":
		return CNY
	case "usdt", "USDT":
		return USDT
	case "usd", "USD":
		return USD
	case "usdc", "USDC":
		return USDC
	case "pax", "PAX":
		return PAX
	case "jpy", "JPY":
		return JPY
	case "krw", "KRW":
		return KRW
	case "eur", "EUR":
		return EUR
	case "btc", "BTC":
		return BTC
	case "xbt", "XBT":
		return XBT
	case "bch", "BCH":
		return BCH
	case "bcc", "BCC":
		return BCC
	case "ltc", "LTC":
		return LTC
	case "sc", "SC":
		return SC
	case "ans", "ANS":
		return ANS
	case "neo", "NEO":
		return NEO
	case "okb", "OKB":
		return OKB
	case "ht", "HT":
		return HT
	case "bnb", "BNB":
		return BNB
	case "trx", "TRX":
		return TRX
	default:
		return Currency{strings.ToUpper(symbol), desc}
	}
}

func NewCurrencyPair(currencyA Currency, currencyB Currency) CurrencyPair {
	return CurrencyPair{currencyA, currencyB}
}

func NewCurrencyPair2(currencyPairSymbol string) CurrencyPair {
	return NewCurrencyPair3(currencyPairSymbol, "_")
}

func NewCurrencyPair3(currencyPairSymbol string, sep string) CurrencyPair {
	currencys := strings.Split(currencyPairSymbol, sep)
	if len(currencys) >= 2 {
		return CurrencyPair{NewCurrency(currencys[0], ""),
			NewCurrency(currencys[1], "")}
	}
	return UNKNOWN_PAIR
}

func (pair CurrencyPair) ToSymbol(joinChar string) string {
	return strings.Join([]string{pair.CurrencyA.Symbol, pair.CurrencyB.Symbol}, joinChar)
}

func (pair CurrencyPair) ToSymbol2(joinChar string) string {
	return strings.Join([]string{pair.CurrencyB.Symbol, pair.CurrencyA.Symbol}, joinChar)
}

func (pair CurrencyPair) AdaptUsdtToUsd() CurrencyPair {
	CurrencyB := pair.CurrencyB
	if pair.CurrencyB.Eq(USDT) {
		CurrencyB = USD
	}
	return CurrencyPair{pair.CurrencyA, CurrencyB}
}

func (pair CurrencyPair) AdaptUsdToUsdt() CurrencyPair {
	CurrencyB := pair.CurrencyB
	if pair.CurrencyB.Eq(USD) {
		CurrencyB = USDT
	}
	return CurrencyPair{pair.CurrencyA, CurrencyB}
}

//It is currently applicable to binance and zb
func (pair CurrencyPair) AdaptBchToBcc() CurrencyPair {
	CurrencyA := pair.CurrencyA
	if pair.CurrencyA.Eq(BCH) {
		CurrencyA = BCC
	}
	return CurrencyPair{CurrencyA, pair.CurrencyB}
}

func (pair CurrencyPair) AdaptBccToBch() CurrencyPair {
	if pair.CurrencyA.Eq(BCC) {
		return CurrencyPair{BCH, pair.CurrencyB}
	}
	return pair
}

//for to symbol lower , Not practical '==' operation method
func (pair CurrencyPair) ToLower() CurrencyPair {
	return CurrencyPair{Currency{strings.ToLower(pair.CurrencyA.Symbol), pair.CurrencyA.Desc},
		Currency{strings.ToLower(pair.CurrencyB.Symbol), pair.CurrencyB.Desc}}
}

func (pair CurrencyPair) Reverse() CurrencyPair {
	return CurrencyPair{pair.CurrencyB, pair.CurrencyA}
}

// -------------------------------------------------------------------------------------------------
// CurrencyMap 币种字典
//var CurrencyMap = map[string]Currency{
//	"CNY":  CNY,
//	"USD":  USD,
//	"USDT": USDT,
//	"PAX":  PAX,
//	"USDC": USDC,
//	"EUR":  EUR,
//	"KRW":  KRW,
//	"JPY":  JPY,
//	"BTC":  BTC,
//	"XBT":  XBT,
//	"BCC":  BCC,
//	"BCH":  BCH,
//	"BCX":  BCX,
//	"LTC":  LTC,
//	"ETH":  ETH,
//	"ETC":  ETC,
//	"EOS":  EOS,
//	"BTS":  BTS,
//	"QTUM": QTUM,
//	"SC":   SC,
//	"ANS":  ANS,
//	"ZEC":  ZEC,
//	"DCR":  DCR,
//	"XRP":  XRP,
//	"BTG":  BTG,
//	"BCD":  BCD,
//	"NEO":  NEO,
//	"HSR":  HSR,
//	"BSV":  BSV,
//	"OKB":  OKB,
//	"HT":   HT,
//	"BNB":  BNB,
//	"TRX":  TRX,
//}

// CurrencyPairMap 币种对字典，key支持无间隔和下划线间隔
var CurrencyPairMap = map[string]CurrencyPair{
	"BTCCNY":  BTC_CNY,
	"LTCCNY":  LTC_CNY,
	"BCCCNY":  BCC_CNY,
	"ETHCNY":  ETH_CNY,
	"ETCCNY":  ETC_CNY,
	"EOSCNY":  EOS_CNY,
	"BTSCNY":  BTS_CNY,
	"QTUMCNY": QTUM_CNY,
	"SCCNY":   SC_CNY,
	"ANSCNY":  ANS_CNY,
	"ZECCNY":  ZEC_CNY,
	"BTCKRW":  BTC_KRW,
	"ETHKRW":  ETH_KRW,
	"ETCKRW":  ETC_KRW,
	"LTCKRW":  LTC_KRW,
	"BCHKRW":  BCH_KRW,
	"BTCUSD":  BTC_USD,
	"LTCUSD":  LTC_USD,
	"ETHUSD":  ETH_USD,
	"ETCUSD":  ETC_USD,
	"BCHUSD":  BCH_USD,
	"BCCUSD":  BCC_USD,
	"XRPUSD":  XRP_USD,
	"BCDUSD":  BCD_USD,
	"EOSUSD":  EOS_USD,
	"BTGUSD":  BTG_USD,
	"BSVUSD":  BSV_USD,
	"BCDUSDT": BCD_USDT,
	"HSRUSDT": HSR_USDT,
	"BSVUSDT": BSV_USDT,
	"OKBUSDT": OKB_USDT,
	"HTUSDT":  HT_USDT,
	"XRPEUR":  XRP_EUR,
	"BTCJPY":  BTC_JPY,
	"LTCJPY":  LTC_JPY,
	"ETHJPY":  ETH_JPY,
	"ETCJPY":  ETC_JPY,
	"BCHJPY":  BCH_JPY,
	"LTCBTC":  LTC_BTC,
	"ETHBTC":  ETH_BTC,
	"ETCBTC":  ETC_BTC,
	"BCCBTC":  BCC_BTC,
	"BCHBTC":  BCH_BTC,
	"DCRBTC":  DCR_BTC,
	"XRPBTC":  XRP_BTC,
	"BTGBTC":  BTG_BTC,
	"BCDBTC":  BCD_BTC,
	"NEOBTC":  NEO_BTC,
	"EOSBTC":  EOS_BTC,
	"HSRBTC":  HSR_BTC,
	"BSVBTC":  BSV_BTC,
	"OKBBTC":  OKB_BTC,
	"HTBTC":   HT_BTC,
	"BNBBTC":  BNB_BTC,
	"TRXBTC":  TRX_BTC,
	"ETCETH":  ETC_ETH,
	"EOSETH":  EOS_ETH,
	"ZECETH":  ZEC_ETH,
	"NEOETH":  NEO_ETH,
	"HSRETH":  HSR_ETH,
	"LTCETH":  LTC_ETH,

	"ADADOWNUSDT":  ADADOWN_USDT,
	"ADAUPUSDT":    ADAUP_USDT,
	"ADAUSDT":      ADA_USDT,
	"AIONUSDT":     AION_USDT,
	"ALGOUSDT":     ALGO_USDT,
	"ANKRUSDT":     ANKR_USDT,
	"ARDRUSDT":     ARDR_USDT,
	"ARPAUSDT":     ARPA_USDT,
	"ATOMUSDT":     ATOM_USDT,
	"AUDUSDT":      AUD_USDT,
	"BALUSDT":      BAL_USDT,
	"BANDUSDT":     BAND_USDT,
	"BATUSDT":      BAT_USDT,
	"BCCUSDT":      BCC_USDT,
	"BCHABCUSDT":   BCHABC_USDT,
	"BCHSVUSDT":    BCHSV_USDT,
	"BCHUSDT":      BCH_USDT,
	"BEAMUSDT":     BEAM_USDT,
	"BEARUSDT":     BEAR_USDT,
	"BKRWUSDT":     BKRW_USDT,
	"BLZUSDT":      BLZ_USDT,
	"BNBBEARUSDT":  BNBBEAR_USDT,
	"BNBBULLUSDT":  BNBBULL_USDT,
	"BNBDOWNUSDT":  BNBDOWN_USDT,
	"BNBUPUSDT":    BNBUP_USDT,
	"BNBUSDT":      BNB_USDT,
	"BNTUSDT":      BNT_USDT,
	"BTCDOWNUSDT":  BTCDOWN_USDT,
	"BTCUPUSDT":    BTCUP_USDT,
	"BTCUSDT":      BTC_USDT,
	"BTSUSDT":      BTS_USDT,
	"BTTUSDT":      BTT_USDT,
	"BULLUSDT":     BULL_USDT,
	"BUSDUSDT":     BUSD_USDT,
	"CELRUSDT":     CELR_USDT,
	"CHRUSDT":      CHR_USDT,
	"CHZUSDT":      CHZ_USDT,
	"COCOSUSDT":    COCOS_USDT,
	"COMPUSDT":     COMP_USDT,
	"COSUSDT":      COS_USDT,
	"COTIUSDT":     COTI_USDT,
	"CTSIUSDT":     CTSI_USDT,
	"CTXCUSDT":     CTXC_USDT,
	"CVCUSDT":      CVC_USDT,
	"DAIUSDT":      DAI_USDT,
	"DASHUSDT":     DASH_USDT,
	"DATAUSDT":     DATA_USDT,
	"DCRUSDT":      DCR_USDT,
	"DENTUSDT":     DENT_USDT,
	"DGBUSDT":      DGB_USDT,
	"DOCKUSDT":     DOCK_USDT,
	"DOGEUSDT":     DOGE_USDT,
	"DOTUSDT":      DOT_USDT,
	"DREPUSDT":     DREP_USDT,
	"DUSKUSDT":     DUSK_USDT,
	"ENJUSDT":      ENJ_USDT,
	"EOSBEARUSDT":  EOSBEAR_USDT,
	"EOSBULLUSDT":  EOSBULL_USDT,
	"EOSUSDT":      EOS_USDT,
	"ERDUSDT":      ERD_USDT,
	"ETCUSDT":      ETC_USDT,
	"ETHBEARUSDT":  ETHBEAR_USDT,
	"ETHBULLUSDT":  ETHBULL_USDT,
	"ETHDOWNUSDT":  ETHDOWN_USDT,
	"ETHUPUSDT":    ETHUP_USDT,
	"ETHUSDT":      ETH_USDT,
	"EURUSDT":      EUR_USDT,
	"FETUSDT":      FET_USDT,
	"FTMUSDT":      FTM_USDT,
	"FTTUSDT":      FTT_USDT,
	"FUNUSDT":      FUN_USDT,
	"GBPUSDT":      GBP_USDT,
	"GTOUSDT":      GTO_USDT,
	"GXSUSDT":      GXS_USDT,
	"HBARUSDT":     HBAR_USDT,
	"HCUSDT":       HC_USDT,
	"HIVEUSDT":     HIVE_USDT,
	"HOTUSDT":      HOT_USDT,
	"ICXUSDT":      ICX_USDT,
	"IOSTUSDT":     IOST_USDT,
	"IOTAUSDT":     IOTA_USDT,
	"IOTXUSDT":     IOTX_USDT,
	"IRISUSDT":     IRIS_USDT,
	"JSTUSDT":      JST_USDT,
	"KAVAUSDT":     KAVA_USDT,
	"KEYUSDT":      KEY_USDT,
	"KMDUSDT":      KMD_USDT,
	"KNCUSDT":      KNC_USDT,
	"LENDUSDT":     LEND_USDT,
	"LINKDOWNUSDT": LINKDOWN_USDT,
	"LINKUPUSDT":   LINKUP_USDT,
	"LINKUSDT":     LINK_USDT,
	"LRCUSDT":      LRC_USDT,
	"LSKUSDT":      LSK_USDT,
	"LTCUSDT":      LTC_USDT,
	"LTOUSDT":      LTO_USDT,
	"MANAUSDT":     MANA_USDT,
	"MATICUSDT":    MATIC_USDT,
	"MBLUSDT":      MBL_USDT,
	"MCOUSDT":      MCO_USDT,
	"MDTUSDT":      MDT_USDT,
	"MFTUSDT":      MFT_USDT,
	"MITHUSDT":     MITH_USDT,
	"MKRUSDT":      MKR_USDT,
	"MTLUSDT":      MTL_USDT,
	"NANOUSDT":     NANO_USDT,
	"NEOUSDT":      NEO_USDT,
	"NKNUSDT":      NKN_USDT,
	"NPXSUSDT":     NPXS_USDT,
	"NULSUSDT":     NULS_USDT,
	"OGNUSDT":      OGN_USDT,
	"OMGUSDT":      OMG_USDT,
	"ONEUSDT":      ONE_USDT,
	"ONGUSDT":      ONG_USDT,
	"ONTUSDT":      ONT_USDT,
	"PAXUSDT":      PAX_USDT,
	"PERLUSDT":     PERL_USDT,
	"PNTUSDT":      PNT_USDT,
	"QTUMUSDT":     QTUM_USDT,
	"RENUSDT":      REN_USDT,
	"REPUSDT":      REP_USDT,
	"RLCUSDT":      RLC_USDT,
	"RVNUSDT":      RVN_USDT,
	"SCUSDT":       SC_USDT,
	"SNXUSDT":      SNX_USDT,
	"SOLUSDT":      SOL_USDT,
	"SRMUSDT":      SRM_USDT,
	"STMXUSDT":     STMX_USDT,
	"STORJUSDT":    STORJ_USDT,
	"STORMUSDT":    STORM_USDT,
	"STPTUSDT":     STPT_USDT,
	"STRATUSDT":    STRAT_USDT,
	"STXUSDT":      STX_USDT,
	"SXPUSDT":      SXP_USDT,
	"TCTUSDT":      TCT_USDT,
	"TFUELUSDT":    TFUEL_USDT,
	"THETAUSDT":    THETA_USDT,
	"TOMOUSDT":     TOMO_USDT,
	"TROYUSDT":     TROY_USDT,
	"TRXUSDT":      TRX_USDT,
	"TUSDUSDT":     TUSD_USDT,
	"USDCUSDT":     USDC_USDT,
	"USDSBUSDT":    USDSB_USDT,
	"USDSUSDT":     USDS_USDT,
	"VENUSDT":      VEN_USDT,
	"VETUSDT":      VET_USDT,
	"VITEUSDT":     VITE_USDT,
	"VTHOUSDT":     VTHO_USDT,
	"WANUSDT":      WAN_USDT,
	"WAVESUSDT":    WAVES_USDT,
	"WINUSDT":      WIN_USDT,
	"WRXUSDT":      WRX_USDT,
	"WTCUSDT":      WTC_USDT,
	"XLMUSDT":      XLM_USDT,
	"XMRUSDT":      XMR_USDT,
	"XRPBEARUSDT":  XRPBEAR_USDT,
	"XRPBULLUSDT":  XRPBULL_USDT,
	"XRPUSDT":      XRP_USDT,
	"XTZDOWNUSDT":  XTZDOWN_USDT,
	"XTZUPUSDT":    XTZUP_USDT,
	"XTZUSDT":      XTZ_USDT,
	"XZCUSDT":      XZC_USDT,
	"YFIUSDT":      YFI_USDT,
	"ZECUSDT":      ZEC_USDT,
	"ZENUSDT":      ZEN_USDT,
	"ZILUSDT":      ZIL_USDT,
	"ZRXUSDT":      ZRX_USDT,

	"BTC_CNY":  BTC_CNY,
	"LTC_CNY":  LTC_CNY,
	"BCC_CNY":  BCC_CNY,
	"ETH_CNY":  ETH_CNY,
	"ETC_CNY":  ETC_CNY,
	"EOS_CNY":  EOS_CNY,
	"BTS_CNY":  BTS_CNY,
	"QTUM_CNY": QTUM_CNY,
	"SC_CNY":   SC_CNY,
	"ANS_CNY":  ANS_CNY,
	"ZEC_CNY":  ZEC_CNY,
	"BTC_KRW":  BTC_KRW,
	"ETH_KRW":  ETH_KRW,
	"ETC_KRW":  ETC_KRW,
	"LTC_KRW":  LTC_KRW,
	"BCH_KRW":  BCH_KRW,
	"BTC_USD":  BTC_USD,
	"LTC_USD":  LTC_USD,
	"ETH_USD":  ETH_USD,
	"ETC_USD":  ETC_USD,
	"BCH_USD":  BCH_USD,
	"BCC_USD":  BCC_USD,
	"XRP_USD":  XRP_USD,
	"BCD_USD":  BCD_USD,
	"EOS_USD":  EOS_USD,
	"BTG_USD":  BTG_USD,
	"BSV_USD":  BSV_USD,
	"BCD_USDT": BCD_USDT,
	"HSR_USDT": HSR_USDT,
	"BSV_USDT": BSV_USDT,
	"OKB_USDT": OKB_USDT,
	"HT_USDT":  HT_USDT,
	"XRP_EUR":  XRP_EUR,
	"BTC_JPY":  BTC_JPY,
	"LTC_JPY":  LTC_JPY,
	"ETH_JPY":  ETH_JPY,
	"ETC_JPY":  ETC_JPY,
	"BCH_JPY":  BCH_JPY,
	"LTC_BTC":  LTC_BTC,
	"ETH_BTC":  ETH_BTC,
	"ETC_BTC":  ETC_BTC,
	"BCC_BTC":  BCC_BTC,
	"BCH_BTC":  BCH_BTC,
	"DCR_BTC":  DCR_BTC,
	"XRP_BTC":  XRP_BTC,
	"BTG_BTC":  BTG_BTC,
	"BCD_BTC":  BCD_BTC,
	"NEO_BTC":  NEO_BTC,
	"EOS_BTC":  EOS_BTC,
	"HSR_BTC":  HSR_BTC,
	"BSV_BTC":  BSV_BTC,
	"OKB_BTC":  OKB_BTC,
	"HT_BTC":   HT_BTC,
	"BNB_BTC":  BNB_BTC,
	"TRX_BTC":  TRX_BTC,
	"ETC_ETH":  ETC_ETH,
	"EOS_ETH":  EOS_ETH,
	"ZEC_ETH":  ZEC_ETH,
	"NEO_ETH":  NEO_ETH,
	"HSR_ETH":  HSR_ETH,
	"LTC_ETH":  LTC_ETH,

	"ADADOWN_USDT":  ADADOWN_USDT,
	"ADAUP_USDT":    ADAUP_USDT,
	"ADA_USDT":      ADA_USDT,
	"AION_USDT":     AION_USDT,
	"ALGO_USDT":     ALGO_USDT,
	"ANKR_USDT":     ANKR_USDT,
	"ARDR_USDT":     ARDR_USDT,
	"ARPA_USDT":     ARPA_USDT,
	"ATOM_USDT":     ATOM_USDT,
	"AUD_USDT":      AUD_USDT,
	"BAL_USDT":      BAL_USDT,
	"BAND_USDT":     BAND_USDT,
	"BAT_USDT":      BAT_USDT,
	"BCC_USDT":      BCC_USDT,
	"BCHABC_USDT":   BCHABC_USDT,
	"BCHSV_USDT":    BCHSV_USDT,
	"BCH_USDT":      BCH_USDT,
	"BEAM_USDT":     BEAM_USDT,
	"BEAR_USDT":     BEAR_USDT,
	"BKRW_USDT":     BKRW_USDT,
	"BLZ_USDT":      BLZ_USDT,
	"BNBBEAR_USDT":  BNBBEAR_USDT,
	"BNBBULL_USDT":  BNBBULL_USDT,
	"BNBDOWN_USDT":  BNBDOWN_USDT,
	"BNBUP_USDT":    BNBUP_USDT,
	"BNB_USDT":      BNB_USDT,
	"BNT_USDT":      BNT_USDT,
	"BTCDOWN_USDT":  BTCDOWN_USDT,
	"BTCUP_USDT":    BTCUP_USDT,
	"BTC_USDT":      BTC_USDT,
	"BTS_USDT":      BTS_USDT,
	"BTT_USDT":      BTT_USDT,
	"BULL_USDT":     BULL_USDT,
	"BUSD_USDT":     BUSD_USDT,
	"CELR_USDT":     CELR_USDT,
	"CHR_USDT":      CHR_USDT,
	"CHZ_USDT":      CHZ_USDT,
	"COCOS_USDT":    COCOS_USDT,
	"COMP_USDT":     COMP_USDT,
	"COS_USDT":      COS_USDT,
	"COTI_USDT":     COTI_USDT,
	"CTSI_USDT":     CTSI_USDT,
	"CTXC_USDT":     CTXC_USDT,
	"CVC_USDT":      CVC_USDT,
	"DAI_USDT":      DAI_USDT,
	"DASH_USDT":     DASH_USDT,
	"DATA_USDT":     DATA_USDT,
	"DCR_USDT":      DCR_USDT,
	"DENT_USDT":     DENT_USDT,
	"DGB_USDT":      DGB_USDT,
	"DOCK_USDT":     DOCK_USDT,
	"DOGE_USDT":     DOGE_USDT,
	"DOT_USDT":      DOT_USDT,
	"DREP_USDT":     DREP_USDT,
	"DUSK_USDT":     DUSK_USDT,
	"ENJ_USDT":      ENJ_USDT,
	"EOSBEAR_USDT":  EOSBEAR_USDT,
	"EOSBULL_USDT":  EOSBULL_USDT,
	"EOS_USDT":      EOS_USDT,
	"ERD_USDT":      ERD_USDT,
	"ETC_USDT":      ETC_USDT,
	"ETHBEAR_USDT":  ETHBEAR_USDT,
	"ETHBULL_USDT":  ETHBULL_USDT,
	"ETHDOWN_USDT":  ETHDOWN_USDT,
	"ETHUP_USDT":    ETHUP_USDT,
	"ETH_USDT":      ETH_USDT,
	"EUR_USDT":      EUR_USDT,
	"FET_USDT":      FET_USDT,
	"FTM_USDT":      FTM_USDT,
	"FTT_USDT":      FTT_USDT,
	"FUN_USDT":      FUN_USDT,
	"GBP_USDT":      GBP_USDT,
	"GTO_USDT":      GTO_USDT,
	"GXS_USDT":      GXS_USDT,
	"HBAR_USDT":     HBAR_USDT,
	"HC_USDT":       HC_USDT,
	"HIVE_USDT":     HIVE_USDT,
	"HOT_USDT":      HOT_USDT,
	"ICX_USDT":      ICX_USDT,
	"IOST_USDT":     IOST_USDT,
	"IOTA_USDT":     IOTA_USDT,
	"IOTX_USDT":     IOTX_USDT,
	"IRIS_USDT":     IRIS_USDT,
	"JST_USDT":      JST_USDT,
	"KAVA_USDT":     KAVA_USDT,
	"KEY_USDT":      KEY_USDT,
	"KMD_USDT":      KMD_USDT,
	"KNC_USDT":      KNC_USDT,
	"LEND_USDT":     LEND_USDT,
	"LINKDOWN_USDT": LINKDOWN_USDT,
	"LINKUP_USDT":   LINKUP_USDT,
	"LINK_USDT":     LINK_USDT,
	"LRC_USDT":      LRC_USDT,
	"LSK_USDT":      LSK_USDT,
	"LTC_USDT":      LTC_USDT,
	"LTO_USDT":      LTO_USDT,
	"MANA_USDT":     MANA_USDT,
	"MATIC_USDT":    MATIC_USDT,
	"MBL_USDT":      MBL_USDT,
	"MCO_USDT":      MCO_USDT,
	"MDT_USDT":      MDT_USDT,
	"MFT_USDT":      MFT_USDT,
	"MITH_USDT":     MITH_USDT,
	"MKR_USDT":      MKR_USDT,
	"MTL_USDT":      MTL_USDT,
	"NANO_USDT":     NANO_USDT,
	"NEO_USDT":      NEO_USDT,
	"NKN_USDT":      NKN_USDT,
	"NPXS_USDT":     NPXS_USDT,
	"NULS_USDT":     NULS_USDT,
	"OGN_USDT":      OGN_USDT,
	"OMG_USDT":      OMG_USDT,
	"ONE_USDT":      ONE_USDT,
	"ONG_USDT":      ONG_USDT,
	"ONT_USDT":      ONT_USDT,
	"PAX_USDT":      PAX_USDT,
	"PERL_USDT":     PERL_USDT,
	"PNT_USDT":      PNT_USDT,
	"QTUM_USDT":     QTUM_USDT,
	"REN_USDT":      REN_USDT,
	"REP_USDT":      REP_USDT,
	"RLC_USDT":      RLC_USDT,
	"RVN_USDT":      RVN_USDT,
	"SC_USDT":       SC_USDT,
	"SNX_USDT":      SNX_USDT,
	"SOL_USDT":      SOL_USDT,
	"SRM_USDT":      SRM_USDT,
	"STMX_USDT":     STMX_USDT,
	"STORJ_USDT":    STORJ_USDT,
	"STORM_USDT":    STORM_USDT,
	"STPT_USDT":     STPT_USDT,
	"STRAT_USDT":    STRAT_USDT,
	"STX_USDT":      STX_USDT,
	"SXP_USDT":      SXP_USDT,
	"TCT_USDT":      TCT_USDT,
	"TFUEL_USDT":    TFUEL_USDT,
	"THETA_USDT":    THETA_USDT,
	"TOMO_USDT":     TOMO_USDT,
	"TROY_USDT":     TROY_USDT,
	"TRX_USDT":      TRX_USDT,
	"TUSD_USDT":     TUSD_USDT,
	"USDC_USDT":     USDC_USDT,
	"USDSB_USDT":    USDSB_USDT,
	"USDS_USDT":     USDS_USDT,
	"VEN_USDT":      VEN_USDT,
	"VET_USDT":      VET_USDT,
	"VITE_USDT":     VITE_USDT,
	"VTHO_USDT":     VTHO_USDT,
	"WAN_USDT":      WAN_USDT,
	"WAVES_USDT":    WAVES_USDT,
	"WIN_USDT":      WIN_USDT,
	"WRX_USDT":      WRX_USDT,
	"WTC_USDT":      WTC_USDT,
	"XLM_USDT":      XLM_USDT,
	"XMR_USDT":      XMR_USDT,
	"XRPBEAR_USDT":  XRPBEAR_USDT,
	"XRPBULL_USDT":  XRPBULL_USDT,
	"XRP_USDT":      XRP_USDT,
	"XTZDOWN_USDT":  XTZDOWN_USDT,
	"XTZUP_USDT":    XTZUP_USDT,
	"XTZ_USDT":      XTZ_USDT,
	"XZC_USDT":      XZC_USDT,
	"YFI_USDT":      YFI_USDT,
	"ZEC_USDT":      ZEC_USDT,
	"ZEN_USDT":      ZEN_USDT,
	"ZIL_USDT":      ZIL_USDT,
	"ZRX_USDT":      ZRX_USDT,
}

// SplitSymbol 把品种拆分为交易币和锚定币
func SplitSymbol(symbol string) (string, string) {
	cp := GetCurrencyPair(symbol)
	if cp.CurrencyA.Symbol != "" {
		return strings.ToLower(cp.CurrencyA.Symbol), strings.ToLower(cp.CurrencyB.Symbol)
	}

	symbol = strings.ToUpper(symbol)
	v, ok := CurrencyPairMap[symbol]
	if !ok {
		return strings.ToLower(symbol), strings.ToLower(symbol)
	}
	return strings.ToLower(v.CurrencyA.Symbol), strings.ToLower(v.CurrencyB.Symbol)
}

// GetCurrency 获取交易币
func GetTradeCurrency(symbol string) string {
	tradeCurrency, _ := SplitSymbol(symbol)
	return tradeCurrency
}

// GetAnchorCurrency 获取锚定币
func GetAnchorCurrency(symbol string) string {
	_, anchorCurrency := SplitSymbol(symbol)
	return anchorCurrency
}

// CurrencyPairMap 交易对
var CurrencyPairSyncMap = new(sync.Map)

// GetCurrencyPair 获取限制值
func GetCurrencyPair(symbol string) *CurrencyPair {
	symbol = strings.ToLower(symbol)
	if value, ok := CurrencyPairSyncMap.Load(symbol); ok {
		return value.(*CurrencyPair)
	}

	return &CurrencyPair{}
}
