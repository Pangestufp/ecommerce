package service

import (
	"backend/dto"
	"backend/entity"
	"backend/errorhandler"
	"backend/helper"
	"backend/repository"
)

type SalesOrderService interface {
	GetAllPaginate(cursor *dto.Paginate, statuses []string, limit int) ([]dto.SalesOrderResponse, *dto.Paginate, error)
	GetAllByUserPaginate(cursor *dto.Paginate, userID string, statuses []string, limit int) ([]dto.SalesOrderResponse, *dto.Paginate, error)
	GetByCode(salesOrderCode string) (*dto.SalesOrderDetailResponse, error)
}

type salesOrderService struct {
	repository        repository.SalesOrderRepository
	historyRepository repository.OrderStatusHistoryRepository
}

func NewSalesOrderService(
	repository repository.SalesOrderRepository,
	historyRepository repository.OrderStatusHistoryRepository,
) SalesOrderService {
	return &salesOrderService{
		repository:        repository,
		historyRepository: historyRepository,
	}
}

func (s *salesOrderService) GetAllPaginate(cursor *dto.Paginate, statuses []string, limit int) ([]dto.SalesOrderResponse, *dto.Paginate, error) {
	orders, err := s.repository.GetAllPaginate(cursor, statuses, limit)
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: "Gagal mengambil data order"}
	}

	return buildPaginateResponse(orders, cursor, limit)
}

func (s *salesOrderService) GetAllByUserPaginate(cursor *dto.Paginate, userID string, statuses []string, limit int) ([]dto.SalesOrderResponse, *dto.Paginate, error) {
	orders, err := s.repository.GetAllByUserPaginate(cursor, userID, statuses, limit)
	if err != nil {
		return nil, nil, &errorhandler.InternalServerError{Message: "Gagal mengambil data order"}
	}

	return buildPaginateResponse(orders, cursor, limit)
}

func (s *salesOrderService) GetByCode(salesOrderCode string) (*dto.SalesOrderDetailResponse, error) {
	order, details, err := s.repository.GetByCode(salesOrderCode)
	if err != nil {
		return nil, err
	}

	histories, err := s.historyRepository.FindBySalesOrderID(order.SalesOrderID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: "Gagal mengambil riwayat status"}
	}

	actions, err := helper.GetAvailableActions(order.Status)
	if err != nil {
		actions = []string{}
	}

	detailResponses := make([]dto.SalesOrderItemResponse, 0, len(details))
	for _, d := range details {
		detailResponses = append(detailResponses, mapDetailToResponse(d))
	}

	historyResponses := make([]dto.OrderStatusHistoryResponse, 0, len(histories))
	for _, h := range histories {
		historyResponses = append(historyResponses, dto.OrderStatusHistoryResponse{
			HistoryID:      h.HistoryID,
			SalesOrderID:   h.SalesOrderID,
			PreviousStatus: h.PreviousStatus,
			Status:         h.Status,
			Note:           h.Note,
			CreatedBy:      h.CreatedBy,
			CreatedName:    h.CreatedName,
			CreatedAt:      h.CreatedAt,
		})
	}

	return &dto.SalesOrderDetailResponse{
		SalesOrder: mapSalesOrderToResponse(*order),
		Details:    detailResponses,
		Histories:  historyResponses,
		Actions:    actions,
	}, nil
}

// ─── shared pagination builder ────────────────────────────────────────────────

