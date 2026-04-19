package request

import (
	"context"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Deprecated: use Atomic ops instead
	Create(ctx context.Context, request *models.Request) error
	Update(ctx context.Context, request *models.Request) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Atomic ops
	CreateWithHistory(ctx context.Context, request *models.Request, history *models.RequestStatusHistory) error
	UpdateWithHistory(ctx context.Context, request *models.Request, history *models.RequestStatusHistory) error
	UpdateWithCategoryChange(ctx context.Context, request *models.Request, change *models.AwardCategoryChange, history *models.RequestStatusHistory) error

	// query ops
	FindByID(ctx context.Context, id uuid.UUID) (*models.Request, error)
	FindByIDWithStudent(ctx context.Context, id uuid.UUID) (*models.Request, error)
	FindByStudentID(ctx context.Context, studentID uuid.UUID) ([]*models.Request, error)
	FindByStatus(ctx context.Context, status string) ([]*models.Request, error)
	FindByAcademicTerm(ctx context.Context, termID int) ([]*models.Request, error)

	// get request with all relations
	FindByIDWithRelations(ctx context.Context, id uuid.UUID) (*models.Request, error)

	// search with filters + pagination
	FindAllWithFilters(ctx context.Context, studentID *uuid.UUID, awardID *int, termID *int, status *string, page, limit int) ([]*models.Request, int64, error)

	// same as FindAllWithFilters but with org scope (faculty/department)
	FindAllWithFiltersAndScope(ctx context.Context,
		studentID *uuid.UUID, awardID *int, termID *int, status *string,
		allowedStatuses []string,
		facultyID *uuid.UUID, departmentID *uuid.UUID,
		page, limit int,
	) ([]*models.Request, int64, error)

	// Find approved winners for honor roll display
	FindApprovedWinners(ctx context.Context,
		academicYear *int, semester *string, campusID *uuid.UUID, awardID *int,
	) ([]*models.Request, error)

	ExistsByStudentAndTerm(ctx context.Context, studentID uuid.UUID, termID int) (bool, error)

	// status ops
	AddStatusHistory(ctx context.Context, history *models.RequestStatusHistory) error
	GetStatusHistory(ctx context.Context, requestID uuid.UUID) ([]*models.RequestStatusHistory, error)

	// doc ops
	AddDocument(ctx context.Context, document *models.RequestDocument) error
	GetDocuments(ctx context.Context, requestID uuid.UUID) ([]*models.RequestDocument, error)
	DeleteDocument(ctx context.Context, documentID uuid.UUID) error

	// award cat change ops
	AddCategoryChange(ctx context.Context, change *models.AwardCategoryChange) error
	GetCategoryChanges(ctx context.Context, requestID uuid.UUID) ([]*models.AwardCategoryChange, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, request *models.Request) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r *repository) Update(ctx context.Context, request *models.Request) error {
	return r.db.WithContext(ctx).Save(request).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Request{}, id).Error
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Request, error) {
	var request models.Request
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&request).Error
	if err != nil {
		return nil, HandleError(err)
	}
	return &request, nil
}

// preload Student for scope check (approve/reject)
func (r *repository) FindByIDWithStudent(ctx context.Context, id uuid.UUID) (*models.Request, error) {
	var request models.Request
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Campus").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Where("id = ?", id).
		First(&request).Error
	if err != nil {
		return nil, HandleError(err)
	}
	return &request, nil
}

func (r *repository) FindByIDWithRelations(ctx context.Context, id uuid.UUID) (*models.Request, error) {
	var request models.Request
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Campus").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Preload("Award").
		Preload("AcademicTerm").
		Preload("StatusHistory").
		Preload("Documents").
		Preload("CategoryChanges").
		Where("id = ?", id).
		First(&request).Error
	if err != nil {
		return nil, HandleError(err)
	}
	return &request, nil
}

func (r *repository) FindByStudentID(ctx context.Context, studentID uuid.UUID) ([]*models.Request, error) {
	var requests []*models.Request
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Campus").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Preload("Award").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&requests).Error
	return requests, err
}

func (r *repository) FindByStatus(ctx context.Context, status string) ([]*models.Request, error) {
	var requests []*models.Request
	err := r.db.WithContext(ctx).
		Where("current_status = ?", status).
		Order("created_at DESC").
		Find(&requests).Error
	return requests, err
}

func (r *repository) FindByAcademicTerm(ctx context.Context, termID int) ([]*models.Request, error) {
	var requests []*models.Request
	err := r.db.WithContext(ctx).
		Where("academic_term_id = ?", termID).
		Order("created_at DESC").
		Find(&requests).Error
	return requests, err
}

func (r *repository) ExistsByStudentAndTerm(ctx context.Context, studentID uuid.UUID, termID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Request{}).
		Where("student_id = ? AND academic_term_id = ?", studentID, termID).
		Count(&count).Error
	return count > 0, err
}

// เดิม: ไม่มี org scope
func (r *repository) FindAllWithFilters(ctx context.Context, studentID *uuid.UUID, awardID *int, termID *int, status *string, page, limit int) ([]*models.Request, int64, error) {
	return r.FindAllWithFiltersAndScope(ctx, studentID, awardID, termID, status, nil, nil, nil, page, limit)
}

