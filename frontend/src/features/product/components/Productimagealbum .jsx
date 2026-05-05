import { useState, useEffect } from "react";

const MAX_VISIBLE = 1;

export default function ProductImageAlbum({ images }) {
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