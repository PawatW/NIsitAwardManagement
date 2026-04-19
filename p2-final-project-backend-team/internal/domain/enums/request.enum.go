package enums
//enum cuzthe workflow of request status is fixed 
type RequestStatus string

const (
	//under review by diff roles
	PendingHOD            RequestStatus = "PENDING_HOD"
	PendingViceDean       RequestStatus = "PENDING_VICE_DEAN"
	PendingDean           RequestStatus = "PENDING_DEAN"
	PendingCommittee      RequestStatus = "PENDING_COMMITTEE" //for the award type check, committee can chnange the award type
	Approved    		  RequestStatus = "APPROVED_PENDING_PRESIDENT"
	
	//rejects by different role
	RejectedByHOD       RequestStatus = "REJECTED_BY_HOD"
	RejectedByViceDean  RequestStatus = "REJECTED_BY_VICE_DEAN"
	RejectedByDean      RequestStatus = "REJECTED_BY_DEAN"
	RejectedByCommittee RequestStatus = "REJECTED_BY_COMMITTEE"
	
	//award cate change need re-review
	AwardChanged RequestStatus = "AWARD_CHANGED"
)
