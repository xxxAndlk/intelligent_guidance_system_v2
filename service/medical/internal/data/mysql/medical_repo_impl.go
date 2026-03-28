package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"intelligent_guidance_system_v2/service/medical/internal/domain/aggregate"
	"intelligent_guidance_system_v2/service/medical/internal/domain/entity"
	"intelligent_guidance_system_v2/service/medical/internal/domain/event"
	"intelligent_guidance_system_v2/service/medical/internal/domain/repository"
	"intelligent_guidance_system_v2/service/medical/internal/domain/vo"
)

type medicalRecordRepo struct {
	db *gorm.DB
}

func NewMedicalRecordRepository(db *gorm.DB) repository.MedicalRecordRepository {
	return &medicalRecordRepo{db: db}
}

func (r *medicalRecordRepo) Save(ctx context.Context, record *aggregate.MedicalRecord) error {
	po := r.toPO(record)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&po).Error; err != nil {
			return err
		}
		record.SetID(po.ID)

		for _, d := range record.Diseases() {
			diseasePO := r.diseaseToPO(po.ID, d)
			if err := tx.Create(&diseasePO).Error; err != nil {
				return err
			}
		}

		for _, p := range record.Prescriptions() {
			prescriptionPO := r.prescriptionToPO(po.ID, p)
			if err := tx.Create(&prescriptionPO).Error; err != nil {
				return err
			}
		}

		for _, s := range record.Timeline() {
			statusPO := r.statusChangeToPO(po.ID, s)
			if err := tx.Create(&statusPO).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (r *medicalRecordRepo) FindByID(ctx context.Context, id int64) (*aggregate.MedicalRecord, error) {
	var po MedicalRecordPO
	err := r.db.WithContext(ctx).
		Preload("Diseases").
		Preload("Prescriptions").
		Preload("StatusChanges").
		First(&po, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toAggregate(&po), nil
}

func (r *medicalRecordRepo) FindByMedicalNumber(ctx context.Context, medicalNumber string) (*aggregate.MedicalRecord, error) {
	var po MedicalRecordPO
	err := r.db.WithContext(ctx).
		Preload("Diseases").
		Preload("Prescriptions").
		Preload("StatusChanges").
		Where("medical_number = ?", medicalNumber).
		First(&po).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toAggregate(&po), nil
}

func (r *medicalRecordRepo) FindByPatientID(ctx context.Context, patientID int64, limit, offset int) ([]*aggregate.MedicalRecord, int64, error) {
	var records []MedicalRecordPO
	var total int64

	query := r.db.WithContext(ctx).Model(&MedicalRecordPO{}).Where("patient_id = ?", patientID)
	query.Count(&total)

	err := query.
		Preload("Diseases").
		Preload("Prescriptions").
		Preload("StatusChanges").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, err
	}

	result := make([]*aggregate.MedicalRecord, len(records))
	for i, po := range records {
		result[i] = r.toAggregate(&po)
	}

	return result, total, nil
}

func (r *medicalRecordRepo) FindByDoctorID(ctx context.Context, doctorID int64, limit, offset int) ([]*aggregate.MedicalRecord, int64, error) {
	var records []MedicalRecordPO
	var total int64

	query := r.db.WithContext(ctx).Model(&MedicalRecordPO{}).Where("doctor_id = ?", doctorID)
	query.Count(&total)

	err := query.
		Preload("Diseases").
		Preload("Prescriptions").
		Preload("StatusChanges").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, err
	}

	result := make([]*aggregate.MedicalRecord, len(records))
	for i, po := range records {
		result[i] = r.toAggregate(&po)
	}

	return result, total, nil
}

func (r *medicalRecordRepo) FindByDepartmentID(ctx context.Context, departmentID int64, limit, offset int) ([]*aggregate.MedicalRecord, int64, error) {
	var records []MedicalRecordPO
	var total int64

	query := r.db.WithContext(ctx).Model(&MedicalRecordPO{}).Where("department_id = ?", departmentID)
	query.Count(&total)

	err := query.
		Preload("Diseases").
		Preload("Prescriptions").
		Preload("StatusChanges").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, err
	}

	result := make([]*aggregate.MedicalRecord, len(records))
	for i, po := range records {
		result[i] = r.toAggregate(&po)
	}

	return result, total, nil
}

func (r *medicalRecordRepo) FindByStatus(ctx context.Context, status int, limit, offset int) ([]*aggregate.MedicalRecord, int64, error) {
	var records []MedicalRecordPO
	var total int64

	query := r.db.WithContext(ctx).Model(&MedicalRecordPO{}).Where("status = ?", status)
	query.Count(&total)

	err := query.
		Preload("Diseases").
		Preload("Prescriptions").
		Preload("StatusChanges").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, err
	}

	result := make([]*aggregate.MedicalRecord, len(records))
	for i, po := range records {
		result[i] = r.toAggregate(&po)
	}

	return result, total, nil
}

