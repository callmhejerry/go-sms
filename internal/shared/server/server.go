package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	academicsession "github.com/callmhejerry/sms/internal/academic/academic_session"
	"github.com/callmhejerry/sms/internal/academic/classes"
	"github.com/callmhejerry/sms/internal/academic/student"
	"github.com/callmhejerry/sms/internal/admission"
	"github.com/callmhejerry/sms/internal/billing"
	"github.com/callmhejerry/sms/internal/grading/grading"
	"github.com/callmhejerry/sms/internal/grading/subjects"
	"github.com/callmhejerry/sms/internal/identity"
	"github.com/callmhejerry/sms/internal/inventory"
	"github.com/callmhejerry/sms/internal/shared/auth"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/tenant"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

type Handlers struct {
	Tenant          *tenant.Handler
	Identity        *identity.Handler
	AcademicSession *academicsession.Handler
	Classes         *classes.Handler
	Admission       *admission.Handler
	Student         *student.Handler
	Billing         *billing.BillingHandler
	Subject         *subjects.SubjectHandler
	Grading         *grading.GradingHandler
	Inventory       *inventory.InventoryHandler
}

func New(
	port string,
	pool *pgxpool.Pool,
	logger *slog.Logger,
	handlers Handlers,
	jwtManger *auth.JWTManager,
	roleChecker middleware.RoleChecker,
) *Server {
	mux := http.NewServeMux()

	// --------------------------
	// Public routes (no auth)
	// --------------------------
	healthHanler := NewHealthHandler(pool)
	mux.HandleFunc("GET /healthz", healthHanler.Healthz)
	mux.HandleFunc("GET /readyz", healthHanler.Readyz)

	mux.HandleFunc("POST /api/v1/tenants", handlers.Tenant.CreateTenant)
	mux.HandleFunc("POST /api/v1/login", handlers.Identity.Login)

	// --------------------------
	// Protected routes
	// --------------------------

	adminOnly := middleware.RequireRole(roleChecker, "admin")

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /api/v1/tenants", handlers.Tenant.ListTenants)
	protectedMux.HandleFunc("GET /api/v1/tenants/{id}", handlers.Tenant.GetTenant)

	// IDENTITY ROUTE

	protectedMux.Handle("POST /api/v1/users", adminOnly(http.HandlerFunc(handlers.Identity.CreateUser)))
	protectedMux.HandleFunc("GET /api/v1/tenants/{tenant_id}/users/{id}", handlers.Identity.GetUser)

	// ACADEMIC ROUTE
	protectedMux.Handle("POST /api/v1/academic-sessions", adminOnly(http.HandlerFunc(handlers.AcademicSession.CreateAcademicSession)))
	protectedMux.HandleFunc("GET /api/v1/academic-sessions", handlers.AcademicSession.ListAcademicSessions)
	protectedMux.HandleFunc("GET /api/v1/academic-sessions/current", handlers.AcademicSession.GetCurrentSession)

	// CLASSES ROUTES
	protectedMux.Handle("POST /api/v1/classes", adminOnly(http.HandlerFunc(handlers.Classes.CreateClass)))
	protectedMux.HandleFunc("GET /api/v1/classes", handlers.Classes.ListClasses)
	protectedMux.Handle("POST /api/v1/class-arms", adminOnly(http.HandlerFunc(handlers.Classes.CreateClassArm)))
	protectedMux.HandleFunc("GET /api/v1/classes/{class_id}/arms", handlers.Classes.ListClassArms)

	// STUDENTS ROUTES
	protectedMux.Handle("POST /api/v1/students", adminOnly(http.HandlerFunc(handlers.Student.CreateStudent)))
	protectedMux.HandleFunc("GET /api/v1/students-page", handlers.Student.ListStudentsPage)
	protectedMux.HandleFunc("GET /api/v1/students/{id}", handlers.Student.GetStudent)
	protectedMux.HandleFunc("GET /api/v1/students/{id}/parents", handlers.Student.GetStudentParents)
	protectedMux.HandleFunc("GET /api/v1/students/{id}/profile", handlers.Student.GetStudentProfile)
	protectedMux.HandleFunc("GET /api/v1/students/search-page", handlers.Student.SearchStudentsPage)
	protectedMux.HandleFunc("GET /api/v1/students/search-cursor", handlers.Student.SearchStudentsCursor)
	protectedMux.HandleFunc("PATCH /api/v1/students/{id}", handlers.Student.UpdateStudent)

	// ADMISSIONS ROUTE
	protectedMux.Handle("POST /api/v1/admissions", adminOnly(http.HandlerFunc(handlers.Admission.CreateAdmission)))
	protectedMux.Handle("GET /api/v1/admissions", adminOnly(http.HandlerFunc(handlers.Admission.ListAdmissions)))
	protectedMux.Handle("POST /api/v1/admissions/{id}/accept", adminOnly(http.HandlerFunc(handlers.Admission.AcceptAdmission)))

	// Fee Types
	protectedMux.HandleFunc("POST /api/v1/fee-types", handlers.Billing.CreateFeeType)
	protectedMux.HandleFunc("GET /api/v1/fee-types", handlers.Billing.ListFeeTypes)

	// Fee Structures
	protectedMux.HandleFunc("POST /api/v1/fee-structures", handlers.Billing.CreateFeeStructure)
	protectedMux.HandleFunc("GET /api/v1/fee-structures", handlers.Billing.ListFeeStructures)

	// Student Fees
	protectedMux.HandleFunc("POST /api/v1/student-fees/assign", handlers.Billing.AssignFeesToStudent)
	protectedMux.HandleFunc("GET /api/v1/students/{student_id}/fees", handlers.Billing.ListStudentFees)
	protectedMux.HandleFunc("GET /api/v1/students/{student_id}/fees/summary", handlers.Billing.GetStudentFeeSummary)
	protectedMux.HandleFunc("GET /api/v1/fees/outstanding", handlers.Billing.ListOutstandingFees)

	// Payments
	protectedMux.HandleFunc("POST /api/v1/payments", handlers.Billing.RecordPayment)
	protectedMux.HandleFunc("GET /api/v1/students/{student_id}/payments", handlers.Billing.ListStudentPayments)

	// Subjects
	protectedMux.HandleFunc("POST /api/v1/subjects", handlers.Subject.CreateSubject)
	protectedMux.HandleFunc("GET /api/v1/subjects", handlers.Subject.ListSubjects)
	protectedMux.HandleFunc("POST /api/v1/class-subjects", handlers.Subject.AddSubjectToClass)

	// Teacher assignments
	protectedMux.HandleFunc("POST /api/v1/teacher-assignments", handlers.Subject.AssignTeacher)

	// assessment types
	protectedMux.HandleFunc("POST /api/v1/assessment-types", handlers.Grading.CreateAssessmentType)
	protectedMux.HandleFunc("GET /api/v1/assessment-types", handlers.Grading.ListAssessmentTypes)

	// Results
	protectedMux.HandleFunc("POST /api/v1/results/compute", handlers.Grading.ComputeResults)
	protectedMux.HandleFunc("GET /api/v1/students/{student_id}/results", handlers.Grading.GetStudentResults)
	protectedMux.HandleFunc("GET /api/v1/students/{student_id}/report-card", handlers.Grading.GetStudentReportCard)

	// Inventory Categories
	protectedMux.HandleFunc("POST /api/v1/inventory/categories", handlers.Inventory.CreateCategory)
	protectedMux.HandleFunc("GET /api/v1/inventory/categories", handlers.Inventory.ListCategories)

	// Inventory Items
	protectedMux.HandleFunc("POST /api/v1/inventory/items", handlers.Inventory.CreateItem)
	protectedMux.HandleFunc("GET /api/v1/inventory/items", handlers.Inventory.ListItems)

	// MIDDLEWARE CHAIN
	protectedHandler := middleware.AuthMiddleware(jwtManger)(protectedMux)
	mux.Handle("/", protectedHandler)

	var handler http.Handler = mux

	handler = middleware.RequestID(handler)
	handler = middleware.Recovery(logger)(handler)

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	server := &Server{
		httpServer: httpServer,
		logger:     logger,
	}
	return server
}

func (server *Server) Start() error {
	server.logger.Info("Starting HTTP SERVER", slog.String("addr", server.httpServer.Addr))
	return server.httpServer.ListenAndServe()
}

func (server *Server) ShutDown(ctx context.Context) error {
	server.logger.Info("Shutting down HTTP server")
	return server.httpServer.Shutdown(ctx)
}
