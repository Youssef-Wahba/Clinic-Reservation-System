package models

import (
	"time"

	"clinic-reservation-system.com/back-end/inits"
)

type Appointment struct {
	DoctorID        uint      `json:"doctor_id"`
	PatientID       uint      `json:"patient_id"`
	AppointmentTime time.Time `json:"appointment_time"`
}

func (a Appointment) InitTable() bool {
	query := `
	CREATE TABLE IF NOT EXISTS appointments(
		id int NOT NULL AUTO_INCREMENT,
		doctor_id int NOT NULL,
		patient_id int,
		appointment_time timestamp	NOT NULL,
		PRIMARY KEY (id),
		KEY doctor_id_fk (doctor_id) ,
		KEY patient_id_fk (patient_id)
		);
		`

	_, err := inits.DB.Exec(query)

	return err == nil

}

func (a Appointment) Create() bool {
	if !a.CheckIfViable() {
		return false
	}

	query := `
	INSERT INTO appointments(doctor_id,appointment_time) VALUES(?,?);
	`
	_, err := inits.DB.Exec(query, a.DoctorID, a.AppointmentTime)

	return err == nil
}

func (a Appointment) Reserve() bool{
	query := `
	UPDATE appointments SET patient_id=? WHERE doctor_id=? AND appointment_time=?;
	`
	_, err := inits.DB.Exec(query, a.PatientID, a.DoctorID, a.AppointmentTime)

	return err == nil
}

func (a Appointment) Delete() bool {
	query := `
	DELETE FROM appointments WHERE doctor_id=? AND appointment_time=?;
	`
	_, err := inits.DB.Exec(query, a.DoctorID, a.AppointmentTime)

	return err == nil
}

func (a Appointment) GetAll(userType string) []Appointment {
	var query string
	var id uint

	if userType == "doctor" {
		query = `
		SELECT * FROM appointments WHERE doctor_id=?;
		`
		id = a.DoctorID
	} else {
		query = `
		SELECT * FROM appointments WHERE patient_id=?;
		`
		id = a.PatientID
	}

	rows, err := inits.DB.Query(query, id)
	var appointments []Appointment

	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var appointment Appointment
		err = rows.Scan(&appointment.DoctorID, &appointment.PatientID, &appointment.AppointmentTime)
		if err != nil {
			panic(err.Error())
		}
		appointments = append(appointments, appointment)
	}

	return appointments

}

func (a Appointment) CheckIfViable() bool {
	
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM appointments
		WHERE doctor_id=? and ABS(TIMESTAMPDIFF(HOUR, appointment_time, ?)) != 1
	) AS result;
	`
	
	var isInvalid bool

	err := inits.DB.QueryRow(query, a.DoctorID, a.AppointmentTime).Scan(&isInvalid)

	if err != nil || isInvalid {
		return false;
	}

	return true;
}
