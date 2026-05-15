import { useState } from "react";
import { addCabinetItem, normalizeError, removeCabinetItem, updateCabinetQty } from "../api/client";
import { CabinetResDto } from "../api/types";
import { useAuth } from "../store/authStore";
import { useProcess } from "../store/processStore";

export function useCabinetViewModel() {
  const auth = useAuth();
  const process = useProcess();
  const [isLoading, setIsLoading] = useState(false);
  const [items, setItems] = useState<CabinetResDto[]>([]);

  const medicineOptions: Array<{ value: string; label: string }> = [];
  const itemOptions: Array<{ value: string; label: string }> = [];

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
      const response = await addCabinetItem(auth.token, {
        medicine_id: medicineId,
        quantity: Number(quantity),
        date_of_manufacture: manufactureDate
      });
      setItems((prev) => [...prev, response]);
      process.setCabinetItemsCount(items.length + 1);
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
      const response = await updateCabinetQty(auth.token, {
        id: itemId,
        qty: Number(newQty)
      });
      setItems((prev) => prev.map((item) => (item.id === response.id ? response : item)));
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
      const updated = items.filter((item) => item.id !== itemId);
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
