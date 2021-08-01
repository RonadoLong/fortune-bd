package account

type SubUserManagementResponse struct {
	Code int `json:"code"`
	Data *SubUserManagement
}
type SubUserManagement struct {
	SubUid    int64  `json:"subUid"`
	UserState string `json:"userState"`
}
