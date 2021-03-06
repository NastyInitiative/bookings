package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/NastyInitiative/bookings/internal/config"
	"github.com/NastyInitiative/bookings/internal/driver"
	"github.com/NastyInitiative/bookings/internal/forms"
	"github.com/NastyInitiative/bookings/internal/helpers"
	"github.com/NastyInitiative/bookings/internal/models"
	"github.com/NastyInitiative/bookings/internal/render"
	"github.com/NastyInitiative/bookings/internal/repository"
	"github.com/NastyInitiative/bookings/internal/repository/dbrepo"

	// "bookings/pkg/render"
	"net/http"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo - the repository used by the handlers
var Repo *Repository

// Creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// Sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	m.DB.AllUsers()
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation - renders the make reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation - handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, "During parsing form", err)
		return
	}

	startDateValue := r.Form.Get("start_date")
	endDateValue := r.Form.Get("end_date")

	// Go time format reference -> 2021-01-01 or 01/01/2021 => 01/02 03:04:05PM '06 -0700

	layout := "01/02/2006"
	startDateParsed, err := time.Parse(layout, startDateValue)
	if err != nil {
		helpers.ServerError(w, "During parsing start date string value", err)
		return
	}

	endDateParsed, err := time.Parse(layout, endDateValue)
	if err != nil {
		helpers.ServerError(w, "During parsing end date string value", err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, "During converting room_id to int", err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, "During saving a new reservation on database", err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDateParsed,
		EndDate:       endDateParsed,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, "During inserting a new room restrinction on database", err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals - renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors - renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability - renders the search availability page
func (m *Repository) Availabilty(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability - renders the search availability page
func (m *Repository) PostAvailabilty(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	// w.Write([]byte("Posted to search availabilty"))
	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// PostAvailability - handles request for availability and sends JSON response
func (m *Repository) AvailabiltyJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, "During marshaling resp to JSON", err)
		return
	}

	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")

	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, req *http.Request) {

	reservation, ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(req.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(req.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, req, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
