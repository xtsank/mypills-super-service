import { useEffect, useState } from "react";
import { normalizeError, updateProfile } from "../api/client";
import { useAuth } from "../store/authStore";
import { useProcess } from "../store/processStore";

export function useProfileViewModel() {
  const auth = useAuth();
  const process = useProcess();
  const [isLoading, setIsLoading] = useState(false);

  const [age, setAge] = useState("");
  const [weight, setWeight] = useState("");
  const [sex, setSex] = useState(false);
  const [isPregnant, setIsPregnant] = useState(false);
  const [isDriver, setIsDriver] = useState(false);
  const [allergies, setAllergies] = useState<string[]>([]);
  const [illnesses, setIllnesses] = useState<string[]>([]);

  useEffect(() => {
    if (auth.user) {
      setAge(String(auth.user.age ?? ""));
      setWeight(String(auth.user.weight ?? ""));
      setSex(Boolean(auth.user.sex));
      setIsPregnant(Boolean(auth.user.is_pregnant));
      setIsDriver(Boolean(auth.user.is_driver));
      setAllergies(auth.user.allergies ?? []);
      setIllnesses(auth.user.illnesses ?? []);
    }
  }, [auth.user]);

  const handleUpdate = async () => {
    if (!auth.token) {
      process.setStatus("error", "Нужно войти");
      return;
    }

    setIsLoading(true);
    try {
      const updated = await updateProfile(auth.token, {
        age: age ? Number(age) : undefined,
        weight: weight ? Number(weight) : undefined,
        sex,
        allergies,
        illnesses,
        is_driver: isDriver,
        is_pregnant: isPregnant
      });
      auth.setAuth(auth.token, updated);
      process.setStatus("success", "Профиль обновлен");
      process.setLastAction("обновление профиля");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return {
    isLoading,
    isAuthenticated: auth.isAuthenticated,
    login: auth.user?.login ?? "-",
    age,
    weight,
    sex,
    isPregnant,
    isDriver,
    allergies,
    illnesses,
    setAge,
    setWeight,
    setSex,
    setIsPregnant,
    setIsDriver,
    setAllergies,
    setIllnesses,
    handleUpdate
  };
}
