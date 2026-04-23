export function ResellerSettings() {
  return (
    <div className="p-4 flex flex-col gap-4 max-w-2xl">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Branding</h3>
        </div>
        <div className="p-4 space-y-4">
          {[
            { label: "Store name", value: "ProxyVN" },
            { label: "Custom domain", value: "proxyvn.io" },
            { label: "Support email", value: "support@proxyvn.io" },
            { label: "Accent color", value: "#D50C2D" },
          ].map(({ label, value }) => (
            <div key={label} className="flex items-center gap-4">
              <label className="text-[12px] text-gray-500 w-36 shrink-0">{label}</label>
              <input
                defaultValue={value}
                className="flex-1 h-8 p-4.5 border border-gray-300 rounded-[3px] text-[13px] font-[inherit] text-gray-800 bg-white outline-none focus:border-[#D50C2D]"
              />
            </div>
          ))}
        </div>
        <div className="p-4 p-4 border-t border-gray-100 flex justify-end">
          <button className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm">
            Save branding
          </button>
        </div>
      </div>
    </div>
  );
}
