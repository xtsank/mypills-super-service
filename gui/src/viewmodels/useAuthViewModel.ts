import { useState } from "react";
import { login, normalizeError, register } from "../api/client";
import { useAuth } from "../store/authStore";
import { useProcess } from "../store/processStore";

export function useAuthViewModel() {
  const auth = useAuth();
  const process = useProcess();
  const [isLoading, setIsLoading] = useState(false);
  const [loginValue, setLoginValue] = useState("");
  const [passwordValue, setPasswordValue] = useState("");

  const [age, setAge] = useState("");
  const [weight, setWeight] = useState("");
  const [sex, setSex] = useState(false);
  const [isPregnant, setIsPregnant] = useState(false);
  const [isDriver, setIsDriver] = useState(false);
  const [allergies, setAllergies] = useState<string[]>([]);
  const [illnesses, setIllnesses] = useState<string[]>([]);

  const handleLogin = async () => {
    setIsLoading(true);
    try {
      const response = await login({ login: loginValue, password: passwordValue });
      auth.setAuth(response.token, response.user);
      process.setStatus("success", "Вход выполнен");
      process.setLastAction("вход");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRegister = async () => {
    setIsLoading(true);
    try {
      const response = await register({
        login: loginValue,
        password: passwordValue,
        age: Number(age),
        weight: Number(weight),
        sex,
        allergies,
        illnesses,
        is_driver: isDriver,
        is_pregnant: isPregnant
      });
      auth.setAuth(response.token, response.user);
      process.setStatus("success", "Регистрация завершена");
      process.setLastAction("регистрация");
    } catch (error) {
      const appError = normalizeError(error);
      process.setStatus("error", `${appError.code}: ${appError.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return {
    isLoading,
    loginValue,
    passwordValue,
    age,
    weight,
    sex,
    isPregnant,
    isDriver,
    allergies,
    illnesses,
    setLoginValue,
    setPasswordValue,
    setAge,
    setWeight,
    setSex,
    setIsPregnant,
    setIsDriver,
    setAllergies,
    setIllnesses,
    handleLogin,
    handleRegister
  };
}
