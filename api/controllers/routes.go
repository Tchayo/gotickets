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

	//Priorities routes
	s.Router.HandleFunc("/priorities", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreatePriority))).Methods("POST")
	s.Router.HandleFunc("/priorities", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetPriorities))).Methods("GET")
	s.Router.HandleFunc("/priorities/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetPriority))).Methods("GET")
	s.Router.HandleFunc("/priorities/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePriority))).Methods("PUT")
	s.Router.HandleFunc("/priorities/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePriority)).Methods("DELETE")

	//Statuses routes
	s.Router.HandleFunc("/statuses", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateStatus))).Methods("POST")
	s.Router.HandleFunc("/statuses", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetStatuses))).Methods("GET")
	s.Router.HandleFunc("/statuses/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetStatus))).Methods("GET")
	s.Router.HandleFunc("/statuses/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateStatus))).Methods("PUT")
	s.Router.HandleFunc("/statuses/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteStatus)).Methods("DELETE")

	//Categories routes
	s.Router.HandleFunc("/categories", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateCategory))).Methods("POST")
	s.Router.HandleFunc("/categories", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCategories))).Methods("GET")
	s.Router.HandleFunc("/categories/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCategory))).Methods("GET")
	s.Router.HandleFunc("/categories/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCategory))).Methods("PUT")
	s.Router.HandleFunc("/categories/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCategory)).Methods("DELETE")

	//Subs routes
	s.Router.HandleFunc("/subs", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateSub))).Methods("POST")
	s.Router.HandleFunc("/subs", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetSubs))).Methods("GET")
	s.Router.HandleFunc("/subs/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetSub))).Methods("GET")
	s.Router.HandleFunc("/subs/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateSub))).Methods("PUT")
	s.Router.HandleFunc("/subs/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteSub)).Methods("DELETE")

}
