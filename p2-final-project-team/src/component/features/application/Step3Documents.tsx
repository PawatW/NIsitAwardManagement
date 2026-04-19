// components/application-form/Step3Documents.tsx
import { useState, useRef, ChangeEvent, DragEvent, useEffect } from "react";
import { Upload, FileText, X, ExternalLink, Eye } from "lucide-react";

export default function Step3Documents({ data, updateData }: any) {

  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isDragging, setIsDragging] = useState(false);

  // เปิดหน้าต่างเลือกไฟล์ 
  const triggerFileInput = () => {
    fileInputRef.current?.click();
  };


  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      processFiles(files);
    }

    if (fileInputRef.current) {
        fileInputRef.current.value = ""; 
    }
  };

  const processFiles = (files: FileList) => {
    const newFiles = Array.from(files).map((file) => ({
      name: file.name,
      size: formatFileSize(file.size),
      type: file.type || "Unknown", 
      fileObj: file, 

      previewUrl: URL.createObjectURL(file) 
    }));

    updateData([...data, ...newFiles]);
  };

  // Drag & Drop Handlers 
  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
    const files = e.dataTransfer.files;
    if (files && files.length > 0) {
      processFiles(files);
    }
  };


  const removeFile = (index: number) => {
    const newData = [...data];

    if (newData[index].previewUrl) {
      URL.revokeObjectURL(newData[index].previewUrl);
    }
    newData.splice(index, 1);
    updateData(newData);
  };



  return (
    <div className="bg-white p-8 rounded-xl border border-gray-200 shadow-sm space-y-6">
      <h2 className="text-xl font-bold text-gray-800 border-b pb-4">3. เอกสารประกอบที่เกี่ยวข้อง</h2>
      
      <div className="bg-blue-50 p-4 rounded-lg text-sm text-blue-700 mb-4">
        กรุณาอัปโหลดเอกสารที่เกี่ยวข้อง เช่น ใบรับรองผลการเรียน, จดหมายรับรอง, หรือผลงานที่เกี่ยวข้องกับรางวัลนี้
      </div>

      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        className="hidden"
        multiple 
        accept=".pdf,.jpg,.jpeg,.png" 
      />

      {/* Upload Area  */}
      <div 
        onClick={triggerFileInput}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        className={`border-2 border-dashed rounded-xl p-10 text-center cursor-pointer transition-all ${
            isDragging 
            ? "border-green-500 bg-green-50" 
            : "border-gray-300 hover:bg-gray-50 hover:border-green-400"
        }`}
      >
        <div className="w-16 h-16 bg-green-100 text-green-600 rounded-full flex items-center justify-center mx-auto mb-4">
          <Upload size={32} />
        </div>
        <p className="font-medium text-gray-700">
            {isDragging ? "Drop files here" : "Click to upload or drag and drop"}
        </p>
        <p className="text-xs text-gray-400 mt-2">PDF, JPG, PNG (Max 10MB)</p>
      </div>

      {/* File List  */}
      {data.length > 0 && (
        <div className="space-y-3">
            {data.map((file: any, index: number) => (
            <div key={index} className="flex items-center justify-between p-4 bg-gray-50 border border-gray-200 rounded-lg animate-in fade-in slide-in-from-bottom-2 group hover:bg-white hover:shadow-sm transition-all">
                <div className="flex items-center gap-4 overflow-hidden flex-1">
                  <div className="w-10 h-10 bg-red-100 text-red-500 rounded-lg flex items-center justify-center flex-shrink-0">
                      <FileText size={20} />
                  </div>
                  <div className="min-w-0 flex-1">
                      <p className="text-sm font-medium text-gray-800 truncate">{file.name}</p>
                      <p className="text-xs text-gray-500">{file.size} • {file.type.split('/')[1]?.toUpperCase() || 'FILE'}</p>
                  </div>
                </div>

                <div className="flex items-center gap-2">

                  {file.previewUrl && (
                    <a 
                      href={file.previewUrl} 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-gray-400 hover:text-blue-500 p-2 rounded-lg hover:bg-blue-50 flex items-center gap-1 transition-colors"
                      title="ดูตัวอย่างไฟล์"
                      onClick={(e) => e.stopPropagation()} // กันการกระตุ้นปุ่มอื่น
                    >
                      <Eye size={18} />
                      <span className="text-xs font-medium hidden sm:inline">ดูไฟล์</span>
                    </a>
                  )}

                  {/* ปุ่มลบไฟล์ */}
                  <button 
                    onClick={(e) => { e.stopPropagation(); removeFile(index); }} 
                    className="text-gray-400 hover:text-red-500 p-2 rounded-lg hover:bg-red-50 transition-colors"
                    title="ลบไฟล์นี้"
                  >
                    <X size={20} />
                  </button>
                </div>
            </div>
            ))}
        </div>
      )}
    </div>
  );
}