const KEY = "auth_token";
// Remove the legacy persistent credential. Sessions now end with the browser tab.
localStorage.removeItem(KEY);
export function saveToken(token: string): void { sessionStorage.setItem(KEY, token); window.dispatchEvent(new Event("auth-changed")); }
export function getStoredToken(): string | null {
 const token = sessionStorage.getItem(KEY);
 if (!token) return null;
 try {
  const payload = JSON.parse(atob(token.split(".")[1].replace(/-/g,"+").replace(/_/g,"/"))) as {exp?:number};
  if (!payload.exp || payload.exp * 1000 <= Date.now()) { removeToken(); return null; }
 } catch { removeToken(); return null; }
 return token;
}
export function removeToken(): void { sessionStorage.removeItem(KEY); }
export function hasToken(): boolean { return getStoredToken() !== null; }
