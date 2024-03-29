package handlers

import (
	"database/sql"
	"strconv"
	"time"

	"clinic-reservation-system.com/back-end/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type DoctorHandler struct{}

func (handler DoctorHandler) AddAppointment(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	Tid := claims["id"].(float64)
	id := int64(Tid)
	nullableID := sql.NullInt64{Int64: id, Valid: true}

	timestamp := ctx.Query("timestamp")

	if timestamp == "" {
		return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "Date or time missing"})
	}

	// Parse the date string into a time.Time object
	layout := "2006-01-02 15:04"
	date, err := time.Parse(layout, timestamp)

	if err != nil || date.Before(time.Now()) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date or time"})
	}

	var sqlTime sql.NullString = sql.NullString{String: timestamp, Valid: true}
	appointment := models.Appointment{DoctorID: nullableID, AppointmentTime: sqlTime}

	if appointment.Create() {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": true,
		})
	}

	return ctx.SendStatus(fiber.StatusInternalServerError)

}

func (handler DoctorHandler) GetAppointment(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	Tid := claims["id"].(float64)
	nullableID := sql.NullInt64{Int64: int64(Tid), Valid: true}

	appointment := models.Appointment{DoctorID: nullableID}

	appointments := appointment.GetReserved("doctor")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"appointments": appointments,
	})
}

func (handler DoctorHandler) DeleteAppointment(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	Tid := claims["id"].(float64)
	nullableID := sql.NullInt64{Int64: int64(Tid), Valid: true}

	appointmentID, err := strconv.Atoi(ctx.Query("appointment_id"))

	if err != nil {
		return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{"error": "appointment id missing"})
	}

	appointment := models.Appointment{ID: sql.NullInt64{Int64: int64(appointmentID), Valid: true}, DoctorID: nullableID}

	if appointment.Delete() {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": true,
		})
	}

	return ctx.SendStatus(fiber.StatusInternalServerError)
}
