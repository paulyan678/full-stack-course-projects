import { LogoutOutlined } from "@ant-design/icons";
import { Button, Layout, Spin, message } from "antd";
import { useEffect, useState } from "react";
import { getAuthStatus, logout } from "./api/client";
import FoodList from "./components/FoodList";
import LoginForm from "./components/LoginForm";
import MyCart from "./components/MyCart";
import SignupModal from "./components/SignupModal";

const { Content, Header } = Layout;

export default function App() {
  const [session, setSession] = useState(null);
  const [cartRevision, setCartRevision] = useState(0);
  const [loggingOut, setLoggingOut] = useState(false);

  useEffect(() => {
    let active = true;
    getAuthStatus()
      .then((status) => active && setSession(status))
      .catch(() => active && setSession({ authenticated: false }))
    return () => {
      active = false;
    };
  }, []);

  async function signOut() {
    setLoggingOut(true);
    try {
      await logout();
      setSession({ authenticated: false });
      setCartRevision(0);
    } catch (error) {
      message.error(error.message);
    } finally {
      setLoggingOut(false);
    }
  }

  return (
    <Layout className="app-shell">
      <Header className="site-header">
        <div className="brand" aria-label="Lai Food home">
          <span className="brand-mark">L</span>
          <span>Lai Food</span>
        </div>
        <div className="header-actions">
          {session?.authenticated ? (
            <>
              <span className="account-email">{session.email}</span>
              <MyCart refreshKey={cartRevision} />
              <Button
                aria-label="Sign out"
                className="logout-button"
                icon={<LogoutOutlined />}
                loading={loggingOut}
                onClick={signOut}
                type="text"
              />
            </>
          ) : session ? (
            <SignupModal />
          ) : null}
        </div>
      </Header>
      <Content className="site-content">
        {session === null ? (
          <div className="page-loading"><Spin size="large" tip="Loading Lai Food…" /></div>
        ) : session.authenticated ? (
          <FoodList onCartChanged={() => setCartRevision((value) => value + 1)} />
        ) : (
          <main className="auth-page">
            <section className="hero-copy">
              <span className="eyebrow">Local kitchens, one table</span>
              <h1>Good food finds its way home.</h1>
              <p>
                Browse neighborhood favorites, build your order, and check out in a few easy steps.
              </p>
              <div className="hero-facts">
                <span>3 local restaurants</span>
                <span>Fresh menus</span>
                <span>Simple checkout</span>
              </div>
            </section>
            <section className="login-card">
              <div>
                <span className="eyebrow">Welcome back</span>
                <h2>Sign in to start your order</h2>
                <p>New here? Use “Create account” in the header.</p>
              </div>
              <LoginForm onSuccess={(email) => setSession({ authenticated: true, email })} />
            </section>
          </main>
        )}
      </Content>
    </Layout>
  );
}
