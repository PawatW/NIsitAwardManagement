# Database Seeding

## Overview
ไฟล์สำหรับสร้าง Mock Data เข้า Database สำหรับการทดสอบระบบ

## ข้อมูลที่จะถูกสร้าง

### 📅 Academic Terms (4 semesters จาก 2 ปีการศึกษา)
- ปีการศึกษา 2024 ภาคเรียนที่ 1 (ปิดแล้ว)
- ปีการศึกษา 2024 ภาคเรียนที่ 2 (ปิดแล้ว)
- ปีการศึกษา 2025 ภาคเรียนที่ 1 (ปิดแล้ว)
- ปีการศึกษา 2025 ภาคเรียนที่ 2 (เปิดรับสมัคร) ✅

### 🏫 Organization Structure
- **3 วิทยาเขต**: กรุงเทพฯ, ปทุมธานี, ชลบุรี
- **7 คณะ**:
  - วิทยาศาสตร์
  - วิศวกรรมศาสตร์
  - ศิลปศาสตร์
  - เทคโนโลยีสารสนเทศ
  - บริหารธุรกิจ
  - สังคมศาสตร์
  - ครุศาสตร์
- **14 ภาควิชา** กระจายตามคณะต่างๆ

### 🏆 Award Categories (4 ประเภท)
1. ด้านกิจกรรมเสริมหลักสูตร
2. ด้านความคิดสร้างสรรค์และนวัตกรรม
3. ด้านผลการเรียนดีเด่น
4. ด้านคุณธรรมจริยธรรม

### 👥 Users (รวม 60+ accounts)
- **1 Admin** - admin@university.ac.th
- **1 Committee Chair** - committee@university.ac.th
- **9 Deans** - dean1@university.ac.th, dean2@university.ac.th, ...
- **9 Vice Deans** - vicedean1@university.ac.th, ...
- **9 Head of Departments** - hod1@university.ac.th, ...
- **30 Students** - student1@university.ac.th, student2@university.ac.th, ...
  - มีข้อมูล: ชื่อ-นามสกุล, เลขนิสิต (65########), GPA, คณะ, ภาควิชา

### 📝 Requests (รวม 50+ requests)

#### ✅ Approved Requests (สำหรับแสดงใน Honor Roll)
- **ปีการศึกษา 2024 ภาค 1**: 12 รายการ (APPROVED_PENDING_PRESIDENT)
- **ปีการศึกษา 2024 ภาค 2**: 10 รายการ (APPROVED_PENDING_PRESIDENT)
- **ปีการศึกษา 2025 ภาค 1**: 8 รายการ (APPROVED_PENDING_PRESIDENT)

#### 🔄 Pending/In-Progress Requests (ปีการศึกษา 2025 ภาค 2)
- กระจายสถานะ: PENDING_HOD, PENDING_VICE_DEAN, PENDING_DEAN, PENDING_COMMITTEE, REJECTED_BY_HOD

## วิธีใช้งาน

### 1. Migrate Schema ก่อน
```bash
make migrate-schema
```

### 2. Seed ข้อมูล
```bash
make seed-data
```

**หมายเหตุ**: คำสั่ง `seed-data` จะลบข้อมูลเก่าออกทั้งหมดก่อนสร้างข้อมูลใหม่

### 3. ถ้าต้องการเก็บข้อมูลเก่าไว้
แก้ไขไฟล์ `cmd/seed/main.go` โดยคอมเมนต์บรรทัด:
```go
// clearData(db)
```

## การทดสอบ Honor Roll

หลังจาก seed ข้อมูลแล้ว สามารถทดสอบ honor roll API ได้:

```bash
# ดูทั้งหมด
GET http://localhost:8080/api/v1/honor-roll

# กรองตามปีการศึกษา 2024
GET http://localhost:8080/api/v1/honor-roll?academic_year=2024

# กรองตามวิทยาเขต
GET http://localhost:8080/api/v1/honor-roll?campus_id=<uuid>

# กรองหลายเงื่อนไข
GET http://localhost:8080/api/v1/honor-roll?academic_year=2024&semester=first
```

## Test Users

### Admin
- Email: admin@university.ac.th
- Role: Admin

### Committee
- Email: committee@university.ac.th
- Role: CommitteeChair

### Students
- Email: student1@university.ac.th ถึง student30@university.ac.th
- เลขนิสิต: 6510000001 ถึง 6510000030
- GPA: 3.55 - 3.95

## ข้อมูลตัวอย่าง

### นิสิตที่ได้รับรางวัลแล้ว (บางส่วน)
- **สมชาย ใจดี** (6510000001) - คณิตศาสตร์, GPA 3.85
- **นันทิดา ฉลาด** (6510000004) - ฟิสิกส์, GPA 3.90
- **ชยพล ก้าวหน้า** (6510000006) - วิศวกรรมคอมพิวเตอร์, GPA 3.95
- **ธัญญารัตน์ ชาญฉลาด** (6510000013) - วิทยาการคอมพิวเตอร์, GPA 3.92
- **จิราภรณ์ ใฝ่ดี** (6510000026) - ครุศาสตร์, GPA 3.91

## Structure
```
cmd/
  seed/
    main.go          # Seed script หลัก
  migration/
    main.go          # Migration script
```

## Dependencies
- GORM (database ORM)
- UUID (สำหรับสร้าง primary keys)
- Zerolog (logging)

## Notes
- ข้อมูลทั้งหมดเป็น Mock Data เท่านั้น
- Password ยังไม่มีในระบบ (ใช้ Google OAuth)
- ข้อมูล AdditionalData ถูกสร้างตาม FormStructure ของแต่ละ Award Category
