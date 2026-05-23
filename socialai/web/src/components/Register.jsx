import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { api } from "../lib/api";

export default function Register() {
  const [form, setForm] = useState({ username: "", password: "", confirm: "" });
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const navigate = useNavigate();

  async function submit(event) {
    event.preventDefault();
    if (form.password !== form.confirm) {
      setError("Passwords do not match.");
      return;
    }
    setBusy(true);
    setError("");
    try {
      await api.signup({ username: form.username, password: form.password });
      navigate("/login", { replace: true, state: { registered: true } });
    } catch (err) {
      setError(err.message || "Could not create the account.");
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="auth-layout register-layout">
      <div className="auth-story">
        <p className="eyebrow">Join the studio</p>
        <h1>Your next idea<br /><em>starts here.</em></h1>
        <p>One account unlocks image generation, media sharing, search, and your personal collection.</p>
        <div className="story-orb orb-three" /><div className="story-orb orb-four" />
      </div>
      <div className="auth-panel">
        <form className="auth-form" onSubmit={submit}>
          <p className="eyebrow">Create an account</p>
          <h2>Start making</h2>
          <label>Username<input autoComplete="username" required minLength="3" maxLength="32" pattern="[A-Za-z0-9][A-Za-z0-9_-]{2,31}" title="3–32 letters, numbers, underscores, or dashes" value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} /></label>
          <label>Password<input type="password" autoComplete="new-password" required minLength="10" maxLength="128" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} /><small>Use at least 10 characters.</small></label>
          <label>Confirm password<input type="password" autoComplete="new-password" required value={form.confirm} onChange={(e) => setForm({ ...form, confirm: e.target.value })} /></label>
          {error && <p className="form-error" role="alert">{error}</p>}
          <button className="button primary wide" disabled={busy}>{busy ? "Creating account…" : "Create account"}</button>
          <p className="form-switch">Already a member? <Link to="/login">Sign in</Link></p>
        </form>
      </div>
    </section>
  );
}
