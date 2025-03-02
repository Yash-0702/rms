package server

import (
	"fmt"
	"net/http"
	"rms/handlers"
	"rms/middlewares"
	"rms/models"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Router chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to RMS API"))
	})

	router.Route("/v1", func(v1 chi.Router) {
		v1.Post("/login", handlers.LoginUser)
		// v1.Post("/register", handlers.RegisterUser)

		v1.Group(func(r chi.Router) {
			r.Use(middlewares.Authenticate)

			r.Post("/logout", handlers.LogoutUser)

			// admin routes

			r.Route("/admin", func(admin chi.Router) {
				admin.Use(middlewares.ShouldHaveRole(models.RoleAdmin))
				admin.Post("/create-sub-admin", handlers.CreateSubAdmin)
				admin.Get("/all-sub-admin", handlers.GetAllSubAdmins)
				admin.Post("/create-user", handlers.CreateUser)
				admin.Get("/all-users", handlers.GetAllUsersByAdminAndSubAdmin)
				admin.Post("/create-restaurant", handlers.CreateRestaurant)
				admin.Get("/all-restaurants", handlers.GetAllRestaurants)

				admin.Route("/{restaurantId}", func(restaurantIDRoute chi.Router) {
					restaurantIDRoute.Post("/create-dish", handlers.CreateDish)
					restaurantIDRoute.Get("/all-dishes", handlers.GetAllDishesFromSpecificRestaurant)
					restaurantIDRoute.Get("/", handlers.GetSpecificRestaurant)

					restaurantIDRoute.Route("/{dishId}", func(dishIDRoute chi.Router) {
						dishIDRoute.Put("/update-dish", handlers.UpdateDish)
						dishIDRoute.Delete("/delete-dish", handlers.DeleteDish)
						dishIDRoute.Get("/", handlers.GetSpecificDish)
					})
				})

				admin.Get("/all-dishes", handlers.GetAllDishesFromAllRestaurants)
			})

			// sub-admin routes

			r.Route("/sub-admin", func(subAdmin chi.Router) {
				subAdmin.Use(middlewares.ShouldHaveRole(models.RoleSubAdmin))
				subAdmin.Post("/create-user", handlers.CreateUser)
				subAdmin.Get("/all-users", handlers.GetAllUsersByAdminAndSubAdmin)
				subAdmin.Post("/create-restaurant", handlers.CreateRestaurant)
				subAdmin.Get("/all-restaurants", handlers.GetAllRestaurants)

				subAdmin.Route("/{restaurantId}", func(restaurantIDRoute chi.Router) {
					restaurantIDRoute.Post("/create-dish", handlers.CreateDish)
					restaurantIDRoute.Get("/all-dishes", handlers.GetAllDishesFromSpecificRestaurant)
					restaurantIDRoute.Get("/", handlers.GetSpecificRestaurant)

					restaurantIDRoute.Route("/{dishId}", func(dishIDRoute chi.Router) {
						dishIDRoute.Put("/update-dish", handlers.UpdateDish)
						dishIDRoute.Delete("/delete-dish", handlers.DeleteDish)
						dishIDRoute.Get("/", handlers.GetSpecificDish)
					})

				})

				subAdmin.Get("/all-dishes", handlers.GetAllDishesFromAllRestaurants)
			})

			// user routes

			r.Route("/users", func(user chi.Router) {
				user.Use(middlewares.ShouldHaveRole(models.RoleUser))
				user.Get("/profile", handlers.GetUser)
				user.Post("/add-address", handlers.AddAddress)
				user.Get("/all-addresses", handlers.GetAllAddress)

				user.Route("/address/{addressId}", func(addressIDRoute chi.Router) {
					addressIDRoute.Delete("/delete-address", handlers.DeleteAddress)
					addressIDRoute.Put("/update-address", handlers.UpdateAddress)
					addressIDRoute.Get("/", handlers.GetSpecificAddress)
				})

				user.Get("/all-restaurants", handlers.GetAllRestaurants)
				user.Get("/all-dishes", handlers.GetAllDishesFromAllRestaurants)
				user.Get("/calculate-distance", handlers.CalculateDistance)

				user.Route("/{restaurantId}", func(restaurantIDRoute chi.Router) {
					restaurantIDRoute.Get("/", handlers.GetSpecificRestaurant)
					restaurantIDRoute.Get("/all-dishes", handlers.GetAllDishesFromSpecificRestaurant)

					restaurantIDRoute.Route("/{dishId}", func(dishIDRoute chi.Router) {
						dishIDRoute.Get("/", handlers.GetSpecificDish)
					})

				})
			})

		})
	})

	return &Server{Router: router}
}

func (svc *Server) Run(port string) error {

	svc.server = &http.Server{
		// Addr:    ":" + port,
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}

	fmt.Println("Server running on port " + port)

	return svc.server.ListenAndServe()
}
