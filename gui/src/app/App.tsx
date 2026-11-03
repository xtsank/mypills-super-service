import { useState } from "react";
import { AuthProvider, useAuth } from "../store/authStore";
import { AdminView } from "../views/AdminView";
import { AuthView } from "../views/AuthView";
import { CabinetView } from "../views/CabinetView";
import { MedicineView } from "../views/MedicineView";
import { ProfileView } from "../views/ProfileView";
import { RestockView } from "../views/RestockView";

const pages = [
  { id: "cabinets", label: "Аптечки", icon: "▦" },
  { id: "recommendations", label: "Рекомендации", icon: "✦" },
  { id: "restock", label: "Пополнение", icon: "⊕" },
  { id: "profile", label: "Профиль", icon: "◉" },
  { id: "admin", label: "Каталог · админ", icon: "⚙" }
] as const;
type Page = (typeof pages)[number]["id"];

function AppShell() {
  const auth = useAuth();
  const [page, setPage] = useState<Page>("cabinets");
  if (!auth.isAuthenticated)
    return (
      <main className="auth-layout">
        <div className="auth-art">
          <div className="logo big">
            <span className="logo-mark">✚</span> MyPills
          </div>
          <h2>Всё важное для здоровья — под рукой.</h2>
          <p>
            Следите за запасами, находите подходящие лекарства и планируйте
            покупки.
          </p>
          <div className="art-circle" />
        </div>
        <AuthView />
      </main>
    );
  return (
    <div className="shell">
      <aside className="sidebar">
        <div className="logo">
          <span className="logo-mark">✚</span> MyPills
        </div>
        <div className="sidebar-caption">МОЙ КАБИНЕТ</div>
        <nav>
          {pages
            .filter((item) => item.id !== "admin" || auth.isAdmin)
            .map((item) => (
              <button
                key={item.id}
                className={page === item.id ? "nav-item active" : "nav-item"}
                onClick={() => setPage(item.id)}
              >
                <span>{item.icon}</span>
                {item.label}
              </button>
            ))}
        </nav>
        <div className="sidebar-user">
          <div className="avatar">
            {auth.user?.login.slice(0, 1).toUpperCase()}
          </div>
          <div>
            <strong>{auth.user?.login}</strong>
            <small>{auth.user?.email}</small>
          </div>
        </div>
        <button className="logout" onClick={auth.logout}>
          Выйти из аккаунта
        </button>
      </aside>
      <main className="content">
        <header className="topbar">
          <div>
            <span className="eyebrow">ЛИЧНОЕ ПРОСТРАНСТВО</span>
            <h1>{pages.find((item) => item.id === page)?.label}</h1>
          </div>
          <span className="topbar-pill">Всё под контролем</span>
        </header>
        {page === "cabinets" && <CabinetView />}
        {page === "recommendations" && <MedicineView />}
        {page === "restock" && <RestockView />}
        {page === "profile" && <ProfileView />}
        {page === "admin" && <AdminView />}
      </main>
    </div>
  );
}
export function App() {
  return (
    <AuthProvider>
      <AppShell />
    </AuthProvider>
  );
}
