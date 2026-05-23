import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../App";

export default function Header() {
  const { session, logout } = useAuth();
  const navigate = useNavigate();
  return (
    <header className="site-header">
      <NavLink className="brand" to={session ? "/create" : "/login"} aria-label="SocialAI home">
        <span className="brand-mark" aria-hidden="true">S</span>
        <span>SocialAI</span>
      </NavLink>
      {session && (
        <>
          <nav aria-label="Primary navigation">
            <NavLink to="/create">Create</NavLink>
            <NavLink to="/collection">Collection</NavLink>
          </nav>
          <div className="account-actions">
            <span className="username">@{session.username}</span>
            <button className="button ghost small" onClick={() => { logout(); navigate("/login"); }}>Log out</button>
          </div>
        </>
      )}
    </header>
  );
}