func buildPaginateResponse(orders []entity.SalesOrder, cursor *dto.Paginate, limit int) ([]dto.SalesOrderResponse, *dto.Paginate, error) {
	var paginate *dto.Paginate

	if len(orders) > 0 {
		isNext := cursor == nil || cursor.Direction == nil || *cursor.Direction == "next"
		isPrev := cursor != nil && cursor.Direction != nil && *cursor.Direction == "prev"

		hasNext := "false"
		hasPrev := "false"

		if isNext {
			if len(orders) > limit {
				hasNext = "true"
				orders = orders[:limit]
			}
			if cursor != nil && cursor.LastID != nil {
				hasPrev = "true"
			}
		} else if isPrev {
			if len(orders) > limit {
				hasPrev = "true"
				orders = orders[1:]
			}
			hasNext = "true"
		}

		direction := "next"
		if isPrev {
			direction = "prev"
		}

		first := orders[0]
		last := orders[len(orders)-1]
		paginate = &dto.Paginate{
			FirstID:        &first.SalesOrderID,
			FirstCreatedAt: &first.CreatedAt,
			LastID:         &last.SalesOrderID,
			LastCreatedAt:  &last.CreatedAt,
			HasNext:        &hasNext,
			HasPrev:        &hasPrev,
			Direction:      &direction,
		}
	}

	responses := make([]dto.SalesOrderResponse, 0, len(orders))
	for _, o := range orders {
		responses = append(responses, mapSalesOrderToResponse(o))
	}

	return responses, paginate, nil
}

// ─── mappers ──────────────────────────────────────────────────────────────────

func mapSalesOrderToResponse(o entity.SalesOrder) dto.SalesOrderResponse {
	return dto.SalesOrderResponse{
		SalesOrderID:            o.SalesOrderID,
		SalesOrderCode:          o.SalesOrderCode,
		UserID:                  o.UserID,
		CustomerName:            o.CustomerName,
		CustomerEmail:           o.CustomerEmail,
		CustomerPhone:           o.CustomerPhone,
		CustomerAddressSnapshot: o.CustomerAddressSnapshot,
		OriginAddressSnapshot:   o.OriginAddressSnapshot,

		TrackingNumber: o.TrackingNumber,

		ShippingCourier:       o.ShippingCourier,
		ShippingDisplayName:   o.ShippingDisplayName,
		ShippingService:       o.ShippingService,
		ShippingDescription:   o.ShippingDescription,
		ShippingETD:           o.ShippingETD,
		OriginDistrictID:      o.OriginDistrictID,
		DestinationDistrictID: o.DestinationDistrictID,

		TotalWeight: o.TotalWeight,

		ShippingFee:       o.ShippingFee,
		ShippingFeeFormat: helper.FormatRupiah(o.ShippingFee),

		TotalBeforeDiscount:       o.TotalBeforeDiscount,
		TotalBeforeDiscountFormat: helper.FormatRupiah(o.TotalBeforeDiscount),

		DiscountAmount:       o.DiscountAmount,
		DiscountAmountFormat: helper.FormatRupiah(o.DiscountAmount),

		FinalTotal:       o.FinalTotal,
		FinalTotalFormat: helper.FormatRupiah(o.FinalTotal),

		Status: o.Status,
		Notes:  o.Notes,

		ConfirmedAt: o.ConfirmedAt,
		PaidAt:      o.PaidAt,
		DeliveredAt: o.DeliveredAt,
		FinishAt:    o.FinishAt,
		CancelledAt: o.CancelledAt,

		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func mapDetailToResponse(d entity.SalesOrderDetail) dto.SalesOrderItemResponse {
	return dto.SalesOrderItemResponse{
		DetailID:     d.DetailID,
		SalesOrderID: d.SalesOrderID,

		ProductID:   d.ProductID,
		ProductName: d.ProductName,
		ProductCode: d.ProductCode,
		ImageURL:    d.ImageURL,

		Quantity:    d.Quantity,
		Weight:      d.Weight,
		TotalWeight: d.TotalWeight,

		PriceBeforeDiscount:       d.PriceBeforeDiscount,
		PriceBeforeDiscountFormat: helper.FormatRupiah(d.PriceBeforeDiscount),

		DiscountID: d.DiscountID,

		DiscountAmount:       d.DiscountAmount,
		DiscountAmountFormat: helper.FormatRupiah(d.DiscountAmount),

		FinalUnitPrice:       d.FinalUnitPrice,
		FinalUnitPriceFormat: helper.FormatRupiah(d.FinalUnitPrice),

		FinalTotal:       d.FinalTotal,
		FinalTotalFormat: helper.FormatRupiah(d.FinalTotal),

		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
