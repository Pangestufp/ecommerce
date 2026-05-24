/**
 * Format angka ke format Rupiah
 * @param {string|number} n
 * @returns {string} e.g. "Rp 1.300.000"
 */
export function formatRupiah(n) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(parseFloat(n));
}

/**
 * Pilih diskon dengan harga akhir (final_amount) paling murah
 * @param {Array} discounts
 * @returns {object|null}
 */
export function getBestDiscount(discounts) {
  if (!discounts || discounts.length === 0) return null;
  return discounts.reduce((best, d) =>
    parseFloat(d.final_amount) < parseFloat(best.final_amount) ? d : best
  );
}

/**
 * Hitung harga satuan produk setelah diskon
 * @param {object} item  - product_price item dari API
 * @param {object|null} selectedDiscount
 * @returns {number}
 */
export function getUnitPrice(item, selectedDiscount) {
  if (selectedDiscount) return parseFloat(selectedDiscount.final_amount);
  return parseFloat(item.product_price);
}