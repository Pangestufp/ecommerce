const CART_KEY = "myecommercedev-cart";

export function saveToCart(product) {
  const current = getCart();
  const existingIdx = current.findIndex(i => i.id === product.id);

  if (existingIdx !== -1) {
    const newQty = current[existingIdx].qty + product.qty;
    current[existingIdx].qty = Math.min(newQty, product.availableStock);
  } else {
    current.push({
      id: product.id,
      productName: product.productName,
      qty: Math.min(product.qty, product.availableStock),
      availableStock: product.availableStock,
      image: product.image,
    });
  }

  localStorage.setItem(CART_KEY, JSON.stringify(current));
}

export function getCart() {
  try {
    return JSON.parse(localStorage.getItem(CART_KEY)) || [];
  } catch {
    return [];
  }
}

// validatedItems dari backend pakai snake_case
// di-convert ke camelCase biar konsisten dengan saveToCart
export function setCart(validatedItems) {
  const normalized = validatedItems.map(item => ({
    id: item.id,
    productName: item.product_name,
    qty: item.qty,
    image: item.image,
  }));
  localStorage.setItem(CART_KEY, JSON.stringify(normalized));
}

export function clearCart() {
  localStorage.removeItem(CART_KEY);
}