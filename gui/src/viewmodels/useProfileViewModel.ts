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
  const [email, setEmail] = useState("");
  const [sex, setSex] = useState(false);
  const [isPregnant, setIsPregnant] = useState(false);
  const [isDriver, setIsDriver] = useState(false);
  const [notifyEnabled, setNotifyEnabled] = useState(false);
  const [notifyIntervalMinutes, setNotifyIntervalMinutes] = useState("1440");
  const [allergies, setAllergies] = useState<string[]>([]);
  const [illnesses, setIllnesses] = useState<string[]>([]);

  const notifyIntervalOptions = [
    { value: "1", label: "1 мин" },
    { value: "30", label: "30 мин" },
    { value: "60", label: "1 час" },
    { value: "720", label: "12 часов" },
    { value: "1440", label: "1 сутки" }
  ];

  useEffect(() => {
    if (auth.user) {
      setAge(String(auth.user.age ?? ""));
      setWeight(String(auth.user.weight ?? ""));
      setEmail(auth.user.email ?? "");
      setSex(Boolean(auth.user.sex));
      setIsPregnant(Boolean(auth.user.is_pregnant));
      setIsDriver(Boolean(auth.user.is_driver));
      setNotifyEnabled(Boolean(auth.user.notify_enabled));
      setNotifyIntervalMinutes(String(auth.user.notify_interval_minutes ?? 1440));
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
        email: email || undefined,
        age: age ? Number(age) : undefined,
        weight: weight ? Number(weight) : undefined,
        sex,
        allergies,
        illnesses,
        is_driver: isDriver,
        is_pregnant: isPregnant,
        notify_enabled: notifyEnabled,
        notify_interval_minutes: Number(notifyIntervalMinutes)
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
    email,
    age,
    weight,
    sex,
    isPregnant,
    isDriver,
    notifyEnabled,
    notifyIntervalMinutes,
    notifyIntervalOptions,
    allergies,
    illnesses,
    setEmail,
    setAge,
    setWeight,
    setSex,
    setIsPregnant,
    setIsDriver,
    setNotifyEnabled,
    setNotifyIntervalMinutes,
    setAllergies,
    setIllnesses,
    handleUpdate
  };
}
