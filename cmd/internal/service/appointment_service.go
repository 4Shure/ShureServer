package service

import (
	"4shure/cmd/internal/domain/entity"
	"4shure/cmd/internal/utils"
	"4shure/cmd/internal/utils/apierror"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"time"
)

type AppointmentRepository interface {
	Save(appointment *entity.Appointment) error
	FindAll() ([]*entity.Appointment, error)
	IsAvailable(begin, end int64) (bool, error)
	FindByUserID(id int) ([]*entity.Appointment, error)
	FindByID(id int) (*entity.Appointment, error)
	FindMonthAppointments(monthStart, monthEnd int64) ([]*entity.Appointment, error)
	Delete(appointment *entity.Appointment) error
}

type AppointmentRequest struct {
	Title    string `json:"title" validate:"max=128"`
	BeginsAt string `json:"begins_at" validate:"required,iso8601"`
}

type AppointmentResponse struct {
	ID        int    `json:"id"`
	BeginsAt  string `json:"begins_at"`
	EndsAt    string `json:"ends_at"`
	UserID    int    `json:"user_id"`
	IsDeleted bool   `json:"is_deleted"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Title     string `json:"title"`
}

type ScheduledDay struct {
	BeginsAt string `json:"begins_at"`
	EndsAt   string `json:"ends_at"`
}

type CalendarResponse struct {
	ScheduledDays []*ScheduledDay `json:"scheduled_days"`
}

type DefaultAppointmentService struct {
	AppointmentRepo AppointmentRepository
	UserRepo        UserRepository
	Validate        *validator.Validate
}

func NewAppointmentService(apptRepo AppointmentRepository, userRepo UserRepository, validate *validator.Validate) *DefaultAppointmentService {
	return &DefaultAppointmentService{AppointmentRepo: apptRepo, UserRepo: userRepo, Validate: validate}
}

func (a *DefaultAppointmentService) GetAppointments(subId string) ([]*AppointmentResponse, apierror.ErrorResponse) {
	caller, err := a.UserRepo.FindBySub(subId)
	if err != nil {
		log.Errorf("failed to check if user %s is admin: %v", subId, err)
		return nil, apierror.InternalServerError
	}

	var appts []*entity.Appointment
	if caller.IsAdmin {
		appts, err = a.AppointmentRepo.FindAll()
	} else {
		appts, err = a.AppointmentRepo.FindByUserID(caller.ID)
	}

	if err != nil {
		log.Errorf("failed to find appointments for user %d: %v", caller.ID, err)
		return nil, apierror.InternalServerError
	}

	response := make([]*AppointmentResponse, len(appts))
	for i, appt := range appts {
		response[i] = toAppointmentResponse(appt)
	}
	return response, nil
}

func (a *DefaultAppointmentService) CreateAppointment(req *AppointmentRequest, subId string) (*AppointmentResponse, apierror.ErrorResponse) {
	caller, err := a.UserRepo.FindBySub(subId)
	if err != nil {
		log.Errorf("failed to fetch user %s: %v", subId, err)
		return nil, apierror.InternalServerError
	}

	utils.Sanitize(req)
	if valerr := a.Validate.Struct(req); valerr != nil {
		return nil, apierror.FromValidationError(valerr)
	}

	begin, err := utils.FromEpoch(req.BeginsAt)
	if err != nil {
		return nil, apierror.MalformedBodyError
	}

	if !utils.IsHourExact(begin) {
		return nil, apierror.HourNotExactError
	}

	if !isFuture(begin) {
		return nil, apierror.AppointmentInPastError
	}

	end := begin + time.Hour.Milliseconds() - 1
	now := utils.NowUTC()

	available, err := a.AppointmentRepo.IsAvailable(begin, end)
	if err != nil {
		log.Errorf("failed to check if time %d is available: %v", begin, err)
		return nil, apierror.InternalServerError
	}

	if !available {
		return nil, apierror.MomentNotAvailable
	}

	appointment := &entity.Appointment{
		BeginsAt:  begin,
		EndsAt:    end,
		UserID:    caller.ID,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: now,
		Title:     req.Title,
	}

	err = a.AppointmentRepo.Save(appointment)
	if err != nil {
		log.Errorf("failed to save appointment: %v", err)
		return nil, apierror.InternalServerError
	}
	return toAppointmentResponse(appointment), nil
}

func (a *DefaultAppointmentService) DeleteAppointment(id int, issuerSub string) apierror.ErrorResponse {
	caller, err := a.UserRepo.FindBySub(issuerSub)
	if err != nil {
		log.Errorf("failed to check if user %s is admin: %v", issuerSub, err)
		return apierror.InternalServerError
	}

	appt, err := a.AppointmentRepo.FindByID(id)
	if err != nil {
		log.Errorf("failed to fetch appointment by id %d: %v", id, err)
		return apierror.InternalServerError
	}

	if caller == nil || appt == nil || appt.IsDeleted || appt.UserID != caller.ID {
		return apierror.NotFoundError
	}

	err = a.AppointmentRepo.Delete(appt)
	if err != nil {
		log.Errorf("failed to delete appointment by id %d: %v", id, err)
		return apierror.InternalServerError
	}
	return nil
}

func (a *DefaultAppointmentService) GetCalendar(monthStart, monthEnd int64) (*CalendarResponse, apierror.ErrorResponse) {
	appts, err := a.AppointmentRepo.FindMonthAppointments(monthStart, monthEnd)
	if err != nil {
		log.Errorf("failed to fetch appointments availability [%d - %d]: %v", monthStart, monthEnd, err)
		return nil, apierror.InternalServerError
	}

	schedDays := make([]*ScheduledDay, len(appts))
	for i, appt := range appts {
		schedDays[i] = toScheduledDay(appt)
	}

	calendar := &CalendarResponse{
		ScheduledDays: schedDays,
	}
	return calendar, nil
}

func isFuture(millis int64) bool {
	now := utils.NowUTC()
	return millis > now
}

func toScheduledDay(appt *entity.Appointment) *ScheduledDay {
	return &ScheduledDay{
		BeginsAt: utils.FormatEpoch(appt.BeginsAt),
		EndsAt:   utils.FormatEpoch(appt.EndsAt),
	}
}

func toAppointmentResponse(appt *entity.Appointment) *AppointmentResponse {
	return &AppointmentResponse{
		ID:        appt.ID,
		UserID:    appt.UserID,
		IsDeleted: appt.IsDeleted,
		Title:     appt.Title,
		BeginsAt:  utils.FormatEpoch(appt.BeginsAt),
		EndsAt:    utils.FormatEpoch(appt.EndsAt),
		CreatedAt: utils.FormatEpoch(appt.CreatedAt),
		UpdatedAt: utils.FormatEpoch(appt.UpdatedAt),
	}
}
