package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/configs"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {
	config := configs.NewConfig()
	db := database.NewPostgrest(config)

	log.Info().Msg("seeding database")

	clearData(db)

	campuses := createCampuses(db)
	log.Info().Msgf("Created %d campuses", len(campuses))

	faculties := createFaculties(db, campuses)
	log.Info().Msgf("Created %d faculties", len(faculties))

	departments := createDepartments(db, faculties)
	log.Info().Msgf("Created %d departments", len(departments))

	awards := createAwardCategories(db)
	log.Info().Msgf("Created %d award categories", len(awards))

	terms := createAcademicTerms(db)
	log.Info().Msgf("Created %d academic terms", len(terms))

	users := createUsers(db, campuses, faculties, departments)
	log.Info().Msgf("Created %d users", len(users))

	requests := createRequests(db, users, awards, terms)
	log.Info().Msgf("Created %d requests", len(requests))

	log.Info().Msg("seed completed")
}

func clearData(db *gorm.DB) {
	db.Exec("DELETE FROM award_category_changes")
	db.Exec("DELETE FROM request_documents")
	db.Exec("DELETE FROM request_status_history")
	db.Exec("DELETE FROM requests")
	db.Exec("DELETE FROM academic_terms")
	db.Exec("DELETE FROM award_categories")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM departments")
	db.Exec("DELETE FROM faculties")
	db.Exec("DELETE FROM campus")
}

// ──────────────────────────────────────────────
// Campuses – วิทยาเขตจริงของมหาวิทยาลัยเกษตรศาสตร์
// ──────────────────────────────────────────────

func createCampuses(db *gorm.DB) []models.Campus {
	campuses := []models.Campus{
		{ID: uuid.New(), Name: "บางเขน"},
		{ID: uuid.New(), Name: "กำแพงแสน"},
		{ID: uuid.New(), Name: "ศรีราชา"},
		{ID: uuid.New(), Name: "เฉลิมพระเกียรติ จังหวัดสกลนคร"},
	}
	for _, c := range campuses {
		db.Create(&c)
	}
	return campuses
}

// ──────────────────────────────────────────────
// Faculties – คณะจริงของ มก. แบ่งตามวิทยาเขต
// ──────────────────────────────────────────────

func createFaculties(db *gorm.DB, campuses []models.Campus) []models.Faculty {
	faculties := []models.Faculty{
		// บางเขน (index 0-9)
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะวิศวกรรมศาสตร์"}, // 0
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะวิทยาศาสตร์"},    // 1
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะเกษตร"},          // 2
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะมนุษยศาสตร์"},    // 3
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะสังคมศาสตร์"},    // 4
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะบริหารธุรกิจ"},   // 5
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะเศรษฐศาสตร์"},    // 6
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะสัตวแพทยศาสตร์"}, // 7
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะศึกษาศาสตร์"},    // 8
		{ID: uuid.New(), CampusID: campuses[0].ID, Name: "คณะประมง"},          // 9

		// กำแพงแสน (index 10-12)
		{ID: uuid.New(), CampusID: campuses[1].ID, Name: "คณะวิศวกรรมศาสตร์ กำแพงแสน"},  // 10
		{ID: uuid.New(), CampusID: campuses[1].ID, Name: "คณะศิลปศาสตร์และวิทยาศาสตร์"}, // 11
		{ID: uuid.New(), CampusID: campuses[1].ID, Name: "คณะศึกษาศาสตร์และพัฒนศาสตร์"}, // 12

		// ศรีราชา (index 13-15)
		{ID: uuid.New(), CampusID: campuses[2].ID, Name: "คณะวิทยาการจัดการ"},        // 13
		{ID: uuid.New(), CampusID: campuses[2].ID, Name: "คณะวิศวกรรมศาสตร์ศรีราชา"}, // 14
		{ID: uuid.New(), CampusID: campuses[2].ID, Name: "คณะเศรษฐศาสตร์ ศรีราชา"},   // 15

		// สกลนคร (index 16-17)
		{ID: uuid.New(), CampusID: campuses[3].ID, Name: "คณะทรัพยากรธรรมชาติและอุตสาหกรรมเกษตร"}, // 16
		{ID: uuid.New(), CampusID: campuses[3].ID, Name: "คณะวิทยาศาสตร์และวิศวกรรมศาสตร์"},       // 17
	}
	for _, f := range faculties {
		db.Create(&f)
	}
	return faculties
}

// ──────────────────────────────────────────────
// Departments – ภาควิชาจริงของแต่ละคณะ
// ──────────────────────────────────────────────

func createDepartments(db *gorm.DB, faculties []models.Faculty) []models.Department {
	departments := []models.Department{
		// คณะวิศวกรรมศาสตร์ บางเขน (fac 0)
		{ID: uuid.New(), FacultyID: faculties[0].ID, Name: "ภาควิชาวิศวกรรมคอมพิวเตอร์"}, // 0
		{ID: uuid.New(), FacultyID: faculties[0].ID, Name: "ภาควิชาวิศวกรรมไฟฟ้า"},       // 1
		{ID: uuid.New(), FacultyID: faculties[0].ID, Name: "ภาควิชาวิศวกรรมเครื่องกล"},   // 2
		{ID: uuid.New(), FacultyID: faculties[0].ID, Name: "ภาควิชาวิศวกรรมอุตสาหการ"},   // 3
		{ID: uuid.New(), FacultyID: faculties[0].ID, Name: "ภาควิชาวิศวกรรมเคมี"},        // 4
		{ID: uuid.New(), FacultyID: faculties[0].ID, Name: "ภาควิชาวิศวกรรมโยธา"},        // 5

		// คณะวิทยาศาสตร์ (fac 1)
		{ID: uuid.New(), FacultyID: faculties[1].ID, Name: "ภาควิชาวิทยาการคอมพิวเตอร์"}, // 6
		{ID: uuid.New(), FacultyID: faculties[1].ID, Name: "ภาควิชาคณิตศาสตร์"},          // 7
		{ID: uuid.New(), FacultyID: faculties[1].ID, Name: "ภาควิชาฟิสิกส์"},             // 8
		{ID: uuid.New(), FacultyID: faculties[1].ID, Name: "ภาควิชาเคมี"},                // 9
		{ID: uuid.New(), FacultyID: faculties[1].ID, Name: "ภาควิชาสถิติ"},               // 10
		{ID: uuid.New(), FacultyID: faculties[1].ID, Name: "ภาควิชาชีววิทยา"},            // 11

		// คณะเกษตร (fac 2)
		{ID: uuid.New(), FacultyID: faculties[2].ID, Name: "ภาควิชาพืชไร่นา"},  // 12
		{ID: uuid.New(), FacultyID: faculties[2].ID, Name: "ภาควิชาสัตวบาล"},   // 13
		{ID: uuid.New(), FacultyID: faculties[2].ID, Name: "ภาควิชากีฏวิทยา"},  // 14
		{ID: uuid.New(), FacultyID: faculties[2].ID, Name: "ภาควิชาปฐพีวิทยา"}, // 15

		// คณะมนุษยศาสตร์ (fac 3)
		{ID: uuid.New(), FacultyID: faculties[3].ID, Name: "ภาควิชาภาษาอังกฤษ"},     // 16
		{ID: uuid.New(), FacultyID: faculties[3].ID, Name: "ภาควิชาภาษาไทย"},        // 17
		{ID: uuid.New(), FacultyID: faculties[3].ID, Name: "ภาควิชาปรัชญาและศาสนา"}, // 18
		{ID: uuid.New(), FacultyID: faculties[3].ID, Name: "ภาควิชาภาษาตะวันออก"},   // 19

		// คณะสังคมศาสตร์ (fac 4)
		{ID: uuid.New(), FacultyID: faculties[4].ID, Name: "ภาควิชารัฐศาสตร์และรัฐประศาสนศาสตร์"}, // 20
		{ID: uuid.New(), FacultyID: faculties[4].ID, Name: "ภาควิชาสังคมวิทยาและมานุษยวิทยา"},     // 21
		{ID: uuid.New(), FacultyID: faculties[4].ID, Name: "ภาควิชาจิตวิทยา"},                     // 22
		{ID: uuid.New(), FacultyID: faculties[4].ID, Name: "ภาควิชานิติศาสตร์"},                   // 23

		// คณะบริหารธุรกิจ (fac 5)
		{ID: uuid.New(), FacultyID: faculties[5].ID, Name: "ภาควิชาการจัดการ"}, // 24
		{ID: uuid.New(), FacultyID: faculties[5].ID, Name: "ภาควิชาการตลาด"},   // 25
		{ID: uuid.New(), FacultyID: faculties[5].ID, Name: "ภาควิชาการเงิน"},   // 26
		{ID: uuid.New(), FacultyID: faculties[5].ID, Name: "ภาควิชาการบัญชี"},  // 27

		// คณะเศรษฐศาสตร์ (fac 6)
		{ID: uuid.New(), FacultyID: faculties[6].ID, Name: "ภาควิชาเศรษฐศาสตร์"},                 // 28
		{ID: uuid.New(), FacultyID: faculties[6].ID, Name: "ภาควิชาเศรษฐศาสตร์เกษตรและทรัพยากร"}, // 29
		{ID: uuid.New(), FacultyID: faculties[6].ID, Name: "ภาควิชาสหกรณ์"},                      // 30

		// คณะสัตวแพทยศาสตร์ (fac 7)
		{ID: uuid.New(), FacultyID: faculties[7].ID, Name: "ภาควิชาเวชศาสตร์คลินิกสัตว์เลี้ยง"}, // 31
		{ID: uuid.New(), FacultyID: faculties[7].ID, Name: "ภาควิชาพยาธิวิทยา"},                 // 32
		{ID: uuid.New(), FacultyID: faculties[7].ID, Name: "ภาควิชาสรีรวิทยา"},                  // 33

		// คณะศึกษาศาสตร์ (fac 8)
		{ID: uuid.New(), FacultyID: faculties[8].ID, Name: "ภาควิชาพลศึกษา"},           // 34
		{ID: uuid.New(), FacultyID: faculties[8].ID, Name: "ภาควิชาเทคโนโลยีการศึกษา"}, // 35

		// คณะประมง (fac 9)
		{ID: uuid.New(), FacultyID: faculties[9].ID, Name: "ภาควิชาเพาะเลี้ยงสัตว์น้ำ"}, // 36
		{ID: uuid.New(), FacultyID: faculties[9].ID, Name: "ภาควิชาวิทยาศาสตร์ทางทะเล"}, // 37

		// คณะวิศวกรรมศาสตร์ กำแพงแสน (fac 10)
		{ID: uuid.New(), FacultyID: faculties[10].ID, Name: "สาขาวิศวกรรมเกษตร"}, // 38
		{ID: uuid.New(), FacultyID: faculties[10].ID, Name: "สาขาวิศวกรรมอาหาร"}, // 39

		// คณะศิลปศาสตร์และวิทยาศาสตร์ (fac 11)
		{ID: uuid.New(), FacultyID: faculties[11].ID, Name: "สาขาวิทยาศาสตร์ชีวภาพ"},            // 40
		{ID: uuid.New(), FacultyID: faculties[11].ID, Name: "สาขาการจัดการโรงแรมและท่องเที่ยว"}, // 41

		// คณะศึกษาศาสตร์และพัฒนศาสตร์ (fac 12)
		{ID: uuid.New(), FacultyID: faculties[12].ID, Name: "สาขาการศึกษา"}, // 42

		// คณะวิทยาการจัดการ (fac 13)
		{ID: uuid.New(), FacultyID: faculties[13].ID, Name: "สาขาการจัดการธุรกิจ"},     // 43
		{ID: uuid.New(), FacultyID: faculties[13].ID, Name: "สาขาการจัดการโลจิสติกส์"}, // 44

		// คณะวิศวกรรมศาสตร์ศรีราชา (fac 14)
		{ID: uuid.New(), FacultyID: faculties[14].ID, Name: "สาขาวิศวกรรมคอมพิวเตอร์และสารสนเทศศาสตร์"}, // 45
		{ID: uuid.New(), FacultyID: faculties[14].ID, Name: "สาขาวิศวกรรมเครื่องกลและการออกแบบ"},        // 46

		// คณะเศรษฐศาสตร์ ศรีราชา (fac 15)
		{ID: uuid.New(), FacultyID: faculties[15].ID, Name: "สาขาเศรษฐศาสตร์"}, // 47

		// คณะทรัพยากรธรรมชาติฯ สกลนคร (fac 16)
		{ID: uuid.New(), FacultyID: faculties[16].ID, Name: "สาขาทรัพยากรเกษตร"},   // 48
		{ID: uuid.New(), FacultyID: faculties[16].ID, Name: "สาขาอุตสาหกรรมเกษตร"}, // 49

		// คณะวิทยาศาสตร์และวิศวกรรมศาสตร์ สกลนคร (fac 17)
		{ID: uuid.New(), FacultyID: faculties[17].ID, Name: "สาขาวิศวกรรมโยธา"},        // 50
		{ID: uuid.New(), FacultyID: faculties[17].ID, Name: "สาขาวิทยาการคอมพิวเตอร์"}, // 51
	}
	for _, d := range departments {
		db.Create(&d)
	}
	return departments
}

