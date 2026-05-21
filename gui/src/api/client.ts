import { request } from "./http";
import {
  AddDosageRuleDto,
  AddItemDto,
  AddMedicineDto,
  AppError,
  AuthResDto,
  CabinetResDto,
  CreateUserDto,
  LoginUserDto,
  MedicineResDto,
  RemoveDosageRuleDto,
  RemoveItemDto,
  RemoveMedicineDto,
  SelectMedicineDto,
  SuccessResDto,
  UpdateCompositionDto,
  UpdateLinksDto,
  UpdateMedicineDto,
  UpdateProfileDto,
  UpdateQtyDto,
  ProfileResDto
} from "./types";

export async function login(payload: LoginUserDto): Promise<AuthResDto> {
  return request<AuthResDto>("/auth/login", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function register(payload: CreateUserDto): Promise<AuthResDto> {
  return request<AuthResDto>("/auth/register", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function updateProfile(token: string, payload: UpdateProfileDto): Promise<ProfileResDto> {
  return request<ProfileResDto>("/profile/me", {
    method: "PATCH",
    token,
    body: JSON.stringify(payload)
  });
}

export async function addCabinetItem(token: string, payload: AddItemDto): Promise<CabinetResDto> {
  return request<CabinetResDto>("/cabinet/items", {
    method: "POST",
    token,
    body: JSON.stringify(payload)
  });
}

export async function updateCabinetQty(token: string, payload: UpdateQtyDto): Promise<CabinetResDto> {
  return request<CabinetResDto>("/cabinet/items", {
    method: "PATCH",
    token,
    body: JSON.stringify(payload)
  });
}

export async function removeCabinetItem(token: string, payload: RemoveItemDto): Promise<SuccessResDto> {
  return request<SuccessResDto>("/cabinet/items", {
    method: "DELETE",
    token,
    body: JSON.stringify(payload)
  });
}

export async function selectMedicine(token: string, payload: SelectMedicineDto): Promise<MedicineResDto> {
  return request<MedicineResDto>("/medicine/select", {
    method: "POST",
    token,
    body: JSON.stringify(payload)
  });
}

export async function adminAddMedicine(token: string, payload: AddMedicineDto) {
  return request<unknown>("/admin/medicine", {
    method: "POST",
    token,
    body: JSON.stringify(payload)
  });
}

export async function adminUpdateMedicine(token: string, payload: UpdateMedicineDto) {
  return request<unknown>("/admin/medicine", {
    method: "PATCH",
    token,
    body: JSON.stringify(payload)
  });
}

export async function adminRemoveMedicine(token: string, payload: RemoveMedicineDto) {
  return request<SuccessResDto>("/admin/medicine", {
    method: "DELETE",
    token,
    body: JSON.stringify(payload)
  });
}

export async function adminUpdateComposition(token: string, payload: UpdateCompositionDto) {
  return request<SuccessResDto>("/admin/medicine/composition", {
    method: "PATCH",
    token,
    body: JSON.stringify(payload)
  });
}

export async function adminUpdateContraindications(token: string, payload: UpdateLinksDto) {
  return request<SuccessResDto>("/admin/medicine/contraindications", {
    method: "PATCH",
    token,
    body: JSON.stringify(payload)
  });
}

export async function adminUpdateIndications(token: string, payload: UpdateLinksDto) {
  return request<SuccessResDto>("/admin/medicine/indications", {
    method: "PATCH",
    token,
    body: JSON.stringify(payload)
  });
}

export async function adminAddDosageRule(token: string, payload: AddDosageRuleDto) {
  return request<SuccessResDto>("/admin/medicine/dosage", {
    method: "POST",
    token,
    body: JSON.stringify(payload)
  });
}

export async function adminRemoveDosageRule(token: string, payload: RemoveDosageRuleDto) {
  return request<SuccessResDto>("/admin/medicine/dosage", {
    method: "DELETE",
    token,
    body: JSON.stringify(payload)
  });
}

export function normalizeError(error: unknown): AppError {
  if (typeof error === "object" && error && "code" in error && "message" in error) {
    return error as AppError;
  }

  if (error instanceof Error) {
    return {
      code: "UNKNOWN",
      message: error.message
    };
  }

  return {
    code: "UNKNOWN",
    message: "Неожиданная ошибка"
  };
}
