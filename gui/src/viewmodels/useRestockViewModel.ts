import { useCallback, useEffect, useState } from "react";
import {
  createRestock,
  deleteRestock,
  listMedicines,
  listRestock,
  normalizeError,
  updateRestock
} from "../api/client";
import type { Medicine, RestockItem } from "../api/types";
import { useAuth } from "../store/authStore";

export function useRestockViewModel() {
  const { token } = useAuth();
  const [items, setItems] = useState<RestockItem[]>([]);
  const [medicines, setMedicines] = useState<Medicine[]>([]);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const refresh = useCallback(async () => {
    if (!token) return;
    try {
      const [restock, catalog] = await Promise.all([
        listRestock(token),
        listMedicines(token)
      ]);
      setItems(restock.items);
      setMedicines(catalog.items);
    } catch (error) {
      setMessage(normalizeError(error).message);
    }
  }, [token]);
  useEffect(() => {
    void refresh();
  }, [refresh]);
  async function run(action: () => Promise<unknown>) {
    setLoading(true);
    setMessage("");
    try {
      await action();
      await refresh();
    } catch (error) {
      setMessage(normalizeError(error).message);
    } finally {
      setLoading(false);
    }
  }
  return {
    items,
    medicines,
    message,
    loading,
    add: (id: string, quantity: number) =>
      token && run(() => createRestock(token, id, quantity)),
    update: (id: string, quantity: number) =>
      token && run(() => updateRestock(token, id, quantity)),
    remove: (id: string) => token && run(() => deleteRestock(token, id))
  };
}
