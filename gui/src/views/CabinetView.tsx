import { useState } from "react";
import { useCabinetViewModel } from "../viewmodels/useCabinetViewModel";

export function CabinetView() {
  const vm = useCabinetViewModel();
  const [name, setName] = useState("");
  const [medicineId, setMedicineId] = useState("");
  const [date, setDate] = useState(new Date().toISOString().slice(0, 10));
  const [quantity, setQuantity] = useState(1);
  const [minimum, setMinimum] = useState(0);
  return (
    <div className="stack">
      <section className="panel">
        <div className="panel-head">
          <div>
            <div className="eyebrow">МОИ ЗАПАСЫ</div>
            <h2>Аптечки</h2>
          </div>
          <span className="badge">{vm.cabinets.length} шт.</span>
        </div>
        <div className="toolbar">
          <select
            value={vm.selectedId}
            onChange={(e) => vm.setSelectedId(e.target.value)}
          >
            <option value="">Выберите аптечку</option>
            {vm.cabinets.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name}
              </option>
            ))}
          </select>
          <input
            placeholder="Название аптечки"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <button
            onClick={() => {
              void vm.addCabinet(name);
              setName("");
            }}
            disabled={!name.trim() || vm.loading}
          >
            Добавить
          </button>
          <button
            onClick={() => void vm.renameCabinet(name)}
            disabled={!vm.selectedId || !name.trim()}
          >
            Переименовать
          </button>
          <button
            className="danger"
            onClick={() => void vm.removeCabinet()}
            disabled={!vm.selectedId}
          >
            Удалить
          </button>
        </div>
        {vm.message && <p className="notice">{vm.message}</p>}
      </section>
      {vm.selectedId && (
        <>
          <section className="panel">
            <div className="panel-head">
              <div>
                <div className="eyebrow">СОДЕРЖИМОЕ</div>
                <h2>Лекарства в аптечке</h2>
              </div>
              <div className="checks compact">
                <label>
                  <input
                    type="checkbox"
                    checked={vm.lowStock}
                    onChange={(e) => vm.setLowStock(e.target.checked)}
                  />{" "}
                  Мало
                </label>
                <label>
                  <input
                    type="checkbox"
                    checked={vm.expired}
                    onChange={(e) => vm.setExpired(e.target.checked)}
                  />{" "}
                  Просрочено
                </label>
              </div>
            </div>
            {vm.items.length === 0 ? (
              <p className="empty">Пока нет подходящих позиций.</p>
            ) : (
              <div className="item-list">
                {vm.items.map((item) => (
                  <CabinetItemRow
                    key={item.id}
                    item={item}
                    save={(q, m) => void vm.editItem(item.id, q, m)}
                    remove={() => void vm.removeItem(item.id)}
                  />
                ))}
              </div>
            )}
          </section>
          <section className="panel">
            <div className="panel-head">
              <div>
                <div className="eyebrow">НОВАЯ ПОЗИЦИЯ</div>
                <h2>Добавить лекарство</h2>
              </div>
            </div>
            <div className="form-grid">
              <label>
                Лекарство
                <select
                  value={medicineId}
                  onChange={(e) => setMedicineId(e.target.value)}
                >
                  <option value="">Выберите из каталога</option>
                  {vm.medicines.map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.name}
                    </option>
                  ))}
                </select>
              </label>
              <label>
                Дата изготовления
                <input
                  type="date"
                  value={date}
                  onChange={(e) => setDate(e.target.value)}
                />
              </label>
              <label>
                Количество
                <input
                  type="number"
                  min="0.01"
                  step="any"
                  value={quantity}
                  onChange={(e) => setQuantity(Number(e.target.value))}
                />
              </label>
              <label>
                Минимальный остаток
                <input
                  type="number"
                  min="0"
                  step="any"
                  value={minimum}
                  onChange={(e) => setMinimum(Number(e.target.value))}
                />
              </label>
            </div>
            <button
              className="primary"
              disabled={!medicineId || quantity <= 0}
              onClick={() =>
                void vm.addItem(medicineId, date, quantity, minimum)
              }
            >
              Добавить в аптечку
            </button>
          </section>
        </>
      )}
    </div>
  );
}

function CabinetItemRow({
  item,
  save,
  remove
}: {
  item: {
    medicine_name: string;
    date_of_manufacture: string;
    expires_at: string;
    quantity: number;
    minimum_quantity: number;
    low_stock: boolean;
  };
  save: (q: number, m: number) => void;
  remove: () => void;
}) {
  const [quantity, setQuantity] = useState(item.quantity);
  const [minimum, setMinimum] = useState(item.minimum_quantity);
  return (
    <div className="item-row">
      <div>
        <strong>{item.medicine_name}</strong>
        <p className="muted">
          Изготовлено {item.date_of_manufacture} · Годно до {item.expires_at}
        </p>
      </div>
      {item.low_stock && <span className="badge warning">Мало</span>}
      <label>
        Остаток
        <input
          type="number"
          min="0.01"
          step="any"
          value={quantity}
          onChange={(e) => setQuantity(Number(e.target.value))}
        />
      </label>
      <label>
        Минимум
        <input
          type="number"
          min="0"
          step="any"
          value={minimum}
          onChange={(e) => setMinimum(Number(e.target.value))}
        />
      </label>
      <button onClick={() => save(quantity, minimum)}>Сохранить</button>
      <button className="danger" onClick={remove}>
        Удалить
      </button>
    </div>
  );
}
