package academic

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func (service *Service) CreateAcademicSession(ctx context.Context, tenantId uuid.UUID, request CreateAcademicSessionRequest) (*store.AcademicSession, error) {

	name := strings.TrimSpace(request.Name)

	if name == "" {
		return nil, apierror.Validation("name is required")
	}
	if request.StartDate.IsZero() || request.EndDate.IsZero() {
		return nil, apierror.Validation("start_date and end_date is required")
	}

	if request.EndDate.Before(request.StartDate) {
		return nil, apierror.Validation("end_date cannot be before start_date")
	}
	if request.IsCurrent {
		_ = service.queries.SetCurrentAcademicSession(ctx, pgtype.UUID{
			Bytes: tenantId,
			Valid: true,
		})
	}

	session, err := service.queries.CreateAcademicSession(ctx, store.CreateAcademicSessionParams{
		TenantID: pgtype.UUID{
			Bytes: tenantId,
			Valid: true,
		},
		Name: name,
		StartDate: pgtype.Date{
			Time:  request.StartDate,
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  request.EndDate,
			Valid: true,
		},
		IsCurrent: request.IsCurrent,
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return &session, nil
}

func (service *Service) ListAcademicSession(ctx context.Context, tenantId uuid.UUID) ([]store.AcademicSession, error) {
	sessions, err := service.queries.ListAcademicSessions(ctx, pgtype.UUID{
		Bytes: tenantId,
		Valid: true,
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return sessions, nil
}

func (service *Service) GetCurrentAcademicSession(ctx context.Context, tenantId uuid.UUID) (*store.AcademicSession, error) {
	currentSession, err := service.queries.GetCurrentAcademicSession(ctx, pgtype.UUID{
		Bytes: tenantId,
		Valid: true,
	})
	if err != nil {
		return nil, translateAcademicError(err)
	}
	return &currentSession, nil
}

func (service *Service) GetAcademicSessionById(ctx context.Context, tenantId, academicSessionId uuid.UUID) (*store.AcademicSession, error) {
	academicSession, err := service.queries.GetAcademicSessionById(ctx, store.GetAcademicSessionByIdParams{
		ID:       pgtype.UUID{Bytes: academicSessionId, Valid: true},
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			message := fmt.Sprintf("Academic session with id %s. Not found", academicSessionId.String())
			return nil, apierror.NotFound(message)
		}
		return nil, translateAcademicError(err)
	}

	return &academicSession, nil
}

func (service *Service) CreateClass(ctx context.Context, tenantId uuid.UUID, request CreateClassRequest) (*store.Class, error) {
	name := strings.TrimSpace(strings.ToLower(request.Name))

	if name == "" {
		return nil, apierror.Validation("class name is required")
	}

	class, err := service.queries.CreateClass(ctx, store.CreateClassParams{
		TenantID:   pgtype.UUID{Bytes: tenantId, Valid: true},
		Name:       name,
		LevelOrder: int32(request.LevelOrder),
	})
	if err != nil {
		return nil, translateAcademicError(err)
	}
	return &class, nil
}

func (service *Service) GetClassById(ctx context.Context, tenantId, classId uuid.UUID) (*store.Class, error) {
	class, err := service.queries.GetClassByID(ctx, store.GetClassByIDParams{
		ID:       pgtype.UUID{Bytes: classId, Valid: true},
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			message := fmt.Sprintf("class with id %s. Not found", classId.String())
			return nil, apierror.NotFound(message)
		}
		return nil, translateAcademicError(err)
	}
	return &class, nil
}

func (service *Service) ListClasses(ctx context.Context, tenantId uuid.UUID) ([]store.Class, error) {
	classes, err := service.queries.ListClasses(ctx, pgtype.UUID{Bytes: tenantId, Valid: true})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return classes, nil
}

func (service *Service) CreateClassArm(ctx context.Context, tenantId uuid.UUID, request CreateClassArmRequest) (*store.ClassArm, error) {
	name := strings.TrimSpace(strings.ToLower(request.Name))

	if name == "" {
		return nil, apierror.Validation("Class arm name is required")
	}

	class, err := service.queries.GetClassByID(ctx, store.GetClassByIDParams{
		ID:       pgtype.UUID{Bytes: request.ClassID, Valid: true},
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}

	classArm, err := service.queries.CreateClassArm(ctx, store.CreateClassArmParams{
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
		ClassID:  class.ID,
		Name:     name,
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return &classArm, nil
}

func (service *Service) ListClassArms(ctx context.Context, tenantId uuid.UUID, classId uuid.UUID) ([]store.ClassArm, error) {
	classArms, err := service.queries.ListClassArms(ctx, store.ListClassArmsParams{
		TenantID: pgtype.UUID{Bytes: tenantId, Valid: true},
		ClassID:  pgtype.UUID{Bytes: classId, Valid: true},
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return classArms, nil
}

func (service *Service) CreateStudent(ctx context.Context, tenantId uuid.UUID, request CreateStudentRequest) (*store.Student, error) {
	admissionNumber := strings.TrimSpace(request.AdmissionNumber)
	firstName := strings.TrimSpace(request.FirstName)
	lastName := strings.TrimSpace(request.LastName)
	gender := strings.TrimSpace(strings.ToLower(request.Gender))

	if admissionNumber == "" {
		return nil, apierror.Validation("admission_number is required")
	}
	if firstName == "" {
		return nil, apierror.Validation("first_name is required")
	}
	if lastName == "" {
		return nil, apierror.Validation("last_name is required")
	}
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
			DateOfBirth:        pgtype.Date{Time: request.DateOfBirth, Valid: true},
			Status:             string(constants.Active),
			CurrentClassArmID:  classArmId,
			AdmissionSessionID: academicSessionId,
		})

		if err != nil {
			return translateAcademicError(err)
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
				return translateAcademicError(err)
			}

			if err := q.LinkStudentParent(ctx, store.LinkStudentParentParams{
				StudentID:    student.ID,
				ParentID:     parent.ID,
				Relationship: p.Relationship,
				IsPrimary:    p.IsPrimary,
			}); err != nil {
				return translateAcademicError(err)
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
		return nil, translateAcademicError(err)
	}

	return &student, nil
}

func (service *Service) ListStudents(ctx context.Context, tenantId uuid.UUID) ([]store.Student, error) {
	students, err := service.queries.ListStudents(ctx, pgtype.UUID{
		Bytes: tenantId,
		Valid: true,
	})
	if err != nil {
		return nil, translateAcademicError(err)
	}

	return students, nil
}

func (service *Service) GetStudentParents(ctx context.Context, tenantId, studentId uuid.UUID) ([]store.GetParentsByStudentRow, error) {
	parents, err := service.queries.GetParentsByStudent(ctx, store.GetParentsByStudentParams{
		StudentID: pgtype.UUID{Bytes: studentId, Valid: true},
		TenantID:  pgtype.UUID{Bytes: tenantId, Valid: true},
	})

	if err != nil {
		return nil, translateAcademicError(err)
	}
	return parents, nil
}
