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

// CreateSub : description
func (server *Server) CreateSub(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	sub := models.Sub{}
	err = json.Unmarshal(body, &sub)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	sub.Prepare()
	err = sub.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	subCreated, err := sub.SaveSub(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, subCreated.ID))
	responses.JSON(w, http.StatusCreated, subCreated)
}

// GetSubs : description
func (server *Server) GetSubs(w http.ResponseWriter, r *http.Request) {

	sub := models.Sub{}
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

	subs, err := sub.FindAllSubs(server.DB, &f)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, subs)
}

// GetSub : description
func (server *Server) GetSub(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tmid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	sub := models.Sub{}

	subReceived, err := sub.FindSubByID(server.DB, tmid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, subReceived)
}

// UpdateSub : description
func (server *Server) UpdateSub(w http.ResponseWriter, r *http.Request) {

	// check if auth token is valid and get the user id from it
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }

	vars := mux.Vars(r)

	// check if sub id is valid
	tmid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check if the sub exists
	sub := models.Sub{}
	err = server.DB.Debug().Model(models.Sub{}).Where("id = ?", tmid).Take(&sub).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Sub not found"))
		return
	}

	// read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// start processing request data
	subUpdate := models.Sub{}
	err = json.Unmarshal(body, &subUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	subUpdate.Prepare()
	err = subUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//this is important to tell the model the post id to update, the other update field are set above
	subUpdate.ID = sub.ID

	// process necessary userIDs in accoradance to update
	// to be done later

	subUpdated, err := subUpdate.UpdateASub(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, subUpdated)

}

// DeleteSub : delete a sub
func (server *Server) DeleteSub(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	sub := models.Sub{}

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
	_, err = sub.DeleteASub(server.DB, uint32(tmid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", tmid))
	responses.JSON(w, http.StatusNoContent, "")
}