// ──────────────────────────────────────────────
// Award Categories – ประเภทรางวัล + แบบฟอร์ม
// ──────────────────────────────────────────────

func createAwardCategories(db *gorm.DB) []models.AwardCategory {

	activityForm := models.FormStructure{
		{ID: "activity_name", Label: "ชื่อกิจกรรม", Type: "text", Required: true},
		{ID: "organizer", Label: "หน่วยงานที่จัด", Type: "text", Required: true},
		{ID: "activity_level", Label: "ระดับกิจกรรม", Type: "select", Required: true,
			Options: []string{"ภาควิชา", "คณะ", "มหาวิทยาลัย", "จังหวัด", "ประเทศ", "นานาชาติ"}},
		{ID: "role", Label: "บทบาทในกิจกรรม", Type: "select", Required: true,
			Options: []string{"ประธานโครงการ", "รองประธาน", "กรรมการ", "ผู้เข้าร่วม"}},
		{ID: "hours", Label: "จำนวนชั่วโมงกิจกรรม", Type: "number", Required: true},
		{ID: "description", Label: "รายละเอียดผลงาน", Type: "textarea", Required: true},
	}

	innovationForm := models.FormStructure{
		{ID: "title", Label: "ชื่อผลงาน/สิ่งประดิษฐ์", Type: "text", Required: true},
		{ID: "competition", Label: "ชื่อการแข่งขัน/เวที", Type: "text", Required: true},
		{ID: "award", Label: "รางวัลที่ได้รับ", Type: "text", Required: true},
		{ID: "level", Label: "ระดับการแข่งขัน", Type: "select", Required: true,
			Options: []string{"คณะ", "มหาวิทยาลัย", "ประเทศ", "นานาชาติ"}},
		{ID: "team_size", Label: "จำนวนสมาชิกทีม", Type: "number", Required: true},
		{ID: "description", Label: "รายละเอียดผลงาน", Type: "textarea", Required: true},
	}

	academicForm := models.FormStructure{
		{ID: "gpa", Label: "เกรดเฉลี่ยสะสม (GPAX)", Type: "number", Required: true},
		{ID: "semester_count", Label: "จำนวนภาคการศึกษาที่รักษาผลการเรียน", Type: "number", Required: true},
		{ID: "publications", Label: "ผลงานวิจัย/ตีพิมพ์", Type: "textarea"},
		{ID: "awards", Label: "รางวัลทางวิชาการ", Type: "textarea"},
		{ID: "description", Label: "ข้อมูลเพิ่มเติม", Type: "textarea"},
	}

	ethicsForm := models.FormStructure{
		{ID: "project", Label: "ชื่อโครงการ/กิจกรรม", Type: "text", Required: true},
		{ID: "community", Label: "ชุมชนหรือหน่วยงานที่เกี่ยวข้อง", Type: "text", Required: true},
		{ID: "duration", Label: "ระยะเวลาดำเนินโครงการ (เดือน)", Type: "number", Required: true},
		{ID: "participants", Label: "จำนวนผู้เข้าร่วม", Type: "number", Required: true},
		{ID: "impact", Label: "ผลกระทบต่อสังคม/ชุมชน", Type: "textarea", Required: true},
		{ID: "description", Label: "รายละเอียดกิจกรรม", Type: "textarea", Required: true},
	}

	leadershipForm := models.FormStructure{
		{ID: "position", Label: "ตำแหน่งผู้นำ", Type: "text", Required: true},
		{ID: "organization", Label: "ชื่อองค์กร/ชมรม", Type: "text", Required: true},
		{ID: "period", Label: "ระยะเวลาดำรงตำแหน่ง", Type: "text", Required: true},
		{ID: "achievements", Label: "ผลงานสำคัญ", Type: "textarea", Required: true},
		{ID: "description", Label: "รายละเอียดเพิ่มเติม", Type: "textarea", Required: true},
	}

	awards := []models.AwardCategory{
		{
			Name:          "ด้านกิจกรรมเสริมหลักสูตร",
			Description:   "นิสิตที่มีผลงานดีเด่นด้านกิจกรรมเสริมหลักสูตร เช่น กีฬา ศิลปะ การบริการสังคม ค่ายอาสา",
			FormStructure: activityForm,
			IsActive:      true,
		},
		{
			Name:          "ด้านความคิดสร้างสรรค์และนวัตกรรม",
			Description:   "นิสิตที่มีผลงานสร้างสรรค์นวัตกรรมหรือสิ่งประดิษฐ์ที่เป็นประโยชน์ต่อสังคม",
			FormStructure: innovationForm,
			IsActive:      true,
		},
		{
			Name:          "ด้านผลการเรียนดีเด่น",
			Description:   "นิสิตที่มีผลการเรียนดีเด่นต่อเนื่อง เกรดเฉลี่ยสะสม 3.50 ขึ้นไป",
			FormStructure: academicForm,
			IsActive:      true,
		},
		{
			Name:          "ด้านคุณธรรมจริยธรรม",
			Description:   "นิสิตที่เป็นแบบอย่างที่ดีด้านคุณธรรมจริยธรรม มีจิตสาธารณะ และบำเพ็ญประโยชน์ต่อสังคม",
			FormStructure: ethicsForm,
			IsActive:      true,
		},
		{
			Name:          "ด้านความเป็นผู้นำ",
			Description:   "นิสิตที่มีภาวะผู้นำโดดเด่น มีบทบาทสำคัญในองค์กรนิสิตหรือกิจกรรมระดับมหาวิทยาลัย",
			FormStructure: leadershipForm,
			IsActive:      true,
		},
	}

	var result []models.AwardCategory
	for _, a := range awards {
		db.Create(&a)
		result = append(result, a)
	}
	return result
}

