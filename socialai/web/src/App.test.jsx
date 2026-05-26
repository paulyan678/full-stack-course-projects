import { MemoryRouter } from "react-router-dom";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";
import App from "./App";

function response(payload, status = 200) {
  return new Response(JSON.stringify(payload), { status, headers: { "Content-Type": "application/json" } });
}

describe("SocialAI app", () => {
  beforeEach(() => { global.fetch = vi.fn(); });

  it("signs in and enters the protected studio", async () => {
    global.fetch.mockResolvedValue(response({ token: "signed-token", username: "nova" }));
    const user = userEvent.setup();
    render(<MemoryRouter initialEntries={["/login"]} future={{ v7_startTransition: true, v7_relativeSplatPath: true }}><App /></MemoryRouter>);
    await user.type(screen.getByLabelText("Username"), "nova");
    await user.type(screen.getByLabelText("Password"), "correct horse battery staple");
    await user.click(screen.getByRole("button", { name: "Sign in" }));
    expect(await screen.findByText(/Turn a thought into/i)).toBeInTheDocument();
    expect(screen.getByText("@nova")).toBeInTheDocument();
    expect(JSON.parse(sessionStorage.getItem("socialai.session"))).toEqual({ token: "signed-token", username: "nova" });
  });

  it("generates and publishes a local AI image", async () => {
    sessionStorage.setItem("socialai.session", JSON.stringify({ token: "signed-token", username: "nova" }));
    global.fetch
      .mockResolvedValueOnce(response({ id: "image-1", url: "http://assets.test/media/image.svg", prompt: "A moon garden" }, 201))
      .mockResolvedValueOnce(response({ id: "image-1", type: "image", message: "A moon garden" }, 201));
    const user = userEvent.setup();
    render(<MemoryRouter initialEntries={["/create"]} future={{ v7_startTransition: true, v7_relativeSplatPath: true }}><App /></MemoryRouter>);
    await user.type(screen.getByLabelText("What should we imagine?"), "A moon garden");
    await user.click(screen.getByRole("button", { name: /Generate image/i }));
    expect(await screen.findByAltText("A moon garden")).toBeInTheDocument();
    await user.click(screen.getByRole("button", { name: "Share image" }));
    expect(await screen.findByText(/Shared to your collection/i)).toBeInTheDocument();
    expect(global.fetch).toHaveBeenCalledTimes(2);
  });

  it("discards an unpublished image through the API", async () => {
    sessionStorage.setItem("socialai.session", JSON.stringify({ token: "signed-token", username: "nova" }));
    global.fetch
      .mockResolvedValueOnce(response({ id: "image-1", url: "http://assets.test/media/image.svg", prompt: "A moon garden" }, 201))
      .mockResolvedValueOnce(new Response(null, { status: 204 }));
    const user = userEvent.setup();
    render(<MemoryRouter initialEntries={["/create"]} future={{ v7_startTransition: true, v7_relativeSplatPath: true }}><App /></MemoryRouter>);
    await user.type(screen.getByLabelText("What should we imagine?"), "A moon garden");
    await user.click(screen.getByRole("button", { name: /Generate image/i }));
    expect(await screen.findByAltText("A moon garden")).toBeInTheDocument();
    await user.click(screen.getByRole("button", { name: "Discard" }));
    expect(await screen.findByText(/Try a direction/i)).toBeInTheDocument();
    expect(global.fetch).toHaveBeenLastCalledWith("/api/ai/images/image-1", expect.objectContaining({ method: "DELETE" }));
  });
});
