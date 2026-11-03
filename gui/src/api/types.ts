export type AppError = { code: string; message: string };
export type Page<T> = {
  items: T[];
  total: number;
  limit: number;
  offset: number;
};
export type ReferenceItem = { id: string; name: string };
export type ProfileResDto = {
  id: string;
  login: string;
  email: string;
  age: number;
  weight: number;
  sex?: boolean;
  is_pregnant: boolean;
  is_driver: boolean;
  allergies: string[];
  illnesses: string[];
};
export type AuthResDto = { token: string; user: ProfileResDto };
export type RegisterUserDto = {
  login: string;
  email: string;
  password: string;
  age: number;
  weight: number;
  sex?: boolean;
  is_pregnant?: boolean;
  is_driver?: boolean;
  allergies?: string[];
  illnesses?: string[];
};
export type ProfilePatchDto = Partial<Omit<ProfileResDto, "id" | "login">>;
export type Medicine = {
  id: string;
  name: string;
  expire_time: number;
  is_prescription: boolean;
  method_of_application: string;
  effect_on_pregnant: boolean;
  effect_on_driver: boolean;
  form_id: string;
  unit_id: string;
  indications: string[];
  contraindications: string[];
  substances: { id: string; concentration: number }[];
  dosages: {
    id: string;
    type: "age" | "weight";
    value_from: number;
    value_to: number;
    dosage_value: number;
    number_of_doses_per_day: number;
  }[];
};
export type MedicineWrite = Omit<Medicine, "id" | "dosages"> & {
  dosages: Omit<Medicine["dosages"][number], "id">[];
};
export type Cabinet = { id: string; name: string };
export type CabinetItem = {
  id: string;
  cabinet_id: string;
  medicine_id: string;
  medicine_name: string;
  date_of_manufacture: string;
  expires_at: string;
  quantity: number;
  minimum_quantity: number;
  low_stock: boolean;
};
export type Recommendation = {
  medicine_id: string;
  cabinet_id: string;
  cabinet_item_id: string;
  name: string;
  method_of_application: string;
  dosage: number;
  unit_name: string;
  frequency_per_day: number;
  quantity_in_cabinet: number;
};
export type RestockItem = {
  id: string;
  medicine_id: string;
  medicine_name: string;
  desired_quantity: number;
};
