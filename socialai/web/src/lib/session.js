const KEY = "socialai.session";

export function readSession() {
  for (const storage of [sessionStorage, localStorage]) {
    try {
      const raw = storage.getItem(KEY);
      if (raw) {
        const session = JSON.parse(raw);
        if (session?.token && session?.username) return session;
      }
    } catch {
      storage.removeItem(KEY);
    }
  }
  return null;
}

export function saveSession(session, remember = false) {
  clearSession();
  (remember ? localStorage : sessionStorage).setItem(KEY, JSON.stringify(session));
}

export function clearSession() {
  sessionStorage.removeItem(KEY);
  localStorage.removeItem(KEY);
}
