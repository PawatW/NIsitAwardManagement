"use client";

import { useState, useEffect, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { CheckCircle2, ChevronRight, ChevronLeft, Loader2, AlertTriangle } from "lucide-react";
import { useAuth } from "../../../../context/AuthContext"; 
import api from "../../../../lib/axios"; 

import Step1Personal from "../../../../component/features/application/Step1Personal";
import Step2Dynamic from "../../../../component/features/application/Step2Activities"; 
import Step3Documents from "../../../../component/features/application/Step3Documents";
import Step4Statement from "../../../../component/features/application/Step4Statement";
import Step5Review from "../../../../component/features/application/Step5Review";

function ApplicationWizardContent() {
  const { user } = useAuth();
  const searchParams = useSearchParams();
  const router = useRouter();
  
  const awardId = searchParams.get('award_id'); 

  const [currentStep, setCurrentStep] = useState(1);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [configError, setConfigError] = useState<string | null>(null);

  const [awardConfig, setAwardConfig] = useState<any | null>(null);
  const [activeTermId, setActiveTermId] = useState<number | null>(null);
  
  const [isRepeatable, setIsRepeatable] = useState(false);

  const [formData, setFormData] = useState({
    personal: {
      firstName: "", lastName: "", studentId: "", phone: "", email: "",
      faculty: "", major: "", facultyId: "", majorId: "",
      studyYear: "", 
      gpa: "",       
      advisor: "",
      dateOfBirth: "",
      address: ""
    },
    dynamicAnswers: {} as any, 
    documents: [] as any[], 
    statement: { topic: "Personal Achievement & Vision", content: "" }
  });

  useEffect(() => {
    const initializePage = async () => {
      if (!user || !awardId) {
        if (!awardId) setConfigError("Missing Award ID.");
        setIsLoading(false);
        return;
      }

      try {
        setIsLoading(true);
        
        const [userRes, awardRes, termRes] = await Promise.all([
          api.get(`/api/v1/user/me`),
          api.get(`/api/v1/award-categories/${awardId}`),
          api.get(`/api/v1/academic-terms`) 
        ]);

        const u = userRes.data.data || userRes.data; 
        const a = awardRes.data.data || awardRes.data; 
        const termsArray = termRes.data.data || termRes.data || []; 

        const formFields = Array.isArray(a.form_structure) ? a.form_structure : [];
        const checkIsRepeatable = a.is_repeatable === true;
        
        setIsRepeatable(checkIsRepeatable);
        setAwardConfig({ ...a, form_structure: formFields });
        
        const initialDynamicData: any = {};
        formFields.forEach((field: any) => {
          if (field.type === "checkbox") {
            initialDynamicData[field.id] = false;
          } else {
            initialDynamicData[field.id] = "";
          }
        });

        setFormData(prev => ({
          ...prev,
          personal: {
            ...prev.personal, 
            firstName: u.first_name || "",
            lastName: u.last_name || "",
            studentId: u.nisit_id || "-",
            email: u.email || "",
            phone: u.phone_number || "-",
            faculty: u.faculty?.name || "-",
            major: u.department?.name || "-",
            facultyId: u.faculty?.id || "",
            majorId: u.department?.id || ""
          },
          dynamicAnswers: initialDynamicData 
        }));

        const currentTerm = termsArray.find((t: any) => t.is_open === true) || termsArray[0];
        setActiveTermId(currentTerm?.id || 1);

      } catch (error) {
        console.error("Initialization Error:", error);
        setConfigError("Failed to load application configuration. Please check your connection.");
      } finally {
        setIsLoading(false);
      }
    };

    initializePage();
  }, [user, awardId]);

  const validateCurrentStep = () => {
    if (currentStep === 1) {
      const { gpa, studyYear } = formData.personal;
      if (!gpa || !studyYear) {
        alert("กรุณากรอกข้อมูลให้ครบถ้วน");
        return false;
      }
      return true;
    }

    if (currentStep === 2) {
      const schemaFields = getFieldsArray();
      
      if (isRepeatable) {
        const items = Array.isArray(formData.dynamicAnswers) ? formData.dynamicAnswers : [];
        if (items.length === 0) return true; 
        
        for (let i = 0; i < items.length; i++) {
          for (const field of schemaFields) {
            if (field.required && (!items[i][field.id] || items[i][field.id].toString().trim() === "")) {
              alert(`กรุณากรอกข้อมูลให้ครบถ้วน`);
              return false; 
            }
          }
        }
      } else {
        const answers = formData.dynamicAnswers || {};
        for (const field of schemaFields) {
          if (field.required && (!answers[field.id] || answers[field.id].toString().trim() === "")) {
            alert(`กรุณากรอก ${field.label} ให้ครบถ้วน`);
            return false; 
          }
        }
      }
      return true;
    }

    if (currentStep === 3 || currentStep === 4) {
      return true; 
    }

    return true; 
  };

  const handleSubmit = async () => {
      if (!confirm("คุณแน่ใจหรือไม่ว่าต้องการส่งใบสมัครนี้?")) return;

      try {
        setIsSubmitting(true);

        const processedAnswers = formData.dynamicAnswers;

        const requestPayload = {
          award_id: parseInt(awardId!),
          academic_term_id: activeTermId,
          additional_data: {
            ...(isRepeatable ? { items: processedAnswers } : processedAnswers),
            statement_content: formData.statement.content 
          },
          snapshot_study_year: formData.personal.studyYear
          ? parseInt(formData.personal.studyYear): null,

        snapshot_gpa: formData.personal.gpa
          ? parseFloat(formData.personal.gpa) : null,
          snapshot_advisor: formData.personal.advisor || "Not Specified",
          snapshot_phone: formData.personal.phone || null,
          snapshot_date_of_birth: formData.personal.dateOfBirth || null,
          snapshot_address: formData.personal.address || null
        };

        const createRes = await api.post("/api/v1/requests", requestPayload);
        const requestId = createRes.data.id;

        if (formData.documents.length > 0) {
          for (const doc of formData.documents) {
            if (doc.fileObj) {
              const fileData = new FormData();
              fileData.append("file", doc.fileObj);
              fileData.append("description", doc.name);
              await api.post(`/api/v1/requests/${requestId}/documents`, fileData, {
                headers: { "Content-Type": "multipart/form-data" }
              });
            }
          }
        }

        alert("ใบสมัครถูกส่งเรียบร้อยแล้ว!");
        router.push("/dashboard");

      } catch (error: any) {
        alert(`Error: ${error.response?.data?.message || "Submission failed"}`);
      } finally {
        setIsSubmitting(false);
      }
  };

  const updateFormData = (section: string, data: any) => {
    setFormData((prev) => ({ ...prev, [section]: data }));
  };

  const getFieldsArray = () => {
    if (!awardConfig?.form_structure) return [];
    if (Array.isArray(awardConfig.form_structure)) {
      return awardConfig.form_structure;
    }
    return awardConfig.form_structure.fields || [];
  };

  if (isLoading) return <div className="flex h-screen items-center justify-center"><Loader2 className="animate-spin text-green-600" size={48} /></div>;
  if (configError) return <div className="flex h-screen items-center justify-center text-red-500 flex-col gap-4"><AlertTriangle size={48}/><p>{configError}</p></div>;

  return (
    <div className="max-w-5xl mx-auto pb-20 px-4 pt-8">
      <button 
        onClick={() => router.push("/dashboard")}
        className="flex items-center text-gray-500 hover:text-gray-900 transition-colors font-bold text-sm mb-6 w-fit"
      >
        <ChevronLeft size={20} className="mr-1" />
        กลับสู่หน้าหลัก
      </button>

      <div className="mb-10">
        <h1 className="text-3xl font-black text-gray-900 uppercase tracking-tighter">ส่งใบสมัคร</h1>
        <p className="text-gray-500 font-bold mt-1 uppercase text-xs tracking-[0.2em]">
          หมวดหมู่: <span className="text-green-600">{awardConfig?.name}</span>
        </p>
      </div>

      <div className="grid grid-cols-1 gap-8">
        {currentStep === 1 && (
          <Step1Personal 
            data={formData.personal} 
            updateData={(d: any) => updateFormData('personal', d)} 
          />
        )}
        
        {currentStep === 2 && (
          <Step2Dynamic 
            schema={getFieldsArray()} 
            isRepeatable={isRepeatable}
            title={isRepeatable ? "Add Your Achievements" : "ข้อมูลเพิ่มเติม"}
            data={formData.dynamicAnswers} 
            updateData={(d: any) => updateFormData('dynamicAnswers', d)} 
          />
        )}
        
        {currentStep === 3 && <Step3Documents data={formData.documents} updateData={(d: any) => updateFormData('documents', d)} />}
        {currentStep === 4 && <Step4Statement data={formData.statement} updateData={(d: any) => updateFormData('statement', d)} />}
        
        {currentStep === 5 && (
            <Step5Review 
                formData={formData} 
                awardType={awardConfig?.name} 
                schema={getFieldsArray()} 
                isRepeatable={isRepeatable} 
            />
        )}
      </div>

      <div className="flex justify-between mt-12 pt-8 border-t border-gray-100">
        <button onClick={() => setCurrentStep(prev => prev - 1)} className={`px-8 py-4 rounded-2xl font-bold uppercase tracking-widest text-xs transition-all ${currentStep === 1 ? 'invisible' : 'bg-gray-100 text-gray-400 hover:bg-gray-200'}`}>Back</button>
        {currentStep < 5 ? (
          <button 
            onClick={() => {
              if (validateCurrentStep()) {
                setCurrentStep(prev => prev + 1);
              }
            }} 
            className="px-10 py-4 bg-gray-900 text-white rounded-2xl font-bold uppercase tracking-widest text-xs shadow-xl shadow-gray-200 hover:bg-green-600 transition-all"
          >
            Continue
          </button>
        ) : (
          <button onClick={handleSubmit} disabled={isSubmitting} className="px-10 py-4 bg-green-600 text-white rounded-2xl font-bold uppercase tracking-widest text-xs shadow-xl shadow-green-100 hover:bg-green-700 transition-all disabled:bg-gray-300">{isSubmitting ? "Submitting..." : "Confirm & Submit"}</button>
        )}
      </div>
    </div>
  );
}

export default function ApplicationWizardPage() {
    return (
        <Suspense fallback={<Loader2 className="animate-spin text-green-600 mx-auto mt-20" size={40} />}>
            <ApplicationWizardContent />
        </Suspense>
    );
}