package routes

import (
	"4shure/cmd/internal/service"
	"4shure/cmd/internal/utils"
	"4shure/cmd/internal/utils/apierror"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type AppointmentService interface {
	GetAppointments(subId string) ([]*service.AppointmentResponse, apierror.ErrorResponse)
	CreateAppointment(req *service.AppointmentRequest, subId string) (*service.AppointmentResponse, apierror.ErrorResponse)
	DeleteAppointment(id int, sub string) apierror.ErrorResponse
	GetCalendar(monthStart, monthEnd int64) (*service.CalendarResponse, apierror.ErrorResponse)
}

type DefaultAppointmentRoute struct {
	AppointmentService AppointmentService
}

func NewAppointmentDefault(apptService AppointmentService) *DefaultAppointmentRoute {
	return &DefaultAppointmentRoute{AppointmentService: apptService}
}

func (a *DefaultAppointmentRoute) GetAppointments(c echo.Context) error {
	data, err := utils.ParseTokenDataCtx(c)
	if err != nil {
		return c.JSON(401, apierror.InvalidAuthTokenError)
	}

	appts, apierr := a.AppointmentService.GetAppointments(data.Sub)
	if apierr != nil {
		return c.JSON(apierr.Code(), apierr)
	}

	resp := echo.Map{"appointments": appts}
	return c.JSON(http.StatusOK, &resp)
}

func (a *DefaultAppointmentRoute) CreateAppointment(c echo.Context) error {
	var req service.AppointmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, apierror.MalformedBodyError)
	}

	data, err := utils.ParseTokenDataCtx(c)
	if err != nil {
		return c.JSON(401, apierror.InvalidAuthTokenError)
	}

	appt, apierr := a.AppointmentService.CreateAppointment(&req, data.Sub)
	if apierr != nil {
		return c.JSON(apierr.Code(), apierr)
	}
	return c.JSON(http.StatusCreated, appt)
}

func (a *DefaultAppointmentRoute) DeleteAppointment(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		errResp := apierror.NewSimple(400, "ID is not a number")
		return c.JSON(errResp.Code(), errResp)
	}

	data, err := utils.ParseTokenDataCtx(c)
	if err != nil {
		return c.JSON(401, apierror.InvalidAuthTokenError)
	}

	serr := a.AppointmentService.DeleteAppointment(id, data.Sub)
	if serr != nil {
		return c.JSON(serr.Code(), serr)
	}
	return c.NoContent(http.StatusOK)
}

func (a *DefaultAppointmentRoute) GetCalendar(c echo.Context) error {
	monthStr := c.QueryParam("month") // "2025-08"
	if monthStr == "" {
		return c.JSON(400, apierror.NewMissingParamError("month"))
	}

	monthStartMillis, monthEndMillis, err := parseMonthString(monthStr)
	if err != nil {
		apierr := apierror.NewSimple(400, "Could not understand month format")
		return c.JSON(apierr.Code(), apierr)
	}

	calendar, apierr := a.AppointmentService.GetCalendar(monthStartMillis, monthEndMillis)
	if apierr != nil {
		return c.JSON(apierr.Code(), apierr)
	}
	return c.JSON(http.StatusOK, &calendar)
}

// parseMonthString takes "YYYY-MM" (e.g., "2025-08") and returns
// the start of that month and the start of the next month as epoch millis.
func parseMonthString(monthString string) (int64, int64, error) {
	t, err := time.Parse("2006-01", monthString)
	if err != nil {
		return 0, 0, errors.New("invalid month format, expected YYYY-MM")
	}

	monthStart := t.UTC() // Ensure UTC always
	monthEnd := monthStart.AddDate(0, 1, 0)
	return monthStart.UnixMilli(), monthEnd.UnixMilli(), nil
}
