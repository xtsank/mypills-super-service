import { useState } from "react";
import { useAuthViewModel } from "../viewmodels/useAuthViewModel";

export function AuthView() {
  const vm = useAuthViewModel();
  const [mode, setMode] = useState<"login" | "register">("login");
  return (
    <div className="auth-card">
      <div className="eyebrow">MY PILLS · ЛИЧНЫЙ КАБИНЕТ</div>
      <h1>{mode === "login" ? "С возвращением" : "Создать аккаунт"}</h1>
      <p className="muted">
        Ваши лекарства, рекомендации и список покупок в одном месте.
      </p>
      <div className="segmented">
        <button
          className={mode === "login" ? "active" : ""}
          onClick={() => setMode("login")}
        >
          Вход
        </button>
        <button
          className={mode === "register" ? "active" : ""}
          onClick={() => setMode("register")}
        >
          Регистрация
        </button>
      </div>
      <form
        onSubmit={(event) => {
          event.preventDefault();
          void vm.submit(mode);
        }}
      >
        <label>
          Логин
          <input
            value={vm.loginValue}
            onChange={(event) => vm.setLoginValue(event.target.value)}
            required
            minLength={mode === "register" ? 8 : 1}
          />
        </label>
        <label>
          Пароль
          <input
            type="password"
            value={vm.password}
            onChange={(event) => vm.setPassword(event.target.value)}
            required
            minLength={mode === "register" ? 8 : 1}
          />
        </label>
        {mode === "register" && (
          <>
            <label>
              Электронная почта
              <input
                type="email"
                value={vm.email}
                onChange={(event) => vm.setEmail(event.target.value)}
                required
              />
            </label>
            <div className="field-row">
              <label>
                Возраст
                <input
                  type="number"
                  min="1"
                  max="120"
                  value={vm.age}
                  onChange={(event) => vm.setAge(Number(event.target.value))}
                />
              </label>
              <label>
                Вес, кг
                <input
                  type="number"
                  min="1"
                  max="500"
                  value={vm.weight}
                  onChange={(event) => vm.setWeight(Number(event.target.value))}
                />
              </label>
            </div>
          </>
        )}
        {vm.message && <p className="error">{vm.message}</p>}
        <button className="primary wide" disabled={vm.loading}>
          {vm.loading
            ? "Подождите…"
            : mode === "login"
              ? "Войти"
              : "Зарегистрироваться"}
        </button>
      </form>
    </div>
  );
}
