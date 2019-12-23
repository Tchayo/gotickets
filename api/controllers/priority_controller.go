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

// CreatePriority : description
func (server *Server) CreatePriority(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	priority := models.Priority{}
	err = json.Unmarshal(body, &priority)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	priority.Prepare()
	err = priority.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	priorityCreated, err := priority.SavePriority(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, priorityCreated.ID))
	responses.JSON(w, http.StatusCreated, priorityCreated)
}

// GetPriorities : description
func (server *Server) GetPriorities(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	priority := models.Priority{}

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

	priorities, err := priority.FindAllPriorities(server.DB, &f)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, priorities)
}

// GetPriority : description
func (server *Server) GetPriority(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tmid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	priority := models.Priority{}

	priorityReceived, err := priority.FindPriorityByID(server.DB, tmid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, priorityReceived)
}

// UpdatePriority : description
func (server *Server) UpdatePriority(w http.ResponseWriter, r *http.Request) {

	// check if auth token is valid and get the user id from it
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }

	vars := mux.Vars(r)

	// check if priority id is valid
	tmid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check if the priority exists
	priority := models.Priority{}
	err = server.DB.Debug().Model(models.Priority{}).Where("id = ?", tmid).Take(&priority).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Priority not found"))
		return
	}

	// read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// start processing request data
	priorityUpdate := models.Priority{}
	err = json.Unmarshal(body, &priorityUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	priorityUpdate.Prepare()
	err = priorityUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//this is important to tell the model the post id to update, the other update field are set above
	priorityUpdate.ID = priority.ID

	// process necessary userIDs in accoradance to update
	// to be done later

	priorityUpdated, err := priorityUpdate.UpdateAPriority(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, priorityUpdated)

}

// DeletePriority : delete a priority
func (server *Server) DeletePriority(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	priority := models.Priority{}

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
	_, err = priority.DeleteAPriority(server.DB, uint32(tmid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", tmid))
	responses.JSON(w, http.StatusNoContent, "")
}
