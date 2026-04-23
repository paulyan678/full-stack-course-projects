let csrfToken;
let csrfHeader = "X-XSRF-TOKEN";
let csrfRequest;

export class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function loadCsrfToken() {
  const response = await fetch("/auth/csrf", {
    credentials: "same-origin",
    headers: { Accept: "application/json" }
  });
  if (!response.ok) {
    throw new ApiError("Could not initialize a secure session", response.status);
  }
  const payload = await response.json();
  if (typeof payload.token !== "string" || typeof payload.header_name !== "string") {
    throw new ApiError("The server returned an invalid security token", 500);
  }
  csrfToken = payload.token;
  csrfHeader = payload.header_name;
}

async function ensureCsrfToken() {
  if (csrfToken) return;
  if (!csrfRequest) {
    csrfRequest = loadCsrfToken().finally(() => {
      csrfRequest = undefined;
    });
  }
  await csrfRequest;
}

async function errorMessage(response) {
  const fallback = response.status === 401
    ? "Your session has expired. Please sign in again."
    : `Request failed (${response.status})`;
  const contentType = response.headers.get("content-type") || "";
  if (!contentType.includes("application/json") && !contentType.includes("application/problem+json")) {
    return fallback;
  }
  try {
    const payload = await response.json();
    const fieldError = payload.errors && Object.values(payload.errors)[0];
    return fieldError || payload.detail || payload.message || fallback;
  } catch {
    return fallback;
  }
}

async function request(path, options = {}) {
  const method = options.method || "GET";
  const mutating = !["GET", "HEAD", "OPTIONS"].includes(method.toUpperCase());
  if (mutating && !csrfToken) {
    await ensureCsrfToken();
  }

  const response = await fetch(path, {
    credentials: "same-origin",
    ...options,
    headers: {
      Accept: "application/json",
      ...(mutating ? { [csrfHeader]: csrfToken } : {}),
      ...options.headers
    }
  });

  if (!response.ok) {
    throw new ApiError(await errorMessage(response), response.status);
  }
  if (response.status === 204) {
    return undefined;
  }
  const contentType = response.headers.get("content-type") || "";
  return contentType.includes("json") ? response.json() : undefined;
}

export async function login({ email, password }) {
  const body = new URLSearchParams({ email, password });
  await request("/login", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8" },
    body
  });
  csrfToken = undefined;
}

export async function logout() {
  await request("/logout", { method: "POST" });
  csrfToken = undefined;
}

export function signup(data) {
  return request("/signup", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data)
  });
}

export function getAuthStatus() {
  return request("/auth/me");
}

export function getRestaurants() {
  return request("/restaurants/menu");
}

export function getMenus(restaurantId) {
  return request(`/restaurant/${encodeURIComponent(restaurantId)}/menu`);
}

export function getCart() {
  return request("/cart");
}

export function addItemToCart(menuId) {
  return request("/cart", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ menu_id: menuId })
  });
}

export function checkout() {
  return request("/cart/checkout", { method: "POST" });
}

export function resetClientForTests() {
  csrfToken = undefined;
  csrfHeader = "X-XSRF-TOKEN";
  csrfRequest = undefined;
}
