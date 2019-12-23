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

// CreateCategory : description
func (server *Server) CreateCategory(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	category := models.Category{}
	err = json.Unmarshal(body, &category)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	category.Prepare()
	err = category.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	categoryCreated, err := category.SaveCategory(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, categoryCreated.ID))
	responses.JSON(w, http.StatusCreated, categoryCreated)
}

// GetCategories : description
func (server *Server) GetCategories(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	category := models.Category{}
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

	categories, err := category.FindAllCategories(server.DB, &f)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, categories)
}

// GetCategory : description
func (server *Server) GetCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tmid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	category := models.Category{}

	categoryReceived, err := category.FindCategoryByID(server.DB, tmid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, categoryReceived)
}

// UpdateCategory : description
func (server *Server) UpdateCategory(w http.ResponseWriter, r *http.Request) {

	// check if auth token is valid and get the user id from it
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }

	vars := mux.Vars(r)

	// check if category id is valid
	tmid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check if the category exists
	category := models.Category{}
	err = server.DB.Debug().Model(models.Category{}).Where("id = ?", tmid).Take(&category).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Category not found"))
		return
	}

	// read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// start processing request data
	categoryUpdate := models.Category{}
	err = json.Unmarshal(body, &categoryUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	categoryUpdate.Prepare()
	err = categoryUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//this is important to tell the model the post id to update, the other update field are set above
	categoryUpdate.ID = category.ID

	// process necessary userIDs in accoradance to update
	// to be done later

	categoryUpdated, err := categoryUpdate.UpdateACategory(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, categoryUpdated)

}

// DeleteCategory : delete a category
func (server *Server) DeleteCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	category := models.Category{}

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
	_, err = category.DeleteACategory(server.DB, uint32(tmid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", tmid))
	responses.JSON(w, http.StatusNoContent, "")
}
