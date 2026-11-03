import { useState } from "react";
import type { RestockItem } from "../api/types";
import { useRestockViewModel } from "../viewmodels/useRestockViewModel";

export function RestockView() {
  const vm = useRestockViewModel();
  const [medicineId, setMedicineId] = useState("");
  const [quantity, setQuantity] = useState(1);
  return (
    <div className="stack">
      <section className="panel">
        <div className="panel-head">
          <div>
            <div className="eyebrow">ПЛАН ПОКУПОК</div>
            <h2>Список пополнения</h2>
          </div>
          <span className="badge">{vm.items.length} позиций</span>
        </div>
        {vm.message && <p className="notice">{vm.message}</p>}
        {vm.items.length === 0 ? (
          <p className="empty">Список пока пуст. Добавьте лекарство ниже.</p>
        ) : (
          <div className="item-list">
            {vm.items.map((item) => (
              <RestockRow
                key={item.id}
                item={item}
                save={(value) => void vm.update(item.id, value)}
                remove={() => void vm.remove(item.id)}
              />
            ))}
          </div>
        )}
      </section>
      <section className="panel">
        <div className="panel-head">
          <h2>Добавить в список</h2>
        </div>
        <div className="toolbar">
          <select
            value={medicineId}
            onChange={(e) => setMedicineId(e.target.value)}
          >
            <option value="">Лекарство</option>
            {vm.medicines.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name}
              </option>
            ))}
          </select>
          <input
            type="number"
            min="0.01"
            step="any"
            value={quantity}
            onChange={(e) => setQuantity(Number(e.target.value))}
          />
          <button
            className="primary"
            disabled={!medicineId || quantity <= 0 || vm.loading}
            onClick={() => void vm.add(medicineId, quantity)}
          >
            Добавить
          </button>
        </div>
      </section>
    </div>
  );
}
function RestockRow({
  item,
  save,
  remove
}: {
  item: RestockItem;
  save: (value: number) => void;
  remove: () => void;
}) {
  const [quantity, setQuantity] = useState(item.desired_quantity);
  return (
    <div className="item-row">
      <strong>{item.medicine_name}</strong>
      <label>
        Нужно
        <input
          type="number"
          min="0.01"
          step="any"
          value={quantity}
          onChange={(e) => setQuantity(Number(e.target.value))}
        />
      </label>
      <button onClick={() => save(quantity)}>Сохранить</button>
      <button className="danger" onClick={remove}>
        Удалить
      </button>
    </div>
  );
}
