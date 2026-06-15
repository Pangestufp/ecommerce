import React from "react";
import Table from "../../../shared/table/Table";
import Pagination from "./Pagination ";
import useBatchTransaction from "../UseBatchTransaction";

const transactionColumns = [
  { 
    key: "type", 
    label: "Tipe", 
    align: "left", 
    render: (val) => (
      <span className={`font-bold px-2 py-1 rounded text-xs ${val === "IN" ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"}`}>
        {val}
      </span>
    )
  },
  { key: "quantity", label: "Jumlah", align: "left" },
  { key: "reference_type", label: "Jenis Ref", align: "left" },
  { key: "note", label: "Catatan Aktivitas", align: "left" },
];

export default function BatchTransactionSubTable({ batchID }) {
  // Panggil logika dari custom hook
  const { 
    transactions, 
    loading, 
    paginate, 
    page, 
    handleNext, 
    handlePrev 
  } = useBatchTransaction(batchID);

  return (
    <div className="p-4 bg-gray-50 rounded-lg border border-gray-200 my-2 shadow-inner">
      <div className="text-xs font-bold text-gray-500 uppercase tracking-wider mb-3">
        📜 Riwayat Keluar Masuk Stok Batch
      </div>
      
      <Table 
        columns={transactionColumns} 
        data={transactions} 
        rowKey="transaction_id"
        isLoading={loading}
      />
      
      <div className="mt-2">
        <Pagination
          page={page}
          onPrev={handlePrev}
          onNext={handleNext}
          disabledPrev={loading || paginate?.has_prev !== "true"}
          disabledNext={loading || paginate?.has_next !== "true"}
        />
      </div>
    </div>
  );
}