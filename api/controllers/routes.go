package controllers

import "github.com/Tchayo/gotickets/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Tickets routes
	s.Router.HandleFunc("/tickets", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateTicket))).Methods("POST")
	s.Router.HandleFunc("/tickets", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetTickets))).Methods("GET")
	s.Router.HandleFunc("/tickets/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetTicket))).Methods("GET")
	s.Router.HandleFunc("/tickets/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateTicket))).Methods("PUT")

	//Teams routes
	s.Router.HandleFunc("/teams", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateTeam))).Methods("POST")
	s.Router.HandleFunc("/teams", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetTeams))).Methods("GET")
	s.Router.HandleFunc("/teams/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetTeam))).Methods("GET")
	s.Router.HandleFunc("/teams/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateTeam))).Methods("PUT")
	s.Router.HandleFunc("/teams/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteTeam)).Methods("DELETE")

}
