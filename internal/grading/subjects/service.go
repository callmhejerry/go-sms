package subjects

import (
	"context"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/google/uuid"
)

type SubjectService struct {
	subjectRepository SubjectRepository
}

func NewSubjectService(subjectRepository SubjectRepository) *SubjectService {
	return &SubjectService{
		subjectRepository: subjectRepository,
	}
}

func (service *SubjectService) CreateSubject(
	ctx context.Context,
	tenantId uuid.UUID,
	request CreateSubjectRequest,
) (*store.Subject, *apierror.AppError) {
	name := strings.TrimSpace(strings.ToLower(request.Name))
	code := strings.TrimSpace(strings.ToLower(request.Code))

	if name == "" {
		return nil, apierror.Validation("Name cannot be empty")
	}
	if code == "" {
		return nil, apierror.Validation("Code cannot be empty")
	}

	return service.subjectRepository.CreateSubject(ctx, tenantId, request.Name, request.Code)
}

func (service *SubjectService) ListSubjects(
	ctx context.Context,
	tenantId uuid.UUID,
) ([]store.Subject, *apierror.AppError) {
	return service.subjectRepository.ListSubjects(ctx, tenantId)
}

func (service *SubjectService) AddSubjectToClass(
	ctx context.Context,
	tenantId, subjectId, classId uuid.UUID,
) (*store.ClassSubject, *apierror.AppError) {
	return service.subjectRepository.AddSubjectToClass(ctx, tenantId, classId, subjectId)
}

func (service *SubjectService) AssignTeacherToSubject(
	ctx context.Context,
	tenantId uuid.UUID,
	request AssignTeacherRequest,
) (*store.TeacherAssignment, *apierror.AppError) {
	return service.subjectRepository.AssignTeacherToSubject(
		ctx, tenantId,
		request.UserID,
		request.SubjectID,
		request.AcademicSessionID,
		request.ClassID,
		request.ClassArmID,
	)
}
