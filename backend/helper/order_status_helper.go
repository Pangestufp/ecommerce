package helper

import "fmt"

const (
	StatusPendingPayment  = "PENDING_PAYMENT"
	StatusPayed           = "PAYED"
	StatusExpired         = "EXPIRED"
	StatusCancelledUnpaid = "CANCELLED_UNPAID"
	StatusAccepted        = "ACCEPTED"
	StatusPacked          = "PACKED"
	StatusShipped         = "SHIPPED"
	StatusCancelled       = "CANCELLED"
	StatusFinished        = "FINISHED"
)

func GetPendingPaymentStatus() string  { return StatusPendingPayment }
func GetPayedStatus() string           { return StatusPayed }
func GetExpiredStatus() string         { return StatusExpired }
func GetCancelledUnpaidStatus() string { return StatusCancelledUnpaid }
func GetAcceptedStatus() string        { return StatusAccepted }
func GetPackedStatus() string          { return StatusPacked }
func GetShippedStatus() string         { return StatusShipped }
func GetCancelledStatus() string       { return StatusCancelled }
func GetFinishedStatus() string        { return StatusFinished }

const (
	ActionPay             = "PAY"              // webhook payment gateway konfirmasi bayar
	ActionExpire          = "EXPIRE"           // cron tengah malam, order tidak dibayar
	ActionCancelUnpaid    = "CANCEL_UNPAID"    // seller cancel manual, order belum dibayar (>2 jam)
	ActionAccept          = "ACCEPT"           // seller terima order
	ActionPack            = "PACK"             // seller selesai kemas barang
	ActionShip            = "SHIP"             // seller input resi, serahkan ke kurir
	ActionCancel          = "CANCEL"           // seller cancel setelah dibayar (PAYED/ACCEPTED/PACKED)
	ActionConfirmReceived = "CONFIRM_RECEIVED" // pembeli klik "pesanan diterima"
	ActionAutoFinish      = "AUTO_FINISH"      // cron, X hari sejak ShippedAt
	ActionForceFinish     = "FORCE_FINISH"     // seller/admin paksa selesai (dispute, dsb)
)

// GetPayAction dst — helper getter per aksi.
func GetPayAction() string             { return ActionPay }
func GetExpireAction() string          { return ActionExpire }
func GetCancelUnpaidAction() string    { return ActionCancelUnpaid }
func GetAcceptAction() string          { return ActionAccept }
func GetPackAction() string            { return ActionPack }
func GetShipAction() string            { return ActionShip }
func GetCancelAction() string          { return ActionCancel }
func GetConfirmReceivedAction() string { return ActionConfirmReceived }
func GetAutoFinishAction() string      { return ActionAutoFinish }
func GetForceFinishAction() string     { return ActionForceFinish }

var transitions = map[string]map[string]string{
	StatusPendingPayment: {
		ActionPay:          StatusPayed,
		ActionExpire:       StatusExpired,
		ActionCancelUnpaid: StatusCancelledUnpaid,
	},
	StatusPayed: {
		ActionAccept: StatusAccepted,
		ActionCancel: StatusCancelled,
	},
	StatusAccepted: {
		ActionPack:   StatusPacked,
		ActionCancel: StatusCancelled,
	},
	StatusPacked: {
		ActionShip:   StatusShipped,
		ActionCancel: StatusCancelled,
	},
	StatusShipped: {
		ActionConfirmReceived: StatusFinished,
		ActionAutoFinish:      StatusFinished,
		ActionForceFinish:     StatusFinished,
	},

	// Status akhir — tidak ada aksi lanjutan yang valid.
	StatusExpired:         {},
	StatusCancelledUnpaid: {},
	StatusCancelled:       {},
	StatusFinished:        {},
}

func ApplyAction(currentStatus string, action string) (string, error) {
	actionsForStatus, statusExists := transitions[currentStatus]
	if !statusExists {
		return "", fmt.Errorf("status tidak dikenal: %s", currentStatus)
	}

	nextStatus, actionValid := actionsForStatus[action]
	if !actionValid {
		return "", fmt.Errorf("aksi %s tidak dapat dilakukan pada status %s", action, currentStatus)
	}

	return nextStatus, nil
}

func GetAvailableActions(currentStatus string) ([]string, error) {
	actionsForStatus, statusExists := transitions[currentStatus]
	if !statusExists {
		return nil, fmt.Errorf("status tidak dikenal: %s", currentStatus)
	}

	actions := make([]string, 0, len(actionsForStatus))
	for action := range actionsForStatus {
		actions = append(actions, action)
	}

	return actions, nil
}

func IsFinalStatus(currentStatus string) bool {
	actionsForStatus, statusExists := transitions[currentStatus]
	if !statusExists {
		return false
	}
	return len(actionsForStatus) == 0
}

var validStatuses = map[string]struct{}{
	StatusPendingPayment:  {},
	StatusPayed:           {},
	StatusExpired:         {},
	StatusCancelledUnpaid: {},
	StatusAccepted:        {},
	StatusPacked:          {},
	StatusShipped:         {},
	StatusCancelled:       {},
	StatusFinished:        {},
}

func ValidateStatuses(statuses []string) error {
	for _, s := range statuses {
		if _, ok := validStatuses[s]; !ok {
			return fmt.Errorf("status tidak dikenal: %s", s)
		}
	}
	return nil
}

// Aksi yang boleh dilakukan customer — hanya PAY dan CONFIRM_RECEIVED.
var customerActions = map[string]struct{}{
	ActionPay:             {},
	ActionConfirmReceived: {},
}

var adminExcludedActions = map[string]struct{}{
	ActionPay:             {},
	ActionConfirmReceived: {},
}

func GetAvailableActionsForCustomer(currentStatus string) ([]string, error) {
	all, err := GetAvailableActions(currentStatus)
	if err != nil {
		return nil, err
	}

	filtered := make([]string, 0)
	for _, action := range all {
		if _, ok := customerActions[action]; ok {
			filtered = append(filtered, action)
		}
	}
	return filtered, nil
}

func GetAvailableActionsForAdmin(currentStatus string) ([]string, error) {
	all, err := GetAvailableActions(currentStatus)
	if err != nil {
		return nil, err
	}

	filtered := make([]string, 0)
	for _, action := range all {
		if _, ok := adminExcludedActions[action]; !ok {
			filtered = append(filtered, action)
		}
	}
	return filtered, nil
}
