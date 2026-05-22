import { useEffect, useState } from "react";
import { listIllnesses, listSubstances, login, normalizeError, register } from "../api/client";
import { useAuth } from "../store/authStore";
import { useProcess } from "../store/processStore";

export function useAuthViewModel() {
  const auth = useAuth();
  const process = useProcess();
  const [isLoading, setIsLoading] = useState(false);
  const [loginValue, setLoginValue] = useState("");
  const [emailValue, setEmailValue] = useState("");
  const [passwordValue, setPasswordValue] = useState("");

  const [registerLogin, setRegisterLogin] = useState("");
  const [registerPassword, setRegisterPassword] = useState("");

  const [age, setAge] = useState("");
  const [weight, setWeight] = useState("");
  const [sex, setSex] = useState(false);
  const [isPregnant, setIsPregnant] = useState(false);
  const [isDriver, setIsDriver] = useState(false);
  const [allergies, setAllergies] = useState<string[]>([]);
  const [illnesses, setIllnesses] = useState<string[]>([]);

  const [allergyOptions, setAllergyOptions] = useState<Array<{ value: string; label: string }>>([]);
  const [illnessOptions, setIllnessOptions] = useState<Array<{ value: string; label: string }>>([]);

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
<<<<<<< HEAD
        login: registerLogin,
        email: emailValue,
        password: registerPassword,
=======
        login: loginValue,
        email: emailValue,
        password: passwordValue,
>>>>>>> 1f83dea7bd71d6b52bdd54933e14f6e23c6bc04a
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

  useEffect(() => {
    const load = async () => {
      try {
        const [illnessList, substanceList] = await Promise.all([listIllnesses(), listSubstances()]);
        setIllnessOptions(illnessList.map((item) => ({ value: item.id, label: item.name })));
        setAllergyOptions(substanceList.map((item) => ({ value: item.id, label: item.name })));
      } catch (error) {
        const appError = normalizeError(error);
        process.setStatus("error", `${appError.code}: ${appError.message}`);
      }
    };
    load();
  }, [process]);

  return {
    isLoading,
    loginValue,
    emailValue,
    passwordValue,
    registerLogin,
    registerPassword,
    age,
    weight,
    sex,
    isPregnant,
    isDriver,
    allergies,
    illnesses,
    allergyOptions,
    illnessOptions,
    setLoginValue,
    setEmailValue,
    setPasswordValue,
    setRegisterLogin,
    setRegisterPassword,
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
