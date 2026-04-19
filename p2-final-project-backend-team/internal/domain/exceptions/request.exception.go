package exceptions

import "errors"

var (
	// Auth/context errors
	ErrUnauthenticated     = errors.New("unauthenticated")
	ErrInvalidUserIDType   = errors.New("invalid user id type")
	ErrRoleNotFound        = errors.New("role not found")
	ErrInvalidUserRoleType = errors.New("invalid user role type")
	ErrViewerNotFound      = errors.New("viewer not found")

	// Domain/business errors
	ErrRequestNotFound                  = errors.New("request not found")
	ErrForbidden                        = errors.New("forbidden")
	ErrDuplicateRequestInTerm           = errors.New("you have already applied for this academic term")
	ErrAcademicTermNotOpen              = errors.New("academic term is not open for applications")
	ErrAcademicTermNotFound             = errors.New("academic term not found")
	ErrRoleNotAllowedToApprove          = errors.New("role is not authorized to approve requests")
	ErrRoleNotAllowedToReject           = errors.New("role is not authorized to reject requests")
	ErrInvalidRejectionRole             = errors.New("invalid role for rejection")
	ErrApprovalStepNotConfigured        = errors.New("approval step not configured for this role")
	ErrInvalidApprovalTransition        = errors.New("request is not in correct status for approval")
	ErrInvalidRejectionTransition       = errors.New("request is not in correct status for rejection")
	ErrInvalidAwardCategoryChangeStatus = errors.New("request is not in correct status for award category change")

	// Input/document errors
	ErrInvalidRequestID       = errors.New("invalid request id")
	ErrInvalidRequestBody     = errors.New("invalid request body")
	ErrInvalidQueryParameters = errors.New("invalid query parameters")
	ErrFileRequired           = errors.New("file is required")
	ErrFileTooLarge           = errors.New("file size exceeds maximum allowed size of 10 MB")
	ErrInvalidFileType        = errors.New("invalid file type. only .png, .jpeg, .jpg, and .pdf files are allowed")
	ErrFailedToOpenFile       = errors.New("failed to open file")
	ErrFailedToUploadFile     = errors.New("failed to upload file to storage")
)
