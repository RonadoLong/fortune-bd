package getrequest

const (
	DEPTH_SIZE_FIVE   = 5
	DEPTH_SIZE_TEN    = 10
	DEPTH_SIZE_TWENTY = 20
)

const (
	STEP0 = "step0"
	STEP1 = "step1"
	STEP2 = "step2"
	STEP3 = "step3"
	STEP4 = "step4"
	STEP5 = "step5"
)

type GetDepthOptionalRequest struct {
	Size int
}
