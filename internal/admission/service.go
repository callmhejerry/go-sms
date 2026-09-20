package admission

import (
	"context"
	"fmt"
	"strings"
	"time"

	academicsession "github.com/callmhejerry/sms/internal/academic/academic_session"
	"github.com/callmhejerry/sms/internal/academic/classes"
	"github.com/callmhejerry/sms/internal/academic/student"
	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/constants"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries         *store.Queries
	studentService  *student.Service
	academicService *academicsession.Service
	classService    *classes.Service
}

func NewService(
	queries *store.Queries,
	studentService *student.Service,
	academicSessionService *academicsession.Service,
	classService *classes.Service,
) *Service {
	return &Service{
		queries:         queries,
		studentService:  studentService,
		academicService: academicSessionService,
	}
}
func (service *Service) CreateAdmission(ctx context.Context, tenantId uuid.UUID, request CreateAdmissionRequest) (*store.Admission, error) {
	firstName := strings.TrimSpace(request.FirstName)
	lastName := strings.TrimSpace(request.LastName)
	gender := strings.TrimSpace(strings.ToLower(request.Gender))
	dateOfBirth, err := time.Parse("2026-01-30", request.DateOfBirth)

	parentFirstName := strings.TrimSpace(request.ParentFirstName)
	parentLastName := strings.TrimSpace(request.ParentLastName)
	parentEmail := strings.TrimSpace(strings.ToLower(request.ParentEmail))
	parentPhone := strings.TrimSpace(request.ParentPhoneNumber)

	if firstName == "" || lastName == "" {
		return nil, apierror.Validation("first_name and last_name are required")
	}

	if gender != string(constants.Male) && gender != string(constants.Female) {
		return nil, apierror.Validation("Invalid gender, gender must be either male or female")
	}

	academicSession, err := service.academicService.GetAcademicSessionById(
		ctx, tenantId, request.AcademicSessionID,
	)

	if err != nil {
		return nil, err
	}

	preferredClass, err := service.classService.GetClassById(ctx, tenantId, request.PreferredClassID)

	if err != nil {
		return nil, err
	}

	admission, err := service.queries.CreateAdmission(ctx, store.CreateAdmissionParams{
		TenantID:           pgtype.UUID{Bytes: tenantId, Valid: true},
		AcademicSessionID:  academicSession.ID,
		FirstName:          firstName,
		LastName:           lastName,
		MiddleName:         request.MiddleName,
		Gender:             gender,
		DateOfBirth:        pgtype.Date{Time: dateOfBirth, Valid: true},
		PreferredClassID:   preferredClass.ID,
		ParentFirstName:    parentFirstName,
		ParentLastName:     parentLastName,
		ParentPhoneNumber:  parentPhone,
		ParentEmail:        parentEmail,
		ParentRelationship: request.Relationship,
	})

	if err != nil {
		return nil, database.TranslateError(err)
	}

	return &admission, nil
}

func (service *Service) ListAdmissions(ctx context.Context, tenantId uuid.UUID) ([]store.Admission, error) {
	admissions, err := service.queries.ListAdmissions(ctx, pgtype.UUID{Bytes: tenantId, Valid: true})

	if err != nil {
		return nil, database.TranslateError(err)
	}
	return admissions, nil
}

func (service *Service) AcceptAdmission(ctx context.Context, tenantId, admissionId, reviewedBy uuid.UUID, classArmId *uuid.UUID) (*store.Admission, *store.Student, error) {

	admission, err := service.queries.GetAdmissionById(ctx, store.GetAdmissionByIdParams{
		ID:       pgtype.UUID{Bytes: admissionId, Valid: true},
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
	})
	if err != nil {
		return nil, nil, database.TranslateError(err)
	}

	if admission.Status != string(constants.AdmissionPending) && admission.Status != string(constants.AdmissionUnderReview) {
		return nil, nil, apierror.Validation("Only pending or under_review admission can be accepted")
	}
	admissionNumber := fmt.Sprintf("ADM/%s/%d", time.Now().Format("2006"), time.Now().Unix()%10000)

	student, err := service.studentService.CreateStudent(ctx, tenantId, student.CreateStudentRequest{
		AdmissionNumber:  admissionNumber,
		FirstName:        admission.FirstName,
		LastName:         admission.LastName,
		MiddleName:       admission.MiddleName,
		Gender:           admission.Gender,
		DateOfBirth:      admission.DateOfBirth.Time.Format("2026-01-30"),
		CurrentClassArm:  classArmId,
		AdmissionSession: (*uuid.UUID)(&admission.AcademicSessionID.Bytes),
	})

	if err != nil {
		return nil, nil, err
	}

	updatedAdmission, err := service.queries.UpdateAdmissionStatus(ctx, store.UpdateAdmissionStatusParams{
		ID:              admission.ID,
		TenantID:        pgtype.UUID{Bytes: tenantId, Valid: true},
		Status:          string(constants.AdmissionAccepted),
		ReviewedBy:      pgtype.UUID{Bytes: reviewedBy, Valid: true},
		AdmissionNumber: &admissionNumber,
		StudentID:       student.ID,
	})

	if err != nil {
		return nil, nil, database.TranslateError(err)
	}
	return &updatedAdmission, student, nil
}
