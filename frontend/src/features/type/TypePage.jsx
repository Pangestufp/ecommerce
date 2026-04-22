import { useState } from "react";
import { useType } from "./useType";
import Table from "../../shared/table/Table";
import Button from "../../shared/ui/Button";
import UpdateTypeModal from "./components/UpdateTypeModal";
import CreateTypeModal from "./components/CreateTypeModal";
import SearchBar from "../../shared/ui/SearchBar";

const columns = [
  { key: "type_name", label: "Nama Tipe", align: "left" },
  { key: "type_code", label: "Kode Tipe", align: "left" },
  { key: "type_desc", label: "Deskripsi", align: "left" },
];

export default function TypePage() {
  const { types, loading, page, hasNext, hasPrev, setSearch, create, update, del, next, prev } = useType();
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
      
      <div className="mb-4">
        <div className="flex justify-between items-center">
          <div className="w-72">
            <SearchBar
              placeholder="Cari tipe..."
              onChange={(value) => {
                setSearch(value);
              }}
            />
          </div>

          <Button onClick={() => setShowCreate(true)}>
            + Tambah
          </Button>
        </div>
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
        <Button variant="secondary" onClick={prev} disabled={loading||!hasPrev}>
          Prev
        </Button>
        <span className="text-sm text-gray-500">
            Page {page} 
        </span>
        <Button variant="secondary" onClick={next} disabled={loading||!hasNext}>
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
