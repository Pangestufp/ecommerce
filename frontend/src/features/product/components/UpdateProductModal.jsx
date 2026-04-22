import { useEffect, useState } from "react";
import TextField from "../../../shared/ui/TextField";
import Button from "../../../shared/ui/Button";
import Textarea from "../../../shared/ui/TextArea";
import SearchableSelect from "../../../shared/ui/SearchableSelect";

export default function UpdateProductModal({
  productId,
  onSubmit,
  onClose,
  loading,
  types,
  getById,
  generatePresignedURLs,
  uploadToPresignedURL,
}) {
  const [form, setForm] = useState({
    product_code: "",
    product_name: "",
    weight_gram: "",
    type_id: "",
    description: "",
  });

  // imageItems: [{ id, previewURL, objectName, uploading, isExisting }]
  // isExisting: true = gambar lama yang sudah di-upload ulang ke MinIO
  const [imageItems, setImageItems] = useState([]);
  const [fetching, setFetching] = useState(true);

  useEffect(() => {
    const loadProduct = async () => {
      setFetching(true);
      try {
        const product = await getById(productId);

        setForm({
          product_code: product.product_code,
          product_name: product.product_name,
          weight_gram: product.weight_gram,
          type_id: product.type_id,
          description: product.description,
        });

        if (product.images?.length > 0) {
          const placeholders = product.images.map((img) => ({
            id: img.image_id,
            previewURL: img.picture_path,
            objectName: null,
            uploading: true,
            isExisting: true,
          }));
          setImageItems(placeholders);

          // Download semua gambar lama sebagai Blob
          const blobs = await Promise.all(
            product.images.map(async (img) => {
              const res = await fetch(img.picture_path);
              const blob = await res.blob();
              const ext =
                img.picture_path.split(".").pop().split("?")[0] || "jpg";
              return new File([blob], `${img.image_id}.${ext}`, {
                type: blob.type || "image/jpeg",
              });
            }),
          );

          // Generate presigned upload URLs untuk semua gambar lama sekaligus
          const uploads = await generatePresignedURLs(blobs);

          // Upload ulang ke MinIO
          await Promise.all(
            uploads.map((u, idx) =>
              uploadToPresignedURL(u.upload_url, blobs[idx]),
            ),
          );

          // Update objectName
          setImageItems(
            product.images.map((img, idx) => ({
              id: img.image_id,
              previewURL: img.picture_path,
              objectName: uploads[idx].object_name,
              uploading: false,
              isExisting: true,
            })),
          );
        }
      } catch (err) {
        // gagal load
      } finally {
        setFetching(false);
      }
    };

    loadProduct();
  }, [productId]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleAddImages = async (e) => {
    const files = Array.from(e.target.files);
    if (!files.length) return;

    const tempItems = files.map((file) => ({
      id: crypto.randomUUID(),
      file,
      previewURL: URL.createObjectURL(file),
      objectName: null,
      uploading: true,
      isExisting: false,
    }));

    setImageItems((prev) => [...prev, ...tempItems]);

    try {
      const uploads = await generatePresignedURLs(files);
      await Promise.all(
        uploads.map((u, idx) => uploadToPresignedURL(u.upload_url, files[idx])),
      );

      setImageItems((prev) =>
        prev.map((item) => {
          const tempIndex = tempItems.findIndex((t) => t.id === item.id);
          if (tempIndex === -1) return item;
          return {
            ...item,
            objectName: uploads[tempIndex].object_name,
            uploading: false,
          };
        }),
      );
    } catch (err) {
      setImageItems((prev) =>
        prev.filter((item) => !tempItems.find((t) => t.id === item.id)),
      );
    }

    e.target.value = "";
  };

  const handleRemoveImage = (id) => {
    setImageItems((prev) => prev.filter((item) => item.id !== id));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const stillUploading = imageItems.some((i) => i.uploading);
    if (stillUploading) return;
    await onSubmit(productId, form, imageItems);
  };

  if (fetching) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-black/40 z-50">
        <div className="bg-white rounded-xl p-6 w-[480px] shadow-lg flex items-center justify-center h-40">
          <span className="text-sm text-gray-500">Memuat data produk...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/40 z-50">
      <div className="bg-white rounded-xl p-6 w-[480px] max-h-[90vh] overflow-y-auto shadow-lg">
        <h2 className="text-base font-semibold text-gray-800 mb-4">
          Update Produk
        </h2>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <TextField
            label="Kode Produk"
            name="product_code"
            placeholder="Masukkan kode produk"
            value={form.product_code}
            onChange={handleChange}
          />
          <TextField
            label="Nama Produk"
            name="product_name"
            placeholder="Masukkan nama produk"
            value={form.product_name}
            onChange={handleChange}
          />
          <TextField
            label="Berat (gram)"
            name="weight_gram"
            type="number"
            placeholder="Masukkan berat dalam gram"
            value={form.weight_gram}
            onChange={handleChange}
          />

          <SearchableSelect
            label="Tipe"
            name="type_id"
            value={form.type_id}
            onChange={handleChange}
            options={types}
            placeholder="Cari tipe..."
            required
            getValue={(t) => t.type_id}
            getLabel={(t) => {
              const a = t.type_code;
              const b = t.type_name;

              if (a === "-") return b;
              if (b === "-") return a;
              return `${a} - ${b}`;
            }}
          />

          <Textarea
            label="Deskripsi"
            name="description"
            placeholder="Masukkan deskripsi produk"
            value={form.description}
            onChange={handleChange}
            rows={3}
          />

          {/* Gambar */}
          <div className="flex flex-col gap-2">
            <label className="text-sm font-medium text-gray-700">Gambar</label>

            {imageItems.length > 0 && (
              <div className="flex flex-wrap gap-2">
                {imageItems.map((item) => (
                  <div key={item.id} className="relative w-20 h-20">
                    <img
                      src={item.previewURL}
                      alt="preview"
                      className="w-20 h-20 object-cover rounded-lg border border-gray-200"
                    />
                    {item.uploading && (
                      <div className="absolute inset-0 bg-black/40 rounded-lg flex items-center justify-center">
                        <span className="text-white text-xs">Upload...</span>
                      </div>
                    )}
                    {!item.uploading && (
                      <button
                        type="button"
                        onClick={() => handleRemoveImage(item.id)}
                        className="absolute -top-1 -right-1 bg-red-500 text-white rounded-full w-4 h-4 flex items-center justify-center text-xs"
                      >
                        ×
                      </button>
                    )}
                  </div>
                ))}
              </div>
            )}

            <label className="cursor-pointer inline-flex items-center gap-2 px-3 py-2 border border-dashed border-gray-300 rounded-lg text-sm text-gray-500 hover:border-blue-400 hover:text-blue-500 transition w-fit">
              <span>+ Tambah Gambar</span>
              <input
                type="file"
                accept="image/*"
                multiple
                className="hidden"
                onChange={handleAddImages}
              />
            </label>
          </div>

          <div className="flex justify-end gap-2 mt-2">
            <Button variant="secondary" onClick={onClose} disabled={loading}>
              Batal
            </Button>
            <Button
              type="submit"
              loading={loading || imageItems.some((i) => i.uploading)}
            >
              Simpan
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
