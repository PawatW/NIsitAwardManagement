package request

import (
	"testing"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
)

func TestCanRoleViewStatus_AwardChanged(t *testing.T) {
	if !canRoleViewStatus(enums.Dean, string(enums.AwardChanged)) {
		t.Fatalf("dean must be able to view AWARD_CHANGED requests")
	}

	if canRoleViewStatus(enums.HeadOfDepartment, string(enums.AwardChanged)) {
		t.Fatalf("HOD must not be able to view AWARD_CHANGED requests")
	}

	if canRoleViewStatus(enums.ViceDean, string(enums.AwardChanged)) {
		t.Fatalf("ViceDean must not be able to view AWARD_CHANGED requests")
	}

	if canRoleViewStatus(enums.CommitteeChair, string(enums.AwardChanged)) {
		t.Fatalf("CommitteeChair must not be able to view AWARD_CHANGED requests")
	}
}

func TestApprovalStepForRole_DeanIncludesAwardChanged(t *testing.T) {
	step, err := approvalStepForRole(enums.Dean)
	if err != nil {
		t.Fatalf("approvalStepForRole(dean) error = %v", err)
	}

	found := false
	for _, status := range step.validFrom {
		if status == enums.AwardChanged {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("dean approval step must allow AWARD_CHANGED")
	}
}

func TestApprovalStepForRole_HODStartsFromPendingHODOnly(t *testing.T) {
	step, err := approvalStepForRole(enums.HeadOfDepartment)
	if err != nil {
		t.Fatalf("approvalStepForRole(hod) error = %v", err)
	}

	if len(step.validFrom) != 1 {
		t.Fatalf("expected exactly 1 validFrom status, got %d", len(step.validFrom))
	}
	if step.validFrom[0] != enums.PendingHOD {
		t.Fatalf("expected validFrom[0] = %s, got %s", enums.PendingHOD, step.validFrom[0])
	}
}

func TestApprovalStepForRole_CommitteeNextIsApproved(t *testing.T) {
	step, err := approvalStepForRole(enums.CommitteeChair)
	if err != nil {
		t.Fatalf("approvalStepForRole(committee) error = %v", err)
	}

	if step.nextStatus != enums.Approved {
		t.Fatalf("expected committee next status = %s, got %s", enums.Approved, step.nextStatus)
	}
}

func TestCanRoleViewStatus_CommitteeChairCanViewApproved(t *testing.T) {
	if !canRoleViewStatus(enums.CommitteeChair, string(enums.Approved)) {
		t.Fatalf("CommitteeChair must be able to view approved status")
	}
}

func TestCanRoleViewStatus_HODCannotViewSubmittedLiteral(t *testing.T) {
	if canRoleViewStatus(enums.HeadOfDepartment, "SUBMITTED") {
		t.Fatalf("HOD must not be able to view SUBMITTED after enum cleanup")
	}
}

func TestAllowedStatusesForRole_AdminIsUnrestricted(t *testing.T) {
	statuses := allowedStatusesForRole(enums.Admin)
	if len(statuses) != 0 {
		t.Fatalf("admin should be unrestricted, got %v", statuses)
	}
}
