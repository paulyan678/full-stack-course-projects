import { useCallback, useEffect, useRef, useState } from "react";
import { useAuth } from "../App";
import { api } from "../lib/api";

export default function Collection() {
  const { session, logout } = useAuth();
  const [posts, setPosts] = useState([]);
  const [filter, setFilter] = useState("all");
  const [searchType, setSearchType] = useState("keywords");
  const [query, setQuery] = useState("");
  const [search, setSearch] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [selected, setSelected] = useState(null);
  const uploadDialog = useRef(null);

  const loadPosts = useCallback(async (params = search) => {
    setLoading(true); setError("");
    try {
      setPosts(await api.posts(session.token, params));
    } catch (err) {
      if (err.status === 401) logout();
      setError(err.message || "Could not load the collection.");
    } finally {
      setLoading(false);
    }
  }, [logout, search, session.token]);

  useEffect(() => { loadPosts(search); }, [loadPosts, search]);

  function submitSearch(event) {
    event.preventDefault();
    const value = query.trim();
    setSearch(value ? { [searchType]: value } : {});
  }

  async function remove(post) {
    if (!window.confirm("Delete this post? This cannot be undone.")) return;
    try {
      await api.deletePost(session.token, post.id);
      setPosts((current) => current.filter((item) => item.id !== post.id));
      setSelected(null); setNotice("Post deleted.");
    } catch (err) {
      if (err.status === 401) logout();
      setError(err.message || "Could not delete the post.");
    }
  }

  async function upload(event) {
    event.preventDefault();
    const form = event.currentTarget;
    const formData = new FormData(form);
    setError("");
    try {
      await api.upload(session.token, formData);
      form.reset(); uploadDialog.current?.close(); setNotice("Post shared.");
      setSearch({}); setQuery(""); await loadPosts({});
    } catch (err) {
      if (err.status === 401) logout();
      setError(err.message || "Could not upload the post.");
    }
  }

  const visible = posts.filter((post) => filter === "all" || post.type === filter);
  return (
    <section className="collection-page">
      <div className="collection-heading">
        <div><p className="eyebrow">The collection</p><h1>Made by people<br /><em>and possibility.</em></h1></div>
        <button className="button primary" onClick={() => uploadDialog.current?.showModal()}>Upload a post <span aria-hidden="true">＋</span></button>
      </div>
      <form className="search-bar" onSubmit={submitSearch}>
        <select aria-label="Search type" value={searchType} onChange={(e) => setSearchType(e.target.value)}><option value="keywords">Caption</option><option value="user">Creator</option></select>
        <input aria-label="Search collection" placeholder={searchType === "user" ? "Search by username" : "Search captions"} value={query} onChange={(e) => setQuery(e.target.value)} />
        <button aria-label="Search" className="search-button">⌕</button>
      </form>
      <div className="collection-controls">
        <div className="filter-tabs" aria-label="Filter media">
          {["all", "image", "video"].map((type) => <button key={type} className={filter === type ? "active" : ""} onClick={() => setFilter(type)}>{type === "all" ? "All work" : `${type}s`}</button>)}
        </div>
        <span>{visible.length} {visible.length === 1 ? "piece" : "pieces"}</span>
      </div>
      {error && <p className="page-message error" role="alert">{error}</p>}
      {notice && <p className="page-message success" role="status">{notice}</p>}
      {loading ? <div className="media-grid loading-grid" aria-label="Loading posts">{[1, 2, 3, 4].map((item) => <div className="media-skeleton" key={item} />)}</div> : visible.length === 0 ? (
        <div className="empty-state"><span aria-hidden="true">✦</span><h2>Nothing here yet</h2><p>Try a different search, or share the first piece in this view.</p></div>
      ) : (
        <div className="media-grid">{visible.map((post) => <MediaCard key={post.id} post={post} own={post.user === session.username} onOpen={setSelected} onDelete={remove} />)}</div>
      )}

      <dialog className="upload-dialog" ref={uploadDialog} onClick={(event) => { if (event.target === uploadDialog.current) uploadDialog.current.close(); }}>
        <form className="upload-form" onSubmit={upload}>
          <div className="dialog-heading"><div><p className="eyebrow">New post</p><h2>Share from your world</h2></div><button type="button" className="icon-button" aria-label="Close" onClick={() => uploadDialog.current?.close()}>×</button></div>
          <label>Caption<textarea name="message" required maxLength="500" rows="3" placeholder="Tell us about this moment…" /></label>
          <label className="file-drop">Image or video<input name="media_file" type="file" required accept="image/jpeg,image/png,image/gif,image/webp,video/mp4,video/webm,video/quicktime" /><span>JPEG, PNG, GIF, WebP, MP4, WebM, or MOV · up to 25 MB</span></label>
          <div className="action-row"><button type="button" className="button ghost" onClick={() => uploadDialog.current?.close()}>Cancel</button><button className="button primary">Share post</button></div>
        </form>
      </dialog>

      {selected && <div className="lightbox" role="dialog" aria-modal="true" aria-label="Media preview" onClick={() => setSelected(null)}><button className="lightbox-close" aria-label="Close preview">×</button><div className="lightbox-content" onClick={(e) => e.stopPropagation()}>{selected.type === "video" ? <video src={selected.url} controls autoPlay /> : <img src={selected.url} alt={selected.message} />}<div><p className="eyebrow">@{selected.user}</p><h2>{selected.message}</h2>{selected.user === session.username && <button className="button danger" onClick={() => remove(selected)}>Delete post</button>}</div></div></div>}
    </section>
  );
}

function MediaCard({ post, own, onOpen, onDelete }) {
  return (
    <article className="media-card">
      <button className="media-preview" onClick={() => onOpen(post)} aria-label={`Open ${post.message}`}>
        {post.type === "video" ? <video src={post.url} muted preload="metadata" /> : <img src={post.url} alt="" loading="lazy" />}
        {post.type === "video" && <span className="play-mark">▶</span>}
      </button>
      <div className="media-meta"><div><p>{post.message}</p><span>@{post.user} · {formatDate(post.created_at)}</span></div>{own && <button className="icon-button" aria-label={`Delete ${post.message}`} onClick={() => onDelete(post)}>×</button>}</div>
    </article>
  );
}

function formatDate(value) {
  if (!value) return "Just now";
  return new Intl.DateTimeFormat(undefined, { month: "short", day: "numeric" }).format(new Date(value));
}
