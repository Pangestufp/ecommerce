import { useParams } from "react-router-dom";
import { useState } from "react";
import { useProductDetail } from "./useProductDetail";

import CreateDiscountModal from "./components/CreateDiscountModal";
import UpdateInventoryModal from "./components/UpdateInventoryModal";
import CreateInventoryModal from "./components/CreateInventoryModal ";
import CreatePriceModal from "./components/CreatepriceModal ";
import ProductImageAlbum from "./components/Productimagealbum ";
import Section from "./components/Section ";
import Pagination from "./components/Pagination ";

import SearchBar from "../../shared/ui/SearchBar";
import Button from "../../shared/ui/Button";
import Table from "../../shared/table/Table";

export default function ProductDetailPage() {
  const { id } = useParams();

  const {
    product,
    types,
    loading,
    prices,
    pricePage,
    hasNextPrice,
    hasPrevPrice,
    createPrice,
    nextPrice,
    prevPrice,
    discounts,
    pageDiscount,
    hasNextDiscount,
    hasPrevDiscount,
    createDiscount,
    delDiscount,
    nextDiscount,
    prevDiscount,
    setSearchDiscount,
    inventories,
    pageInventory,
    hasNextInventory,
    hasPrevInventory,
    createInventory,
    updateInventory,
    nextInventory,
    prevInventory,
    setSearchInventory,
  } = useProductDetail(id);

  const [showCreateDiscount, setShowCreateDiscount] = useState(false);
  const [showCreateInventory, setShowCreateInventory] = useState(false);
  const [showCreatePrice, setShowCreatePrice] = useState(false);
  const [selectedInventory, setSelectedInventory] = useState(null);

  // inventoryColumns di sini karena butuh setSelectedInventory
  const inventoryColumns = [
    { key: "batch_code", label: "Kode Batch", align: "left" },
    { key: "cost_price_format", label: "Harga Modal", align: "left" },
    { key: "stock", label: "Stok", align: "left" },
    { key: "reserved_stock", label: "Direservasi", align: "left" },
    {
      key: "actions",
      label: "",
      align: "right",
      render: (_, row) => (
        <Button variant="secondary" onClick={() => setSelectedInventory(row)}>
          Edit
        </Button>
      ),
    },
  ];

  const priceColumns = [
    {
      key: "product_price_format",
      label: "Harga",
      align: "left",
    },
    {
      key: "created_at",
      label: "Ditambahkan Pada",
      align: "left",
      render: (value) =>
        new Date(value).toLocaleString("id-ID", {
          dateStyle: "medium",
          timeStyle: "short",
        }),
    },
  ];

  const discountColumns = [
    { key: "discount_name", label: "Nama Diskon", align: "left" },
    { key: "discount_type", label: "Tipe", align: "left" },
    { key: "discount_value_format", label: "Nilai", align: "left" },
    { key: "discount_Amount_format", label: "Jumlah", align: "left" },
    {
      key: "final_value",
      label: "Harga setelah diskon",
      align: "left",
      render: (value) => (
        <span
          className={
            value === "Harga belum diatur" ? "text-red-500 font-medium" : ""
          }
        >
          {value}
        </span>
      ),
    },
    {
      key: "status_format",
      label: "Status",
      align: "left",
      render: (value) => (
        <span
          className={`font-medium ${value === "Aktif" ? "text-green-600" : "text-red-500"}`}
        >
          {value}
        </span>
      ),
    },
    { key: "start_at_format", label: "Mulai", align: "left" },
    { key: "expired_at_format", label: "Berakhir", align: "left" },
  ];

  const handleCreatePrice = async (form) => {
    try {
      await createPrice(form);
      setShowCreatePrice(false);
    } catch {}
  };

  const handleCreateDiscount = async (form) => {
    try {
      await createDiscount(form);
      setShowCreateDiscount(false);
    } catch {}
  };

  const handleCreateInventory = async (form) => {
    try {
      await createInventory(form);
      setShowCreateInventory(false);
    } catch {}
  };

  const handleUpdateInventory = async (batchId, payload) => {
    try {
      await updateInventory(batchId, payload);
      setSelectedInventory(null);
    } catch {}
  };

  return (
    <div className="p-6 flex flex-col gap-6">
      {/* Info Produk */}
      {product && (
        <Section title="Info Produk">
          <div className="flex flex-col gap-5">
            <div className="flex gap-6 items-start">
              <ProductImageAlbum images={product.images} />

              <div className="flex flex-col gap-3 min-w-0 flex-1">
                <div className="flex items-start gap-2 flex-wrap">
                  <h1 className="text-base font-semibold text-gray-800">
                    {product.product_name}
                  </h1>
                  <span
                    className={`text-xs px-2 py-0.5 rounded-full font-medium flex-shrink-0 ${
                      product.status === 1
                        ? "bg-green-50 text-green-600"
                        : "bg-red-50 text-red-500"
                    }`}
                  >
                    {product.status === 1 ? "Aktif" : "Nonaktif"}
                  </span>
                </div>

                <div className="grid grid-cols-3 gap-x-10 gap-y-2 text-sm w-fit mt-1">
                  <div>
                    <p className="text-xs text-gray-400 mb-0.5">Kode Produk</p>
                    <p className="text-gray-700">{product.product_code}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-400 mb-0.5">Tipe</p>
                    <p className="text-gray-700">
                      {product.type_code !== "-"
                        ? `${product.type_code} - ${product.type_name}`
                        : product.type_name}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-400 mb-0.5">Berat</p>
                    <p className="text-gray-700">{product.weight_gram} g</p>
                  </div>
                </div>
              </div>
            </div>

            {product.description && (
              <div className="max-w-4xl">
                <p className="text-xs text-gray-400 mb-1">Deskripsi</p>
                <p className="text-sm text-gray-600 leading-relaxed text-justify">
                  {product.description}
                </p>
              </div>
            )}
          </div>
        </Section>
      )}

      {/* Harga */}
      <Section
        title="Riwayat Harga"
        action={
          <Button onClick={() => setShowCreatePrice(true)}>
            + Tambah Harga
          </Button>
        }
      >
        <Table columns={priceColumns} data={prices} rowKey="price_id" />
        <Pagination
          page={pricePage}
          onPrev={prevPrice}
          onNext={nextPrice}
          disabledPrev={loading || !hasPrevPrice}
          disabledNext={loading || !hasNextPrice}
        />
      </Section>

      {/* Diskon */}
      <Section
        title="Diskon"
        action={
          <>
            <div className="w-52">
              <SearchBar
                placeholder="Cari diskon..."
                onChange={(value) => setSearchDiscount(value)}
              />
            </div>
            <Button onClick={() => setShowCreateDiscount(true)}>
              + Tambah Diskon
            </Button>
          </>
        }
      >
        <Table
          columns={discountColumns}
          data={discounts}
          rowKey="discount_id"
          onDelete={(id) => delDiscount(id)}
        />
        <Pagination
          page={pageDiscount}
          onPrev={prevDiscount}
          onNext={nextDiscount}
          disabledPrev={loading || !hasPrevDiscount}
          disabledNext={loading || !hasNextDiscount}
        />
      </Section>

      {/* Inventory */}
      <Section
        title="Inventory"
        action={
          <>
            <div className="w-52">
              <SearchBar
                placeholder="Cari inventory..."
                onChange={(value) => setSearchInventory(value)}
              />
            </div>
            <Button onClick={() => setShowCreateInventory(true)}>
              + Tambah Batch
            </Button>
          </>
        }
      >
        <Table
          columns={inventoryColumns}
          data={inventories}
          rowKey="batch_id"
        />
        <Pagination
          page={pageInventory}
          onPrev={prevInventory}
          onNext={nextInventory}
          disabledPrev={loading || !hasPrevInventory}
          disabledNext={loading || !hasNextInventory}
        />
      </Section>

      {/* Modals */}
      {showCreatePrice && (
        <CreatePriceModal
          productID={id}
          onSubmit={handleCreatePrice}
          onClose={() => setShowCreatePrice(false)}
          loading={loading}
        />
      )}

      {showCreateDiscount && (
        <CreateDiscountModal
          productID={id}
          types={types}
          onSubmit={handleCreateDiscount}
          onClose={() => setShowCreateDiscount(false)}
          loading={loading}
        />
      )}

      {showCreateInventory && (
        <CreateInventoryModal
          productID={id}
          onSubmit={handleCreateInventory}
          onClose={() => setShowCreateInventory(false)}
          loading={loading}
        />
      )}

      {selectedInventory && (
        <UpdateInventoryModal
          inventory={selectedInventory}
          onSubmit={handleUpdateInventory}
          onClose={() => setSelectedInventory(null)}
          loading={loading}
        />
      )}
    </div>
  );
}
