// components/application-form/Step4Statement.tsx
export default function Step4Statement({ data, updateData }: any) {
  const wordCount = data.content.trim().split(/\s+/).length;

  return (
    <div className="bg-white p-8 rounded-xl border border-gray-200 shadow-sm space-y-6">
      <h2 className="text-xl font-bold text-gray-800 border-b pb-4">4. คำแถลงส่วนตัว</h2>
      
      <div>
        <label className="block text-sm font-bold text-gray-700 mb-2">
          กรุณาเขียนคำแถลงส่วนตัวที่อธิบายถึงความสำเร็จและวิสัยทัศน์ของคุณ
        </label>
   
        
        <textarea
          value={data.content}
          onChange={(e) => updateData({ ...data, content: e.target.value })}
          className="w-full h-64 p-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 outline-none resize-none leading-relaxed"
          placeholder="Start typing your statement here..."
        ></textarea>
        
        <div className="flex justify-end mt-2">
          <span className={`text-xs font-medium ${wordCount > 500 ? 'text-red-500' : 'text-gray-400'}`}>
            {data.content === "" ? 0 : wordCount} / 500 คำ
          </span>
        </div>
      </div>
    </div>
  );
}