func (r *medicalRecordRepo) List(ctx context.Context, filter *repository.MedicalRecordFilter) ([]*aggregate.MedicalRecord, int64, error) {
	var records []MedicalRecordPO
	var total int64

	query := r.db.WithContext(ctx).Model(&MedicalRecordPO{})

	if filter.PatientID > 0 {
		query = query.Where("patient_id = ?", filter.PatientID)
	}
	if filter.DoctorID > 0 {
		query = query.Where("doctor_id = ?", filter.DoctorID)
	}
	if filter.DepartmentID > 0 {
		query = query.Where("department_id = ?", filter.DepartmentID)
	}
	if filter.Status > 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.MedicalNumber != "" {
		query = query.Where("medical_number = ?", filter.MedicalNumber)
	}
	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.
		Preload("Diseases").
		Preload("Prescriptions").
		Preload("StatusChanges").
		Order("created_at DESC").
		Limit(filter.PageSize).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, err
	}

	result := make([]*aggregate.MedicalRecord, len(records))
	for i, po := range records {
		result[i] = r.toAggregate(&po)
	}

	return result, total, nil
}

func (r *medicalRecordRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&MedicalRecordPO{}, id).Error
}

func (r *medicalRecordRepo) Update(ctx context.Context, record *aggregate.MedicalRecord) error {
	po := r.toPO(record)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&po).Error; err != nil {
			return err
		}

		if len(record.Diseases()) > 0 {
			diseasePOs := make([]DiseaseItemPO, len(record.Diseases()))
			for i, d := range record.Diseases() {
				diseasePOs[i] = r.diseaseToPO(po.ID, d)
			}
			if err := tx.Where("medical_record_id = ?", po.ID).Delete(&DiseaseItemPO{}).Error; err != nil {
				return err
			}
			if err := tx.Create(&diseasePOs).Error; err != nil {
				return err
			}
		}

		if len(record.Prescriptions()) > 0 {
			prescriptionPOs := make([]PrescriptionItemPO, len(record.Prescriptions()))
			for i, p := range record.Prescriptions() {
				prescriptionPOs[i] = r.prescriptionToPO(po.ID, p)
			}
			if err := tx.Where("medical_record_id = ?", po.ID).Delete(&PrescriptionItemPO{}).Error; err != nil {
				return err
			}
			if err := tx.Create(&prescriptionPOs).Error; err != nil {
				return err
			}
		}

		if len(record.Timeline()) > 0 {
			statusPOs := make([]StatusChangePO, len(record.Timeline()))
			for i, s := range record.Timeline() {
				statusPOs[i] = r.statusChangeToPO(po.ID, s)
			}
			if err := tx.Where("medical_record_id = ?", po.ID).Delete(&StatusChangePO{}).Error; err != nil {
				return err
			}
			if err := tx.Create(&statusPOs).Error; err != nil {
				return err
			}
		}

		for _, e := range record.Events() {
			eventPO := r.eventToPO(e)
			if err := tx.Create(&eventPO).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *medicalRecordRepo) toPO(record *aggregate.MedicalRecord) MedicalRecordPO {
	po := MedicalRecordPO{
		ID:            record.ID(),
		MedicalNumber: record.MedicalNumber(),
		PatientID:     record.PatientID(),
		DoctorID:      record.DoctorID(),
		DepartmentID:  record.DepartmentID(),
		Symptom:       record.Symptom(),
		Diagnosis:     record.Diagnosis(),
		TreatmentPlan: record.TreatmentPlan(),
		Status:        int(record.Status()),
		CreatedAt:     record.CreatedAt(),
		UpdatedAt:     record.UpdatedAt(),
	}

	if record.Payment() != nil {
		po.PaymentAmount = record.Payment().Amount().Amount()
		po.PaymentMethod = string(record.Payment().Method())
		po.PaymentStatus = string(record.Payment().Status())
		po.PaymentTime = record.Payment().PayTime()
	}

	return po
}

