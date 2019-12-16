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
	"github.com/Tchayo/gotickets/api/utils/formaterror"
	"github.com/gorilla/mux"
)

// CreateTeam : description
func (server *Server) CreateTeam(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	team := models.Team{}
	err = json.Unmarshal(body, &team)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	team.Prepare()
	err = team.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	teamCreated, err := team.SaveTeam(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, teamCreated.ID))
	responses.JSON(w, http.StatusCreated, teamCreated)
}

// GetTeams : description
func (server *Server) GetTeams(w http.ResponseWriter, r *http.Request) {

	team := models.Team{}

	teams, err := team.FindAllTeams(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, teams)
}

// GetTeam : description
func (server *Server) GetTeam(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tmid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	team := models.Team{}

	teamReceived, err := team.FindTeamByID(server.DB, tmid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, teamReceived)
}

// UpdateTeam : description
func (server *Server) UpdateTeam(w http.ResponseWriter, r *http.Request) {

	// check if auth token is valid and get the user id from it
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }

	vars := mux.Vars(r)

	// check if team id is valid
	tmid, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check if the team exists
	team := models.Team{}
	err = server.DB.Debug().Model(models.Team{}).Where("id = ?", tmid).Take(&team).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Team not found"))
		return
	}

	// read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// start processing request data
	teamUpdate := models.Team{}
	err = json.Unmarshal(body, &teamUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	teamUpdate.Prepare()
	err = teamUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//this is important to tell the model the post id to update, the other update field are set above
	teamUpdate.ID = team.ID

	// process necessary userIDs in accoradance to update
	// to be done later

	teamUpdated, err := teamUpdate.UpdateATeam(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, teamUpdated)

}

// DeleteTeam : delete a team
func (server *Server) DeleteTeam(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	team := models.Team{}

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
	_, err = team.DeleteATeam(server.DB, uint32(tmid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", tmid))
	responses.JSON(w, http.StatusNoContent, "")
}
