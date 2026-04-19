import { useState } from "react";
import { useType } from "./useType";
import CreateTypeModal from "./CreateTypeModal";
import Table from "../../shared/table/Table";
import Button from "../../shared/ui/Button";
import UpdateTypeModal from "./UpdateTypeModal";

const columns = [
  { key: "type_name", label: "Nama Tipe" },
  { key: "type_code", label: "Kode Tipe" },
  { key: "type_desc", label: "Deskripsi" },
];

export default function TypePage() {
  const { types, loading, cursorHistory, create, update, del, next, prev } = useType();
  const [showCreate, setShowCreate] = useState(false);
  const [showUpdate, setShowUpdate] = useState(false);
  const [selectedRow, setSelectedRow] = useState(null);

  const handleCreate = async (form) => {
    try {
      await create(form);
      setShowCreate(false);
    } catch {}
  };

  const handleUpdate = async (id, form) => {
    try {
      await update(id, form);
      setShowUpdate(false);
    } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-lg font-semibold text-gray-800">Tipe</h1>
        <Button onClick={() => setShowCreate(true)}>+ Tambah</Button>
      </div>

      <Table
        columns={columns}
        data={types}
        rowKey="type_id"
        onUpdate={(row) => {
          setSelectedRow(row);
          setShowUpdate(true);
        }}
        onDelete={(id) => del(id)}
      />

      <div className="flex items-center justify-end gap-2 mt-4">
        <Button variant="secondary" onClick={prev} disabled={loading}>
          Prev
        </Button>
        <span className="text-sm text-gray-500">
            Page {cursorHistory.length + 1}
        </span>
        <Button variant="secondary" onClick={next} disabled={loading}>
          Next
        </Button>
      </div>

      {showCreate && (
        <CreateTypeModal
          onSubmit={handleCreate}
          onClose={() => setShowCreate(false)}
          loading={loading}
        />
      )}

      {showUpdate && (
        <UpdateTypeModal
          data={selectedRow}
          onSubmit={handleUpdate}
          onClose={() => setShowUpdate(false)}
          loading={loading}
        />
      )}
    </div>
  );
}
