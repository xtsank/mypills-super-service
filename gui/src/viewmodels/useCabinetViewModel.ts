import { useEffect, useMemo, useState } from "react";
import {
  addCabinetItem,
  listCabinetItems,
  listMedicines,
  normalizeError,
  removeCabinetItem,
  updateCabinetQty
} from "../api/client";
import { CabinetItemDetailsResDto } from "../api/types";
import { useAuth } from "../store/authStore";
import { useProcess } from "../store/processStore";

export function useCabinetViewModel() {
  const auth = useAuth();
  const process = useProcess();
  const [isLoading, setIsLoading] = useState(false);
  const [items, setItems] = useState<CabinetItemDetailsResDto[]>([]);
  const [medicineOptions, setMedicineOptions] = useState<Array<{ value: string; label: string }>>([]);

  const itemOptions = useMemo(
    () => items.map((item) => ({ value: item.id, label: `${item.medicine_name} (${item.quantity})` })),
    [items]
  );

  const [medicineId, setMedicineId] = useState("");
  const [quantity, setQuantity] = useState("");
  const [manufactureDate, setManufactureDate] = useState("");

  const [itemId, setItemId] = useState("");
  const [newQty, setNewQty] = useState("");

  const handleAdd = async () => {
    if (!auth.token) {
      process.setStatus("error", "Нужно войти");
      return;
    }

    setIsLoading(true);
    try {
      await addCabinetItem(auth.token, {
        medicine_id: medicineId,
        quantity: Number(quantity),
        date_of_manufacture: manufactureDate
      });
      const updated = await listCabinetItems(auth.token);
      setItems(updated);
      process.setCabinetItemsCount(updated.length);
      process.setStatus("success", "Предмет добавлен");
      process.setLastAction("добавление предмета");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateQty = async () => {
    if (!auth.token) {
      process.setStatus("error", "Нужно войти");
      return;
    }

    setIsLoading(true);
    try {
      await updateCabinetQty(auth.token, {
        id: itemId,
        qty: Number(newQty)
      });
      const updated = await listCabinetItems(auth.token);
      setItems(updated);
      process.setStatus("success", "Количество обновлено");
      process.setLastAction("обновление количества");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRemove = async () => {
    if (!auth.token) {
      process.setStatus("error", "Нужно войти");
      return;
    }

    setIsLoading(true);
    try {
      await removeCabinetItem(auth.token, { id: itemId });
      const updated = await listCabinetItems(auth.token);
      setItems(updated);
      process.setCabinetItemsCount(updated.length);
      process.setStatus("success", "Предмет удален");
      process.setLastAction("удаление предмета");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const load = async () => {
      try {
        const meds = await listMedicines();
        setMedicineOptions(meds.map((item) => ({ value: item.id, label: item.name })));
      } catch (error) {
        const appError = normalizeError(error);
        process.setStatus("error", `${appError.code}: ${appError.message}`);
      }
    };
    load();
  }, [process]);

  useEffect(() => {
    const load = async () => {
      if (!auth.token) {
        return;
      }
      try {
        const itemsResponse = await listCabinetItems(auth.token);
        setItems(itemsResponse);
        process.setCabinetItemsCount(itemsResponse.length);
      } catch (error) {
        const appError = normalizeError(error);
        process.setStatus("error", `${appError.code}: ${appError.message}`);
      }
    };
    load();
  }, [auth.token, process]);

  return {
    isLoading,
    items,
    medicineOptions,
    itemOptions,
    medicineId,
    quantity,
    manufactureDate,
    itemId,
    newQty,
    setMedicineId,
    setQuantity,
    setManufactureDate,
    setItemId,
    setNewQty,
    handleAdd,
    handleUpdateQty,
    handleRemove
  };
}
