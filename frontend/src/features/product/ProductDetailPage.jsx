import { useParams } from "react-router-dom";
import CreateDiscountModal from "./components/CreateDiscountModal";
import UpdateInventoryModal from "./components/UpdateInventoryModal";
import CreateInventoryModal from "./components/CreateInventoryModal ";
import CreatePriceModal from "./components/CreatepriceModal ";
import { useProductDetail } from "./useProductDetail";
import { useState, useEffect } from "react";
import SearchBar from "../../shared/ui/SearchBar";
import Button from "../../shared/ui/Button";
import Table from "../../shared/table/Table";

//  Image Album + Lightbox
const MAX_VISIBLE = 1;

function ProductImageAlbum({ images }) {
  const [lightboxIndex, setLightboxIndex] = useState(null);

  const close = () => setLightboxIndex(null);
  const prev = (e) => {
    e.stopPropagation();
    setLightboxIndex((i) => (i - 1 + images.length) % images.length);
  };
  const next = (e) => {
    e.stopPropagation();
    setLightboxIndex((i) => (i + 1) % images.length);
  };

  useEffect(() => {
    if (lightboxIndex === null) return;
    const handler = (e) => {
      if (e.key === "ArrowLeft")
        setLightboxIndex((i) => (i - 1 + images.length) % images.length);
      if (e.key === "ArrowRight")
        setLightboxIndex((i) => (i + 1) % images.length);
      if (e.key === "Escape") close();
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [lightboxIndex, images.length]);

  if (!images?.length) {
    return (
      <div className="w-20 h-20 flex-shrink-0 rounded-lg border border-dashed border-gray-300 flex items-center justify-center">
        <span className="text-xs text-gray-400">No image</span>
      </div>
    );
  }

  const visibleImages = images.slice(0, MAX_VISIBLE);
  const remaining = images.length - MAX_VISIBLE;

  return (
    <>
      <div className="flex gap-2 flex-shrink-0 flex-wrap">
        {visibleImages.map((img, idx) => (
          <button
            key={img.image_id ?? idx}
            type="button"
            onClick={() => setLightboxIndex(idx)}
            className="relative w-20 h-20 rounded-lg overflow-hidden border border-gray-200 hover:border-blue-400 hover:ring-2 hover:ring-blue-100 transition flex-shrink-0 focus:outline-none"
          >
            <img
              src={img.picture_path}
              alt={`Gambar ${idx + 1}`}
              className="w-full h-full object-cover"
            />
            {idx === MAX_VISIBLE - 1 && remaining > 0 && (
              <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                <span className="text-white text-sm font-semibold">
                  +{remaining}
                </span>
              </div>
            )}
          </button>
        ))}
      </div>

      {lightboxIndex !== null && (
        <div
          className="fixed inset-0 z-50 bg-black/80 flex items-center justify-center"
          onClick={close}
        >
          <div className="absolute top-4 left-1/2 -translate-x-1/2 bg-black/50 text-white text-xs px-3 py-1 rounded-full select-none">
            {lightboxIndex + 1} / {images.length}
          </div>

          <button
            type="button"
            onClick={close}
            className="absolute top-4 right-4 text-white bg-black/40 hover:bg-black/60 rounded-full w-9 h-9 flex items-center justify-center text-xl leading-none transition"
          >
            ×
          </button>

          {images.length > 1 && (
            <button
              type="button"
              onClick={prev}
              className="absolute left-4 top-1/2 -translate-y-1/2 text-white bg-black/40 hover:bg-black/60 rounded-full w-10 h-10 flex items-center justify-center text-2xl leading-none transition select-none"
            >
              ‹
            </button>
          )}

          <img
            src={images[lightboxIndex].picture_path}
            alt={`Gambar ${lightboxIndex + 1}`}
            className="max-h-[85vh] max-w-[90vw] rounded-lg object-contain"
            onClick={(e) => e.stopPropagation()}
          />

          {images.length > 1 && (
            <button
              type="button"
              onClick={next}
              className="absolute right-4 top-1/2 -translate-y-1/2 text-white bg-black/40 hover:bg-black/60 rounded-full w-10 h-10 flex items-center justify-center text-2xl leading-none transition select-none"
            >
              ›
            </button>
          )}

          {images.length > 1 && (
            <div className="absolute bottom-5 left-1/2 -translate-x-1/2 flex gap-1.5">
              {images.map((_, idx) => (
                <button
                  key={idx}
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    setLightboxIndex(idx);
                  }}
                  className={`w-2 h-2 rounded-full transition ${
                    idx === lightboxIndex ? "bg-white" : "bg-white/40"
                  }`}
                />
              ))}
            </div>
          )}
        </div>
      )}
    </>
  );
}

//  Section wrapper
function Section({ title, action, children }) {
  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div className="flex items-center justify-between px-5 py-3.5 border-b border-gray-100">
        <h2 className="text-sm font-semibold text-gray-700">{title}</h2>
        {action && <div className="flex items-center gap-3">{action}</div>}
      </div>
      <div className="p-5">{children}</div>
    </div>
  );
}

//  Pagination
function Pagination({ page, onPrev, onNext, disabledPrev, disabledNext }) {
  return (
    <div className="flex items-center justify-end gap-2 mt-4">
      <Button variant="secondary" onClick={onPrev} disabled={disabledPrev}>
        Prev
      </Button>
      <span className="text-sm text-gray-500">Page {page}</span>
      <Button variant="secondary" onClick={onNext} disabled={disabledNext}>
        Next
      </Button>
    </div>
  );
}

//  Main Page
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

  //  Columns
  const priceColumns = [
    {
      key: "product_price",
      label: "Harga",
      align: "left",
      render: (value) =>
        new Intl.NumberFormat("id-ID", {
          style: "currency",
          currency: "IDR",
          maximumFractionDigits: 0,
        }).format(value),
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

  const inventoryColumns = [
    { key: "batch_code", label: "Kode Batch", align: "left" },
    {
      key: "cost_price",
      label: "Harga Modal",
      align: "left",
      render: (value) =>
        new Intl.NumberFormat("id-ID", {
          style: "currency",
          currency: "IDR",
          maximumFractionDigits: 0,
        }).format(value),
    },
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

  //  Handlers
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

  //  Render
  return (
    <div className="p-6 flex flex-col gap-6">
      {/*  Info Produk  */}
      {/*  Info Produk  */}
      {product && (
        <Section title="Info Produk">
          <div className="flex flex-col gap-5">
            {/* Top Info */}
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

            {/* Description */}
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

      {/*  Harga  */}
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

      {/*  Diskon  */}
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

      {/*  Inventory  */}
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

      {/*  Modals  */}
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