// ──────────────────────────────────────────────
// Academic Terms
// ──────────────────────────────────────────────

func createAcademicTerms(db *gorm.DB) []models.AcademicTerm {
	now := time.Now()

	terms := []models.AcademicTerm{
		{
			AcademicYear: 2024,
			Semester:     "first",
			StartDate:    time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC),
			IsOpen:       false,
		},
		{
			AcademicYear: 2024,
			Semester:     "second",
			StartDate:    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
			IsOpen:       false,
		},
		{
			AcademicYear: 2025,
			Semester:     "first",
			StartDate:    time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
			IsOpen:       false,
		},
		{
			AcademicYear: 2025,
			Semester:     "second",
			StartDate:    time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
			IsOpen:       true,
		},
	}

	var result []models.AcademicTerm
	for _, t := range terms {
		t.CreatedAt = now
		t.UpdatedAt = now
		db.Create(&t)
		result = append(result, t)
	}
	return result
}

// ──────────────────────────────────────────────
// Advisors – อาจารย์ที่ปรึกษาอ้างอิงชื่อจริง (KU)
// ──────────────────────────────────────────────

var advisors = []string{
	"ผศ.ดร.ภัทร ลีลาพฤทธิ์",
	"รศ.ดร.อานนท์ รุ่งสว่าง",
	"ผศ.ดร.สุรีย์พร อัศวศิลปินท์",
	"รศ.ดร.ณัฐวุฒิ ขวัญแก้ว",
	"ผศ.ดร.วรรณดี สุทธินยชื่อ",
	"รศ.ดร.จักรกฤษ เติมกล้า",
	"ผศ.ดร.ธีรนันท์ นาคทอง",
	"รศ.ดร.สมศักดิ์ ศรีสมบูรณ์",
	"ผศ.ดร.ปิยนุช เวทย์วิวรณ์",
	"รศ.ดร.กฤษณ์ คงเจริญ",
	"ผศ.ดร.นิภาพร กัลยา",
	"รศ.ดร.วิชัย ชินบุตร",
	"ผศ.ดร.ชัยยุทธ ธรรมสาร",
	"รศ.ดร.สิริพร ศศิมณฑลกุล",
	"ผศ.ดร.พรทิพย์ เกิดทรัพย์",
	"รศ.ดร.ประภาส ช่างเรือ",
}

func randomAdvisor() string {
	return advisors[rand.Intn(len(advisors))]
}

// ──────────────────────────────────────────────
// Users
// ──────────────────────────────────────────────

type studentInfo struct {
	firstName  string
	lastName   string
	facultyIdx int
	deptIdx    int
	gpa        float64
	year       int
}

