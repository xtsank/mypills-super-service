import { useEffect, useState } from "react";
import {
  listCabinets,
  listRecommendations,
  listReference,
  normalizeError
} from "../api/client";
import type { Cabinet, Recommendation, ReferenceItem } from "../api/types";
import { useAuth } from "../store/authStore";

export function useMedicineViewModel() {
  const { token } = useAuth();
  const [illnesses, setIllnesses] = useState<ReferenceItem[]>([]);
  const [cabinets, setCabinets] = useState<Cabinet[]>([]);
  const [illnessId, setIllnessId] = useState("");
  const [cabinetId, setCabinetId] = useState("");
  const [items, setItems] = useState<Recommendation[]>([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  useEffect(() => {
    if (!token) return;
    void Promise.all([listReference("illnesses"), listCabinets(token)])
      .then(([illnessPage, cabinetPage]) => {
        setIllnesses(illnessPage.items);
        setCabinets(cabinetPage.items);
      })
      .catch((error) => setMessage(normalizeError(error).message));
  }, [token]);
  async function find() {
    if (!token || !illnessId) return;
    setLoading(true);
    setMessage("");
    try {
      setItems(
        (await listRecommendations(token, illnessId, cabinetId || undefined))
          .items
      );
    } catch (error) {
      setMessage(normalizeError(error).message);
    } finally {
      setLoading(false);
    }
  }
  return {
    illnesses,
    cabinets,
    illnessId,
    setIllnessId,
    cabinetId,
    setCabinetId,
    items,
    loading,
    message,
    find
  };
}
