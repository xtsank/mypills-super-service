import { useState } from "react";
import type { MedicineWrite, ReferenceItem } from "../api/types";
import { useAdminViewModel } from "../viewmodels/useAdminViewModel";

export function AdminView() {
  const vm = useAdminViewModel();
  const d = vm.draft;
  const [jsonError, setJsonError] = useState("");
  const change = <K extends keyof MedicineWrite>(
    key: K,
    value: MedicineWrite[K]
  ) => vm.setDraft({ ...d, [key]: value });
  function parseList<K extends "substances" | "dosages">(
    key: K,
    value: string
  ) {
    try {
      const parsed: unknown = JSON.parse(value);
      if (!Array.isArray(parsed)) throw new Error("Нужен массив JSON");
      change(key, parsed as MedicineWrite[K]);
      setJsonError("");
    } catch (error) {
      setJsonError(error instanceof Error ? error.message : "Ошибка JSON");
    }
  }
  return (
    <section className="panel">
      <div className="panel-head">
        <div>
          <div className="eyebrow">УПРАВЛЕНИЕ КАТАЛОГОМ</div>
          <h2>Лекарства</h2>
        </div>
        <span className="badge">Администратор</span>
      </div>
      <div className="toolbar">
        <select
          value={vm.selectedId}
          onChange={(e) => vm.select(e.target.value)}
        >
          <option value="">Новое лекарство</option>
          {vm.medicines.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
            </option>
          ))}
        </select>
        <button onClick={() => vm.select("")}>Новая запись</button>
      </div>
      <div className="form-grid">
        <label>
          Название
          <input
            value={d.name}
            onChange={(e) => change("name", e.target.value)}
          />
        </label>
        <label>
          Срок хранения, мес
          <input
            type="number"
            min="1"
            value={d.expire_time}
            onChange={(e) => change("expire_time", Number(e.target.value))}
          />
        </label>
        <label>
          Способ применения
          <input
            value={d.method_of_application}
            onChange={(e) => change("method_of_application", e.target.value)}
          />
        </label>
        <label>
          Форма
          <select
            value={d.form_id}
            onChange={(e) => change("form_id", e.target.value)}
          >
            <option value="">Выберите</option>
            {vm.forms.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          Единица измерения
          <select
            value={d.unit_id}
            onChange={(e) => change("unit_id", e.target.value)}
          >
            <option value="">Выберите</option>
            {vm.units.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div className="checks">
        <label>
          <input
            type="checkbox"
            checked={d.is_prescription}
            onChange={(e) => change("is_prescription", e.target.checked)}
          />{" "}
          По рецепту
        </label>
        <label>
          <input
            type="checkbox"
            checked={d.effect_on_pregnant}
            onChange={(e) => change("effect_on_pregnant", e.target.checked)}
          />{" "}
          Ограничение при беременности
        </label>
        <label>
          <input
            type="checkbox"
            checked={d.effect_on_driver}
            onChange={(e) => change("effect_on_driver", e.target.checked)}
          />{" "}
          Ограничение для водителей
        </label>
      </div>
      <div className="form-grid">
        <ReferenceSelect
          label="Показания"
          options={vm.illnesses}
          value={d.indications}
          onChange={(value) => change("indications", value)}
        />
        <ReferenceSelect
          label="Противопоказания"
          options={vm.illnesses}
          value={d.contraindications}
          onChange={(value) => change("contraindications", value)}
        />
      </div>
      <div className="form-grid">
        <label>
          Активные вещества (JSON)
          <textarea
            key={`${vm.selectedId}-substances`}
            defaultValue={JSON.stringify(d.substances, null, 2)}
            onBlur={(e) => parseList("substances", e.target.value)}
          />
        </label>
        <label>
          Правила дозирования (JSON)
          <textarea
            key={`${vm.selectedId}-dosages`}
            defaultValue={JSON.stringify(d.dosages, null, 2)}
            onBlur={(e) => parseList("dosages", e.target.value)}
          />
        </label>
      </div>
      <p className="muted">
        Вещества: id и concentration. Дозировки: type (age/weight), value_from,
        value_to, dosage_value, number_of_doses_per_day.
      </p>
      <p className="muted">
        Доступные вещества:{" "}
        {vm.substances.map((item) => item.name).join(", ") || "нет данных"}
      </p>
      {(jsonError || vm.message) && (
        <p className="notice">{jsonError || vm.message}</p>
      )}
      <div className="toolbar">
        <button
          className="primary"
          disabled={
            vm.loading || !!jsonError || !d.name || !d.form_id || !d.unit_id
          }
          onClick={() => void vm.save()}
        >
          Сохранить
        </button>
        <button
          className="danger"
          disabled={!vm.selectedId || vm.loading}
          onClick={() => void vm.remove()}
        >
          Удалить лекарство
        </button>
      </div>
    </section>
  );
}
function ReferenceSelect({
  label,
  options,
  value,
  onChange
}: {
  label: string;
  options: ReferenceItem[];
  value: string[];
  onChange: (value: string[]) => void;
}) {
  return (
    <label>
      {label}
      <select
        multiple
        value={value}
        onChange={(e) =>
          onChange(
            Array.from(e.target.selectedOptions, (option) => option.value)
          )
        }
      >
        {options.map((item) => (
          <option key={item.id} value={item.id}>
            {item.name}
          </option>
        ))}
      </select>
    </label>
  );
}
