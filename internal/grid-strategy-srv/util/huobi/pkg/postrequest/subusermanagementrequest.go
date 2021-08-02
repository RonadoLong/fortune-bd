package postrequest

type SubUserManagementRequest struct {
	SubUid int64  `json:"subUid"`
	Action string `json:"action"`
}
