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

// CreateStatus : description
func (server *Server) CreateStatus(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	status := models.Status{}
	err = json.Unmarshal(body, &status)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	status.Prepare()
	err = status.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	statusCreated, err := status.SaveStatus(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, statusCreated.ID))
	responses.JSON(w, http.StatusCreated, statusCreated)
}

// GetStatuses : description
func (server *Server) GetStatuses(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	status := models.Status{}
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

	statuses, err := status.FindAllStatuses(server.DB, &f)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, statuses)
}

// GetStatus : description
func (server *Server) GetStatus(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tmid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	status := models.Status{}

	statusReceived, err := status.FindStatusByID(server.DB, tmid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, statusReceived)
}

// UpdateStatus : description
func (server *Server) UpdateStatus(w http.ResponseWriter, r *http.Request) {

	// check if auth token is valid and get the user id from it
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }

	vars := mux.Vars(r)

	// check if status id is valid
	tmid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check if the status exists
	status := models.Status{}
	err = server.DB.Debug().Model(models.Status{}).Where("id = ?", tmid).Take(&status).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Status not found"))
		return
	}

	// read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// start processing request data
	statusUpdate := models.Status{}
	err = json.Unmarshal(body, &statusUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	statusUpdate.Prepare()
	err = statusUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//this is important to tell the model the post id to update, the other update field are set above
	statusUpdate.ID = status.ID

	// process necessary userIDs in accoradance to update
	// to be done later

	statusUpdated, err := statusUpdate.UpdateAStatus(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, statusUpdated)

}

// DeleteStatus : delete a status
func (server *Server) DeleteStatus(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	status := models.Status{}

	tmid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = status.DeleteAStatus(server.DB, uint32(tmid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", tmid))
	responses.JSON(w, http.StatusNoContent, "")
}
