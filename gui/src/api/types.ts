export type AppError = {
  code: string;
  message: string;
};

export type ProfileResDto = {
  id: string;
  login: string;
  age: number;
  weight: number;
  sex: boolean;
  allergies: string[];
  illnesses: string[];
  is_driver: boolean;
  is_pregnant: boolean;
};

export type AuthResDto = {
  token: string;
  user: ProfileResDto;
};

export type CreateUserDto = {
  login: string;
  password: string;
  age: number;
  weight: number;
  sex: boolean;
  allergies: string[];
  illnesses: string[];
  is_driver: boolean;
  is_pregnant: boolean;
};

export type LoginUserDto = {
  login: string;
  password: string;
};

export type UpdateProfileDto = Partial<Omit<CreateUserDto, "login" | "password">>;

export type CabinetResDto = {
  id: string;
  medicine_id: string;
  quantity: number;
};

export type AddItemDto = {
  medicine_id: string;
  quantity: number;
  date_of_manufacture: string;
};

export type UpdateQtyDto = {
  id: string;
  qty: number;
};

export type RemoveItemDto = {
  id: string;
};

export type SelectMedicineDto = {
  illness_id: string;
};

export type MedicineRecommendation = {
  id: string;
  name: string;
  method_of_application: string;
  dosage: number;
  frequency: number;
  quantity_in_cabinet: number;
  unit_name: string;
};

export type MedicineResDto = {
  recommendations: MedicineRecommendation[];
};

export type ActiveSubstanceDto = {
  id: string;
  concentration: number;
};

export type DosageType = "weight" | "age";

export type DosageRuleDto = {
  dosageValue: number;
  numberOfDosesPerDay: number;
  type: DosageType;
  valueFrom: number;
  valueTo: number;
};

export type AddDosageRuleDto = {
  medicine_id: string;
  dosage: DosageRuleDto;
};

export type RemoveDosageRuleDto = {
  rule_id: string;
};

export type AddMedicineDto = {
  name: string;
  form_id: string;
  unit_id: string;
  method_of_application: string;
  expire_time: number;
  is_prescription: boolean;
  effect_on_pregnant: boolean;
  effect_on_driver: boolean;
  contraindications: string[];
  recommendations: string[];
  substances: ActiveSubstanceDto[];
  dosages: DosageRuleDto[];
};

export type UpdateMedicineDto = {
  id: string;
  form_id: string;
  unit_id: string;
  method_of_application: string;
  expire_time: number;
  is_prescription: boolean;
  effect_on_pregnant: boolean;
  effect_on_driver: boolean;
};

export type RemoveMedicineDto = {
  id: string;
};

export type UpdateCompositionDto = {
  medicine_id: string;
  substances: ActiveSubstanceDto[];
};

export type UpdateLinksDto = {
  medicine_id: string;
  ids: string[];
};

export type SuccessResDto = {
  status: string;
};

