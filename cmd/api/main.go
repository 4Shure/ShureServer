package main

import (
	"4shure/cmd/internal/domain/sqlite"
	"4shure/cmd/internal/domain/sqlite/repository"
	cognitoclient "4shure/cmd/internal/integration/aws/cognito"
	"4shure/cmd/internal/routes"
	"4shure/cmd/internal/service"
	"4shure/cmd/internal/utils/validators"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	validate := validator.New()
	registerValidators(validate)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file", err)
	}

	// Init SQLite
	db, err := sqlite.Init()
	if err != nil {
		log.Fatal("failed to initialize database", err)
	}

	// Cognito cliente
	cogClient, err := cognitoclient.InitCognitoClient()
	if err != nil {
		log.Fatal("failed to initialize cognito client", err)
	}

	// Getting repositories
	userRepo := repository.NewUserRepository(db)
	apptRepo := repository.NewAppointmentRepository(db)

	// Getting services
	userService := service.NewUserService(userRepo, validate, cogClient)
	apptService := service.NewAppointmentService(apptRepo, userRepo, validate)

	// Getting routes
	userRoutes := routes.NewUserDefault(userService)
	apptRoutes := routes.NewAppointmentDefault(apptService)

	e := echo.New()
	e.Use(middleware.CORS())

	// Appointments
	e.GET("/api/appointments", apptRoutes.GetAppointments)
	e.POST("/api/appointments", apptRoutes.CreateAppointment)
	e.DELETE("/api/appointments/:id", apptRoutes.DeleteAppointment)

	// Pseudo-entity "Calendar" to check the availability of a new appointment
	e.GET("/api/calendar", apptRoutes.GetCalendar)

	// Users
	e.GET("/api/users", userRoutes.GetUsers)
	e.GET("/api/users/:id", userRoutes.GetUser)
	e.POST("/api/users", userRoutes.CreateUser)
	e.POST("/api/users/login", userRoutes.CreateLogin)
	e.POST("/api/users/verify", userRoutes.VerifySignup)

	err = e.Start(":6060")
	if err != nil {
		e.Logger.Fatal(err)
	}
}

func registerValidators(validate *validator.Validate) {
	_ = validate.RegisterValidation("hasupper", validators.HasUpper)
	_ = validate.RegisterValidation("haslower", validators.HasLower)
	_ = validate.RegisterValidation("hasdigit", validators.HasDigit)
	_ = validate.RegisterValidation("hasspecial", validators.HasSpecial)
	_ = validate.RegisterValidation("nodupes", validators.NoDupes)
	_ = validate.RegisterValidation("nospaces", validators.NoWhiteSpaces)
	_ = validate.RegisterValidation("iso8601", validators.IsIso8601)
}
