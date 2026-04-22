export function AdminSettings() {
  return (
    <div className="p-4 flex flex-col gap-4 max-w-2xl">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Platform settings</h3>
        </div>
        <div className="p-4 space-y-4">
          {[
            { label: "Platform name", value: "HANetwork Billing" },
            { label: "Primary domain", value: "billing.hanetwork.vn" },
            { label: "Support email", value: "support@hanetwork.vn" },
            { label: "Default currency", value: "USD" },
            { label: "Default timezone", value: "Asia/Ho_Chi_Minh" },
          ].map(({ label, value }) => (
            <div key={label} className="flex items-center gap-4">
              <label className="text-[12px] text-gray-500 w-40 shrink-0">{label}</label>
              <input
                defaultValue={value}
                className="flex-1 h-8 p-4.5 border border-gray-300 rounded-[3px] text-[13px] font-[inherit] text-gray-800 bg-white outline-none focus:border-[#D50C2D]"
              />
            </div>
          ))}
        </div>
        <div className="p-4 p-4 border-t border-gray-100 flex justify-end">
          <button className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm">
            Save changes
          </button>
        </div>
      </div>

      <div className="bg-amber-50 border border-amber-200 rounded p-4 text-[12px] text-amber-700">
        Settings changes are audited. All modifications are logged with actor, timestamp, and previous value.
      </div>
    </div>
  );
}
