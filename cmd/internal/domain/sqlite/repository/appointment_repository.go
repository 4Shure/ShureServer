package repository

import (
	"4shure/cmd/internal/domain/entity"
	"errors"
	"gorm.io/gorm"
)

type DefaultAppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *DefaultAppointmentRepository {
	return &DefaultAppointmentRepository{db: db}
}

func (a *DefaultAppointmentRepository) FindByID(id int) (*entity.Appointment, error) {
	var appt entity.Appointment
	err := a.db.First(&appt, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &appt, err
}

func (a *DefaultAppointmentRepository) IsAvailable(begin, end int64) (bool, error) {
	if begin >= end {
		return false, errors.New("start time must be before end time")
	}

	var count int64
	err := a.db.Model(&entity.Appointment{}).
		Where("is_deleted = ?", false).
		Where("begins_at < ?", end).
		Where("ends_at > ?", begin).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (a *DefaultAppointmentRepository) FindAll() ([]*entity.Appointment, error) {
	var appts []*entity.Appointment
	err := a.db.Find(&appts).Error
	return appts, err
}

// FindMonthAppointments finds all appointments that overlap with a given month.
// This method returns PARTIAL appointment entities, having only `BeginsAt` and `EndsAt` fields.
func (a *DefaultAppointmentRepository) FindMonthAppointments(monthStart, monthEnd int64) ([]*entity.Appointment, error) {
	var results []*entity.Appointment

	err := a.db.Model(&entity.Appointment{}).
		Select("begins_at, ends_at").
		Where("is_deleted = ?", false).
		Where("begins_at < ?", monthEnd).
		Where("ends_at > ?", monthStart).
		Order("begins_at asc").
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (a *DefaultAppointmentRepository) FindByUserID(id int) ([]*entity.Appointment, error) {
	var appts []*entity.Appointment
	err := a.db.Where("user_id = ?", id).Find(&appts).Error
	return appts, err
}

func (a *DefaultAppointmentRepository) Save(appointment *entity.Appointment) error {
	return a.db.Save(appointment).Error
}

func (a *DefaultAppointmentRepository) Delete(appointment *entity.Appointment) error {
	return a.db.Delete(appointment).Error
}
