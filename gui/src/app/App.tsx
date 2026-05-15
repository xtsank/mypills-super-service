import { useMemo, useState } from "react";
import { AuthProvider, useAuth } from "../store/authStore";
import { ProcessProvider, useProcess } from "../store/processStore";
import { AuthView } from "../views/AuthView";
import { ProfileView } from "../views/ProfileView";
import { CabinetView } from "../views/CabinetView";
import { MedicineView } from "../views/MedicineView";
import { AdminView } from "../views/AdminView";

const tabs = ["Вход", "Профиль", "Кабинет", "Рекомендации", "Админ"] as const;

type TabKey = (typeof tabs)[number];

function AppShell() {
  const [activeTab, setActiveTab] = useState<TabKey>("Вход");
  const auth = useAuth();
  const process = useProcess();

  const statusText = useMemo(() => {
    const authState = auth.isAuthenticated ? `Пользователь: ${auth.user?.login ?? "-"}` : "Не авторизован";
    const counts = `Кабинет: ${process.cabinetItemsCount}, Рекомендации: ${process.recommendationsCount}`;
    const last = process.lastAction ? `Последнее действие: ${process.lastAction}` : "Последнее действие: -";
    return `${authState} | ${counts} | ${last}`;
  }, [auth.isAuthenticated, auth.user?.login, process.cabinetItemsCount, process.recommendationsCount, process.lastAction]);

  return (
    <div className="app">
      <header className="app__header">
        <div className="app__title">MyPills — интерфейс</div>
        <div className="app__status">{statusText}</div>
        {process.statusMessage && (
          <div className={`app__banner app__banner--${process.statusType}`}>
            {process.statusMessage}
          </div>
        )}
      </header>

      <nav className="app__nav">
        {tabs.map((tab) => (
          <button
            key={tab}
            className={`app__tab ${activeTab === tab ? "app__tab--active" : ""}`}
            onClick={() => setActiveTab(tab)}
          >
            {tab}
          </button>
        ))}
        {auth.isAuthenticated && (
          <button className="app__tab" onClick={auth.logout}>
            Выход
          </button>
        )}
      </nav>

      <main className="app__main">
        {activeTab === "Вход" && <AuthView />}
        {activeTab === "Профиль" && <ProfileView />}
        {activeTab === "Кабинет" && <CabinetView />}
        {activeTab === "Рекомендации" && <MedicineView />}
        {activeTab === "Админ" && <AdminView />}
      </main>
    </div>
  );
}

export function App() {
  return (
    <AuthProvider>
      <ProcessProvider>
        <AppShell />
      </ProcessProvider>
    </AuthProvider>
  );
}