func createUsers(db *gorm.DB, campuses []models.Campus, faculties []models.Faculty, departments []models.Department) map[string]models.User {
	users := make(map[string]models.User)

	// --- Admin ---
	admin := models.User{
		UserID:       uuid.New(),
		FirstName:    "วิชาญ",
		LastName:     "จิตรกร",
		Email:        "admin@ku.th",
		Role:         enums.Admin,
		AuthProvider: "google",
	}
	db.Create(&admin)
	users["admin"] = admin

	// --- Admin 2 ---
	admin2 := models.User{
		UserID:       uuid.New(),
		FirstName:    "My",
		LastName:     "Admin",
		Email:        "game134613@gmail.com",
		Role:         enums.Admin,
		AuthProvider: "google",
	}
	db.Create(&admin2)
	users["admin2"] = admin2

	// --- Committee Chair ---
	// คณะกรรมการอยู่ที่วิทยาเขตบางเขน (วิทยาเขตหลัก)
	committee := models.User{
		UserID:       uuid.New(),
		FirstName:    "ศ.ดร.สุวิทย์",
		LastName:     "เมษินทรีย์",
		Email:        "committee@ku.th",
		Role:         enums.CommitteeChair,
		CampusID:     &campuses[0].ID, // วิทยาเขตบางเขน
		AuthProvider: "google",
	}
	db.Create(&committee)
	users["committee"] = committee

	// --- Deans / Vice Deans / HODs for key faculties ---
	type facultyStaff struct {
		facultyIdx    int
		deptIdx       int
		deanFirst     string
		deanLast      string
		viceDeanFirst string
		viceDeanLast  string
		hodFirst      string
		hodLast       string
	}

	staffList := []facultyStaff{
		{0, 0, "รศ.ดร.วิศิษฐ์", "ลิ้มสกุล", "ผศ.ดร.ธวัชชัย", "สุวรรณคีรี", "ผศ.ดร.ชัยพร", "ใจแก้ว"},
		{0, 1, "", "", "", "", "ผศ.ดร.สุทธิศักดิ์", "พงศ์ธนา"},
		{1, 6, "รศ.ดร.อรินทิพย์", "ศิริไพ", "ผศ.ดร.ปิยะ", "ตั้งศรีวงศ์", "ผศ.ดร.อรรถสิทธิ์", "สุรฤกษ์"},
		{2, 12, "รศ.ดร.ปิยะ", "กิตติภาดากุล", "ผศ.ดร.ธงชัย", "มาลา", "ผศ.ดร.เจริญศักดิ์", "โรจนฤทธิ์พิเชฐ"},
		{3, 16, "รศ.ดร.กิตติชัย", "แสงสว่าง", "ผศ.ดร.จรัลศรี", "พิทักษ์รังสา", "ผศ.ดร.ณัฏฐ์ชุดา", "วิจิตรจามรี"},
		{4, 20, "รศ.ดร.ธนารัตน์", "มุนินทร์", "ผศ.ดร.เฉลิมพล", "เฉลิมดง", "ผศ.ดร.พิริยะ", "ผลพิรุฬห์"},
		{5, 24, "รศ.ดร.ศศิวิมล", "มีอำพล", "ผศ.ดร.ชยันต์", "ตันติวัสดาการ", "ผศ.ดร.หฤทัย", "นำประเสริฐ"},
		{6, 28, "รศ.ดร.วุฒิยา", "สาหร่ายทอง", "ผศ.ดร.อรุณี", "ปัญญสวัสดิ์สุทธิ์", "ผศ.ดร.ฐาปกร", "จิตรถเวช"},
		{7, 31, "รศ.ดร.ปรีดา", "เลิศพงศ์วิภูษณะ", "ผศ.ดร.สุณีรัตน์", "เอี่ยมศิริ", "ผศ.ดร.อุคเดช", "บุสดี"},
		{8, 34, "รศ.ดร.ปัทมาวดี", "เล่ห์มงคล", "ผศ.ดร.สุวรรณา", "สุภิมล", "ผศ.ดร.สมบัติ", "อ่อนศิริ"},
		{13, 43, "รศ.ดร.ศรีอร", "สมบูรณ์ทรัพย์", "ผศ.ดร.ธนัท", "อมาตยกุล", "ผศ.ดร.นุชนาถ", "มั่งเจริญ"},
		{14, 45, "รศ.ดร.สถาพร", "เชื้อเพ็ง", "ผศ.ดร.ปรเมศวร์", "ตั้งจาตุรนต์", "ผศ.ดร.สมชาย", "ลิ่มอรุณ"},
	}

	staffIdx := 0
	for _, s := range staffList {
		if s.deanFirst != "" {
			dean := models.User{
				UserID:       uuid.New(),
				FirstName:    s.deanFirst,
				LastName:     s.deanLast,
				Email:        fmt.Sprintf("dean%d@ku.th", staffIdx+1),
				Role:         enums.Dean,
				FacultyID:    &faculties[s.facultyIdx].ID,
				CampusID:     &faculties[s.facultyIdx].CampusID,
				AuthProvider: "google",
			}
			db.Create(&dean)
			users[fmt.Sprintf("dean%d", staffIdx)] = dean
		}

		if s.viceDeanFirst != "" {
			viceDean := models.User{
				UserID:       uuid.New(),
				FirstName:    s.viceDeanFirst,
				LastName:     s.viceDeanLast,
				Email:        fmt.Sprintf("vicedean%d@ku.th", staffIdx+1),
				Role:         enums.ViceDean,
				FacultyID:    &faculties[s.facultyIdx].ID,
				CampusID:     &faculties[s.facultyIdx].CampusID,
				AuthProvider: "google",
			}
			db.Create(&viceDean)
			users[fmt.Sprintf("vicedean%d", staffIdx)] = viceDean
		}

		hod := models.User{
			UserID:       uuid.New(),
			FirstName:    s.hodFirst,
			LastName:     s.hodLast,
			Email:        fmt.Sprintf("hod%d@ku.th", staffIdx+1),
			Role:         enums.HeadOfDepartment,
			DepartmentID: &departments[s.deptIdx].ID,
			FacultyID:    &faculties[s.facultyIdx].ID,
			CampusID:     &faculties[s.facultyIdx].CampusID,
			AuthProvider: "google",
		}
		db.Create(&hod)
		users[fmt.Sprintf("hod%d", staffIdx)] = hod

		staffIdx++
	}

	// --- Students (50 คน กระจายทุกวิทยาเขต) ---
	students := []studentInfo{
		// วิศวกรรมคอมพิวเตอร์ (fac 0, dept 0)
		{"ณัฐพล", "ศรีวิชัย", 0, 0, 3.92, 4},
		{"ปวริศ", "ธนาพรรณ", 0, 0, 3.78, 3},
		{"กัญญาวีร์", "ลิมปิสวัสดิ์", 0, 0, 3.65, 3},
		{"ธนดล", "เจริญสุข", 0, 0, 3.51, 2},

		// วิศวกรรมไฟฟ้า (fac 0, dept 1)
		{"พีรพัฒน์", "อินทรสุวรรณ", 0, 1, 3.84, 4},
		{"อภิสรา", "คงสุวรรณ", 0, 1, 3.72, 3},

		// วิศวกรรมเครื่องกล (fac 0, dept 2)
		{"สิรวิชญ์", "ตันติพจน์", 0, 2, 3.55, 3},
		{"นภัสสร", "จิตตรง", 0, 2, 3.68, 4},

		// วิศวกรรมโยธา (fac 0, dept 5)
		{"ชนาธิป", "กิจจานนท์", 0, 5, 3.81, 4},

		// วิทยาการคอมพิวเตอร์ (fac 1, dept 6)
		{"ภูมิพัฒน์", "วัฒนเสรี", 1, 6, 3.96, 4},
		{"สโรชา", "ภูมิภาคพัทธ์", 1, 6, 3.88, 3},
		{"เตชินท์", "ศรีโยธิน", 1, 6, 3.73, 3},

		// คณิตศาสตร์ (fac 1, dept 7)
		{"ปุณยนุช", "วิจิตรวงศ์ทอง", 1, 7, 3.90, 4},

		// เคมี (fac 1, dept 9)
		{"ธนภัทร", "เอื้อเฟื้อ", 1, 9, 3.62, 3},
		{"มนัสนันท์", "สมานมิตร", 1, 9, 3.77, 2},

		// ชีววิทยา (fac 1, dept 11)
		{"จิดาภา", "ธรรมชาติ", 1, 11, 3.85, 4},

		// พืชไร่นา (fac 2, dept 12)
		{"ศุภกร", "เกษตรพงษ์", 2, 12, 3.70, 3},
		{"พัชรินทร์", "ทุ่งสว่าง", 2, 12, 3.64, 3},

		// สัตวบาล (fac 2, dept 13)
		{"กิตติธัช", "ศิริสาร", 2, 13, 3.58, 4},

		// ภาษาอังกฤษ (fac 3, dept 16)
		{"ณิชาภัทร", "พงศ์สุชน", 3, 16, 3.87, 3},
		{"ธัญชนก", "สุขประเสริฐ", 3, 16, 3.69, 4},

		// ภาษาไทย (fac 3, dept 17)
		{"กวินท์", "อภิรักษ์", 3, 17, 3.74, 3},

		// รัฐศาสตร์ (fac 4, dept 20)
		{"ภคพล", "ชัยนิมิต", 4, 20, 3.82, 4},
		{"อาทิตยา", "เทียนทอง", 4, 20, 3.76, 3},

		// จิตวิทยา (fac 4, dept 22)
		{"พิมพ์มาดา", "เลิศประสิทธิ์", 4, 22, 3.91, 4},

		// การจัดการ (fac 5, dept 24)
		{"นรวิชญ์", "ธนโชติ", 5, 24, 3.79, 3},
		{"ชลธิชา", "ดีเลิศ", 5, 24, 3.67, 4},

		// การตลาด (fac 5, dept 25)
		{"ปภังกร", "วรรณพฤกษ์", 5, 25, 3.54, 3},

		// การเงิน (fac 5, dept 26)
		{"สิริยากร", "มั่งมี", 5, 26, 3.83, 4},

		// เศรษฐศาสตร์ (fac 6, dept 28)
		{"กันตพัฒน์", "ภูวนาท", 6, 28, 3.86, 3},
		{"ปาณิสรา", "สุนทรวัฒน์", 6, 28, 3.71, 4},

		// เวชศาสตร์คลินิกสัตว์เลี้ยง (fac 7, dept 31)
		{"ภัสสร", "ชุติกุลกรณ์", 7, 31, 3.93, 5},
		{"จิรัฏฐ์", "อภิบาลธรรม", 7, 31, 3.80, 4},

		// พลศึกษา (fac 8, dept 34)
		{"กษิดิศ", "ศิลป์วิสุทธิ์", 8, 34, 3.61, 3},
		{"สุพิชฌาย์", "พงศ์วราภา", 8, 34, 3.75, 4},

		// เพาะเลี้ยงสัตว์น้ำ (fac 9, dept 36)
		{"ธีรภัทร", "สาครประเสริฐ", 9, 36, 3.63, 3},

		// วิศวกรรมเกษตร กำแพงแสน (fac 10, dept 38)
		{"ณัฐนนท์", "สวนทอง", 10, 38, 3.57, 3},
		{"วรรณพร", "ผ่องแผ้ว", 10, 38, 3.72, 4},

		// ศิลปศาสตร์ฯ กำแพงแสน (fac 11, dept 41)
		{"ศิรดา", "ภูวธนชัย", 11, 41, 3.66, 3},

		// ศึกษาศาสตร์ฯ กำแพงแสน (fac 12, dept 42)
		{"ปริชญา", "ทวีสิน", 12, 42, 3.78, 4},

		// วิทยาการจัดการ ศรีราชา (fac 13, dept 43)
		{"ธนกฤต", "ศรศรี", 13, 43, 3.74, 3},
		{"อัจฉริยา", "สัตยาบัน", 13, 43, 3.81, 4},

		// วิศวกรรมศาสตร์ ศรีราชา (fac 14, dept 45)
		{"พชรพล", "เจริญรัตน์", 14, 45, 3.69, 3},
		{"นันท์นภัส", "แก้วมณี", 14, 45, 3.87, 4},

		// เศรษฐศาสตร์ ศรีราชา (fac 15, dept 47)
		{"กฤตภาส", "จรัสพิสิฐ", 15, 47, 3.60, 3},

		// ทรัพยากรเกษตร สกลนคร (fac 16, dept 48)
		{"วรากร", "แดนทอง", 16, 48, 3.56, 3},
		{"ฐิตาภา", "โพธิ์ทอง", 16, 48, 3.73, 4},

		// วิทยาการคอมพิวเตอร์ สกลนคร (fac 17, dept 51)
		{"พัฒนพล", "ดวงดาว", 17, 51, 3.65, 3},
		{"ศิริลักษณ์", "สายน้ำ", 17, 51, 3.82, 4},

		// เพิ่มนิสิตวิทยาการคอมพิวเตอร์ วิทยาเขตบางเขน (fac 1, dept 6) - 50 คนสำหรับครบทุก award × status
		{"ธนกร", "พัฒนากิจ", 1, 6, 3.94, 4},
		{"วรัญญา", "สุขสันต์", 1, 6, 3.89, 3},
		{"ณัฐวุฒิ", "เจริญชัย", 1, 6, 3.76, 4},
		{"กัญญารัตน์", "บุญมี", 1, 6, 3.68, 3},
		{"อภิสิทธิ์", "รุ่งเรือง", 1, 6, 3.82, 4},
		{"ชนิดา", "สว่างจิต", 1, 6, 3.91, 3},
		{"ณัฐกานต์", "ปัญญา", 1, 6, 3.71, 2},
		{"พีรพัศ", "ทองคำ", 1, 6, 3.85, 4},
		{"ปิยะนุช", "สมบูรณ์", 1, 6, 3.79, 3},
		{"รัฐภูมิ", "วิชัยดิษฐ", 1, 6, 3.66, 4},

		// เพิ่ม 40 คนเพื่อให้ครบ 50 คน (แต่ละคนยื่นรางวัลได้ 1 ประเภท/ภาคการศึกษา)
		{"ธนาธิป", "สุขใส", 1, 6, 3.88, 3},
		{"พิชญา", "ดีงาม", 1, 6, 3.75, 4},
		{"กิตติพงษ์", "วงศ์ดี", 1, 6, 3.92, 3},
		{"สุภาพร", "ใจดี", 1, 6, 3.69, 4},
		{"ชัยวัฒน์", "รุ่งเรือง", 1, 6, 3.87, 3},
		{"นภัสวรรณ", "ชัยชนะ", 1, 6, 3.74, 4},
		{"ภูริณัฐ", "สว่างศรี", 1, 6, 3.83, 3},
		{"จิราภรณ์", "มั่งคั่ง", 1, 6, 3.78, 4},
		{"ธีรพงศ์", "เจริญกิจ", 1, 6, 3.90, 3},
		{"ปภาวรินทร์", "สุขสม", 1, 6, 3.67, 4},

		{"วีรภัทร", "ทองดี", 1, 6, 3.86, 3},
		{"นันทิดา", "พูลสวัสดิ์", 1, 6, 3.72, 4},
		{"ศุภชัย", "มีสุข", 1, 6, 3.95, 3},
		{"อรพิน", "แสงสว่าง", 1, 6, 3.70, 2},
		{"ธนพล", "ชาญชัย", 1, 6, 3.84, 4},
		{"พรรณพิมล", "สมบูรณ์", 1, 6, 3.77, 3},
		{"กฤตนัย", "เลิศล้ำ", 1, 6, 3.89, 4},
		{"ชญานี", "วิริยะ", 1, 6, 3.73, 3},
		{"ปฏิพล", "อุดมสุข", 1, 6, 3.81, 4},
		{"สุธีรา", "ผลดี", 1, 6, 3.76, 3},

		{"ธนวัฒน์", "สุขเกษม", 1, 6, 3.93, 4},
		{"ญาณิศา", "ปัญญา", 1, 6, 3.68, 3},
		{"อติชาต", "เจริญดี", 1, 6, 3.85, 4},
		{"มนัสนันท์", "รุ่งโรจน์", 1, 6, 3.79, 3},
		{"ธนากร", "บุญมา", 1, 6, 3.71, 4},
		{"ปิยะธิดา", "ศรีสุข", 1, 6, 3.88, 3},
		{"วรากร", "ไพบูลย์", 1, 6, 3.74, 4},
		{"ชนิกานต์", "สมสุข", 1, 6, 3.82, 3},
		{"ภาณุพงศ์", "เจริญรุ่ง", 1, 6, 3.90, 4},
		{"ศิริลักษณ์", "พิมพ์ดี", 1, 6, 3.67, 3},

		{"ธนโชติ", "สุขสดใส", 1, 6, 3.86, 4},
		{"พิมพ์ชนก", "วัฒนา", 1, 6, 3.75, 3},
		{"ศุภกร", "เจริญศรี", 1, 6, 3.91, 4},
		{"นันท์นภัส", "รุ่งเรือง", 1, 6, 3.69, 3},
		{"ณัฐพงศ์", "ดีเด่น", 1, 6, 3.83, 4},
		{"สุภัสสร", "มั่งมี", 1, 6, 3.78, 3},
		{"อภิวัฒน์", "ชัยชนะ", 1, 6, 3.87, 2},
		{"พัชรพร", "สว่าง", 1, 6, 3.72, 4},
		{"ธนบดี", "เจริญสุข", 1, 6, 3.94, 3},
		{"วริศรา", "ปัญญาดี", 1, 6, 3.80, 4},
	}

	for i, s := range students {
		nisitID := fmt.Sprintf("65%08d", 10340000+i)
		phone := fmt.Sprintf("09%d%04d", 10000+rand.Intn(89999), rand.Intn(10000))

		user := models.User{
			UserID:       uuid.New(),
			FirstName:    s.firstName,
			LastName:     s.lastName,
			Email:        fmt.Sprintf("%s.%s@ku.th", s.firstName, s.lastName),
			NisitID:      &nisitID,
			PhoneNumber:  &phone,
			Role:         enums.Student,
			DepartmentID: &departments[s.deptIdx].ID,
			FacultyID:    &faculties[s.facultyIdx].ID,
			CampusID:     &faculties[s.facultyIdx].CampusID,
			AuthProvider: "google",
		}
		db.Create(&user)
		users[fmt.Sprintf("student%d", i)] = user
	}

	return users
}

