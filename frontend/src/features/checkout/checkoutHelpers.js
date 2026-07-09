export function formatRupiah(n) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(parseFloat(n));
}

export function getBestDiscount(discounts) {
  if (!discounts || discounts.length === 0) return null;
  return discounts.reduce((best, d) =>
    parseFloat(d.final_amount) < parseFloat(best.final_amount) ? d : best
  );
}

export function generateIdempotencyKey() {
  if (typeof crypto !== "undefined" && crypto.randomUUID) {
    return crypto.randomUUID();
  }

  // fallback kalau browser lama / crypto.randomUUID gak ada
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

export function getUnitPrice(item, selectedDiscount) {
  if (selectedDiscount) return parseFloat(selectedDiscount.final_amount);
  return parseFloat(item.product_price);
}