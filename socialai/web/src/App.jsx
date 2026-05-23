import { createContext, useContext, useMemo, useState } from "react";
import { Navigate, Route, Routes, useLocation } from "react-router-dom";
import Header from "./components/Header";
import Login from "./components/Login";
import Register from "./components/Register";
import Create from "./components/Create";
import Collection from "./components/Collection";
import { clearSession, readSession, saveSession } from "./lib/session";

const AuthContext = createContext(null);
export const useAuth = () => useContext(AuthContext);

function Protected({ children }) {
  const { session } = useAuth();
  const location = useLocation();
  return session ? children : <Navigate to="/login" state={{ from: location.pathname }} replace />;
}

export default function App() {
  const [session, setSession] = useState(readSession);
  const auth = useMemo(() => ({
    session,
    login(next, remember) {
      saveSession(next, remember);
      setSession(next);
    },
    logout() {
      clearSession();
      setSession(null);
    },
  }), [session]);

  return (
    <AuthContext.Provider value={auth}>
      <div className="app-shell">
        <Header />
        <main>
          <Routes>
            <Route path="/" element={<Navigate to={session ? "/create" : "/login"} replace />} />
            <Route path="/login" element={session ? <Navigate to="/create" replace /> : <Login />} />
            <Route path="/register" element={session ? <Navigate to="/create" replace /> : <Register />} />
            <Route path="/create" element={<Protected><Create /></Protected>} />
            <Route path="/collection" element={<Protected><Collection /></Protected>} />
            <Route path="*" element={<section className="not-found"><p className="eyebrow">404</p><h1>That page wandered off.</h1></section>} />
          </Routes>
        </main>
      </div>
    </AuthContext.Provider>
  );
}
