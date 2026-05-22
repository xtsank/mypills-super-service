import { useState } from "react";
import { normalizeError, selectMedicine } from "../api/client";
import { MedicineRecommendation } from "../api/types";
import { useAuth } from "../store/authStore";
import { useProcess } from "../store/processStore";

export function useMedicineViewModel() {
  const auth = useAuth();
  const process = useProcess();
  const [isLoading, setIsLoading] = useState(false);
  const [illnessId, setIllnessId] = useState("");
  const [recommendations, setRecommendations] = useState<MedicineRecommendation[]>([]);

  const illnessOptions: Array<{ value: string; label: string }> = [];

  const handleSelect = async () => {
    if (!auth.token) {
      process.setStatus("error", "Нужно войти");
      return;
    }

    setIsLoading(true);
    try {
      const response = await selectMedicine(auth.token, { illness_id: illnessId });
      setRecommendations(response.recommendations ?? []);
      process.setRecommendationsCount(response.recommendations?.length ?? 0);
      process.setStatus("success", "Рекомендации получены");
      process.setLastAction("подбор лекарств");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return {
    isLoading,
    illnessId,
    recommendations,
    illnessOptions,
    setIllnessId,
    handleSelect
  };
}
