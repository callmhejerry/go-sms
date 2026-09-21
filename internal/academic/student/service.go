package student

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/callmhejerry/sms/internal/academic"
	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/constants"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/callmhejerry/sms/internal/shared/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *store.Queries
	pool    *pgxpool.Pool
}

func NewService(queries *store.Queries, pool *pgxpool.Pool) *Service {
	return &Service{
		queries: queries,
		pool:    pool,
	}
}

func (service *Service) CreateStudent(ctx context.Context, tenantId uuid.UUID, request CreateStudentRequest) (*store.Student, error) {
	admissionNumber := strings.TrimSpace(request.AdmissionNumber)
	firstName := strings.TrimSpace(request.FirstName)
	lastName := strings.TrimSpace(request.LastName)
	gender := strings.TrimSpace(strings.ToLower(request.Gender))
	dateOfBirth, _ := time.Parse("2026-01-30", request.DateOfBirth)

	if gender != string(constants.Male) && gender != string(constants.Female) {
		return nil, apierror.Validation("gender must either be male or female")
	}

	var newStudent store.Student

	err := database.WithTx(ctx, service.pool, func(q *store.Queries) error {
		var classArmId uuid.UUID
		var academicSessionId uuid.UUID

		if request.CurrentClassArm != nil {
			classArm, err := q.GetClassArmByID(ctx, store.GetClassArmByIDParams{
				ID:       *request.CurrentClassArm,
				TenantID: tenantId,
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					message := fmt.Sprintf("class arm with id %s. Not found", request.CurrentClassArm.String())
					return apierror.NotFound(message)
				}
				return apierror.Validation("Invalid class arm id")
			}
			classArmId = classArm.ID
		}

		if request.AdmissionSession != nil {
			academicSession, err := q.GetAcademicSessionById(ctx, store.GetAcademicSessionByIdParams{
				ID:       *request.AdmissionSession,
				TenantID: tenantId,
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					message := fmt.Sprintf("Academic session with id %s. Not found", academicSessionId.String())
					return apierror.NotFound(message)
				}
				return apierror.Validation("Invalid academic session")
			}
			academicSessionId = academicSession.ID
		}

		student, err := q.CreateStudent(ctx, store.CreateStudentParams{
			TenantID:           tenantId,
			AdmissionNumber:    admissionNumber,
			FirstName:          firstName,
			LastName:           lastName,
			MiddleName:         request.MiddleName,
			Gender:             gender,
			DateOfBirth:        dateOfBirth,
			Status:             string(constants.Active),
			CurrentClassArmID:  &classArmId,
			AdmissionSessionID: &academicSessionId,
		})

		if err != nil {
			return academic.TranslateAcademicError(err)
		}

		for _, p := range request.Parents {
			parentFirstName := strings.TrimSpace(p.FirstName)
			parentLastName := strings.TrimSpace(p.LastName)
			parentEmail := strings.TrimSpace(strings.ToLower(p.Email))

			if parentFirstName == "" {
				return apierror.Validation("parent_first_name is required")
			}
			if parentLastName == "" {
				return apierror.Validation("parent_last_name is required")
			}
			if parentEmail == "" {
				return apierror.Validation("parent_email is required")
			}
			if !validation.IsValidEmail(parentEmail) {
				return apierror.Validation("Invalid parent_email")
			}

			parent, err := q.CreateParent(ctx, store.CreateParentParams{
				TenantID:    tenantId,
				FirstName:   firstName,
				LastName:    lastName,
				Email:       &parentEmail,
				PhoneNumber: p.PhoneNumber,
				Address:     &p.Address,
			})
			if err != nil {
				return academic.TranslateAcademicError(err)
			}

			if err := q.LinkStudentParent(ctx, store.LinkStudentParentParams{
				StudentID:    student.ID,
				ParentID:     parent.ID,
				Relationship: p.Relationship,
				IsPrimary:    p.IsPrimary,
			}); err != nil {
				return academic.TranslateAcademicError(err)
			}
		}

		newStudent = student
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &newStudent, nil
}

func (service *Service) GetStudent(ctx context.Context, tenantId, studentId uuid.UUID) (*store.Student, error) {
	student, err := service.queries.GetStudentByID(ctx, store.GetStudentByIDParams{
		ID:       studentId,
		TenantID: tenantId,
	})

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}

	return &student, nil
}

func (service *Service) ListStudents(ctx context.Context, tenantId uuid.UUID) ([]store.Student, error) {
	students, err := service.queries.ListStudents(ctx, tenantId)
	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}

	return students, nil
}

func (service *Service) GetStudentParents(ctx context.Context, tenantId, studentId uuid.UUID) ([]store.GetParentsByStudentRow, error) {
	parents, err := service.queries.GetParentsByStudent(ctx, store.GetParentsByStudentParams{
		StudentID: studentId,
		TenantID:  tenantId,
	})

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}
	return parents, nil
}

func (service *Service) GetStudentProfile(ctx context.Context, tenantId, studentId uuid.UUID) (*store.GetStudentProfileRow, error) {

	studentProfile, err := service.queries.GetStudentProfile(ctx, store.GetStudentProfileParams{
		ID:       studentId,
		TenantID: tenantId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NotFound("Student profile not found")
		}
		return nil, academic.TranslateAcademicError(err)
	}
	return &studentProfile, nil
}

func (service *Service) SearchStudent(ctx context.Context, tenantId uuid.UUID, searchRequest SearchStudentRequest) ([]store.Student, error) {

	students, err := service.queries.SearchStudents(ctx, store.SearchStudentsParams{
		TenantID:        tenantId,
		Search:          searchRequest.Query,
		CurrentClassArm: searchRequest.ClassArmId,
		Status:          searchRequest.Status,
	})

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}
	return students, nil
}

func (service *Service) UpdateStudent(ctx context.Context, tenantId, studentId uuid.UUID, request UpdateStudentRequest) (*store.Student, error) {

	var dateOfBirth *time.Time

	if request.DateOfBirth != nil {
		time, _ := time.Parse("2026-01-30", *request.DateOfBirth)
		dateOfBirth = &time
	}

	updatedStudent, err := service.queries.UpdateStudent(ctx, store.UpdateStudentParams{
		ID:                studentId,
		TenantID:          tenantId,
		FirstName:         request.FirstName,
		LastName:          request.LastName,
		MiddleName:        request.MiddleName,
		Gender:            request.Gender,
		DateOfBirth:       dateOfBirth,
		CurrentClassArmID: request.CurrentClassArmID,
		Status:            request.Status,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NotFound("Student profile not found")
		}
		return nil, academic.TranslateAcademicError(err)
	}

	return &updatedStudent, nil
}