func (r *medicalRecordRepo) toAggregate(po *MedicalRecordPO) *aggregate.MedicalRecord {
	diseases := make([]*entity.DiseaseItem, len(po.Diseases))
	for i, d := range po.Diseases {
		diseases[i] = entity.ReconstructDiseaseItem(
			d.ItemID,
			d.DiseaseID,
			d.DiseaseName,
			d.Symptoms,
			d.Diagnosis,
			d.CreatedAt,
		)
	}

	prescriptions := make([]*entity.PrescriptionItem, len(po.Prescriptions))
	for i, p := range po.Prescriptions {
		price, _ := vo.NewMoney(p.Price, "CNY")
		prescriptions[i] = entity.ReconstructPrescriptionItem(
			p.ItemID,
			p.DrugID,
			p.DrugName,
			p.Quantity,
			price,
			p.Usage,
		)
	}

	timeline := make([]*entity.StatusChange, len(po.StatusChanges))
	for i, s := range po.StatusChanges {
		timeline[i] = entity.ReconstructStatusChange(
			entity.MedicalStatus(s.FromStatus),
			entity.MedicalStatus(s.ToStatus),
			s.ChangeTime,
			s.OperatorID,
			s.Reason,
		)
	}

	payment, _ := vo.NewPaymentInfo(vo.Zero("CNY"), vo.PaymentMethod(po.PaymentMethod))
	if payment != nil {
		amount, _ := vo.NewMoney(po.PaymentAmount, "CNY")
		payment.SetAmount(amount)
		payment.Status = vo.PaymentStatus(po.PaymentStatus)
		if po.PaymentTime != nil {
			payment.PayTime = po.PaymentTime
		}
	}

	return aggregate.ReconstructMedicalRecord(
		po.ID,
		po.MedicalNumber,
		po.PatientID,
		po.DoctorID,
		po.DepartmentID,
		po.Symptom,
		po.Diagnosis,
		po.TreatmentPlan,
		diseases,
		prescriptions,
		payment,
		entity.MedicalStatus(po.Status),
		timeline,
		po.CreatedAt,
		po.UpdatedAt,
	)
}

func (r *medicalRecordRepo) diseaseToPO(recordID int64, d *entity.DiseaseItem) DiseaseItemPO {
	return DiseaseItemPO{
		MedicalRecordID: recordID,
		ItemID:          d.ID(),
		DiseaseID:       d.DiseaseID(),
		DiseaseName:     d.DiseaseName(),
		Symptoms:        d.Symptoms(),
		Diagnosis:       d.Diagnosis(),
		CreatedAt:       d.CreatedAt(),
	}
}

func (r *medicalRecordRepo) prescriptionToPO(recordID int64, p *entity.PrescriptionItem) PrescriptionItemPO {
	return PrescriptionItemPO{
		MedicalRecordID: recordID,
		ItemID:          p.ID(),
		DrugID:          p.DrugID(),
		DrugName:        p.DrugName(),
		Quantity:        p.Quantity(),
		Price:           p.Price().Amount(),
		Usage:           p.Usage(),
		CreatedAt:       time.Now(),
	}
}

func (r *medicalRecordRepo) statusChangeToPO(recordID int64, s *entity.StatusChange) StatusChangePO {
	return StatusChangePO{
		MedicalRecordID: recordID,
		FromStatus:      int(s.FromStatus()),
		ToStatus:        int(s.ToStatus()),
		ChangeTime:      s.ChangeTime(),
		OperatorID:      s.OperatorID(),
		Reason:          s.Reason(),
	}
}

func (r *medicalRecordRepo) eventToPO(e *event.MedicalEvent) MedicalEventPO {
	payload, _ := json.Marshal(e.Payload())
	return MedicalEventPO{
		RecordID:    e.RecordID(),
		EventType:   string(e.EventType()),
		OccurredAt:  e.OccurredAt(),
		OperatorID:  e.OperatorID(),
		Payload:     string(payload),
	}
}

type eventRepo struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) repository.EventRepository {
	return &eventRepo{db: db}
}

func (r *eventRepo) Save(ctx context.Context, events ...interface{}) error {
	for _, e := range events {
		if me, ok := e.(*event.MedicalEvent); ok {
			payload, _ := json.Marshal(me.Payload())
			po := MedicalEventPO{
				RecordID:    me.RecordID(),
				EventType:   string(me.EventType()),
				OccurredAt:  me.OccurredAt(),
				OperatorID:  me.OperatorID(),
				Payload:     string(payload),
			}
			if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *eventRepo) FindByRecordID(ctx context.Context, recordID int64) ([]interface{}, error) {
	var events []MedicalEventPO
	err := r.db.WithContext(ctx).
		Where("record_id = ?", recordID).
		Order("occurred_at ASC").
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(events))
	for i, po := range events {
		var payload map[string]interface{}
		json.Unmarshal([]byte(po.Payload), &payload)
		result[i] = event.NewMedicalEvent(
			event.EventType(po.EventType),
			po.RecordID,
			po.OperatorID,
			payload,
		)
	}

	return result, nil
}