import { useCallback, useEffect, useState } from "react";
import {
  createMedicine,
  deleteMedicine,
  listMedicines,
  listReference,
  normalizeError,
  updateMedicine
} from "../api/client";
import type { Medicine, MedicineWrite, ReferenceItem } from "../api/types";
import { useAuth } from "../store/authStore";

const emptyDraft: MedicineWrite = {
  name: "",
  expire_time: 24,
  is_prescription: false,
  method_of_application: "Внутрь",
  effect_on_pregnant: false,
  effect_on_driver: false,
  form_id: "",
  unit_id: "",
  indications: [],
  contraindications: [],
  substances: [],
  dosages: []
};
export function useAdminViewModel() {
  const { token } = useAuth();
  const [medicines, setMedicines] = useState<Medicine[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [draft, setDraft] = useState<MedicineWrite>(emptyDraft);
  const [forms, setForms] = useState<ReferenceItem[]>([]);
  const [units, setUnits] = useState<ReferenceItem[]>([]);
  const [illnesses, setIllnesses] = useState<ReferenceItem[]>([]);
  const [substances, setSubstances] = useState<ReferenceItem[]>([]);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const refresh = useCallback(async () => {
    if (!token) return;
    try {
      const [catalog, formPage, unitPage, illnessPage, substancePage] =
        await Promise.all([
          listMedicines(token),
          listReference("forms"),
          listReference("units"),
          listReference("illnesses"),
          listReference("substances")
        ]);
      setMedicines(catalog.items);
      setForms(formPage.items);
      setUnits(unitPage.items);
      setIllnesses(illnessPage.items);
      setSubstances(substancePage.items);
    } catch (error) {
      setMessage(normalizeError(error).message);
    }
  }, [token]);
  useEffect(() => {
    void refresh();
  }, [refresh]);
  function select(id: string) {
    setSelectedId(id);
    const medicine = medicines.find((item) => item.id === id);
    if (medicine) {
      const { id: _id, dosages, ...rest } = medicine;
      setDraft({
        ...rest,
        dosages: dosages.map(({ id: _ruleId, ...rule }) => rule)
      });
    } else setDraft(emptyDraft);
  }
  async function save() {
    if (!token) return;
    setLoading(true);
    setMessage("");
    try {
      const result = selectedId
        ? await updateMedicine(token, selectedId, draft)
        : await createMedicine(token, draft);
      await refresh();
      setSelectedId(result.id);
      setMessage("Лекарство сохранено");
    } catch (error) {
      setMessage(normalizeError(error).message);
    } finally {
      setLoading(false);
    }
  }
  async function remove() {
    if (!token || !selectedId) return;
    setLoading(true);
    setMessage("");
    try {
      await deleteMedicine(token, selectedId);
      select("");
      await refresh();
      setMessage("Лекарство удалено");
    } catch (error) {
      setMessage(normalizeError(error).message);
    } finally {
      setLoading(false);
    }
  }
  return {
    medicines,
    selectedId,
    draft,
    setDraft,
    forms,
    units,
    illnesses,
    substances,
    message,
    loading,
    select,
    save,
    remove
  };
}
