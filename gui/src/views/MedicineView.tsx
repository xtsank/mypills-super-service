import { Button, SelectInput } from "../components/Controls";
import { useMedicineViewModel } from "../viewmodels/useMedicineViewModel";

export function MedicineView() {
  const vm = useMedicineViewModel();

  return (
    <div className="section">
      <div className="section__title">Подбор лекарств</div>

      <div className="inline">
        <SelectInput
          label="ID болезни"
          value={vm.illnessId}
          onChange={vm.setIllnessId}
          options={vm.illnessOptions}
          placeholder="Пока нет данных"
        />
        <Button onClick={vm.handleSelect}>{vm.isLoading ? "Загрузка..." : "Подобрать"}</Button>
      </div>

      <div className="section">
        <div className="section__title">Результаты</div>
        <div className="list">
          {vm.recommendations.length === 0 && <div className="badge">Нет рекомендаций</div>}
          {vm.recommendations.map((rec) => (
            <div key={rec.id} className="badge">
              {rec.name} | дозировка: {rec.dosage} | частота: {rec.frequency} | ед.: {rec.unit_name} | в кабинете: {rec.quantity_in_cabinet}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
