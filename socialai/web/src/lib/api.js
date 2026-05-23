const API_BASE = (import.meta.env.VITE_API_BASE_URL || "").replace(/\/$/, "");

export class APIError extends Error {
  constructor(message, status, code) {
    super(message);
    this.name = "APIError";
    this.status = status;
    this.code = code;
  }
}

async function request(path, { token, body, headers, ...options } = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    body,
    headers: {
      ...(body && !(body instanceof FormData) ? { "Content-Type": "application/json" } : {}),
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...headers,
    },
  });
  if (response.status === 204) return null;
  const contentType = response.headers.get("content-type") || "";
  const payload = contentType.includes("application/json") ? await response.json() : await response.text();
  if (!response.ok) {
    const message = payload?.error?.message || payload || "Request failed.";
    throw new APIError(message, response.status, payload?.error?.code);
  }
  return payload;
}

export const api = {
  signup: (credentials) => request("/api/auth/signup", { method: "POST", body: JSON.stringify(credentials) }),
  signin: (credentials) => request("/api/auth/signin", { method: "POST", body: JSON.stringify(credentials) }),
  posts: (token, params = {}) => {
    const search = new URLSearchParams();
    if (params.user) search.set("user", params.user);
    if (params.keywords) search.set("keywords", params.keywords);
    return request(`/api/posts${search.size ? `?${search}` : ""}`, { token });
  },
  upload: (token, formData) => request("/api/posts", { method: "POST", token, body: formData }),
  deletePost: (token, id) => request(`/api/posts/${encodeURIComponent(id)}`, { method: "DELETE", token }),
  generate: (token, prompt) => request("/api/ai/images", { method: "POST", token, body: JSON.stringify({ prompt }) }),
  discard: (token, id) => request(`/api/ai/images/${encodeURIComponent(id)}`, { method: "DELETE", token }),
  publish: (token, id, message) => request(`/api/ai/images/${encodeURIComponent(id)}/publish`, {
    method: "POST", token, body: JSON.stringify({ message }),
  }),
};
