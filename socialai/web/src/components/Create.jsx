import { useState } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "../App";
import { api } from "../lib/api";

const suggestions = ["A glass observatory above a lavender cloud sea", "A tiny library hidden inside an ancient tree", "Editorial portrait lit by blue hour and candlelight"];

export default function Create() {
  const { session, logout } = useAuth();
  const [prompt, setPrompt] = useState("");
  const [generated, setGenerated] = useState(null);
  const [message, setMessage] = useState("");
  const [busy, setBusy] = useState(false);
  const [sharing, setSharing] = useState(false);
  const [discarding, setDiscarding] = useState(false);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");

  function handleError(err) {
    if (err.status === 401) logout();
    setError(err.message || "Something went wrong.");
  }

  async function generate(event) {
    event?.preventDefault();
    if (prompt.trim().length < 3) return;
    setBusy(true); setError(""); setNotice("");
    try {
      if (generated) {
        try {
          await api.discard(session.token, generated.id);
        } catch (err) {
          if (err.status !== 404) throw err;
        }
      }
      setGenerated(null);
      const image = await api.generate(session.token, prompt.trim());
      setGenerated(image);
      setMessage(prompt.trim());
    } catch (err) {
      handleError(err);
    } finally {
      setBusy(false);
    }
  }

  async function discard() {
    setDiscarding(true); setError("");
    try {
      await api.discard(session.token, generated.id);
      setGenerated(null);
    } catch (err) {
      if (err.status === 404) setGenerated(null);
      else handleError(err);
    } finally {
      setDiscarding(false);
    }
  }

  async function publish() {
    setSharing(true); setError("");
    try {
      await api.publish(session.token, generated.id, message.trim());
      setNotice("Shared to your collection.");
      setGenerated(null); setPrompt(""); setMessage("");
    } catch (err) {
      handleError(err);
    } finally {
      setSharing(false);
    }
  }

  return (
    <section className="create-page">
      <div className="create-intro">
        <p className="eyebrow">AI image studio</p>
        <h1>Turn a thought into<br /><em>something you can see.</em></h1>
        <p>Describe a scene with light, mood, texture, and point of view. Your server handles generation securely—no API key reaches this browser.</p>
      </div>
      <form className="prompt-card" onSubmit={generate}>
        <label htmlFor="prompt">What should we imagine?</label>
        <textarea id="prompt" rows="4" maxLength="1000" placeholder="A quiet train crossing a salt flat at sunrise, cinematic wide angle…" value={prompt} onChange={(e) => setPrompt(e.target.value)} />
        <div className="prompt-footer"><span>{prompt.length}/1000</span><button className="button primary" disabled={busy || prompt.trim().length < 3}>{busy ? "Imagining…" : "Generate image"}<span aria-hidden="true">→</span></button></div>
      </form>
      {!generated && !busy && (
        <div className="suggestions"><span>Try a direction</span>{suggestions.map((value) => <button key={value} onClick={() => setPrompt(value)}>{value}</button>)}</div>
      )}
      {busy && <div className="generation-skeleton" role="status"><span className="spinner" /><p>Composing your image…</p><small>Local previews are instant; OpenAI can take a little longer.</small></div>}
      {generated && (
        <article className="generated-card">
          <img src={generated.url} alt={generated.prompt} />
          <div className="generated-actions">
            <div><p className="eyebrow">Ready to share</p><h2>Make it part of your collection</h2></div>
            <label>Caption<input maxLength="500" value={message} onChange={(e) => setMessage(e.target.value)} /></label>
            <div className="action-row"><button className="button ghost" disabled={discarding || sharing} onClick={discard}>{discarding ? "Discarding…" : "Discard"}</button><button className="button primary" disabled={sharing || discarding} onClick={publish}>{sharing ? "Sharing…" : "Share image"}</button></div>
          </div>
        </article>
      )}
      {error && <p className="page-message error" role="alert">{error}</p>}
      {notice && <p className="page-message success" role="status">{notice} <Link to="/collection">View it →</Link></p>}
    </section>
  );
}