// ──────────────────────────────────────────────
// Requests – คำขอรางวัล (พร้อม status history)
// ──────────────────────────────────────────────

func createRequests(
	db *gorm.DB,
	users map[string]models.User,
	awards []models.AwardCategory,
	terms []models.AcademicTerm,
) []models.Request {
	now := time.Now()
	var requests []models.Request

	totalStudents := 100

	// Student GPAs (100 students - 50 คนแรก + 50 นิสิตวิทยาการคอมพิวเตอร์)
	studentGPAs := []float64{
		3.92, 3.78, 3.65, 3.51, 3.84, 3.72, 3.55, 3.68, 3.81,
		3.96, 3.88, 3.73, 3.90, 3.62, 3.77, 3.85, 3.70, 3.64, 3.58,
		3.87, 3.69, 3.74, 3.82, 3.76, 3.91, 3.79, 3.67, 3.54, 3.83,
		3.86, 3.71, 3.93, 3.80, 3.61, 3.75, 3.63, 3.57, 3.72, 3.66,
		3.78, 3.74, 3.81, 3.69, 3.87, 3.60, 3.56, 3.73, 3.65, 3.82,
		3.75, // student49
		// นิสิตวิทยาการคอมพิวเตอร์ 50 คน (student50-99)
		3.94, 3.89, 3.76, 3.68, 3.82, 3.91, 3.71, 3.85, 3.79, 3.66, // 50-59
		3.88, 3.75, 3.92, 3.69, 3.87, 3.74, 3.83, 3.78, 3.90, 3.67, // 60-69
		3.86, 3.72, 3.95, 3.70, 3.84, 3.77, 3.89, 3.73, 3.81, 3.76, // 70-79
		3.93, 3.68, 3.85, 3.79, 3.71, 3.88, 3.74, 3.82, 3.90, 3.67, // 80-89
		3.86, 3.75, 3.91, 3.69, 3.83, 3.78, 3.87, 3.72, 3.94, 3.80, // 90-99
	}
	studentYears := []int{
		4, 3, 3, 2, 4, 3, 3, 4, 4,
		4, 3, 3, 4, 3, 2, 4, 3, 3, 4,
		3, 4, 3, 4, 3, 4, 3, 4, 3, 4,
		3, 4, 5, 4, 3, 4, 3, 3, 4, 3,
		4, 3, 4, 3, 4, 3, 3, 4, 3, 4,
		3, // student49
		// นิสิตวิทยาการคอมพิวเตอร์ 50 คน (student50-99)
		4, 3, 4, 3, 4, 3, 2, 4, 3, 4, // 50-59
		3, 4, 3, 4, 3, 4, 3, 4, 3, 4, // 60-69
		3, 4, 3, 2, 4, 3, 4, 3, 4, 3, // 70-79
		4, 3, 4, 3, 4, 3, 4, 3, 4, 3, // 80-89
		4, 3, 4, 3, 4, 3, 2, 4, 3, 4, // 90-99
	}

	// ── helper: สร้าง additional data ตามประเภทรางวัล ──

	activityNames := []string{
		"โครงการค่าย KU ร่วมใจ พัฒนาชนบท",
		"กิจกรรมปลูกป่าชายเลนบางปู",
		"โครงการอาสา KU สอนน้องที่สระบุรี",
		"การแข่งขันกีฬามหาวิทยาลัยแห่งประเทศไทย ครั้งที่ 49",
		"KU Band Music Festival",
		"กิจกรรมรับน้องสร้างสรรค์ Freshy KU",
		"โครงการพัฒนาโรงเรียนตำรวจตระเวนชายแดน",
		"KU Hackathon for Social Good",
		"โครงการ KU Green Campus อาสาพัฒนาสิ่งแวดล้อม",
		"งานนนทรีสีทอง (กิจกรรมรับน้อง)",
		"โครงการรณรงค์ลดขยะพลาสติกในมหาวิทยาลัย",
		"กิจกรรมวิ่ง KU Health Run 2025",
	}

	activityOrganizers := []string{
		"องค์การบริหาร องค์การนิสิต มก.",
		"สโมสรนิสิตคณะวิศวกรรมศาสตร์",
		"สโมสรนิสิตคณะวิทยาศาสตร์",
		"กองกิจการนิสิต มก.",
		"ชมรมอาสาพัฒนาและบำเพ็ญประโยชน์",
		"สภาผู้แทนนิสิต มก.",
		"ชมรมดนตรีสากล มก.",
		"ฝ่ายกิจการนิสิต คณะเกษตร",
		"สำนักบริการคอมพิวเตอร์ มก.",
		"ชมรมอนุรักษ์ธรรมชาติและสิ่งแวดล้อม",
	}

	activityLevels := []string{"ภาควิชา", "คณะ", "มหาวิทยาลัย", "จังหวัด", "ประเทศ", "นานาชาติ"}
	activityRoles := []string{"ประธานโครงการ", "รองประธาน", "กรรมการ", "ผู้เข้าร่วม"}

	innovationTitles := []string{
		"ระบบ IoT ตรวจวัดคุณภาพอากาศสำหรับฟาร์ม",
		"แอปพลิเคชันช่วยเหลือเกษตรกร SmartFarm KU",
		"ระบบ AI วิเคราะห์โรคข้าวจากภาพถ่าย",
		"หุ่นยนต์เก็บเกี่ยวผลไม้อัตโนมัติ",
		"แพลตฟอร์มจัดการขยะรีไซเคิลในมหาวิทยาลัย",
		"ระบบคัดแยกผลไม้ด้วย Computer Vision",
		"อุปกรณ์วัดความชื้นในดินแบบ Real-Time",
		"ระบบแนะนำเส้นทางรถพลังงานไฟฟ้าในวิทยาเขต",
		"แอปพลิเคชัน KU Buddy ช่วยเหลือนิสิตใหม่",
		"ระบบ Blockchain สำหรับตรวจสอบย้อนกลับสินค้าเกษตร",
	}

	innovationCompetitions := []string{
		"Thailand ICT Awards",
		"NSC (National Software Contest)",
		"SCB Challenge",
		"Startup Thailand League",
		"ASEAN Data Science Explorers",
		"True 5G Tech Startup",
		"KU Innovation Awards",
		"Young Technopreneur",
		"Thailand Research Expo",
		"IEEE Region 10 Student Paper Contest",
	}

	innovationAwards := []string{
		"รางวัลชนะเลิศ", "รางวัลรองชนะเลิศอันดับ 1", "รางวัลรองชนะเลิศอันดับ 2",
		"รางวัลชมเชย", "รางวัล Popular Vote", "รางวัล Best Innovation",
	}

	ethicsProjects := []string{
		"โครงการพัฒนาชุมชนเกษตรอินทรีย์บางเขน",
		"โครงการสอนพิเศษนิสิตด้อยโอกาส",
		"โครงการ KU จิตอาสาพัฒนาวัด",
		"โครงการ Big Brother บ้านเด็กกำพร้า",
		"โครงการส่งเสริมการอ่านโรงเรียนชนบท",
		"โครงการทำความสะอาดคลองบางเขน",
		"โครงการ KU ร่วมใจ ช่วยเหลือผู้ประสบอุทกภัย",
		"โครงการดูแลผู้สูงอายุชุมชนรอบมหาวิทยาลัย",
		"โครงการรณรงค์ต่อต้านยาเสพติดในโรงเรียน",
		"โครงการเศรษฐกิจพอเพียงตามรอยพ่อ",
	}

	ethicsCommunities := []string{
		"ชุมชนบางเขน",
		"ชุมชนวัดพระศรีมหาธาตุ",
		"ชุมชนหมู่บ้านเสนานิคม",
		"ศูนย์พัฒนาเด็กเล็กเขตจตุจักร",
		"โรงเรียนวัดเสมียนนารี",
		"มูลนิธิช่วยเหลือเด็กกำพร้า",
		"ชุมชนกำแพงแสน",
		"ชุมชนแหลมฉบัง ศรีราชา",
		"หมู่บ้านเกษตรกรรอบมหาวิทยาลัย",
		"มูลนิธิดวงประทีป",
	}

	leadershipPositions := []string{
		"นายกองค์การบริหาร องค์การนิสิต มก.",
		"ประธานสโมสรนิสิตคณะ",
		"ประธานชมรมอาสาพัฒนา",
		"ประธานสภาผู้แทนนิสิต",
		"หัวหน้าคณะกรรมการจัดงานเกษตรแฟร์",
		"ประธานชมรมดนตรีสากล",
		"ประธานชมรมหุ่นยนต์",
		"ประธานชมรมวิ่ง KU Runner",
	}

	leadershipOrgs := []string{
		"องค์การนิสิต มหาวิทยาลัยเกษตรศาสตร์",
		"สโมสรนิสิตคณะวิศวกรรมศาสตร์",
		"สโมสรนิสิตคณะวิทยาศาสตร์",
		"สภาผู้แทนนิสิต มก.",
		"ชมรมอาสาพัฒนาและบำเพ็ญประโยชน์",
		"ชมรมหุ่นยนต์ มก.",
		"สโมสรนิสิตคณะบริหารธุรกิจ",
		"ชมรม Kasetsart Maker Club",
	}

	makeAdditionalData := func(awardIdx int) models.AdditionalData {
		data := models.AdditionalData{}
		switch awardIdx {
		case 0: // กิจกรรมเสริมหลักสูตร
			data["activity_name"] = activityNames[rand.Intn(len(activityNames))]
			data["organizer"] = activityOrganizers[rand.Intn(len(activityOrganizers))]
			data["activity_level"] = activityLevels[rand.Intn(len(activityLevels))]
			data["role"] = activityRoles[rand.Intn(len(activityRoles))]
			data["hours"] = 40 + rand.Intn(161) // 40-200 hrs
			data["description"] = "เป็นผู้มีส่วนร่วมในการจัดกิจกรรมเพื่อพัฒนาชุมชนและมหาวิทยาลัยอย่างต่อเนื่อง"

		case 1: // ความคิดสร้างสรรค์และนวัตกรรม
			data["title"] = innovationTitles[rand.Intn(len(innovationTitles))]
			data["competition"] = innovationCompetitions[rand.Intn(len(innovationCompetitions))]
			data["award"] = innovationAwards[rand.Intn(len(innovationAwards))]
			data["level"] = activityLevels[2+rand.Intn(4)] // มหาวิทยาลัย, จังหวัด, ประเทศ, นานาชาติ
			data["team_size"] = 2 + rand.Intn(5)           // 2-6 members
			data["description"] = "ผลงานนวัตกรรมที่ตอบโจทย์ปัญหาสังคมและได้รับการยอมรับในระดับการแข่งขัน"

		case 2: // ผลการเรียนดีเด่น
			data["gpa"] = 3.70 + rand.Float64()*0.29  // 3.70-3.99
			data["semester_count"] = 4 + rand.Intn(5) // 4-8 semesters
			pubs := []string{
				"บทความวิจัยในวารสารระดับชาติ 1 ฉบับ",
				"บทความวิจัยในวารสาร TCI กลุ่ม 1 จำนวน 2 ฉบับ",
				"นำเสนอบทความใน National Conference 2025",
				"ผลงานวิจัยร่วมกับอาจารย์ตีพิมพ์ใน SCOPUS",
				"ได้รับเหรียญทองศึกษิตแห่งปี มก.",
			}
			data["publications"] = pubs[rand.Intn(len(pubs))]
			acAwards := []string{
				"ทุนเรียนดี มหาวิทยาลัยเกษตรศาสตร์",
				"Dean's List ต่อเนื่อง 4 ภาคการศึกษา",
				"รางวัลนิสิตผลการเรียนยอดเยี่ยม",
				"ทุนกิตติบัณฑิต ปีการศึกษา 2024",
			}
			data["awards"] = acAwards[rand.Intn(len(acAwards))]
			data["description"] = "มีผลการเรียนดีเด่นต่อเนื่อง พร้อมผลงานวิจัยที่เป็นที่ยอมรับ"

		case 3: // คุณธรรมจริยธรรม
			data["project"] = ethicsProjects[rand.Intn(len(ethicsProjects))]
			data["community"] = ethicsCommunities[rand.Intn(len(ethicsCommunities))]
			data["duration"] = 3 + rand.Intn(10)      // 3-12 months
			data["participants"] = 10 + rand.Intn(91) // 10-100 people
			data["impact"] = "ช่วยยกระดับคุณภาพชีวิตของชุมชนโดยรอบมหาวิทยาลัย สร้างความตระหนักด้านจิตสาธารณะ"
			data["description"] = "อุทิศตนเพื่อสังคมอย่างต่อเนื่อง เป็นแบบอย่างที่ดีแก่เพื่อนนิสิต"

		case 4: // ความเป็นผู้นำ
			data["position"] = leadershipPositions[rand.Intn(len(leadershipPositions))]
			data["organization"] = leadershipOrgs[rand.Intn(len(leadershipOrgs))]
			periods := []string{"1 ปีการศึกษา (2567)", "2 ภาคการศึกษา", "1 ปี 6 เดือน"}
			data["period"] = periods[rand.Intn(len(periods))]
			data["achievements"] = "จัดกิจกรรมสำเร็จตามเป้าหมาย สร้างเครือข่ายความร่วมมือระหว่างคณะ"
			data["description"] = "แสดงภาวะผู้นำโดดเด่น สามารถบริหารจัดการทีมและโครงการได้อย่างมีประสิทธิภาพ"
		}
		return data
	}

	// ── Term 0 (2024 sem 1) — Approved: 15 requests ──
	approvedT0 := []int{0, 3, 5, 9, 12, 16, 19, 22, 25, 29, 31, 34, 37, 40, 45}
	for i, sIdx := range approvedT0 {
		aIdx := i % len(awards)
		student := users[fmt.Sprintf("student%d", sIdx)]

		req := models.Request{
			ID:                uuid.New(),
			StudentID:         student.UserID,
			AwardID:           awards[aIdx].ID,
			AcademicTermID:    terms[0].ID,
			AdditionalData:    makeAdditionalData(aIdx),
			SnapshotStudyYear: studentYears[sIdx],
			SnapshotGPA:       studentGPAs[sIdx],
			SnapshotAdvisor:   randomAdvisor(),
			CurrentStatus:     string(enums.Approved),
			CreatedAt:         terms[0].StartDate.Add(time.Hour * 24 * time.Duration(10+rand.Intn(30))),
			UpdatedAt:         now,
		}
		db.Create(&req)

		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.PendingHOD), UpdatedBy: student.UserID,
			Remark: "ยื่นคำขอรับรางวัล", UpdatedAt: req.CreatedAt,
		})
		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod0"].UserID,
			Remark: "ผ่านการพิจารณาจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 5),
		})
		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.Approved), UpdatedBy: users["committee"].UserID,
			Remark: "อนุมัติตามมติที่ประชุมคณะกรรมการ", UpdatedAt: terms[0].EndDate.Add(-time.Hour * 24 * 7),
		})
		requests = append(requests, req)
	}

	// ── Term 1 (2024 sem 2) — Approved: 12 requests ──
	approvedT1 := []int{1, 4, 7, 10, 14, 17, 20, 23, 27, 30, 35, 42}
	for i, sIdx := range approvedT1 {
		aIdx := i % len(awards)
		student := users[fmt.Sprintf("student%d", sIdx)]

		req := models.Request{
			ID:                uuid.New(),
			StudentID:         student.UserID,
			AwardID:           awards[aIdx].ID,
			AcademicTermID:    terms[1].ID,
			AdditionalData:    makeAdditionalData(aIdx),
			SnapshotStudyYear: studentYears[sIdx],
			SnapshotGPA:       studentGPAs[sIdx],
			SnapshotAdvisor:   randomAdvisor(),
			CurrentStatus:     string(enums.Approved),
			CreatedAt:         terms[1].StartDate.Add(time.Hour * 24 * time.Duration(10+rand.Intn(20))),
			UpdatedAt:         now,
		}
		db.Create(&req)

		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.PendingHOD), UpdatedBy: student.UserID,
			Remark: "ยื่นคำขอรับรางวัล", UpdatedAt: req.CreatedAt,
		})
		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.Approved), UpdatedBy: users["committee"].UserID,
			Remark: "อนุมัติตามมติที่ประชุม", UpdatedAt: terms[1].EndDate.Add(-time.Hour * 24 * 10),
		})
		requests = append(requests, req)
	}

	// ── Term 2 (2025 sem 1) — Approved: 10 requests ──
	approvedT2 := []int{2, 6, 8, 11, 15, 18, 24, 33, 38, 46}
	for i, sIdx := range approvedT2 {
		aIdx := i % len(awards)
		student := users[fmt.Sprintf("student%d", sIdx)]

		req := models.Request{
			ID:                uuid.New(),
			StudentID:         student.UserID,
			AwardID:           awards[aIdx].ID,
			AcademicTermID:    terms[2].ID,
			AdditionalData:    makeAdditionalData(aIdx),
			SnapshotStudyYear: studentYears[sIdx],
			SnapshotGPA:       studentGPAs[sIdx],
			SnapshotAdvisor:   randomAdvisor(),
			CurrentStatus:     string(enums.Approved),
			CreatedAt:         terms[2].StartDate.Add(time.Hour * 24 * time.Duration(15+rand.Intn(25))),
			UpdatedAt:         now,
		}
		db.Create(&req)

		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.PendingHOD), UpdatedBy: student.UserID,
			Remark: "ยื่นคำขอรับรางวัล", UpdatedAt: req.CreatedAt,
		})
		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.Approved), UpdatedBy: users["committee"].UserID,
			Remark: "อนุมัติตามมติคณะกรรมการ", UpdatedAt: terms[2].EndDate.Add(-time.Hour * 24 * 5),
		})
		requests = append(requests, req)
	}

	// ── Term 3 (2025 sem 2 — current / open) — Various statuses: 20 requests ──
	type pendingReq struct {
		studentIdx int
		awardIdx   int
		status     enums.RequestStatus
	}

	currentRequests := []pendingReq{
		{0, 0, enums.PendingHOD},
		{1, 1, enums.PendingViceDean},
		{2, 2, enums.PendingDean},
		{4, 3, enums.PendingCommittee},
		{5, 4, enums.PendingHOD},
		{7, 0, enums.RejectedByHOD},
		{9, 1, enums.PendingViceDean},
		{10, 2, enums.PendingDean},
		{12, 3, enums.PendingCommittee},
		{13, 4, enums.PendingHOD},
		{16, 0, enums.RejectedByViceDean},
		{19, 1, enums.PendingViceDean},
		{22, 2, enums.PendingDean},
		{25, 3, enums.PendingCommittee},
		{28, 4, enums.PendingHOD},
		{31, 0, enums.RejectedByDean},
		{34, 1, enums.PendingViceDean},
		{37, 2, enums.PendingHOD},
		{41, 3, enums.PendingDean},
		{44, 4, enums.PendingCommittee},
	}

	for i, cr := range currentRequests {
		student := users[fmt.Sprintf("student%d", cr.studentIdx)]

		req := models.Request{
			ID:                uuid.New(),
			StudentID:         student.UserID,
			AwardID:           awards[cr.awardIdx].ID,
			AcademicTermID:    terms[3].ID,
			AdditionalData:    makeAdditionalData(cr.awardIdx),
			SnapshotStudyYear: studentYears[cr.studentIdx],
			SnapshotGPA:       studentGPAs[cr.studentIdx],
			SnapshotAdvisor:   randomAdvisor(),
			CurrentStatus:     string(cr.status),
			CreatedAt:         terms[3].StartDate.Add(time.Hour * 24 * time.Duration(i*2+1)),
			UpdatedAt:         now,
		}
		db.Create(&req)

		// Initial submission
		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.PendingHOD), UpdatedBy: student.UserID,
			Remark: "ยื่นคำขอรับรางวัล", UpdatedAt: req.CreatedAt,
		})

		// Intermediate statuses based on current status
		switch cr.status {
		case enums.PendingViceDean:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod0"].UserID,
				Remark: "ผ่านการพิจารณาจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
		case enums.PendingDean:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod0"].UserID,
				Remark: "ผ่านจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingDean), UpdatedBy: users["vicedean0"].UserID,
				Remark: "ผ่านจากรองคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 6),
			})
		case enums.PendingCommittee:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod0"].UserID,
				Remark: "ผ่านจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingDean), UpdatedBy: users["vicedean0"].UserID,
				Remark: "ผ่านจากรองคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 6),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingCommittee), UpdatedBy: users["dean0"].UserID,
				Remark: "ผ่านจากคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 9),
			})
		case enums.RejectedByHOD:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.RejectedByHOD), UpdatedBy: users["hod0"].UserID,
				Remark: "ข้อมูลไม่ครบถ้วน กรุณาแก้ไขและส่งใหม่", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 4),
			})
		case enums.RejectedByViceDean:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod0"].UserID,
				Remark: "ผ่านจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.RejectedByViceDean), UpdatedBy: users["vicedean0"].UserID,
				Remark: "ผลงานยังไม่เพียงพอตามเกณฑ์", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 7),
			})
		case enums.RejectedByDean:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod0"].UserID,
				Remark: "ผ่านจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingDean), UpdatedBy: users["vicedean0"].UserID,
				Remark: "ผ่านจากรองคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 6),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.RejectedByDean), UpdatedBy: users["dean0"].UserID,
				Remark: "ไม่ผ่านเกณฑ์ประเมินระดับคณะ", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 10),
			})
		}

		requests = append(requests, req)
	}

	// ── นิสิตวิทยาการคอมพิวเตอร์ (student50-99): ครบทุกสถานะ × ทุก award category ──
	// แต่ละคนยื่นเพียง 1 ประเภทรางวัลต่อภาคการศึกษา (ตามกฎของระบบ)
	// 5 award categories × 10 statuses = 50 students
	csStudentRequests := []pendingReq{
		// Award 0: กิจกรรมเสริมหลักสูตร - ทุกสถานะ (student50-59)
		{50, 0, enums.PendingHOD},
		{51, 0, enums.PendingViceDean},
		{52, 0, enums.PendingDean},
		{53, 0, enums.PendingCommittee},
		{54, 0, enums.Approved},
		{55, 0, enums.RejectedByHOD},
		{56, 0, enums.RejectedByViceDean},
		{57, 0, enums.RejectedByDean},
		{58, 0, enums.RejectedByCommittee},
		{59, 0, enums.AwardChanged},

		// Award 1: ความคิดสร้างสรรค์และนวัตกรรม - ทุกสถานะ (student60-69)
		{60, 1, enums.PendingHOD},
		{61, 1, enums.PendingViceDean},
		{62, 1, enums.PendingDean},
		{63, 1, enums.PendingCommittee},
		{64, 1, enums.Approved},
		{65, 1, enums.RejectedByHOD},
		{66, 1, enums.RejectedByViceDean},
		{67, 1, enums.RejectedByDean},
		{68, 1, enums.RejectedByCommittee},
		{69, 1, enums.AwardChanged},

		// Award 2: ผลการเรียนดีเด่น - ทุกสถานะ (student70-79)
		{70, 2, enums.PendingHOD},
		{71, 2, enums.PendingViceDean},
		{72, 2, enums.PendingDean},
		{73, 2, enums.PendingCommittee},
		{74, 2, enums.Approved},
		{75, 2, enums.RejectedByHOD},
		{76, 2, enums.RejectedByViceDean},
		{77, 2, enums.RejectedByDean},
		{78, 2, enums.RejectedByCommittee},
		{79, 2, enums.AwardChanged},

		// Award 3: คุณธรรมจริยธรรม - ทุกสถานะ (student80-89)
		{80, 3, enums.PendingHOD},
		{81, 3, enums.PendingViceDean},
		{82, 3, enums.PendingDean},
		{83, 3, enums.PendingCommittee},
		{84, 3, enums.Approved},
		{85, 3, enums.RejectedByHOD},
		{86, 3, enums.RejectedByViceDean},
		{87, 3, enums.RejectedByDean},
		{88, 3, enums.RejectedByCommittee},
		{89, 3, enums.AwardChanged},

		// Award 4: ความเป็นผู้นำ - ทุกสถานะ (student90-99)
		{90, 4, enums.PendingHOD},
		{91, 4, enums.PendingViceDean},
		{92, 4, enums.PendingDean},
		{93, 4, enums.PendingCommittee},
		{94, 4, enums.Approved},
		{95, 4, enums.RejectedByHOD},
		{96, 4, enums.RejectedByViceDean},
		{97, 4, enums.RejectedByDean},
		{98, 4, enums.RejectedByCommittee},
		{99, 4, enums.AwardChanged},
	}

	for i, cr := range csStudentRequests {
		student := users[fmt.Sprintf("student%d", cr.studentIdx)]

		// สำหรับ AwardChanged: เปลี่ยนเป็น award ถัดไป
		actualAwardIdx := cr.awardIdx
		if cr.status == enums.AwardChanged {
			actualAwardIdx = (cr.awardIdx + 1) % len(awards)
		}

		req := models.Request{
			ID:                uuid.New(),
			StudentID:         student.UserID,
			AwardID:           awards[actualAwardIdx].ID,
			AcademicTermID:    terms[3].ID,
			AdditionalData:    makeAdditionalData(cr.awardIdx), // ใช้ award เดิมสำหรับ form data
			SnapshotStudyYear: studentYears[cr.studentIdx],
			SnapshotGPA:       studentGPAs[cr.studentIdx],
			SnapshotAdvisor:   randomAdvisor(),
			CurrentStatus:     string(cr.status),
			CreatedAt:         terms[3].StartDate.Add(time.Hour * 24 * time.Duration(i*3+1)),
			UpdatedAt:         now,
		}
		db.Create(&req)

		// Initial submission
		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.PendingHOD), UpdatedBy: student.UserID,
			Remark: "ยื่นคำขอรับรางวัล", UpdatedAt: req.CreatedAt,
		})

		// Intermediate statuses based on current status
		switch cr.status {
		case enums.PendingViceDean:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod2"].UserID,
				Remark: "ผ่านการพิจารณาจากหัวหน้าภาควิชาวิทยาการคอมพิวเตอร์", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
		case enums.PendingDean:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod2"].UserID,
				Remark: "ผ่านการพิจารณาจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingDean), UpdatedBy: users["vicedean2"].UserID,
				Remark: "ผ่านการพิจารณาจากรองคณบดี คณะวิทยาศาสตร์", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 6),
			})
		case enums.PendingCommittee:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod2"].UserID,
				Remark: "ผ่านการพิจารณาจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingDean), UpdatedBy: users["vicedean2"].UserID,
				Remark: "ผ่านการพิจารณาจากรองคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 6),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingCommittee), UpdatedBy: users["dean2"].UserID,
				Remark: "ผ่านการพิจารณาจากคณบดี ส่งต่อคณะกรรมการ", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 9),
			})
		case enums.Approved:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod2"].UserID,
				Remark: "ผ่านการพิจารณาจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingDean), UpdatedBy: users["vicedean2"].UserID,
				Remark: "ผ่านการพิจารณาจากรองคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 6),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingCommittee), UpdatedBy: users["dean2"].UserID,
				Remark: "ผ่านการพิจารณาจากคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 9),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.Approved), UpdatedBy: users["committee"].UserID,
				Remark: "อนุมัติตามมติที่ประชุมคณะกรรมการ รอการอนุมัติจากอธิการบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 12),
			})
		case enums.RejectedByHOD:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.RejectedByHOD), UpdatedBy: users["hod2"].UserID,
				Remark: "ข้อมูลไม่ครบถ้วน กรุณาแก้ไขและส่งใหม่", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 4),
			})
		case enums.RejectedByViceDean:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod2"].UserID,
				Remark: "ผ่านการพิจารณาจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.RejectedByViceDean), UpdatedBy: users["vicedean2"].UserID,
				Remark: "ผลงานยังไม่เพียงพอตามเกณฑ์ของคณะ", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 7),
			})
		case enums.RejectedByDean:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod2"].UserID,
				Remark: "ผ่านการพิจารณาจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingDean), UpdatedBy: users["vicedean2"].UserID,
				Remark: "ผ่านการพิจารณาจากรองคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 6),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.RejectedByDean), UpdatedBy: users["dean2"].UserID,
				Remark: "ไม่ผ่านเกณฑ์ประเมินระดับคณะวิทยาศาสตร์", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 10),
			})
		case enums.RejectedByCommittee:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod2"].UserID,
				Remark: "ผ่านการพิจารณาจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingDean), UpdatedBy: users["vicedean2"].UserID,
				Remark: "ผ่านการพิจารณาจากรองคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 6),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingCommittee), UpdatedBy: users["dean2"].UserID,
				Remark: "ผ่านการพิจารณาจากคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 9),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.RejectedByCommittee), UpdatedBy: users["committee"].UserID,
				Remark: "คณะกรรมการเห็นว่าผลงานไม่ผ่านเกณฑ์ในการพิจารณา", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 13),
			})
		case enums.AwardChanged:
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingViceDean), UpdatedBy: users["hod2"].UserID,
				Remark: "ผ่านการพิจารณาจากหัวหน้าภาค", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 3),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingDean), UpdatedBy: users["vicedean2"].UserID,
				Remark: "ผ่านการพิจารณาจากรองคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 6),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.PendingCommittee), UpdatedBy: users["dean2"].UserID,
				Remark: "ผ่านการพิจารณาจากคณบดี", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 9),
			})
			db.Create(&models.RequestStatusHistory{
				RequestID: req.ID, Status: string(enums.AwardChanged), UpdatedBy: users["committee"].UserID,
				Remark: "คณะกรรมการเปลี่ยนประเภทรางวัล ต้องรออนุมัติใหม่", UpdatedAt: req.CreatedAt.Add(time.Hour * 24 * 12),
			})

			// สร้าง AwardCategoryChange record เพื่อบันทึกการเปลี่ยนประเภทรางวัล
			// เปลี่ยนจาก award ที่นิสิตยื่น (cr.awardIdx) ไปเป็น award ใหม่ (req.AwardID)
			db.Create(&models.AwardCategoryChange{
				RequestID:     req.ID,
				OldCategoryID: awards[cr.awardIdx].ID,
				NewCategoryID: req.AwardID, // ใช้ AwardID ที่อัปเดทแล้ว
				ChangedBy:     users["committee"].UserID,
				ChangedAt:     req.CreatedAt.Add(time.Hour * 24 * 12),
			})
		}

		requests = append(requests, req)
	}

	// ── Extra: นิสิตที่ยังไม่เคยยื่น (เพื่อให้มี active students ใน term ปัจจุบัน) ──
	unusedStudents := []int{3, 6, 8, 11, 14, 15, 17, 20, 21, 23, 24, 26, 27, 29, 30, 32, 33, 35, 36, 38, 39, 40, 42, 43, 45, 46, 47, 48, 49}
	for i := 0; i < 10; i++ {
		sIdx := unusedStudents[rand.Intn(len(unusedStudents))]
		aIdx := rand.Intn(len(awards))
		student := users[fmt.Sprintf("student%d", sIdx)]

		req := models.Request{
			ID:                uuid.New(),
			StudentID:         student.UserID,
			AwardID:           awards[aIdx].ID,
			AcademicTermID:    terms[3].ID,
			AdditionalData:    makeAdditionalData(aIdx),
			SnapshotStudyYear: studentYears[sIdx],
			SnapshotGPA:       studentGPAs[sIdx],
			SnapshotAdvisor:   randomAdvisor(),
			CurrentStatus:     string(enums.PendingHOD),
			CreatedAt:         now.Add(-time.Hour * 24 * time.Duration(rand.Intn(14))),
			UpdatedAt:         now,
		}
		db.Create(&req)

		db.Create(&models.RequestStatusHistory{
			RequestID: req.ID, Status: string(enums.PendingHOD), UpdatedBy: student.UserID,
			Remark: "ยื่นคำขอรับรางวัล", UpdatedAt: req.CreatedAt,
		})
		requests = append(requests, req)
	}

	_ = totalStudents
	return requests
}
