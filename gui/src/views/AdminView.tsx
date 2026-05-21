import { Button, Checkbox, MultiSelectInput, NumberInput, SelectInput, TagInput, TextInput } from "../components/Controls";
import { useAdminViewModel } from "../viewmodels/useAdminViewModel";

export function AdminView() {
  const vm = useAdminViewModel();

  return (
    <div className="section">
      <div className="section__title">Администрирование</div>

      <div className="section">
        <div className="section__title">Добавить лекарство</div>
        <div className="inline">
          <TextInput label="Название" value={vm.name} onChange={vm.setName} />
          <SelectInput
            label="ID формы"
            value={vm.formId}
            onChange={vm.setFormId}
            options={vm.formOptions}
            placeholder="Пока нет данных"
          />
          <SelectInput
            label="ID единицы"
            value={vm.unitId}
            onChange={vm.setUnitId}
            options={vm.unitOptions}
            placeholder="Пока нет данных"
          />
        </div>
        <div className="inline">
          <TextInput label="Способ применения" value={vm.methodOfApplication} onChange={vm.setMethodOfApplication} />
          <NumberInput label="Срок годности" value={vm.expireTime} onChange={vm.setExpireTime} />
        </div>
        <div className="inline">
          <Checkbox label="По рецепту" checked={vm.isPrescription} onChange={vm.setIsPrescription} />
          <Checkbox label="Влияние на беременных" checked={vm.effectOnPregnant} onChange={vm.setEffectOnPregnant} />
          <Checkbox label="Влияние на водителей" checked={vm.effectOnDriver} onChange={vm.setEffectOnDriver} />
        </div>
        <div className="inline">
          <TagInput label="Противопоказания" value={vm.contraindications} onChange={vm.setContraindications} />
          <TagInput label="Рекомендации" value={vm.recommendations} onChange={vm.setRecommendations} />
        </div>

        <div className="section">
          <div className="section__title">Действующие вещества</div>
          {vm.substances.map((item, index) => (
            <div key={index} className="inline">
              <SelectInput
                label={`ID вещества ${index + 1}`}
                value={item.id}
                onChange={(value) => vm.updateSubstance(index, { id: value })}
                options={vm.substanceOptions}
                placeholder="Пока нет данных"
              />
              <NumberInput
                label="Концентрация"
                value={String(item.concentration)}
                onChange={(value) => vm.updateSubstance(index, { concentration: Number(value) })}
              />
              <Button variant="secondary" onClick={() => vm.removeSubstance(index)}>
                Удалить
              </Button>
            </div>
          ))}
          <Button variant="secondary" onClick={vm.addSubstance}>
            Добавить вещество
          </Button>
        </div>

        <div className="section">
          <div className="section__title">Правила дозировки</div>
          {vm.dosages.map((item, index) => (
            <div key={index} className="inline">
              <NumberInput
                label="Дозировка"
                value={String(item.dosageValue)}
                onChange={(value) => vm.updateDosage(index, { dosageValue: Number(value) })}
              />
              <NumberInput
                label="Приемов в день"
                value={String(item.numberOfDosesPerDay)}
                onChange={(value) => vm.updateDosage(index, { numberOfDosesPerDay: Number(value) })}
              />
              <SelectInput
                label="Тип"
                value={item.type}
                options={[
                  { value: "weight", label: "Вес" },
                  { value: "age", label: "Возраст" }
                ]}
                onChange={(value) => vm.updateDosage(index, { type: value as "weight" | "age" })}
              />
              <NumberInput
                label="От"
                value={String(item.valueFrom)}
                onChange={(value) => vm.updateDosage(index, { valueFrom: Number(value) })}
              />
              <NumberInput
                label="До"
                value={String(item.valueTo)}
                onChange={(value) => vm.updateDosage(index, { valueTo: Number(value) })}
              />
              <Button variant="secondary" onClick={() => vm.removeDosage(index)}>
                Удалить
              </Button>
            </div>
          ))}
          <Button variant="secondary" onClick={vm.addDosage}>
            Добавить правило
          </Button>
        </div>

        <Button onClick={vm.handleAddMedicine}>{vm.isLoading ? "Загрузка..." : "Добавить лекарство"}</Button>
      </div>

      <div className="section">
        <div className="section__title">Обновить лекарство</div>
        <div className="inline">
          <SelectInput
            label="ID лекарства"
            value={vm.updateId}
            onChange={vm.setUpdateId}
            options={vm.medicineOptions}
            placeholder="Пока нет данных"
          />
          <SelectInput
            label="ID формы"
            value={vm.updateFormId}
            onChange={vm.setUpdateFormId}
            options={vm.formOptions}
            placeholder="Пока нет данных"
          />
          <SelectInput
            label="ID единицы"
            value={vm.updateUnitId}
            onChange={vm.setUpdateUnitId}
            options={vm.unitOptions}
            placeholder="Пока нет данных"
          />
        </div>
        <div className="inline">
          <TextInput label="Способ применения" value={vm.updateMethod} onChange={vm.setUpdateMethod} />
          <NumberInput label="Срок годности" value={vm.updateExpireTime} onChange={vm.setUpdateExpireTime} />
        </div>
        <div className="inline">
          <Checkbox label="По рецепту" checked={vm.updatePrescription} onChange={vm.setUpdatePrescription} />
          <Checkbox label="Влияние на беременных" checked={vm.updatePregnant} onChange={vm.setUpdatePregnant} />
          <Checkbox label="Влияние на водителей" checked={vm.updateDriver} onChange={vm.setUpdateDriver} />
        </div>
        <Button onClick={vm.handleUpdateMedicine}>{vm.isLoading ? "Загрузка..." : "Обновить"}</Button>
      </div>

      <div className="section">
        <div className="section__title">Удалить лекарство</div>
        <div className="inline">
          <SelectInput
            label="ID лекарства"
            value={vm.deleteId}
            onChange={vm.setDeleteId}
            options={vm.medicineOptions}
            placeholder="Пока нет данных"
          />
          <Button variant="secondary" onClick={vm.handleRemoveMedicine}>
            {vm.isLoading ? "Загрузка..." : "Удалить"}
          </Button>
        </div>
      </div>

      <div className="section">
        <div className="section__title">Обновить состав</div>
        <div className="inline">
          <SelectInput
            label="ID лекарства"
            value={vm.compositionMedicineId}
            onChange={vm.setCompositionMedicineId}
            options={vm.medicineOptions}
            placeholder="Пока нет данных"
          />
        </div>
        {vm.compositionSubstances.map((item, index) => (
          <div key={index} className="inline">
            <SelectInput
              label={`ID вещества ${index + 1}`}
              value={item.id}
              onChange={(value) => vm.updateCompositionSubstance(index, { id: value })}
              options={vm.substanceOptions}
              placeholder="Пока нет данных"
            />
            <NumberInput
              label="Концентрация"
              value={String(item.concentration)}
              onChange={(value) => vm.updateCompositionSubstance(index, { concentration: Number(value) })}
            />
            <Button variant="secondary" onClick={() => vm.removeCompositionSubstance(index)}>
              Удалить
            </Button>
          </div>
        ))}
        <div className="inline">
          <Button variant="secondary" onClick={vm.addCompositionSubstance}>
            Добавить вещество
          </Button>
          <Button onClick={vm.handleUpdateComposition}>{vm.isLoading ? "Загрузка..." : "Обновить состав"}</Button>
        </div>
      </div>

      <div className="section">
        <div className="section__title">Обновить противопоказания</div>
        <div className="inline">
          <SelectInput
            label="ID лекарства"
            value={vm.contraMedicineId}
            onChange={vm.setContraMedicineId}
            options={vm.medicineOptions}
            placeholder="Пока нет данных"
          />
          <MultiSelectInput
            label="ID противопоказаний"
            values={vm.contraIds}
            onChange={vm.setContraIds}
            options={vm.contraindicationOptions}
            placeholder="Пока нет данных"
          />
        </div>
        <Button onClick={vm.handleUpdateContraindications}>{vm.isLoading ? "Загрузка..." : "Обновить"}</Button>
      </div>

      <div className="section">
        <div className="section__title">Обновить показания</div>
        <div className="inline">
          <SelectInput
            label="ID лекарства"
            value={vm.indicationMedicineId}
            onChange={vm.setIndicationMedicineId}
            options={vm.medicineOptions}
            placeholder="Пока нет данных"
          />
          <MultiSelectInput
            label="ID показаний"
            values={vm.indicationIds}
            onChange={vm.setIndicationIds}
            options={vm.indicationOptions}
            placeholder="Пока нет данных"
          />
        </div>
        <Button onClick={vm.handleUpdateIndications}>{vm.isLoading ? "Загрузка..." : "Обновить"}</Button>
      </div>

      <div className="section">
        <div className="section__title">Правила дозировки</div>
        <div className="inline">
          <SelectInput
            label="ID лекарства"
            value={vm.dosageMedicineId}
            onChange={vm.setDosageMedicineId}
            options={vm.medicineOptions}
            placeholder="Пока нет данных"
          />
          <NumberInput
            label="Дозировка"
            value={String(vm.dosageRule.dosageValue)}
            onChange={(value) => vm.updateDosageRule({ dosageValue: Number(value) })}
          />
          <NumberInput
            label="Приемов в день"
            value={String(vm.dosageRule.numberOfDosesPerDay)}
            onChange={(value) => vm.updateDosageRule({ numberOfDosesPerDay: Number(value) })}
          />
          <SelectInput
            label="Тип"
            value={vm.dosageRule.type}
            options={[
              { value: "weight", label: "Вес" },
              { value: "age", label: "Возраст" }
            ]}
            onChange={(value) => vm.setDosageType(value as "weight" | "age")}
          />
          <NumberInput
            label="От"
            value={String(vm.dosageRule.valueFrom)}
            onChange={(value) => vm.updateDosageRule({ valueFrom: Number(value) })}
          />
          <NumberInput
            label="До"
            value={String(vm.dosageRule.valueTo)}
            onChange={(value) => vm.updateDosageRule({ valueTo: Number(value) })}
          />
        </div>
        <div className="inline">
          <Button onClick={vm.handleAddDosageRule}>{vm.isLoading ? "Загрузка..." : "Добавить правило"}</Button>
        </div>
        <div className="inline">
          <SelectInput
            label="ID правила"
            value={vm.removeRuleId}
            onChange={vm.setRemoveRuleId}
            options={vm.ruleOptions}
            placeholder="Пока нет данных"
          />
          <Button variant="secondary" onClick={vm.handleRemoveDosageRule}>
            {vm.isLoading ? "Загрузка..." : "Удалить правило"}
          </Button>
        </div>
      </div>
    </div>
  );
}
