package helper

type DocumentStatus string
type DocumentAction string

const (
	// Status
	StatusOpened   DocumentStatus = "opened"
	StatusApproved DocumentStatus = "approved"
	StatusReviewed DocumentStatus = "reviewed"
	StatusReturned DocumentStatus = "returned"
	StatusClosed   DocumentStatus = "closed"
	StatusRejected DocumentStatus = "rejected"

	// Action
	ActionCreate  DocumentAction = "create"
	ActionApprove DocumentAction = "approve"
	ActionReject  DocumentAction = "reject"
	ActionReview  DocumentAction = "review"
	ActionReturn  DocumentAction = "return"
)

func CanDoAction(status DocumentStatus, action DocumentAction, assignJob string) bool {

	if !canByStatus(status, action) {
		return false
	}

	if action != ActionReturn && action != ActionReject {
		if assignJob != string(action) {
			return false
		}
	}

	return true
}

func canByStatus(status DocumentStatus, action DocumentAction) bool {
	switch status {

	case StatusOpened:
		return action == ActionApprove ||
			action == ActionReject ||
			action == ActionReview

	case StatusApproved:
		return action == ActionApprove ||
			action == ActionReject ||
			action == ActionReview

	case StatusReviewed:
		return action == ActionApprove ||
			action == ActionReject ||
			action == ActionReview ||
			action == ActionReturn
	case StatusReturned:
		return action == ActionReview
	case StatusClosed, StatusRejected:
		return false

	default:
		return false
	}
}

func GetNextStatusByAction(action DocumentAction) string {
	switch action {

	case ActionApprove:
		return string(StatusApproved)

	case ActionReject:
		return string(StatusRejected)

	case ActionReview:
		return string(StatusReviewed)

	case ActionReturn:
		return string(StatusReturned)

	case ActionCreate:
		return string(StatusOpened)

	default:
		return "error"
	}
}

func ToDocumentStatus(status string) DocumentStatus {
	ds := DocumentStatus(status)

	switch ds {
	case StatusOpened,
		StatusApproved,
		StatusReviewed,
		StatusReturned,
		StatusClosed,
		StatusRejected:
		return ds
	default:
		return ""
	}
}
