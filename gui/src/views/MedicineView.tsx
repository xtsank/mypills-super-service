import { useMedicineViewModel } from "../viewmodels/useMedicineViewModel";

export function MedicineView() {
  const vm = useMedicineViewModel();
  return (
    <section className="panel">
      <div className="panel-head">
        <div>
          <div className="eyebrow">ПОДБОР</div>
          <h2>Рекомендации</h2>
          <p className="muted">
            Показываем лекарства из ваших аптечек с учётом профиля, срока
            годности и дозировки.
          </p>
        </div>
      </div>
      <div className="toolbar">
        <select
          value={vm.illnessId}
          onChange={(e) => vm.setIllnessId(e.target.value)}
        >
          <option value="">Выберите заболевание</option>
          {vm.illnesses.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
            </option>
          ))}
        </select>
        <select
          value={vm.cabinetId}
          onChange={(e) => vm.setCabinetId(e.target.value)}
        >
          <option value="">Все аптечки</option>
          {vm.cabinets.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
            </option>
          ))}
        </select>
        <button
          className="primary"
          disabled={!vm.illnessId || vm.loading}
          onClick={() => void vm.find()}
        >
          {vm.loading ? "Подбор…" : "Подобрать"}
        </button>
      </div>
      {vm.message && <p className="notice">{vm.message}</p>}
      {vm.items.length === 0 ? (
        <p className="empty">
          Подходящие лекарства появятся здесь после поиска.
        </p>
      ) : (
        <div className="cards">
          {vm.items.map((item) => (
            <article className="result-card" key={item.cabinet_item_id}>
              <div className="eyebrow">ДОСТУПНО В АПТЕЧКЕ</div>
              <h3>{item.name}</h3>
              <p>{item.method_of_application}</p>
              <strong>
                {item.dosage} {item.unit_name} · {item.frequency_per_day}{" "}
                раз/сутки
              </strong>
              <p className="muted">Остаток: {item.quantity_in_cabinet}</p>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}
