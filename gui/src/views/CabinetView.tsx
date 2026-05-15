import { Button, NumberInput, SelectInput, TextInput } from "../components/Controls";
import { useCabinetViewModel } from "../viewmodels/useCabinetViewModel";

export function CabinetView() {
  const vm = useCabinetViewModel();

  return (
    <div className="section">
      <div className="section__title">Кабинет</div>

      <div className="section">
        <div className="section__title">Добавить предмет</div>
        <div className="inline">
          <SelectInput
            label="ID лекарства"
            value={vm.medicineId}
            onChange={vm.setMedicineId}
            options={vm.medicineOptions}
            placeholder="Пока нет данных"
          />
          <NumberInput label="Количество" value={vm.quantity} onChange={vm.setQuantity} />
          <TextInput label="Дата производства" value={vm.manufactureDate} onChange={vm.setManufactureDate} placeholder="ГГГГ-ММ-ДД" />
        </div>
        <Button onClick={vm.handleAdd}>{vm.isLoading ? "Загрузка..." : "Добавить"}</Button>
      </div>

      <div className="section">
        <div className="section__title">Обновить количество</div>
        <div className="inline">
          <SelectInput
            label="ID предмета"
            value={vm.itemId}
            onChange={vm.setItemId}
            options={vm.itemOptions}
            placeholder="Пока нет данных"
          />
          <NumberInput label="Новое количество" value={vm.newQty} onChange={vm.setNewQty} />
        </div>
        <Button onClick={vm.handleUpdateQty}>{vm.isLoading ? "Загрузка..." : "Обновить"}</Button>
      </div>

      <div className="section">
        <div className="section__title">Удалить предмет</div>
        <div className="inline">
          <SelectInput
            label="ID предмета"
            value={vm.itemId}
            onChange={vm.setItemId}
            options={vm.itemOptions}
            placeholder="Пока нет данных"
          />
          <Button variant="secondary" onClick={vm.handleRemove}>
            {vm.isLoading ? "Загрузка..." : "Удалить"}
          </Button>
        </div>
      </div>

      <div className="section">
        <div className="section__title">Текущие предметы</div>
        <div className="list">
          {vm.items.length === 0 && <div className="badge">Пусто</div>}
          {vm.items.map((item) => (
            <div key={item.id} className="badge">
              {item.id} | лекарство: {item.medicine_id} | кол-во: {item.quantity}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
