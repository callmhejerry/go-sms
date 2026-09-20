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
	"github.com/jackc/pgx/v5/pgtype"
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
		var classArmId pgtype.UUID
		var academicSessionId pgtype.UUID

		if request.CurrentClassArm != nil {
			classArm, err := q.GetClassArmByID(ctx, store.GetClassArmByIDParams{
				ID:       pgtype.UUID{Bytes: *request.CurrentClassArm, Valid: true},
				TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
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
				ID:       pgtype.UUID{Bytes: *request.AdmissionSession, Valid: true},
				TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
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
			TenantID:           pgtype.UUID{Bytes: tenantId, Valid: true},
			AdmissionNumber:    admissionNumber,
			FirstName:          firstName,
			LastName:           lastName,
			MiddleName:         request.MiddleName,
			Gender:             gender,
			DateOfBirth:        pgtype.Date{Time: dateOfBirth, Valid: true},
			Status:             string(constants.Active),
			CurrentClassArmID:  classArmId,
			AdmissionSessionID: academicSessionId,
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
				TenantID:    pgtype.UUID{Bytes: tenantId, Valid: true},
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
		ID:       pgtype.UUID{Bytes: studentId, Valid: true},
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
	})

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}

	return &student, nil
}

func (service *Service) ListStudents(ctx context.Context, tenantId uuid.UUID) ([]store.Student, error) {
	students, err := service.queries.ListStudents(ctx, pgtype.UUID{
		Bytes: tenantId,
		Valid: true,
	})
	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}

	return students, nil
}

func (service *Service) GetStudentParents(ctx context.Context, tenantId, studentId uuid.UUID) ([]store.GetParentsByStudentRow, error) {
	parents, err := service.queries.GetParentsByStudent(ctx, store.GetParentsByStudentParams{
		StudentID: pgtype.UUID{Bytes: studentId, Valid: true},
		TenantID:  pgtype.UUID{Bytes: tenantId, Valid: true},
	})

	if err != nil {
		return nil, academic.TranslateAcademicError(err)
	}
	return parents, nil
}
