package models

import (
	"database/sql"
	"log"

	"clinic-reservation-system.com/back-end/inits"
)

type Appointment struct {
	ID              sql.NullInt64  `json:"id"`
	DoctorID        sql.NullInt64  `json:"doctor_id"`
	Name 			sql.NullString `json:"doctor_name"`
	PatientID       sql.NullInt64  `json:"patient_id"`
	AppointmentTime sql.NullString `json:"appointment_time"`
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

func (a Appointment) Reserve() bool {
	query := `
	UPDATE appointments SET patient_id=? WHERE id=? AND patient_id IS NULL;
	`
	_, err := inits.DB.Exec(query, a.PatientID, a.ID)

	return err == nil
}

func (a Appointment) Delete() bool {
	query := `
	DELETE FROM appointments WHERE id=? and doctor_id=?;
	`
	_, err := inits.DB.Exec(query, a.ID, a.DoctorID)

	return err == nil
}

func (a Appointment) GetReserved(userType string) []Appointment {
	var query string
	var id sql.NullInt64

	if userType == "doctor" {
		query = `
		select users.name, appointments.id, appointments.doctor_id, appointments.patient_id, appointments.appointment_time from users inner join appointments where users.id = appointments.patient_id and appointments.doctor_id = ?;
		`
		id = a.DoctorID
	} else {
		query = `
		select users.name, appointments.id, appointments.doctor_id, appointments.patient_id, appointments.appointment_time from users inner join appointments where users.id = appointments.doctor_id and appointments.patient_id = ?;
		`
		id = a.PatientID
	}

	var rows *sql.Rows
	var err error

	idVal, err := id.Value()

	if err != nil {
		log.Println(err.Error())
		return nil
	}


	rows, err = inits.DB.Query(query, idVal)

	
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	
	defer rows.Close()
	
	var appointments []Appointment
	
	for rows.Next() {
		var appointment Appointment
		err = rows.Scan(&appointment.Name, &appointment.ID, &appointment.DoctorID, &appointment.PatientID, &appointment.AppointmentTime)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		appointments = append(appointments, appointment)
	}

	return appointments

}

func (a Appointment) GetAll(userType string) []Appointment{
	var query string
	var id sql.NullInt64

	if userType == "doctor" {
		query = `
		select users.name, appointments.id, appointments.doctor_id, appointments.patient_id, appointments.appointment_time from appointments left join users on users.id = appointments.patient_id where appointments.doctor_id = ?;
		`
		id = a.DoctorID
	} else {
		query = `
		select users.name, appointments.id, appointments.doctor_id, appointments.patient_id, appointments.appointment_time from appointments left join users on users.id = appointments.doctor_id where appointments.patient_id = ?;
		`
		id = a.PatientID
	}

	var rows *sql.Rows
	var err error

	idVal, err := id.Value()

	if err != nil {
		log.Println(err.Error())
		return nil
	}


	rows, err = inits.DB.Query(query, idVal)

	
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	
	defer rows.Close()
	
	var appointments []Appointment
	
	for rows.Next() {
		var appointment Appointment
		err = rows.Scan(&appointment.Name, &appointment.ID, &appointment.DoctorID, &appointment.PatientID, &appointment.AppointmentTime)
		if err != nil {
			log.Println(err.Error())
			return nil
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
		WHERE doctor_id=? and ABS(TIMESTAMPDIFF(HOUR, appointment_time, ?)) < 1
	) AS result;
	`

	var isInvalid bool

	time, _ := a.AppointmentTime.Value()
	err := inits.DB.QueryRow(query, a.DoctorID, time).Scan(&isInvalid)

	if err != nil || isInvalid {
		return false
	}

	return true
}

func (a Appointment) UnReserve() bool {
	query := `
	UPDATE appointments SET patient_id=NULL WHERE id=? and patient_id=?;
	`
	_, err := inits.DB.Exec(query, a.ID, a.PatientID)

	return err == nil
}

func (a Appointment) Edit() bool {
	query := `
	UPDATE appointments SET appointment_time=? WHERE id=? and patient_id=?;
	`
	_, err := inits.DB.Exec(query, a.AppointmentTime, a.ID, a.PatientID)

	return err == nil
}
