import { useProfileViewModel } from "../viewmodels/useProfileViewModel";

export function ProfileView() {
  const vm = useProfileViewModel();
  const p = vm.profile;
  if (!p) return <section className="panel">Загрузка профиля…</section>;
  const change = <K extends keyof typeof p>(key: K, value: (typeof p)[K]) =>
    vm.setProfile({ ...p, [key]: value });
  return (
    <section className="panel">
      <div className="panel-head">
        <div>
          <div className="eyebrow">АККАУНТ</div>
          <h2>Мой профиль</h2>
        </div>
        <span className="badge">{p.login}</span>
      </div>
      <div className="form-grid">
        <label>
          Почта
          <input
            type="email"
            value={p.email}
            onChange={(e) => change("email", e.target.value)}
          />
        </label>
        <label>
          Возраст
          <input
            type="number"
            min="1"
            max="120"
            value={p.age}
            onChange={(e) => change("age", Number(e.target.value))}
          />
        </label>
        <label>
          Вес, кг
          <input
            type="number"
            min="1"
            max="500"
            value={p.weight}
            onChange={(e) => change("weight", Number(e.target.value))}
          />
        </label>
      </div>
      <div className="checks">
        <label>
          <input
            type="checkbox"
            checked={p.is_pregnant}
            onChange={(e) => change("is_pregnant", e.target.checked)}
          />{" "}
          Беременность
        </label>
        <label>
          <input
            type="checkbox"
            checked={p.is_driver}
            onChange={(e) => change("is_driver", e.target.checked)}
          />{" "}
          Вожу автомобиль
        </label>
      </div>
      <div className="form-grid">
        <label>
          Заболевания
          <select
            multiple
            value={p.illnesses}
            onChange={(e) =>
              change(
                "illnesses",
                Array.from(e.target.selectedOptions, (o) => o.value)
              )
            }
          >
            {vm.illnesses.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          Аллергии
          <select
            multiple
            value={p.allergies}
            onChange={(e) =>
              change(
                "allergies",
                Array.from(e.target.selectedOptions, (o) => o.value)
              )
            }
          >
            {vm.substances.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <p className="muted">
        Для выбора нескольких значений удерживайте Ctrl или Cmd.
      </p>
      {vm.message && <p className="notice">{vm.message}</p>}
      <button
        className="primary"
        onClick={() => void vm.save()}
        disabled={vm.loading}
      >
        Сохранить изменения
      </button>
    </section>
  );
}
