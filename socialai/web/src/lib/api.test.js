import { describe, expect, it, vi } from "vitest";
import { api, APIError } from "./api";

describe("api client", () => {
  it("sends bearer auth and parses posts", async () => {
    global.fetch = vi.fn().mockResolvedValue(new Response(JSON.stringify([{ id: "one" }]), {
      status: 200, headers: { "Content-Type": "application/json" },
    }));
    await expect(api.posts("secret", { keywords: "night sky" })).resolves.toEqual([{ id: "one" }]);
    const [url, options] = global.fetch.mock.calls[0];
    expect(url).toBe("/api/posts?keywords=night+sky");
    expect(options.headers.Authorization).toBe("Bearer secret");
  });

  it("turns structured backend errors into APIError", async () => {
    global.fetch = vi.fn().mockResolvedValue(new Response(JSON.stringify({ error: { code: "invalid_credentials", message: "Nope" } }), {
      status: 401, headers: { "Content-Type": "application/json" },
    }));
    await expect(api.signin({ username: "x", password: "y" })).rejects.toMatchObject({
      name: "APIError", status: 401, code: "invalid_credentials", message: "Nope",
    });
    expect(APIError).toBeTypeOf("function");
  });

  it("encodes generated image IDs when discarding previews", async () => {
    global.fetch = vi.fn().mockResolvedValue(new Response(null, { status: 204 }));
    await expect(api.discard("secret", "image/one")).resolves.toBeNull();
    expect(global.fetch).toHaveBeenCalledWith("/api/ai/images/image%2Fone", expect.objectContaining({
      method: "DELETE",
      headers: expect.objectContaining({ Authorization: "Bearer secret" }),
    }));
  });
});
