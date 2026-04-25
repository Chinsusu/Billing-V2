"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { clientServiceOrderLabel, clientServiceSourceLabel } from "@/lib/api/clientViewModels";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { CLIENT_SERVICES } from "@/mocks/billingData";

export function ClientDashboard() {
  const wallets = useApiResource(billingApi.listClientWallets);
  const services = useApiResource(billingApi.listClientServices);
  const orders = useApiResource(billingApi.listClientOrders);
  const wallet = wallets.data?.[0];
  const liveServices = services.status === "success" ? services.data ?? [] : null;
  const serviceRows = liveServices
    ? liveServices.map((service) => ({
        id: service.id,
        label: recordLabel(service.display_id, "SVC-"),
        region: clientServiceSourceLabel(service),
        bandwidth: clientServiceOrderLabel(service),
        expiry: compactDateTime(service.term_end),
        status: service.status,
        note: undefined,
      }))
    : CLIENT_SERVICES;
  const suspended = serviceRows.filter((service) => service.status === "suspended");

  return (
    <div className="p-4 flex flex-col gap-4">
      {suspended.map((service) => (
        <div key={service.id} className="bg-amber-50 border border-amber-200 text-amber-700 text-[12px] p-4 p-4.5 rounded flex items-center gap-4">
          <span>!</span>
          <span><strong>{service.label}</strong> is suspended. {service.note} Renew to restore access.</span>
          <button className="ml-auto inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-amber-600 hover:bg-amber-700 text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm">
            Renew now
          </button>
        </div>
      ))}

      <div className="bg-white border border-gray-200 rounded p-4 flex items-center justify-between">
        <div>
          <div className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Wallet balance</div>
          <div className="text-lg font-medium tabular-nums text-gray-900">
            {wallet ? moneyMinor(wallet.available_balance_minor, wallet.currency) : "$128.40"}
          </div>
          <div className="text-[11px] text-gray-400 mt-0.5">{orders.data?.length ?? 0} order records</div>
        </div>
        <button className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm">
          + Top up
        </button>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">My services</h3>
          <a href="#" className="text-[12px] text-[#D50C2D]">View all -&gt;</a>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Label", "Region / Source", "Order", "Expires", "Status"].map((heading) => (
                <th key={heading} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {heading}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {serviceRows.map((service) => (
              <tr key={service.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 font-medium text-gray-900">{service.label}</td>
                <td className="p-4 p-4 text-[12px] text-gray-400">{service.region}</td>
                <td className="p-4 p-4 text-gray-500">{service.bandwidth}</td>
                <td className="p-4 p-4 text-gray-400">{service.expiry}</td>
                <td className="p-4 p-4"><StatusBadge status={service.status} dot /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
