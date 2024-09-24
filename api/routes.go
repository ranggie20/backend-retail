package api

import (
	"database/sql"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/online-bnsp/backend/api/cart"
	"github.com/online-bnsp/backend/api/categories"
	"github.com/online-bnsp/backend/api/courses"
	coursesvideo "github.com/online-bnsp/backend/api/courses_video"
	"github.com/online-bnsp/backend/api/payment"
	paymentmethod "github.com/online-bnsp/backend/api/payment_method"
	"github.com/online-bnsp/backend/api/paymentstatus"
	"github.com/online-bnsp/backend/api/subscriptions"
	"github.com/online-bnsp/backend/api/user"
	"github.com/online-bnsp/backend/api/wishlist"
	"github.com/online-bnsp/backend/middleware"
	"github.com/online-bnsp/backend/middleware/auth"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
	"github.com/online-bnsp/backend/util/buckets"
	"github.com/online-bnsp/backend/util/mailer"
	queue "github.com/online-bnsp/backend/util/queue"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	router *chi.Mux
}

type RoleMiddleware func(http.Handler) http.Handler

var validate *validator.Validate

func New(db *sql.DB, rdb *redis.Client, q queue.Queuer, bucket buckets.Bucket, mail *mailer.Mailer, cors RoleMiddleware) *Handler {
	r := chi.NewMux()
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.BirthTime)
	r.Use(auth.ExtractTokenClaims) // Extract JWT claims into context
	r.Use(cors)                    // CORS Middleware, if needed

	h := &Handler{
		router: r,
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	dbGenerated := repo.New(db)

	r.Get("/ping", h.Ping)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		util.NewResponse(http.StatusNotFound, http.StatusNotFound, "404 Not found!", nil).WriteResponse(w, r)
	})

	//payment Handler
	PaymentHandler := payment.NewHandler(validate, dbGenerated)
	// Routes for payment
	r.Route("/payment", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("student"))

		r.Post("/create-payment", PaymentHandler.CreatePayment)
		r.Get("/get-payment", PaymentHandler.GetPayment)
		r.Get("/get-paymenthistory", PaymentHandler.GetPaymentHistory)
	})

	//paymentmethod Handler
	PaymentMethodHandler := paymentmethod.NewHandler(validate, dbGenerated)
	// Routes for paymentmethod
	r.Route("/payment-method", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("student"))

		r.Post("/create-paymentmethod", PaymentMethodHandler.CreatePaymentMethod)
		r.Get("/get-paymentmethod", PaymentMethodHandler.GetAllPaymentMethod)
		r.Get("/get-paymentmethod/{id}", PaymentMethodHandler.GetPaymentMethodById)
	})

	//paymentstatus Handler
	PaymentStatusHandler := paymentstatus.NewHandler(validate, dbGenerated)
	r.Route("/paymentstatus", func(r chi.Router) {
		r.Use(auth.RequireRole("student"))

		r.Post("/create-paymentsstatus", PaymentStatusHandler.CreatePaymentStatus)
		r.Get("/get-paymentsstatus", PaymentStatusHandler.GetAllPaymentStatus)
		r.Get("/get-paymentsstatus/{id}", PaymentStatusHandler.GetPaymentStatusById)
	})
	// Routes for paymentsstatus

	//subscriptions Handler
	SubscriptionHandler := subscriptions.NewHandler(validate, dbGenerated)
	r.Route("/subscription", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("student"))

		// Routes for subscriptions
		r.Post("/create-subscription", SubscriptionHandler.CreateSubscription)
		r.Get("/get-subscription", SubscriptionHandler.GetAllSubscriptions)
	})

	// Course Handler
	CoursesHandler := courses.NewHandler(validate, dbGenerated)
	// Routes for courses

	r.Route("/my-course", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("student"))

		r.Get("/", CoursesHandler.GetMyCoursePage)
		r.Get("/{course_id}", CoursesHandler.GetMyCourse)
	})
	r.Route("/teacher", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("teacher", "admin"))

		r.Post("/create-course", CoursesHandler.CreateCourses)
		r.Put("/update-course/{id}", CoursesHandler.UpdateCourse)
		r.Delete("/delete-course/{id}", CoursesHandler.DeleteCourse)
	})

	// Cart Handler
	CartHandler := cart.NewHandler(validate, dbGenerated)
	r.Route("/cart", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("student"))

		r.Post("/create-cart", CartHandler.CreateCart)
		r.Get("/getall-cart", CartHandler.GetAllCart)
		r.Delete("/delete-cart/{course_id}", CartHandler.DeleteCart)
		r.Get("/cartpage", CartHandler.GetCartByUserID)
	})

	//course_video handler
	coursesVideo := coursesvideo.NewHandler(validate, dbGenerated)

	r.Get("/course_video", coursesVideo.GetCourseVideoHandler)
	// route course_video
	r.Route("/course_video", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("teacher"))

		r.Post("/create-course_video", coursesVideo.CreateCourseVideo)
		r.Put("/update-course_video/{id}", coursesVideo.UpdateCourseVideo)
		r.Delete("/delete-coursevideo/{id}", coursesVideo.DeleteCourseVideo)
	})

	// Wishlist Handler
	WishlistHandler := wishlist.NewHandler(validate, dbGenerated)
	r.Route("/wishlist", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("student"))

		r.Post("/create-wishlist", WishlistHandler.CreateWishlist)
		r.Get("/wishlist", WishlistHandler.GetAllWishlist)
		r.Delete("/delete-wishlist/{course_id}", WishlistHandler.DeleteWishlist)
	})

	// User Handler
	userHandler := user.NewHandler(validate, dbGenerated)
	r.Route("/my-user", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("student"))

		r.Put("/profile/{id}", userHandler.UpdateUser)
		r.Get("/list-teacher", userHandler.GetAllUserByTeacher)
	})

	//route user
	r.Route("/user", func(r chi.Router) {
		r.Get("/all-user", userHandler.GetAllUser)
		r.Get("/list-teacher", userHandler.GetAllUserByTeacher)
		r.Get("/list-student", userHandler.GetAllUserByStudent)
		r.Get("/user/{id}", userHandler.GetUserByID)
		r.Put("/profile/{id}", userHandler.UpdateUser)

		r.Post("/register", userHandler.Register)
		r.Post("/sign-in", userHandler.Login)

		r.Route("/", func(r chi.Router) {
			r.Use(auth.AuthMiddleware)

			r.Post("/profile", userHandler.UpdateUser)
			r.Get("/user-info", userHandler.UserInfo)
			r.Post("/sign-out", userHandler.Logout)
		})
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("admin"))

		r.Get("/all-user", userHandler.GetAllUser)
		r.Get("/list-teacher", userHandler.GetAllUserByTeacher)
		r.Get("/list-student", userHandler.GetAllUserByStudent)
		r.Get("/list-payment", PaymentHandler.GetAllPayment)
		r.Get("/list-subscription", SubscriptionHandler.GetAllSubscriptions)
	})
	// Category Handler
	CategoryHandler := categories.NewHandler(validate, dbGenerated)
	// Routes for categories
	r.Route("/category", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(auth.RequireRole("teacher", "admin")) // Middleware to require 'teacher' role

		r.Post("/create-category", CategoryHandler.CreateCategory)
		r.Put("/update-category/{id}", CategoryHandler.UpdateCategory)
		r.Delete("/delete-category/{id}", CategoryHandler.DeleteCategory)
	})

	workingDir, _ := os.Getwd()
	fileServer := http.FileServer(http.Dir(path.Join(workingDir, "public")))
	r.Route("/static", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/static/", fileServer))
	})

	r.Route("/public", func(r chi.Router) {
		r.Get("/category", CategoryHandler.GetAllCategories)
		r.Get("/categoryy", CategoryHandler.GetCategory) //getbynamecategory
		r.Get("/category/{id}", CategoryHandler.GetCategoryByID)
		r.Get("/get-category/{category_id}", CategoryHandler.GetCoursesByCategoryID)
		r.Get("/course_video", coursesVideo.GetAllCourseVideos)
		r.Get("/course_video/{id}", coursesVideo.GetCourseVideoByID)
		r.Get("/getall-course", CoursesHandler.GetAllCourses)
		r.Get("/home", CoursesHandler.GetCourseByNew)
		r.Get("/popular", CoursesHandler.GetPopularCourses)
		r.Get("/price", CoursesHandler.GetCoursePrice)
		r.Get("/get-course/{course_id}", CoursesHandler.GetCourseByID)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Get("/get-user-info", userHandler.GetUserInfo)
	})

	return h
}

func (h *Handler) Handler() http.Handler {
	return h.router
}
