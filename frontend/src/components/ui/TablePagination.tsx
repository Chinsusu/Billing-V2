"use client";

import React from "react";

interface Props {
    page: number;
    setPage: (p: number) => void;
    limit: number;
    setLimit: (l: number) => void;
    total: number;
}

export function TablePagination({ page, setPage, limit, setLimit, total }: Props) {
    const totalPages = limit === -1 ? 1 : Math.ceil(total / limit);

    return (
        <div className="p-3 border-t border-gray-100 flex items-center justify-between text-[11px] text-gray-500 bg-gray-50/50 rounded-b">
            <div className="flex items-center gap-2">
                <span>Hiển thị</span>
                <select
                    value={limit}
                    onChange={(e) => { setLimit(Number(e.target.value)); setPage(1); }}
                    className="border border-gray-200 rounded px-1.5 py-0.5 outline-none bg-white font-medium cursor-pointer"
                >
                    <option value={10}>10</option>
                    <option value={20}>20</option>
                    <option value={50}>50</option>
                    <option value={100}>100</option>
                    <option value={-1}>All</option>
                </select>
                <span>dòng</span>
            </div>

            {totalPages > 1 && (
                <div className="flex items-center gap-1.5">
                    <button
                        disabled={page === 1}
                        onClick={() => setPage(page - 1)}
                        className="px-2.5 py-1 border border-gray-200 bg-white hover:bg-gray-50 rounded disabled:opacity-50 transition-colors cursor-pointer"
                    >
                        Trang trước
                    </button>
                    <span className="px-1 text-gray-400 font-medium">{page} / {totalPages}</span>
                    <button
                        disabled={page === totalPages}
                        onClick={() => setPage(page + 1)}
                        className="px-2.5 py-1 border border-gray-200 bg-white hover:bg-gray-50 rounded disabled:opacity-50 transition-colors cursor-pointer"
                    >
                        Trang tiếp
                    </button>
                </div>
            )}
        </div>
    );
}
