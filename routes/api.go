package routes

func (route *Route) RouteApi() {
	api := route.App.Group("/api")

	controllers := route.Controllers
	api.Post("/login", controllers.Login)
	api.Post("register", controllers.Register)
	customer := api.Group("cust", route.Middlewares.AuthIsCustomer)
	customer.Get("/user", controllers.UserInfo)
	customer.Put("/user/:userid", controllers.UpdateUser)

}
