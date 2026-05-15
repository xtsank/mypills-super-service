import { useState } from "react";
import {
  adminAddDosageRule,
  adminAddMedicine,
  adminRemoveDosageRule,
  adminRemoveMedicine,
  adminUpdateComposition,
  adminUpdateContraindications,
  adminUpdateIndications,
  adminUpdateMedicine,
  normalizeError
} from "../api/client";
import { ActiveSubstanceDto, DosageRuleDto, DosageType } from "../api/types";
import { useAuth } from "../store/authStore";
import { useProcess } from "../store/processStore";

const emptySubstance: ActiveSubstanceDto = { id: "", concentration: 0 };
const emptyDosage: DosageRuleDto = {
  dosageValue: 0,
  numberOfDosesPerDay: 0,
  type: "weight",
  valueFrom: 0,
  valueTo: 0
};

export function useAdminViewModel() {
  const auth = useAuth();
  const process = useProcess();
  const [isLoading, setIsLoading] = useState(false);

  const [name, setName] = useState("");
  const [formId, setFormId] = useState("");
  const [unitId, setUnitId] = useState("");
  const [methodOfApplication, setMethodOfApplication] = useState("");
  const [expireTime, setExpireTime] = useState("");
  const [isPrescription, setIsPrescription] = useState(false);
  const [effectOnPregnant, setEffectOnPregnant] = useState(false);
  const [effectOnDriver, setEffectOnDriver] = useState(false);
  const [contraindications, setContraindications] = useState<string[]>([]);
  const [recommendations, setRecommendations] = useState<string[]>([]);
  const [substances, setSubstances] = useState<ActiveSubstanceDto[]>([emptySubstance]);
  const [dosages, setDosages] = useState<DosageRuleDto[]>([emptyDosage]);

  const [updateId, setUpdateId] = useState("");
  const [updateFormId, setUpdateFormId] = useState("");
  const [updateUnitId, setUpdateUnitId] = useState("");
  const [updateMethod, setUpdateMethod] = useState("");
  const [updateExpireTime, setUpdateExpireTime] = useState("");
  const [updatePrescription, setUpdatePrescription] = useState(false);
  const [updatePregnant, setUpdatePregnant] = useState(false);
  const [updateDriver, setUpdateDriver] = useState(false);

  const [deleteId, setDeleteId] = useState("");

  const [compositionMedicineId, setCompositionMedicineId] = useState("");
  const [compositionSubstances, setCompositionSubstances] = useState<ActiveSubstanceDto[]>([emptySubstance]);

  const [contraMedicineId, setContraMedicineId] = useState("");
  const [contraIds, setContraIds] = useState<string[]>([]);

  const [indicationMedicineId, setIndicationMedicineId] = useState("");
  const [indicationIds, setIndicationIds] = useState<string[]>([]);

  const [dosageMedicineId, setDosageMedicineId] = useState("");
  const [dosageRule, setDosageRule] = useState<DosageRuleDto>(emptyDosage);
  const [removeRuleId, setRemoveRuleId] = useState("");

  const requireToken = () => {
    if (!auth.token) {
      process.setStatus("error", "Нужно войти");
      return null;
    }
    return auth.token;
  };

  const handleAddMedicine = async () => {
    const token = requireToken();
    if (!token) {
      return;
    }
    setIsLoading(true);
    try {
      await adminAddMedicine(token, {
        name,
        form_id: formId,
        unit_id: unitId,
        method_of_application: methodOfApplication,
        expire_time: Number(expireTime),
        is_prescription: isPrescription,
        effect_on_pregnant: effectOnPregnant,
        effect_on_driver: effectOnDriver,
        contraindications,
        recommendations,
        substances,
        dosages
      });
      process.setStatus("success", "Лекарство добавлено");
      process.setLastAction("добавление лекарства (админ)");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateMedicine = async () => {
    const token = requireToken();
    if (!token) {
      return;
    }
    setIsLoading(true);
    try {
      await adminUpdateMedicine(token, {
        id: updateId,
        form_id: updateFormId,
        unit_id: updateUnitId,
        method_of_application: updateMethod,
        expire_time: Number(updateExpireTime),
        is_prescription: updatePrescription,
        effect_on_pregnant: updatePregnant,
        effect_on_driver: updateDriver
      });
      process.setStatus("success", "Лекарство обновлено");
      process.setLastAction("обновление лекарства (админ)");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRemoveMedicine = async () => {
    const token = requireToken();
    if (!token) {
      return;
    }
    setIsLoading(true);
    try {
      await adminRemoveMedicine(token, { id: deleteId });
      process.setStatus("success", "Лекарство удалено");
      process.setLastAction("удаление лекарства (админ)");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateComposition = async () => {
    const token = requireToken();
    if (!token) {
      return;
    }
    setIsLoading(true);
    try {
      await adminUpdateComposition(token, {
        medicine_id: compositionMedicineId,
        substances: compositionSubstances
      });
      process.setStatus("success", "Состав обновлен");
      process.setLastAction("обновление состава (админ)");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateContraindications = async () => {
    const token = requireToken();
    if (!token) {
      return;
    }
    setIsLoading(true);
    try {
      await adminUpdateContraindications(token, {
        medicine_id: contraMedicineId,
        ids: contraIds
      });
      process.setStatus("success", "Противопоказания обновлены");
      process.setLastAction("обновление противопоказаний (админ)");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateIndications = async () => {
    const token = requireToken();
    if (!token) {
      return;
    }
    setIsLoading(true);
    try {
      await adminUpdateIndications(token, {
        medicine_id: indicationMedicineId,
        ids: indicationIds
      });
      process.setStatus("success", "Показания обновлены");
      process.setLastAction("обновление показаний (админ)");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddDosageRule = async () => {
    const token = requireToken();
    if (!token) {
      return;
    }
    setIsLoading(true);
    try {
      await adminAddDosageRule(token, {
        medicine_id: dosageMedicineId,
        dosage: dosageRule
      });
      process.setStatus("success", "Правило дозировки добавлено");
      process.setLastAction("добавление дозировки (админ)");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRemoveDosageRule = async () => {
    const token = requireToken();
    if (!token) {
      return;
    }
    setIsLoading(true);
    try {
      await adminRemoveDosageRule(token, { rule_id: removeRuleId });
      process.setStatus("success", "Правило дозировки удалено");
      process.setLastAction("удаление дозировки (админ)");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const updateSubstance = (index: number, patch: Partial<ActiveSubstanceDto>) => {
    setSubstances((prev) => prev.map((item, idx) => (idx === index ? { ...item, ...patch } : item)));
  };

  const addSubstance = () => setSubstances((prev) => [...prev, { ...emptySubstance }]);
  const removeSubstance = (index: number) => setSubstances((prev) => prev.filter((_, idx) => idx !== index));

  const updateDosage = (index: number, patch: Partial<DosageRuleDto>) => {
    setDosages((prev) => prev.map((item, idx) => (idx === index ? { ...item, ...patch } : item)));
  };

  const addDosage = () => setDosages((prev) => [...prev, { ...emptyDosage }]);
  const removeDosage = (index: number) => setDosages((prev) => prev.filter((_, idx) => idx !== index));

  const updateCompositionSubstance = (index: number, patch: Partial<ActiveSubstanceDto>) => {
    setCompositionSubstances((prev) => prev.map((item, idx) => (idx === index ? { ...item, ...patch } : item)));
  };

  const addCompositionSubstance = () => setCompositionSubstances((prev) => [...prev, { ...emptySubstance }]);
  const removeCompositionSubstance = (index: number) =>
    setCompositionSubstances((prev) => prev.filter((_, idx) => idx !== index));

  const updateDosageRule = (patch: Partial<DosageRuleDto>) => {
    setDosageRule((prev) => ({ ...prev, ...patch }));
  };

  const setDosageType = (type: DosageType) => setDosageRule((prev) => ({ ...prev, type }));

  const medicineOptions: Array<{ value: string; label: string }> = [];
  const formOptions: Array<{ value: string; label: string }> = [];
  const unitOptions: Array<{ value: string; label: string }> = [];
  const substanceOptions: Array<{ value: string; label: string }> = [];
  const ruleOptions: Array<{ value: string; label: string }> = [];
  const contraindicationOptions: Array<{ value: string; label: string }> = [];
  const indicationOptions: Array<{ value: string; label: string }> = [];

  return {
    isLoading,
    medicineOptions,
    formOptions,
    unitOptions,
    substanceOptions,
    ruleOptions,
    contraindicationOptions,
    indicationOptions,
    name,
    formId,
    unitId,
    methodOfApplication,
    expireTime,
    isPrescription,
    effectOnPregnant,
    effectOnDriver,
    contraindications,
    recommendations,
    substances,
    dosages,
    updateId,
    updateFormId,
    updateUnitId,
    updateMethod,
    updateExpireTime,
    updatePrescription,
    updatePregnant,
    updateDriver,
    deleteId,
    compositionMedicineId,
    compositionSubstances,
    contraMedicineId,
    contraIds,
    indicationMedicineId,
    indicationIds,
    dosageMedicineId,
    dosageRule,
    removeRuleId,
    setName,
    setFormId,
    setUnitId,
    setMethodOfApplication,
    setExpireTime,
    setIsPrescription,
    setEffectOnPregnant,
    setEffectOnDriver,
    setContraindications,
    setRecommendations,
    setUpdateId,
    setUpdateFormId,
    setUpdateUnitId,
    setUpdateMethod,
    setUpdateExpireTime,
    setUpdatePrescription,
    setUpdatePregnant,
    setUpdateDriver,
    setDeleteId,
    setCompositionMedicineId,
    setContraMedicineId,
    setContraIds,
    setIndicationMedicineId,
    setIndicationIds,
    setDosageMedicineId,
    setRemoveRuleId,
    updateSubstance,
    addSubstance,
    removeSubstance,
    updateDosage,
    addDosage,
    removeDosage,
    updateCompositionSubstance,
    addCompositionSubstance,
    removeCompositionSubstance,
    updateDosageRule,
    setDosageType,
    handleAddMedicine,
    handleUpdateMedicine,
    handleRemoveMedicine,
    handleUpdateComposition,
    handleUpdateContraindications,
    handleUpdateIndications,
    handleAddDosageRule,
    handleRemoveDosageRule
  };
}
