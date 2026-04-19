package requests

import "github.com/google/uuid"
import "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"


//create request dto (student submit)
type CreateRequestDTO struct {
	StudentID        uuid.UUID `json:"student_id" validate:"required,uuid"`
	AwardID          int       `json:"award_id" validate:"required"`
	AcademicTermID   int       `json:"academic_term_id" validate:"required"`
	AdditionalData   models.AdditionalData  `json:"additional_data" validate:"required"` 

	
	//snapshot req at the tiem of submit
	SnapshotStudyYear   int     `json:"snapshot_study_year" validate:"required,min=1,max=8"`
	SnapshotGPA         float64 `json:"snapshot_gpa" validate:"required,min=0,max=4"`
	SnapshotAdvisor     string  `json:"snapshot_advisor" validate:"required"`
	SnapshotDateOfBirth *string `json:"snapshot_date_of_birth,omitempty"` 
	SnapshotPhone       *string `json:"snapshot_phone,omitempty"`        
	SnapshotAddress     *string `json:"snapshot_address,omitempty"`      
}

//approve req for HOD, Vice Dean, Dean, Committee Chair 
type ApproveRequestDTO struct {
	Remark string `json:"remark"` //optional
}

//reject req for HOD, Vice Dean, Dean, Committee 
type RejectRequestDTO struct {
	Remark string `json:"remark" validate:"required"` // require when rejecting
}

//verify req for Student Development Office to verify
type VerifyRequestDTO struct {
	IsApproved bool   `json:"is_approved" validate:"required"`
	Remark     string `json:"remark"`
}

//award category change for กองพัฒนานิสิต hence the student developmetn division
type ChangeAwardCategoryDTO struct {
	NewCategoryID int    `json:"new_category_id" validate:"required"`
	AdditionalData models.AdditionalData `json:"additional_data" validate:"required"` 
	Remark        string `json:"remark" validate:"required"`
}



//upload doc to req
type UploadDocumentDTO struct {
	Description string `json:"description" validate:"required"`
}

//req query dto for filtering/searching requests
type RequestQueryDTO struct {
	StudentID      string `query:"student_id"`
	AwardID        int    `query:"award_id"`
	AcademicTermID int    `query:"academic_term_id"`
	CurrentStatus  string `query:"current_status"`
	Page           int    `query:"page" validate:"min=1"`
	Limit          int    `query:"limit" validate:"min=1,max=100"`
}

//p.s viotign is done manually by the comittee members outside the system