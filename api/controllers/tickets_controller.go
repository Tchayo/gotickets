package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Tchayo/gotickets/api/auth"
	"github.com/Tchayo/gotickets/api/models"
	"github.com/Tchayo/gotickets/api/responses"
	"github.com/Tchayo/gotickets/api/utils/filter"
	"github.com/Tchayo/gotickets/api/utils/formaterror"
	"github.com/gorilla/mux"
)

// CreateTicket : description
func (server *Server) CreateTicket(w http.ResponseWriter, r *http.Request) {

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	ticket := models.Ticket{}
	err = json.Unmarshal(body, &ticket)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	ticket.Prepare()
	err = ticket.Validate("new")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if uid != ticket.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	ticketCreated, err := ticket.SaveTicket(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, ticketCreated.ID))
	responses.JSON(w, http.StatusCreated, ticketCreated)
}

// GetTickets : description
func (server *Server) GetTickets(w http.ResponseWriter, r *http.Request) {

	ticket := models.Ticket{}
	query := r.URL.Query()

	f := filter.Filter{}

	pg, pErr := strconv.Atoi(query.Get("page"))
	lm, lErr := strconv.Atoi(query.Get("limit"))
	sch := query.Get("search")

	if pErr != nil {
		pg = 1
	}

	if lErr != nil {
		lm = 10
	}

	f.Page = pg
	f.Limit = lm
	f.Search = sch

	tickets, err := ticket.FindAllTickets(server.DB, &f)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, tickets)
}

// GetTicket : description
func (server *Server) GetTicket(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	ticket := models.Ticket{}

	ticketReceived, err := ticket.FindTicketByID(server.DB, tid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, ticketReceived)
}

// UpdateTicket : description
func (server *Server) UpdateTicket(w http.ResponseWriter, r *http.Request) {

	// check if auth token is valid and get the user id from it
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }

	vars := mux.Vars(r)

	// check if ticket id is valid
	tid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check if the ticket exists
	ticket := models.Ticket{}
	err = server.DB.Debug().Model(models.Ticket{}).Where("id = ?", tid).Take(&ticket).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Ticket not found"))
		return
	}

	// read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// start processing request data
	ticketUpdate := models.Ticket{}
	err = json.Unmarshal(body, &ticketUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	ticketUpdate.Prepare()
	err = ticketUpdate.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//this is important to tell the model the post id to update, the other update field are set above
	ticketUpdate.ID = ticket.ID

	// process necessary userIDs in accoradance to update
	// to be done later

	ticketUpdated, err := ticketUpdate.UpdateATicket(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, ticketUpdated)

}
