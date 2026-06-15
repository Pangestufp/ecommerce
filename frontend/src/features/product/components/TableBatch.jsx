import React, { useState } from "react";
import Button from "../../../shared/ui/Button";
import { ChevronRight } from "lucide-react";
import BatchTransactionSubTable from "./BatchTransactionSubTable";

export default function TableBatch({ data, onEdit }) {
  const [expandedRows, setExpandedRows] = useState({});

  const toggleRow = (batchID) => {
    setExpandedRows((prev) => ({
      ...prev,
      [batchID]: !prev[batchID], 
    }));
  };


  const handleEditClick = (row) => {
    setExpandedRows((prev) => ({
      ...prev,
      [row.batch_id]: false,
    }));

    if (onEdit) onEdit(row);
  };

  return (
    <div className="w-full overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
      <div className="overflow-x-auto w-full">
        <table className="min-w-full divide-y divide-gray-200 text-sm text-left">
          <thead className="bg-gray-50 text-xs font-semibold uppercase tracking-wider text-gray-600">
            <tr>
              <th className="px-4 py-3.5 w-10"></th> 
              <th className="px-4 py-3.5">Kode Batch</th>
              <th className="px-4 py-3.5">Harga Modal</th>
              <th className="px-4 py-3.5">Stok</th>
              <th className="px-4 py-3.5">Direservasi</th>
              <th className="px-4 py-3.5 text-right">Aksi</th>
            </tr>
          </thead>
          
          <tbody className="divide-y divide-gray-200 bg-white">
            {data && data.length > 0 ? (
              data.map((row) => {
                const isExpanded = !!expandedRows[row.batch_id];
                return (
                  <React.Fragment key={row.batch_id}>
                    <tr className={`hover:bg-gray-50 transition-colors ${isExpanded ? "bg-gray-100/50" : ""}`}>
                     <td className="px-4 py-3 text-center">
                      <button 
                        type="button"
                        onClick={() => toggleRow(row.batch_id)} 
                        className="text-gray-400 hover:text-gray-700 transition-transform duration-200 transform flex items-center justify-center w-full focus:outline-none"
                        style={{ transform: isExpanded ? "rotate(90deg)" : "rotate(0deg)" }}
                      >
                       
                        <ChevronRight size={18} strokeWidth={2.5} />
                      </button>
                    </td>
                      <td className="px-4 py-3 font-medium text-gray-900">{row.batch_code}</td>
                      <td className="px-4 py-3 text-gray-700">{row.cost_price_format}</td>
                      <td className="px-4 py-3 text-gray-700 font-semibold">{row.stock}</td>
                      <td className="px-4 py-3 text-gray-500">{row.reserved_stock}</td>
                      <td className="px-4 py-3 text-right">
                        <Button variant="secondary" onClick={() => handleEditClick(row)}>
                          Edit
                        </Button>
                      </td>
                    </tr>
                    
                    {isExpanded && (
                      <tr>
                        <td colSpan={6} className="bg-gray-50/50 px-6 py-4">
                          {/* Memanggil sub-table komponen yang sudah dipisahkan */}
                          <BatchTransactionSubTable batchID={row.batch_id} />
                        </td>
                      </tr>
                    )}
                  </React.Fragment>
                );
              })
            ) : (
              <tr>
                <td colSpan={6} className="text-center py-8 text-gray-400">
                  Data batch tidak ditemukan
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}