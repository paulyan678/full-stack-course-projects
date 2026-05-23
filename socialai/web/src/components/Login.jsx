import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../App";
import { api } from "../lib/api";

export default function Login() {
  const [form, setForm] = useState({ username: "", password: "", remember: false });
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  async function submit(event) {
    event.preventDefault();
    setBusy(true);
    setError("");
    try {
      const session = await api.signin({ username: form.username, password: form.password });
      login(session, form.remember);
      navigate(location.state?.from || "/create", { replace: true });
    } catch (err) {
      setError(err.message || "Could not sign in.");
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="auth-layout">
      <div className="auth-story">
        <p className="eyebrow">A visual commons</p>
        <h1>Imagine it.<br /><em>Make it visible.</em></h1>
        <p>Generate new worlds with AI, share your own moments, and discover what everyone else is creating.</p>
        <div className="story-orb orb-one" /><div className="story-orb orb-two" />
      </div>
      <div className="auth-panel">
        <form className="auth-form" onSubmit={submit}>
          <p className="eyebrow">Welcome back</p>
          <h2>Sign in to SocialAI</h2>
          {location.state?.registered && <p className="form-success" role="status">Account created. You can sign in now.</p>}
          <label>Username<input autoComplete="username" required value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} /></label>
          <label>Password<input type="password" autoComplete="current-password" required minLength="10" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} /></label>
          <label className="check-row"><input type="checkbox" checked={form.remember} onChange={(e) => setForm({ ...form, remember: e.target.checked })} /> Keep me signed in on this device</label>
          {error && <p className="form-error" role="alert">{error}</p>}
          <button className="button primary wide" disabled={busy}>{busy ? "Signing in…" : "Sign in"}</button>
          <p className="form-switch">New here? <Link to="/register">Create an account</Link></p>
        </form>
      </div>
    </section>
  );
}
