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
	s.Router.HandleFunc("/tickets", middlewares.SetMiddlewareJSON(s.CreateTicket)).Methods("POST")
	s.Router.HandleFunc("/tickets", middlewares.SetMiddlewareJSON(s.GetTickets)).Methods("GET")
	s.Router.HandleFunc("/tickets/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetTicket))).Methods("GET")
	s.Router.HandleFunc("/tickets/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateTicket))).Methods("PUT")

}
