"use client";

import { useState, useEffect, useMemo, useCallback } from "react";
import { useRouter } from "next/navigation";
import { 
  Trophy, Calendar, Sparkles, GraduationCap, Loader2, ChevronLeft, Building2, Award, Search, BookOpen, Layers
} from "lucide-react";
import { useAuth } from "../../context/AuthContext"; 

// TYPES 
interface WinnerResponse {
  student_id: string;
  student_name: string;
  nisit_id: string;
  faculty_name: string;
  department_name: string;
  campus_name: string;
  gpa: number;
}

interface AwardCategoryGroupResponse {
  award_id: number;
  award_name: string;
  winners: WinnerResponse[];
}

interface AcademicYearGroupResponse {
  academic_year: number;
  semester: string;
  categories: AwardCategoryGroupResponse[];
}

export default function HallOfFamePage() {
  const router = useRouter();
  const { user } = useAuth(); 

  const [loading, setLoading] = useState(true);
  

  const [rawHonorRollData, setRawHonorRollData] = useState<AcademicYearGroupResponse[]>([]);
  
  // Selected Filters
  const [selectedYear, setSelectedYear] = useState<string>("");
  const [selectedSemester, setSelectedSemester] = useState<string>("");
  const [selectedCampus, setSelectedCampus] = useState<string>("");
  const [selectedFaculty, setSelectedFaculty] = useState<string>(""); 
  const [selectedAward, setSelectedAward] = useState<string>("");
  const [searchName, setSearchName] = useState<string>("");

  const formatSemester = (sem: string) => {
    const s = sem.toLowerCase();
    if (s === "first" || s === "1") return "ต้น";
    if (s === "second" || s === "2") return "ปลาย";
    return sem; 
  };

  const fetchHonorRoll = useCallback(async () => {
    try {
      setLoading(true);
      
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/honor-roll`);
      if (!response.ok) throw new Error("Failed to fetch honor roll");

      const json = await response.json();
      const responseData = json.data || json; 
      
      if (Array.isArray(responseData)) {
        setRawHonorRollData(responseData);
      } else {
        setRawHonorRollData([]); 
      }
    } catch (error) {
      console.error("Failed to load honor roll data", error);
      setRawHonorRollData([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchHonorRoll();
  }, [fetchHonorRoll]);

  const filterOptions = useMemo(() => {
    const years = new Set<number>();
    const awardsMap = new Map<number, string>();
    const campuses = new Set<string>();
    const faculties = new Set<string>();

    rawHonorRollData.forEach(yearGroup => {
      years.add(yearGroup.academic_year);
      
      yearGroup.categories.forEach(cat => {
        awardsMap.set(cat.award_id, cat.award_name);
        
        cat.winners.forEach(winner => {
          if (winner.campus_name) campuses.add(winner.campus_name);
          if (winner.faculty_name) faculties.add(winner.faculty_name);
        });
      });
    });

    return {
      years: Array.from(years).sort((a, b) => b - a), // เรียงปีล่าสุดขึ้นก่อน
      awards: Array.from(awardsMap.entries()).map(([id, name]) => ({ id, name })),
      campuses: Array.from(campuses).sort(),
      faculties: Array.from(faculties).sort()
    };
  }, [rawHonorRollData]);

  const displayData = useMemo(() => {
    let filtered = rawHonorRollData;

    if (selectedYear) {
      filtered = filtered.filter(yg => yg.academic_year.toString() === selectedYear);
    }
    if (selectedSemester) {
      filtered = filtered.filter(yg => yg.semester.toLowerCase() === selectedSemester.toLowerCase());
    }

    const lowerSearch = searchName.toLowerCase().trim();

    return filtered.map(yearGroup => {
      
      let filteredCats = yearGroup.categories;

      if (selectedAward) {
        filteredCats = filteredCats.filter(cat => cat.award_id.toString() === selectedAward);
      }

      filteredCats = filteredCats.map(cat => {
        let matchedWinners = cat.winners;

        if (selectedCampus) {
          matchedWinners = matchedWinners.filter(w => w.campus_name === selectedCampus);
        }
        if (selectedFaculty) {
          matchedWinners = matchedWinners.filter(w => w.faculty_name === selectedFaculty);
        }
        if (lowerSearch) {
          matchedWinners = matchedWinners.filter(w => 
            w.student_name.toLowerCase().includes(lowerSearch) || 
            w.nisit_id.includes(lowerSearch)
          );
        }

        return { ...cat, winners: matchedWinners };
      }).filter(cat => cat.winners.length > 0);

      return { ...yearGroup, categories: filteredCats };
    }).filter(yearGroup => yearGroup.categories.length > 0); 
    
  }, [rawHonorRollData, selectedYear, selectedSemester, selectedCampus, selectedFaculty, selectedAward, searchName]);

  const handleGoBack = () => {
    if (!user) {
      router.replace("/");
      return;
    }
    const role = user?.role?.toUpperCase() || "";
    switch (role) {
      case "ADMIN":
        router.replace("/admin/app");
        break;
      case "COMMITTEE_CHAIR":
        router.replace("/staff-page/app/review");
        break;
      case "DEAN":
      case "HEAD_OF_DEPARTMENT":
      case "VICE_DEAN":
        router.replace("/approval");
        break;
      case "STUDENT":
      default:
        router.replace("/dashboard");
        break;
    }
  };

  return (
    <div className="flex min-h-screen bg-[#f8f7f4] font-sans text-slate-800">
      <main className="flex-1 p-4 md:p-8 lg:p-12 overflow-y-auto">
        <div className="max-w-6xl mx-auto">
          
          {/* ── Back Button ── */}
          <div className="mb-6">
            <button 
              onClick={handleGoBack}
              className="flex items-center gap-2 text-sm font-bold text-slate-500 hover:text-[#d4af37] hover:border-[#d4af37]/30 transition-all bg-white px-4 py-2.5 rounded-xl border border-[#ede9e0] shadow-sm w-fit"
            >
              <ChevronLeft size={18} /> 
              {user ? "กลับสู่หน้าหลัก" : "กลับหน้าเข้าสู่ระบบ"}
            </button>
          </div>

          {/* ── Header ── */}
          <div className="text-center mb-10 relative mt-4">
            <div className="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2 w-40 h-40 bg-yellow-200 blur-3xl rounded-full opacity-30 z-0"></div>
            <div className="relative z-10">
              <div className="w-16 h-16 bg-white border border-[#ede9e0] rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-sm text-[#d4af37]">
                <Trophy size={32} />
              </div>
              <h1 className="font-serif text-4xl lg:text-5xl font-bold text-[#1a1a2e] mb-3">
                ทำเนียบนิสิตดีเด่น
              </h1>
              <p className="text-slate-500 font-medium max-w-lg mx-auto">
                ประกาศรายชื่อผู้ที่ได้รับรางวัลเชิดชูเกียรติในด้านต่างๆ 
              </p>
            </div>
          </div>

          {/* Filter Card  */}
          <div className="bg-white p-6 rounded-3xl shadow-sm border border-[#ede9e0] mb-12">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">

              {/* ค้นหารายชื่อ */}
              <div>
                <label className="block text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1.5 ml-1">ค้นหารายชื่อ</label>
                <div className="relative">
                  <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                  <input 
                    type="text"
                    placeholder="ชื่อ หรือ รหัสนิสิต..."
                    value={searchName}
                    onChange={(e) => setSearchName(e.target.value)}
                    className="w-full pl-10 pr-4 py-3 bg-[#f8f7f4] border border-transparent rounded-xl text-sm font-medium focus:outline-none focus:border-[#d4af37]/40 focus:bg-white transition-colors"
                  />
                </div>
              </div>

              {/* ปีการศึกษา */}
              <div>
                <label className="block text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1.5 ml-1">ปีการศึกษา</label>
                <div className="relative">
                  <Calendar className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                  <select 
                    value={selectedYear}
                    onChange={(e) => setSelectedYear(e.target.value)}
                    className="w-full pl-10 pr-4 py-3 bg-[#f8f7f4] border border-transparent rounded-xl text-sm font-medium focus:outline-none focus:border-[#d4af37]/40 focus:bg-white transition-colors appearance-none cursor-pointer"
                  >
                    <option value="">ทุกปีการศึกษา</option>
                    {filterOptions.years.map(y => <option key={y} value={y}>{y}</option>)}
                  </select>
                </div>
              </div>

              {/* ภาคเรียน */}
              <div>
                <label className="block text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1.5 ml-1">ภาคเรียน</label>
                <div className="relative">
                  <BookOpen className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                  <select 
                    value={selectedSemester}
                    onChange={(e) => setSelectedSemester(e.target.value)}
                    className="w-full pl-10 pr-4 py-3 bg-[#f8f7f4] border border-transparent rounded-xl text-sm font-medium focus:outline-none focus:border-[#d4af37]/40 focus:bg-white transition-colors appearance-none cursor-pointer"
                  >
                    <option value="">ทุกเทอม</option>
                    <option value="first">เทอมต้น (First)</option>
                    <option value="second">เทอมปลาย (Second)</option>
                  </select>
                </div>
              </div>

              {/* วิทยาเขต */}
              <div>
                <label className="block text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1.5 ml-1">วิทยาเขต</label>
                <div className="relative">
                  <Building2 className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                  <select 
                    value={selectedCampus}
                    onChange={(e) => setSelectedCampus(e.target.value)}
                    className="w-full pl-10 pr-4 py-3 bg-[#f8f7f4] border border-transparent rounded-xl text-sm font-medium focus:outline-none focus:border-[#d4af37]/40 focus:bg-white transition-colors appearance-none cursor-pointer"
                  >
                    <option value="">ทุกวิทยาเขต</option>
                    {filterOptions.campuses.map(c => <option key={c} value={c}>{c}</option>)}
                  </select>
                </div>
              </div>

              {/* คณะ */}
              <div>
                <label className="block text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1.5 ml-1">คณะ</label>
                <div className="relative">
                  <Layers className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                  <select 
                    value={selectedFaculty}
                    onChange={(e) => setSelectedFaculty(e.target.value)}
                    className="w-full pl-10 pr-4 py-3 bg-[#f8f7f4] border border-transparent rounded-xl text-sm font-medium focus:outline-none focus:border-[#d4af37]/40 focus:bg-white transition-colors appearance-none cursor-pointer"
                  >
                    <option value="">ทุกคณะ</option>
                    {filterOptions.faculties.map(f => <option key={f} value={f}>{f}</option>)}
                  </select>
                </div>
              </div>

              {/* ประเภทรางวัล */}
              <div>
                <label className="block text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-1.5 ml-1">ประเภทรางวัล</label>
                <div className="relative">
                  <Award className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                  <select 
                    value={selectedAward}
                    onChange={(e) => setSelectedAward(e.target.value)}
                    className="w-full pl-10 pr-4 py-3 bg-[#f8f7f4] border border-transparent rounded-xl text-sm font-medium focus:outline-none focus:border-[#d4af37]/40 focus:bg-white transition-colors appearance-none cursor-pointer"
                  >
                    <option value="">ทุกรางวัล</option>
                    {filterOptions.awards.map(a => <option key={a.id} value={a.id}>{a.name}</option>)}
                  </select>
                </div>
              </div>

            </div>
          </div>

          {/* Loading State  */}
          {loading ? (
            <div className="py-20 flex flex-col items-center justify-center">
              <Loader2 className="animate-spin text-[#d4af37]" size={40} />
              <p className="mt-4 text-slate-500 font-medium">กำลังโหลดข้อมูล...</p>
            </div>
          ) : 

          /* Content  */
          displayData.length === 0 ? (
            <div className="text-center bg-white border border-dashed border-[#ede9e0] rounded-3xl py-20 px-6 shadow-sm">
                <div className="w-20 h-20 bg-slate-50 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Trophy size={32} className="text-slate-300" />
                </div>
                <p className="text-lg font-bold text-slate-700">ไม่พบข้อมูลที่ค้นหา</p>
                <p className="text-sm text-slate-500 mt-1">ลองปรับเปลี่ยนตัวกรอง หรือค้นหาด้วยคำอื่นดูอีกครั้ง</p>
            </div>
          ) : (
            <div className="space-y-16">
              {/* วนลูปแสดงทีละเทอม */}
              {displayData.map((termGroup) => (
                <div key={`${termGroup.academic_year}-${termGroup.semester}`} className="relative">
                  
                  {/* Divider เทอม */}
                  <div className="flex items-center gap-4 mb-8">
                    <div className="flex items-center gap-2 bg-[#1a1a2e] text-white px-5 py-2.5 rounded-full font-bold text-sm shadow-md">
                      <Calendar size={16} className="text-[#d4af37]" />
                      ปีการศึกษา {termGroup.academic_year} / เทอม{formatSemester(termGroup.semester)}
                    </div>
                    <div className="flex-1 h-px bg-gradient-to-r from-[#ede9e0] to-transparent"></div>
                  </div>

                  <div className="space-y-8 pl-4 md:pl-8 border-l-2 border-[#ede9e0]/50 ml-6 md:ml-10">
                    {termGroup.categories.map((category) => (
                      <div key={category.award_id} className="bg-white rounded-3xl border border-[#ede9e0] overflow-hidden shadow-sm hover:shadow-md transition-shadow relative">
                        
                        <div className="absolute -left-[19px] md:-left-[35px] top-8 w-3 h-3 bg-[#d4af37] rounded-full border-4 border-[#f8f7f4]"></div>

                        {/* Category Header */}
                        <div className="p-6 border-b border-[#ede9e0] bg-gradient-to-br from-white to-[#fcfcfc] flex flex-col md:flex-row md:items-center justify-between gap-4">
                          <div className="flex items-center gap-4">
                              <div className="w-12 h-12 bg-yellow-50 rounded-xl flex items-center justify-center text-yellow-600 shrink-0">
                                <Sparkles size={24} />
                              </div>
                              <div>
                                <h2 className="text-xl font-bold text-[#1a1a2e] leading-tight">{category.award_name}</h2>
                              </div>
                          </div>
                          <div className="bg-slate-50 px-4 py-1.5 rounded-full text-sm font-bold text-slate-600 border border-slate-200 self-start md:self-auto shrink-0">
                              {category.winners.length} รางวัล
                          </div>
                        </div>

                        {/* Students Grid */}
                        <div className="p-6">
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          {category.winners.map((student, index) => (
                              <div key={student.student_id} className="flex items-center gap-4 p-4 rounded-2xl bg-[#faf9f7] border border-[#ede9e0]/60 hover:border-[#d4af37]/30 transition-colors relative group">
                              
                              {/* อันดับ / เหรียญ */}
                              <div className="w-12 h-12 rounded-full flex items-center justify-center font-black text-lg shadow-sm shrink-0 z-10 transition-colors bg-white text-[#1a1a2e] border border-[#ede9e0]">
                                    {index + 1}
                                </div>
                              
                              <div className="z-10 min-w-0">
                                  <p className="font-bold text-[#1a1a2e] text-base leading-tight truncate">
                                    {student.student_name}
                                  </p>
                                  <p className="text-xs font-medium text-slate-500 mt-1 flex items-center gap-1.5 truncate">
                                    <GraduationCap size={14} className="text-slate-400 shrink-0"/>
                                    {student.faculty_name || 'ไม่ระบุคณะ'}
                                  </p>
                                  <p className="text-[10px] text-slate-400 font-medium mt-0.5 ml-5 truncate">
                                    {student.campus_name}
                                  </p>
                              </div>
                              </div>
                          ))}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>

                </div>
              ))}
            </div>
          )}

        </div>
      </main>
    </div>
  );
}