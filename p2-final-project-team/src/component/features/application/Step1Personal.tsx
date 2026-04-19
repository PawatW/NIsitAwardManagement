import { User, Mail, Phone, Book, School, Hash, GraduationCap, TrendingUp, UserCheck, MapPin, Calendar } from "lucide-react";

export default function Step1Personal({ data, updateData }: any) {
  const handleChange = (field: string, value: string) => {
    updateData({ ...data, [field]: value });
  };

  return (
    <div className="bg-white p-8 rounded-xl border border-gray-200 shadow-sm space-y-8">
      
      {/*  ข้อมูล Read-only  */}
      <section>
        <div className="border-b pb-4 mb-4">
          <h2 className="text-xl font-bold text-gray-800">1. ข้อมูลส่วนตัว</h2>
          <p className="text-sm text-gray-500">ข้อมูลส่วนตัวถูกดึงมาจากระบบทะเบียน หากข้อมูลไม่ถูกต้องกรุณาติดต่อเจ้าหน้าที่</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-700">ชื่อ</label>
            <div className="relative">
              <User className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input type="text" value={data.firstName || ""} disabled className="w-full pl-10 pr-4 py-2 bg-gray-100 border border-gray-300 rounded-lg text-gray-500 cursor-not-allowed" />
            </div>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-700">นามสกุล</label>
            <div className="relative">
              <User className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input type="text" value={data.lastName || ""} disabled className="w-full pl-10 pr-4 py-2 bg-gray-100 border border-gray-300 rounded-lg text-gray-500 cursor-not-allowed" />
            </div>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-700">รหัสประจำตัวนักศึกษา</label>
            <div className="relative">
              <Hash className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input type="text" value={data.studentId || ""} disabled className="w-full pl-10 pr-4 py-2 bg-gray-100 border border-gray-300 rounded-lg text-gray-500 cursor-not-allowed" />
            </div>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-700">Email</label>
            <div className="relative">
              <Mail className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input type="email" value={data.email || ""} disabled className="w-full pl-10 pr-4 py-2 bg-gray-100 border border-gray-300 rounded-lg text-gray-500 cursor-not-allowed" />
            </div>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-700">คณะ</label>
            <div className="relative">
              <School className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input type="text" value={data.faculty || "-"} disabled className="w-full pl-10 pr-4 py-2 bg-gray-100 border border-gray-300 rounded-lg text-gray-500 cursor-not-allowed" />
            </div>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-700">สาขา</label>
            <div className="relative">
              <Book className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input type="text" value={data.major || "-"} disabled className="w-full pl-10 pr-4 py-2 bg-gray-100 border border-gray-300 rounded-lg text-gray-500 cursor-not-allowed" />
            </div>
          </div>
        </div>
      </section>

      {/* Editable Snapshot Data  */}
      <section>
        <div className="border-b pb-4 mb-4 mt-6">
          <h2 className="text-lg font-bold text-gray-800">สถานะการศึกษาปัจจุบัน</h2>
          <p className="text-sm text-gray-500">กรุณาระบุข้อมูล ณ วันที่ยื่นใบสมัครให้ถูกต้อง</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          
          <div className="space-y-2">
            <label className="text-sm font-bold text-gray-700">Year of Study (ชั้นปี) <span className="text-red-500">*</span></label>
            <div className="relative">
              <GraduationCap className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <select
                value={data.studyYear || ""}
                onChange={(e) => handleChange("studyYear", e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-white border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-green-500 outline-none"
              >
                <option value="">Select Year</option>
                <option value="1">Year 1</option>
                <option value="2">Year 2</option>
                <option value="3">Year 3</option>
                <option value="4">Year 4</option>
                <option value="5">Year 5+</option>
              </select>
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-bold text-gray-700">GPA สะสม <span className="text-red-500">*</span></label>
            <div className="relative">
              <TrendingUp className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input
                type="text"  min="0" max="4" placeholder="e.g. 3.50"
                value={data.gpa || ""}
                onChange={(e) => handleChange("gpa", e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-white border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-green-500 outline-none"
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-bold text-gray-700">ชื่อ-นามสกุล อาจารย์ที่ปรึกษา <span className="text-red-500">*</span></label>
            <div className="relative">
              <UserCheck className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input
                type="text" placeholder="ชื่อ-นามสกุล อาจารย์ที่ปรึกษา"
                value={data.advisor || ""}
                onChange={(e) => handleChange("advisor", e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-white border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-green-500 outline-none"
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-bold text-gray-700">เบอร์โทรศัพท์ <span className="text-red-500">*</span></label>
            <div className="relative">
              <Phone className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input
                type="text" placeholder="08x-xxx-xxxx"
                value={data.phone || ""}
                onChange={(e) => handleChange("phone", e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-white border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-green-500 outline-none"
              />
            </div>
          </div>


          <div className="space-y-2">
            <label className="text-sm font-bold text-gray-700">วันเดือนปีเกิด</label>
            <div className="relative">
              <Calendar className="absolute left-3 top-2.5 text-gray-400" size={18} />
              <input
                type="date"
                value={data.dateOfBirth || ""}
                onChange={(e) => handleChange("dateOfBirth", e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-white border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-green-500 outline-none"
              />
            </div>
          </div>


          <div className="space-y-2 md:col-span-2">
            <label className="text-sm font-bold text-gray-700">ที่อยู่ปัจจุบัน</label>
            <div className="relative">
              <MapPin className="absolute left-3 top-3 text-gray-400" size={18} />
              <textarea
                rows={3} placeholder="ที่อยู่ปัจจุบันที่สามารถติดต่อได้"
                value={data.address || ""}
                onChange={(e) => handleChange("address", e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-white border border-gray-300 rounded-lg text-gray-900 focus:ring-2 focus:ring-green-500 outline-none resize-none"
              />
            </div>
          </div>

        </div>
      </section>
    </div>
  );
}