// มี org scope โดย join users (student)
func (r *repository) FindAllWithFiltersAndScope(
	ctx context.Context,
	studentID *uuid.UUID, awardID *int, termID *int, status *string,
	allowedStatuses []string,
	facultyID *uuid.UUID, departmentID *uuid.UUID,
	page, limit int,
) ([]*models.Request, int64, error) {
	var requests []*models.Request
	var total int64

	// base
	query := r.db.WithContext(ctx).Model(&models.Request{})

	// join users for org scope (student)
	// assumes: users.user_id, users.faculty_id, users.department_id
	if facultyID != nil || departmentID != nil {
		query = query.Joins("JOIN users u ON u.user_id = requests.student_id")
		if facultyID != nil {
			query = query.Where("u.faculty_id = ?", *facultyID)
		}
		if departmentID != nil {
			query = query.Where("u.department_id = ?", *departmentID)
		}
	}

	if studentID != nil {
		query = query.Where("requests.student_id = ?", *studentID)
	}
	if awardID != nil {
		query = query.Where("requests.award_id = ?", *awardID)
	}
	if termID != nil {
		query = query.Where("requests.academic_term_id = ?", *termID)
	}
	if status != nil {
		query = query.Where("requests.current_status = ?", *status)
	} else if len(allowedStatuses) > 0 {
		query = query.Where("requests.current_status IN ?", allowedStatuses)
	}

	// count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.
		Preload("Student").
		Preload("Student.Campus").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Preload("Award").
		Preload("AcademicTerm").
		Order("requests.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&requests).Error

	return requests, total, err
}

func (r *repository) AddStatusHistory(ctx context.Context, history *models.RequestStatusHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *repository) GetStatusHistory(ctx context.Context, requestID uuid.UUID) ([]*models.RequestStatusHistory, error) {
	var history []*models.RequestStatusHistory
	err := r.db.WithContext(ctx).
		Where("request_id = ?", requestID).
		Order("updated_at ASC").
		Find(&history).Error
	return history, err
}

func (r *repository) AddDocument(ctx context.Context, document *models.RequestDocument) error {
	return r.db.WithContext(ctx).Create(document).Error
}

func (r *repository) GetDocuments(ctx context.Context, requestID uuid.UUID) ([]*models.RequestDocument, error) {
	var documents []*models.RequestDocument
	err := r.db.WithContext(ctx).
		Where("request_id = ?", requestID).
		Find(&documents).Error
	return documents, err
}

func (r *repository) DeleteDocument(ctx context.Context, documentID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.RequestDocument{}, documentID).Error
}

func (r *repository) AddCategoryChange(ctx context.Context, change *models.AwardCategoryChange) error {
	return r.db.WithContext(ctx).Create(change).Error
}

func (r *repository) GetCategoryChanges(ctx context.Context, requestID uuid.UUID) ([]*models.AwardCategoryChange, error) {
	var changes []*models.AwardCategoryChange
	err := r.db.WithContext(ctx).
		Where("request_id = ?", requestID).
		Order("changed_at ASC").
		Find(&changes).Error
	return changes, err
}

func (r *repository) CreateWithHistory(ctx context.Context, request *models.Request, history *models.RequestStatusHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(request).Error; err != nil {
			return err
		}
		if err := tx.Create(history).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *repository) UpdateWithCategoryChange(ctx context.Context, request *models.Request, change *models.AwardCategoryChange, history *models.RequestStatusHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(request).Error; err != nil {
			return err
		}
		if err := tx.Create(change).Error; err != nil {
			return err
		}
		if err := tx.Create(history).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *repository) UpdateWithHistory(ctx context.Context, request *models.Request, history *models.RequestStatusHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(request).Error; err != nil {
			return err
		}
		if err := tx.Create(history).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindApprovedWinners returns all approved requests for honor roll display
// Filters by APPROVED_PENDING_PRESIDENT status and optional parameters
func (r *repository) FindApprovedWinners(
	ctx context.Context,
	academicYear *int,
	semester *string,
	campusID *uuid.UUID,
	awardID *int,
) ([]*models.Request, error) {
	var requests []*models.Request

	// Start with base query filtering approved status
	query := r.db.WithContext(ctx).Model(&models.Request{}).
		Where("requests.current_status = ?", "APPROVED_PENDING_PRESIDENT")

	// Join academic_terms for year/semester filtering
	if academicYear != nil || semester != nil {
		query = query.Joins("JOIN academic_terms at ON at.id = requests.academic_term_id")
		if academicYear != nil {
			query = query.Where("at.academic_year = ?", *academicYear)
		}
		if semester != nil {
			query = query.Where("at.semester = ?", *semester)
		}
	}

	// Join users table for campus filtering
	if campusID != nil {
		query = query.Joins("JOIN users u ON u.user_id = requests.student_id").
			Where("u.campus_id = ?", *campusID)
	}

	// Filter by award category
	if awardID != nil {
		query = query.Where("requests.award_id = ?", *awardID)
	}

	// Execute with preloads and ordering
	err := query.
		Preload("Student").
		Preload("Student.Campus").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Preload("Award").
		Preload("AcademicTerm").
		Order("requests.academic_term_id DESC, requests.award_id ASC, requests.snapshot_gpa DESC").
		Find(&requests).Error

	return requests, err
}
