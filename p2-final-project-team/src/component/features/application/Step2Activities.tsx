import { Plus, Trash2 } from "lucide-react";

interface Field {
  id: string;
  label: string;
  type: "text" | "number" | "textarea" | "select";
  options?: string[];
  required?: boolean;
  placeholder?: string;
}

interface Step2Props {
  schema: Field[];
  data: any;
  updateData: (data: any) => void;
  isRepeatable?: boolean;
  title?: string; 
}

const renderFieldInput = (field: Field, value: any, onChange: (val: any) => void) => {
  const commonClass = "w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 outline-none transition-all";
  
  switch (field.type) {
    case "textarea":
      return (
        <textarea
          rows={4}
          className={commonClass}
          placeholder={field.placeholder}
          value={value || ""}
          onChange={(e) => onChange(e.target.value)}
        />
      );
    case "select":
      return (
        <select
          className={commonClass}
          value={value || ""}
          onChange={(e) => onChange(e.target.value)}
        >
          <option value="">Select {field.label}</option>
          {field.options?.map((opt) => (
            <option key={opt} value={opt}>{opt}</option>
          ))}
        </select>
      );
    case "number":
      return (
        <input
          type="number"
          className={commonClass}
          placeholder={field.placeholder}
          value={value || ""}
          onChange={(e) => onChange(e.target.value)}
        />
      );
    default:
      return (
        <input
          type="text"
          className={commonClass}
          placeholder={field.placeholder}
          value={value || ""}
          onChange={(e) => onChange(e.target.value)}
        />
      );
  }
};

export default function Step2Dynamic({ 
  schema, 
  data = {}, 
  updateData, 
  isRepeatable = false,
  title = "Details" 
}: Step2Props) {

  const handleChange = (id: string, value: any) => {
    updateData({ ...data, [id]: value });
  };

  const handleArrayChange = (index: number, id: string, value: any) => {

    const newData = [...(Array.isArray(data) ? data : [])];
    
    if (!newData[index]) newData[index] = {};
    
    newData[index][id] = value;
    updateData(newData);
  };

  const addItem = () => {
    const newData = [...(Array.isArray(data) ? data : []), {}];
    updateData(newData);
  };

  const removeItem = (index: number) => {
    const newData = [...(Array.isArray(data) ? data : [])];
    newData.splice(index, 1);
    updateData(newData);
  };


  //  Repeater (List)

  if (isRepeatable) {
    const items = Array.isArray(data) && data.length > 0 ? data : [{}];
    
    const singularTitle = title.endsWith('s') ? title.slice(0, -1) : title;

    return (
      <div className="space-y-6">
        <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
           <h2 className="text-xl font-bold text-gray-800 mb-2">{title}</h2>
           <p className="text-gray-500 text-sm">Please provide information for your {singularTitle.toLowerCase()}.</p>
        </div>

        {items.map((item: any, index: number) => (
          <div key={index} className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm relative group">
            <div className="flex justify-between items-center mb-4">
               <h3 className="text-sm font-bold text-gray-600 uppercase tracking-wider">
                 {singularTitle} #{index + 1}
               </h3>
               {items.length > 1 && (
                 <button 
                   onClick={() => removeItem(index)}
                   className="text-gray-400 hover:text-red-500 transition-colors"
                 >
                   <Trash2 size={20} />
                 </button>
               )}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {schema.map((field) => (
                <div key={field.id} className={field.type === 'textarea' || field.type === 'text' ? "col-span-1 md:col-span-2" : ""}>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {field.label} {field.required && <span className="text-red-500">*</span>}
                  </label>
                  {renderFieldInput(
                    field, 
                    item[field.id], 
                    (val) => handleArrayChange(index, field.id, val)
                  )}
                </div>
              ))}
            </div>
          </div>
        ))}

        <button
          onClick={addItem}
          className="w-full py-4 border-2 border-dashed border-green-300 rounded-xl flex items-center justify-center gap-2 text-green-600 font-bold hover:bg-green-50 hover:border-green-500 transition-all"
        >

          <Plus size={20} /> Add Another {singularTitle}
        </button>
      </div>
    );
  }



  return (
    <div className="bg-white p-8 rounded-xl border border-gray-200 shadow-sm space-y-6">

      <h2 className="text-xl font-bold text-gray-800 border-b pb-4">{title}</h2>
      <div className="grid grid-cols-1 gap-6">
        {schema.map((field) => (
          <div key={field.id}>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {field.label} {field.required && <span className="text-red-500">*</span>}
            </label>

            {renderFieldInput(field, data?.[field.id], (val) => handleChange(field.id, val))}
          </div>
        ))}
      </div>
    </div>
  );
}