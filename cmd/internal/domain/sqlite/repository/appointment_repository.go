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

func (a *DefaultAppointmentRepository) Save(appointment *entity.Appointment) error {
	return a.db.Save(appointment).Error
}

func (a *DefaultAppointmentRepository) Delete(appointment *entity.Appointment) error {
	return a.db.Delete(appointment).Error
}
