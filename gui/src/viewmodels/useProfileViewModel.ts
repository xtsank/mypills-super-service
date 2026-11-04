import { useEffect, useState } from "react";
import {
  getProfile,
  listReference,
  normalizeError,
  updateProfile
} from "../api/client";
import type { ProfileResDto, ReferenceItem } from "../api/types";
import { useAuth } from "../store/authStore";

export function useProfileViewModel() {
  const auth = useAuth();
  const [profile, setProfile] = useState<ProfileResDto | null>(auth.user);
  const [illnesses, setIllnesses] = useState<ReferenceItem[]>([]);
  const [substances, setSubstances] = useState<ReferenceItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  useEffect(() => {
    if (!auth.token) return;
    void Promise.all([
      getProfile(auth.token),
      listReference("illnesses"),
      listReference("substances")
    ])
      .then(([user, illnessPage, substancePage]) => {
        setProfile(user);
        setIllnesses(illnessPage.items);
        setSubstances(substancePage.items);
      })
      .catch((error) => setMessage(normalizeError(error).message));
  }, [auth.token]);
  async function save() {
    if (!auth.token || !profile) return;
    setLoading(true);
    try {
      const updated = await updateProfile(auth.token, {
        email: profile.email,
        age: profile.age,
        weight: profile.weight,
        sex: profile.sex,
        is_pregnant: profile.is_pregnant,
        is_driver: profile.is_driver,
        allergies: profile.allergies,
        illnesses: profile.illnesses
      });
      auth.setAuth(auth.token, updated);
      setProfile(updated);
      setMessage("Профиль сохранён");
    } catch (error) {
      setMessage(normalizeError(error).message);
    } finally {
      setLoading(false);
    }
  }
  return { profile, setProfile, illnesses, substances, loading, message, save };
}
