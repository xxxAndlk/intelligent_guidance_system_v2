package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"intelligent-guidance-system/service/payment/internal/domain/aggregate"
	"intelligent-guidance-system/service/payment/internal/domain/entity"
	"intelligent-guidance-system/service/payment/internal/domain/repository"
)

type PaymentRepoImpl struct {
	db *gorm.DB
}

func NewPaymentRepoImpl(db *gorm.DB) repository.PaymentRepository {
	return &PaymentRepoImpl{db: db}
}

func (r *PaymentRepoImpl) Save(ctx context.Context, payment *aggregate.Payment) error {
	po := r.toPO(payment)

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
			return err
		}
		payment.SetID(po.ID)
	} else {
		if err := r.db.WithContext(ctx).Save(&po).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *PaymentRepoImpl) FindByID(ctx context.Context, id int64) (*aggregate.Payment, error) {
	var po PaymentPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toAggregate(&po), nil
}

func (r *PaymentRepoImpl) FindByMedicalID(ctx context.Context, medicalID int64) (*aggregate.Payment, error) {
	var po PaymentPO
	if err := r.db.WithContext(ctx).Where("medical_id = ?", medicalID).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toAggregate(&po), nil
}

func (r *PaymentRepoImpl) FindByPatientID(ctx context.Context, patientID int64) ([]*aggregate.Payment, error) {
	var pos []PaymentPO
	if err := r.db.WithContext(ctx).Where("patient_id = ?", patientID).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *PaymentRepoImpl) FindByStatus(ctx context.Context, status int) ([]*aggregate.Payment, error) {
	var pos []PaymentPO
	statusCode := entity.PaymentStatus(status).Code()
	if err := r.db.WithContext(ctx).Where("status = ?", statusCode).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toAggregates(pos), nil
}

func (r *PaymentRepoImpl) FindByTransactionID(ctx context.Context, transactionID string) (*aggregate.Payment, error) {
	var po PaymentPO
	if err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toAggregate(&po), nil
}

func (r *PaymentRepoImpl) FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Payment, int64, error) {
	var pos []PaymentPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&PaymentPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return r.toAggregates(pos), total, nil
}

func (r *PaymentRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&PaymentPO{}, id).Error
}

func (r *PaymentRepoImpl) toPO(payment *aggregate.Payment) *PaymentPO {
	return &PaymentPO{
		ID:            payment.ID(),
		MedicalID:     payment.MedicalID(),
		PatientID:     payment.PatientID(),
		Amount:        payment.Amount(),
		Method:        payment.Method().Code(),
		Status:        payment.Status().Code(),
		TransactionID: payment.TransactionID(),
		CreatedAt:     payment.CreatedAt(),
		UpdatedAt:     payment.UpdatedAt(),
	}
}

func (r *PaymentRepoImpl) toAggregate(po *PaymentPO) *aggregate.Payment {
	return aggregate.ReconstructPayment(
		po.ID,
		po.MedicalID,
		po.PatientID,
		po.Amount,
		entity.PaymentMethodFromCode(po.Method),
		entity.PaymentStatusFromCode(po.Status),
		po.TransactionID,
		nil,
		po.CreatedAt,
		po.UpdatedAt,
	)
}

func (r *PaymentRepoImpl) toAggregates(pos []PaymentPO) []*aggregate.Payment {
	payments := make([]*aggregate.Payment, 0, len(pos))
	for _, po := range pos {
		payments = append(payments, r.toAggregate(&po))
	}
	return payments
}

type RefundRepoImpl struct {
	db *gorm.DB
}

func NewRefundRepoImpl(db *gorm.DB) repository.RefundRepository {
	return &RefundRepoImpl{db: db}
}

func (r *RefundRepoImpl) Save(ctx context.Context, refund *entity.Refund) error {
	po := RefundPO{
		ID:          refund.ID(),
		PaymentID:   refund.PaymentID(),
		Amount:      refund.Amount(),
		Reason:      refund.Reason(),
		Status:      refund.Status().Code(),
		ProcessedAt: refund.ProcessedAt(),
		CreatedAt:   refund.CreatedAt(),
		UpdatedAt:   refund.UpdatedAt(),
	}

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
			return err
		}
		refund.SetID(po.ID)
	} else {
		if err := r.db.WithContext(ctx).Save(&po).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *RefundRepoImpl) FindByID(ctx context.Context, id int64) (*entity.Refund, error) {
	var po RefundPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return entity.ReconstructRefund(
		po.ID, po.PaymentID, po.Amount, po.Reason,
		entity.RefundStatusFromCode(po.Status),
		po.ProcessedAt, po.CreatedAt, po.UpdatedAt,
	), nil
}

func (r *RefundRepoImpl) FindByPaymentID(ctx context.Context, paymentID int64) ([]*entity.Refund, error) {
	var pos []RefundPO
	if err := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).Find(&pos).Error; err != nil {
		return nil, err
	}

	refunds := make([]*entity.Refund, 0, len(pos))
	for _, po := range pos {
		refunds = append(refunds, entity.ReconstructRefund(
			po.ID, po.PaymentID, po.Amount, po.Reason,
			entity.RefundStatusFromCode(po.Status),
			po.ProcessedAt, po.CreatedAt, po.UpdatedAt,
		))
	}

	return refunds, nil
}

func (r *RefundRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&RefundPO{}, id).Error
}