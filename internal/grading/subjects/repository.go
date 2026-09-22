package subjects

import (
	"context"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type SubjectRepository interface {
	CreateSubject(
		ctx context.Context,
		tenantId uuid.UUID,
		name, code string,
	) (*store.Subject, *apierror.AppError)

	ListSubjects(
		ctx context.Context,
		tenantId uuid.UUID,
	) ([]store.Subject, *apierror.AppError)

	AddSubjectToClass(
		ctx context.Context,
		tenantId, classId, subjectId uuid.UUID,
	) (*store.ClassSubject, *apierror.AppError)

	AssignTeacherToSubject(
		ctx context.Context,
		tenantId, userId, subjectId, academicSessionId, classId uuid.UUID,
		class_arm_id *uuid.UUID,
	) (*store.TeacherAssignment, *apierror.AppError)
}

type SubjectRepositoryImpl struct {
	queries *store.Queries
}

func NewSubjectRepository(queries *store.Queries) *SubjectRepositoryImpl {
	return &SubjectRepositoryImpl{
		queries: queries,
	}
}

func (repo *SubjectRepositoryImpl) CreateSubject(
	ctx context.Context,
	tenantId uuid.UUID,
	name, code string,
) (*store.Subject, *apierror.AppError) {
	subject, err := repo.queries.CreateSubject(
		ctx, store.CreateSubjectParams{
			TenantID: tenantId,
			Name:     name,
			Code:     code,
		},
	)

	if err != nil {
		return nil, translateSubjectError(err)
	}
	return &subject, nil
}

func (repo *SubjectRepositoryImpl) ListSubjects(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.Subject, *apierror.AppError) {
	subject, err := repo.queries.ListSubjects(ctx, tenantId)

	if err != nil {
		return nil, translateSubjectError(err)
	}
	return subject, nil
}

func (repo *SubjectRepositoryImpl) AddSubjectToClass(
	ctx context.Context,
	tenantId, classId, subjectId uuid.UUID,
) (*store.ClassSubject, *apierror.AppError) {
	classSubject, err := repo.queries.AddSubjectToClass(
		ctx, store.AddSubjectToClassParams{
			TenantID:  tenantId,
			ClassID:   classId,
			SubjectID: subjectId,
		},
	)
	if err != nil {
		return nil, translateSubjectError(err)
	}
	return &classSubject, nil
}

func (repo *SubjectRepositoryImpl) AssignTeacherToSubject(
	ctx context.Context,
	tenantId, userId, subjectId, academicSessionId, classId uuid.UUID,
	class_arm_id *uuid.UUID,
) (*store.TeacherAssignment, *apierror.AppError) {
	teacherAssignment, err := repo.queries.AssignTeacher(
		ctx, store.AssignTeacherParams{
			TenantID:          tenantId,
			UserID:            userId,
			SubjectID:         subjectId,
			AcademicSessionID: academicSessionId,
			ClassID:           classId,
			ClassArmID:        class_arm_id,
		},
	)

	if err != nil {
		return nil, translateSubjectError(err)
	}
	return &teacherAssignment, nil
}
