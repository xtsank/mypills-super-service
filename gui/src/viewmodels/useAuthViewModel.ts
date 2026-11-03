import { useState } from "react";
import { login, normalizeError, register } from "../api/client";
import { useAuth } from "../store/authStore";

export function useAuthViewModel() {
  const auth = useAuth();
  const [loginValue, setLoginValue] = useState("");
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");
  const [age, setAge] = useState(30);
  const [weight, setWeight] = useState(70);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  async function submit(mode: "login" | "register") {
    setLoading(true);
    setMessage("");
    try {
      const result =
        mode === "login"
          ? await login({ login: loginValue, password })
          : await register({ login: loginValue, password, email, age, weight });
      auth.setAuth(result.token, result.user);
    } catch (error) {
      setMessage(normalizeError(error).message);
    } finally {
      setLoading(false);
    }
  }
  return {
    loginValue,
    setLoginValue,
    password,
    setPassword,
    email,
    setEmail,
    age,
    setAge,
    weight,
    setWeight,
    loading,
    message,
    submit
  };
}
