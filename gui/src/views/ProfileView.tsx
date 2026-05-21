import { Button, Checkbox, NumberInput, TagInput, TextInput } from "../components/Controls";
import { useProfileViewModel } from "../viewmodels/useProfileViewModel";

export function ProfileView() {
  const vm = useProfileViewModel();

  return (
    <div className="section">
      <div className="section__title">Профиль</div>
      <TextInput label="Логин" value={vm.login} onChange={() => undefined} readOnly />

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
        <TagInput label="Аллергии" value={vm.allergies} onChange={vm.setAllergies} placeholder="через запятую" />
        <TagInput label="Болезни" value={vm.illnesses} onChange={vm.setIllnesses} placeholder="через запятую" />
      </div>

      <div className="inline">
        <Button onClick={vm.handleUpdate}>
          {vm.isLoading ? "Загрузка..." : "Обновить профиль"}
        </Button>
      </div>
    </div>
  );
}
