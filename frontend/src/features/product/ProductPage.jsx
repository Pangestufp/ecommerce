import { useState } from "react";
import { useProduct } from "./useProduct";
import TableProduct from "./components/TableProduct";
import Button from "../../shared/ui/Button";
import SearchBar from "../../shared/ui/SearchBar";
import CreateProductModal from "./components/CreateProductModal";
import UpdateProductModal from "./components/UpdateProductModal";

export default function ProductPage() {
  const {
    products,
    types,
    loading,
    page,
    hasNext,
    hasPrev,
    setSearch,
    create,
    update,
    del,
    next,
    prev,
    getById,
    generatePresignedURLs,
    uploadToPresignedURL,
  } = useProduct();

  const [showCreate, setShowCreate] = useState(false);
  const [showUpdate, setShowUpdate] = useState(false);
  const [selectedProductId, setSelectedProductId] = useState(null);

  const handleCreate = async (form, imageItems) => {
    try {
      await create(form, imageItems);
      setShowCreate(false);
    } catch {}
  };

  const handleUpdate = async (id, form, imageItems) => {
    try {
      await update(id, form, imageItems);
      setShowUpdate(false);
    } catch {}
  };

  return (
    <div className="p-6">
      <div className="mb-4">
        <div className="flex justify-between items-center">
          <div className="w-72">
            <SearchBar
              placeholder="Cari produk..."
              onChange={(value) => setSearch(value)}
            />
          </div>
          <Button onClick={() => setShowCreate(true)}>+ Tambah</Button>
        </div>
      </div>

      <TableProduct
        data={products}
        onUpdate={(row) => {
          setSelectedProductId(row.product_id);
          setShowUpdate(true);
        }}
        onDelete={(id) => del(id)}
      />

      <div className="flex items-center justify-end gap-2 mt-4">
        <Button variant="secondary" onClick={prev} disabled={loading || !hasPrev}>
          Prev
        </Button>
        <span className="text-sm text-gray-500">Page {page}</span>
        <Button variant="secondary" onClick={next} disabled={loading || !hasNext}>
          Next
        </Button>
      </div>

      {showCreate && (
        <CreateProductModal
          onSubmit={handleCreate}
          onClose={() => setShowCreate(false)}
          loading={loading}
          types={types}
          generatePresignedURLs={generatePresignedURLs}
          uploadToPresignedURL={uploadToPresignedURL}
        />
      )}

      {showUpdate && selectedProductId && (
        <UpdateProductModal
          productId={selectedProductId}
          onSubmit={handleUpdate}
          onClose={() => {
            setShowUpdate(false);
            setSelectedProductId(null);
          }}
          loading={loading}
          types={types}
          getById={getById}
          generatePresignedURLs={generatePresignedURLs}
          uploadToPresignedURL={uploadToPresignedURL}
        />
      )}
    </div>
  );
}