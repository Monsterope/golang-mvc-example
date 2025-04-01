package routes

func (route *Route) RouteBo() {
	api := route.App.Group("/api")
	controller := route.Controllers
	middleware := route.Middlewares

	bo := api.Group("bo", middleware.AuthIsAdmin)
	bo.Get("/user", controller.UserInfo)
	bo.Put("/user/:userid", controller.UpdateUser)

	bo.Get("/category", controller.GetCategoryAll)
	bo.Post("/category", controller.CreateCategory)
	bo.Put("/category/:cateid", controller.UpdateCategory)

}
