import { FileText, User, BookOpen } from "lucide-react";

export default function Step5Review({ formData, awardType, schema, isRepeatable }: any) {
  
  const renderSpecificDetails = () => {
    if (!formData.dynamicAnswers || Object.keys(formData.dynamicAnswers).length === 0) {
      return <p className="text-gray-500 italic">ไม่มีรายละเอียดเพิ่มเติม</p>;
    }

    // กรณี(Repeatable)
    if (isRepeatable && Array.isArray(formData.dynamicAnswers)) {
      return (
        <div className="space-y-4">
          {formData.dynamicAnswers.map((item: any, index: number) => (
            <div key={index} className="p-5 bg-gray-50 rounded-xl border border-gray-100">
              <p className="text-xs font-bold text-gray-500 mb-3 border-b border-gray-200 pb-2">Item #{index + 1}</p>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {Object.entries(item).map(([key, value]) => {
                  const fieldLabel = schema.find((f: any) => f.id === key)?.label || key;
                  return (
                    <div key={key}>
                      <p className="text-[10px] text-gray-400 uppercase font-black tracking-widest mb-1">{fieldLabel}</p>
                      <p className="text-gray-900 font-bold text-sm break-words">{String(value)}</p>
                    </div>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      );
    }

    // กรณี non repeat
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {Object.entries(formData.dynamicAnswers).map(([key, value]) => {
           const fieldLabel = schema.find((f: any) => f.id === key)?.label || key;
           return (
             <div key={key} className="p-4 bg-gray-50 rounded-xl border border-gray-100">
                <p className="text-[10px] text-gray-400 uppercase font-black tracking-widest mb-1">{fieldLabel}</p>
                <p className="text-gray-900 font-bold text-sm break-words">
                   {String(value)}
                </p>
             </div>
           );
        })}
      </div>
    );
  };

  return (
    <div className="space-y-8">
      <div className="bg-green-50 p-6 rounded-2xl border border-green-200 text-center shadow-sm">
        <h2 className="text-xl font-black text-green-800 uppercase tracking-tight">ยืนยันการส่งข้อมูล</h2>
        <p className="text-green-700 text-sm mt-1 font-medium">กรุณาตรวจสอบข้อมูลของคุณให้ถูกต้องก่อนส่ง</p>
      </div>

      <div className="bg-white p-8 rounded-[2rem] border border-gray-100 shadow-sm space-y-10">
        
        {/* Personal Info & Snapshot  */}
        <div>
           <h3 className="text-lg font-bold text-gray-800 border-b pb-3 mb-6 flex items-center gap-2">
             <User size={20} className="text-green-600"/> 1. ข้อมูลส่วนตัว
           </h3>
           
           <div className="bg-gray-50 p-6 rounded-2xl border border-gray-100 space-y-6">

              <div className="grid grid-cols-1 md:grid-cols-2 gap-y-5 gap-x-8 text-sm">
                <div><p className="text-[10px] text-gray-400 font-black tracking-widest uppercase mb-1">Full Name</p><p className="text-gray-900 font-bold">{formData.personal.firstName} {formData.personal.lastName}</p></div>
                <div><p className="text-[10px] text-gray-400 font-black tracking-widest uppercase mb-1">Student ID</p><p className="text-gray-900 font-bold">{formData.personal.studentId}</p></div>
                <div><p className="text-[10px] text-gray-400 font-black tracking-widest uppercase mb-1">Faculty & Major</p><p className="text-gray-900 font-bold">{formData.personal.faculty} <br/><span className="text-gray-500 font-medium">{formData.personal.major}</span></p></div>
                <div><p className="text-[10px] text-gray-400 font-black tracking-widest uppercase mb-1">Contact</p><p className="text-gray-900 font-bold">{formData.personal.phone || "-"}<br/><span className="text-gray-500 font-medium">{formData.personal.email}</span></p></div>
                  <div><p className="text-[10px] text-gray-400 font-black tracking-widest uppercase mb-1">Study Year</p><p className="text-gray-900 font-bold">{formData.personal.studyYear ? `Year ${formData.personal.studyYear}` : "-"}</p></div>
                  <div><p className="text-[10px] text-gray-400 font-black tracking-widest uppercase mb-1">Cumulative GPA</p><p className="text-gray-900 font-bold">{formData.personal.gpa || "-"}</p></div>
                  <div><p className="text-[10px] text-gray-400 font-black tracking-widest uppercase mb-1">Academic Advisor</p><p className="text-gray-900 font-bold">{formData.personal.advisor || "-"}</p></div>
                  <div>
                    <p className="text-[10px] text-gray-400 font-black tracking-widest uppercase mb-1">Date of Birth</p>
                    <p className="text-gray-900 font-bold">{formData.personal.dateOfBirth ? new Date(formData.personal.dateOfBirth).toLocaleDateString() : "-"}</p>
                  </div>
                  <div className="md:col-span-2">
                    <p className="text-[10px] text-gray-400 font-black tracking-widest uppercase mb-1">Current Address</p>
                    <p className="text-gray-900 font-bold">{formData.personal.address || "-"}</p>
                  </div>
                
              </div>
           </div>
        </div>

        {/*  Specific Details (Dynamic)  */}
        <div>
           <h3 className="text-lg font-bold text-gray-800 border-b pb-3 mb-6 flex items-center gap-2">
             <BookOpen size={20} className="text-green-600"/> 2. รายละเอียดเฉพาะ ({awardType})
           </h3>
           {renderSpecificDetails()}
        </div>

        {/* Documents  */}
        <div>
            <h3 className="text-lg font-bold text-gray-800 border-b pb-3 mb-6 flex items-center gap-2">
              <FileText size={20} className="text-green-600"/> 3. เอกสารประกอบ ({formData.documents?.length || 0})
            </h3>
            {formData.documents && formData.documents.length > 0 ? (
                <ul className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    {formData.documents.map((doc: any, i: number) => (
                        <li key={i} className="flex items-center gap-3 text-sm text-gray-700 bg-gray-50 p-4 rounded-xl border border-gray-100">
                            <div className="w-10 h-10 bg-white rounded-lg flex items-center justify-center text-red-500 shadow-sm"><FileText size={18}/></div>
                            <div className="overflow-hidden">
                              <p className="truncate font-bold text-gray-800">{doc.name}</p>
                              <p className="text-[10px] text-gray-400 font-black tracking-widest uppercase">{doc.size}</p>
                            </div>
                        </li>
                    ))}
                </ul>
            ) : (
                <p className="text-gray-500 italic text-sm">No documents attached.</p>
            )}
        </div>

        {/*  Statement  */}
        <div>
           <h3 className="text-lg font-bold text-gray-800 border-b pb-3 mb-6 flex items-center gap-2">
             <FileText size={20} className="text-green-600"/> 4. คำแถลงส่วนตัว
           </h3>
           <p className="font-bold text-sm mb-3 text-gray-900 bg-gray-50 inline-block px-4 py-1.5 rounded-lg">{formData.statement.topic || "No Topic"}</p>
           <div className="text-sm text-gray-700 leading-relaxed p-6 bg-blue-50/30 rounded-2xl whitespace-pre-wrap border border-blue-100/50">
              {formData.statement.content || <span className="text-gray-400 italic">No content provided.</span>}
           </div>
        </div>

      </div>
    </div>
  );
}