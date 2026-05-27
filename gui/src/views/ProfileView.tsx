<<<<<<< HEAD
import { Button, Checkbox, MultiSelectInput, NumberInput, SelectInput, TextInput } from "../components/Controls";
=======
import { Button, Checkbox, NumberInput, SelectInput, TagInput, TextInput } from "../components/Controls";
>>>>>>> 1f83dea7bd71d6b52bdd54933e14f6e23c6bc04a
import { useProfileViewModel } from "../viewmodels/useProfileViewModel";

export function ProfileView() {
  const vm = useProfileViewModel();

  return (
    <div className="section">
      <div className="section__title">Профиль</div>
      <TextInput label="Логин" value={vm.login} onChange={() => undefined} readOnly />
      <TextInput label="Почта" value={vm.email} onChange={vm.setEmail} />

      <div className="inline">
        <NumberInput label="Возраст" value={vm.age} onChange={vm.setAge} />
        <NumberInput label="Вес" value={vm.weight} onChange={vm.setWeight} />
      </div>

      <div className="inline">
        <label className="field">
          <span className="field__label">Пол</span>
          <select
            className="select"
            value={vm.sex ? "male" : "female"}
            onChange={(event) => vm.setSex(event.target.value === "male")}
          >
            <option value="male">Мужской</option>
            <option value="female">Женский</option>
          </select>
        </label>
        <Checkbox label="Беременность" checked={vm.isPregnant} onChange={vm.setIsPregnant} />
        <Checkbox label="Водитель" checked={vm.isDriver} onChange={vm.setIsDriver} />
      </div>

      <div className="inline">
        <MultiSelectInput
          label="Аллергии"
          values={vm.allergies}
          onChange={vm.setAllergies}
          options={vm.allergyOptions}
          placeholder="Пока нет данных"
        />
        <MultiSelectInput
          label="Болезни"
          values={vm.illnesses}
          onChange={vm.setIllnesses}
          options={vm.illnessOptions}
          placeholder="Пока нет данных"
        />
      </div>

      <div className="inline">
        <Checkbox label="Отправлять уведомления" checked={vm.notifyEnabled} onChange={vm.setNotifyEnabled} />
        <SelectInput
          label="Интервал"
          value={vm.notifyIntervalMinutes}
          options={vm.notifyIntervalOptions}
          onChange={vm.setNotifyIntervalMinutes}
        />
      </div>

      <div className="inline">
        <Checkbox label="Отправлять уведомления" checked={vm.notifyEnabled} onChange={vm.setNotifyEnabled} />
        <SelectInput
          label="Интервал"
          value={vm.notifyIntervalMinutes}
          options={vm.notifyIntervalOptions}
          onChange={vm.setNotifyIntervalMinutes}
        />
      </div>

      <div className="inline">
        <Button onClick={vm.handleUpdate}>
          {vm.isLoading ? "Загрузка..." : "Обновить профиль"}
        </Button>
      </div>
    </div>
  );
}
