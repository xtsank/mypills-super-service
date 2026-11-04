import { useCallback, useEffect, useState } from "react";
import {
  createCabinet,
  createCabinetItem,
  deleteCabinet,
  deleteCabinetItem,
  listCabinetItems,
  listCabinets,
  listMedicines,
  normalizeError,
  renameCabinet,
  updateCabinetItem
} from "../api/client";
import type { Cabinet, CabinetItem, Medicine } from "../api/types";
import { useAuth } from "../store/authStore";

export function useCabinetViewModel() {
  const { token } = useAuth();
  const [cabinets, setCabinets] = useState<Cabinet[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [items, setItems] = useState<CabinetItem[]>([]);
  const [medicines, setMedicines] = useState<Medicine[]>([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [lowStock, setLowStock] = useState(false);
  const [expired, setExpired] = useState(false);
  const refresh = useCallback(async () => {
    if (!token) return;
    setLoading(true);
    try {
      const [groups, catalog] = await Promise.all([
        listCabinets(token),
        listMedicines(token)
      ]);
      setCabinets(groups.items);
      setMedicines(catalog.items);
      setSelectedId((current) =>
        groups.items.some((item) => item.id === current)
          ? current
          : (groups.items[0]?.id ?? "")
      );
    } catch (error) {
      setMessage(normalizeError(error).message);
    } finally {
      setLoading(false);
    }
  }, [token]);
  useEffect(() => {
    void refresh();
  }, [refresh]);
  const refreshItems = useCallback(async () => {
    if (!token || !selectedId) {
      setItems([]);
      return;
    }
    try {
      setItems(
        (
          await listCabinetItems(token, selectedId, {
            low_stock: lowStock,
            expired
          })
        ).items
      );
    } catch (error) {
      setMessage(normalizeError(error).message);
    }
  }, [token, selectedId, lowStock, expired]);
  useEffect(() => {
    void refreshItems();
  }, [refreshItems]);
  async function run(action: () => Promise<unknown>) {
    setLoading(true);
    setMessage("");
    try {
      await action();
      await refresh();
      await refreshItems();
    } catch (error) {
      setMessage(normalizeError(error).message);
    } finally {
      setLoading(false);
    }
  }
  return {
    cabinets,
    selectedId,
    setSelectedId,
    items,
    medicines,
    loading,
    message,
    lowStock,
    setLowStock,
    expired,
    setExpired,
    refreshItems,
    addCabinet: (name: string) =>
      token && run(() => createCabinet(token, name)),
    renameCabinet: (name: string) =>
      token && selectedId && run(() => renameCabinet(token, selectedId, name)),
    removeCabinet: () =>
      token && selectedId && run(() => deleteCabinet(token, selectedId)),
    addItem: (
      medicineId: string,
      date: string,
      quantity: number,
      minimum: number
    ) =>
      token &&
      selectedId &&
      run(() =>
        createCabinetItem(token, selectedId, {
          medicine_id: medicineId,
          date_of_manufacture: date,
          quantity,
          minimum_quantity: minimum
        })
      ),
    editItem: (id: string, quantity: number, minimum: number) =>
      token &&
      selectedId &&
      run(() =>
        updateCabinetItem(token, selectedId, id, {
          quantity,
          minimum_quantity: minimum
        })
      ),
    removeItem: (id: string) =>
      token && selectedId && run(() => deleteCabinetItem(token, selectedId, id))
  };
}
