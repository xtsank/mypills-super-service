import { request } from "./http";
import type {
  AppError,
  AuthResDto,
  Cabinet,
  CabinetItem,
  Medicine,
  MedicineWrite,
  Page,
  ProfilePatchDto,
  ProfileResDto,
  Recommendation,
  ReferenceItem,
  RegisterUserDto,
  RestockItem
} from "./types";

const API_PREFIX = "/api/v2";

function query(values: Record<string, string | number | undefined>): string {
  const params = new URLSearchParams();
  Object.entries(values).forEach(([key, value]) => {
    if (value !== undefined && value !== "") {
      params.set(key, String(value));
    }
  });
  return params.size ? `?${params}` : "";
}

export async function login(payload: {
  login: string;
  password: string;
}): Promise<AuthResDto> {
  return request<AuthResDto>(`${API_PREFIX}/sessions`, {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function register(payload: RegisterUserDto): Promise<AuthResDto> {
  return request<AuthResDto>(`${API_PREFIX}/users`, {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function getProfile(token: string): Promise<ProfileResDto> {
  return request<ProfileResDto>(`${API_PREFIX}/users/me`, { token });
}

export async function updateProfile(
  token: string,
  payload: ProfilePatchDto
): Promise<ProfileResDto> {
  return request<ProfileResDto>(`${API_PREFIX}/users/me`, {
    method: "PATCH",
    token,
    body: JSON.stringify(payload)
  });
}

export async function listReference(
  kind: "illnesses" | "substances" | "forms" | "units"
): Promise<Page<ReferenceItem>> {
  return request<Page<ReferenceItem>>(
    `${API_PREFIX}/reference/${kind}?limit=100`
  );
}

export async function listMedicines(
  token: string,
  q?: string
): Promise<Page<Medicine>> {
  return request<Page<Medicine>>(
    `${API_PREFIX}/medicines${query({ q, limit: 100 })}`,
    { token }
  );
}

export async function createMedicine(
  token: string,
  payload: MedicineWrite
): Promise<Medicine> {
  return request<Medicine>(`${API_PREFIX}/medicines`, {
    method: "POST",
    token,
    body: JSON.stringify(payload)
  });
}

export async function updateMedicine(
  token: string,
  id: string,
  payload: Partial<MedicineWrite>
): Promise<Medicine> {
  return request<Medicine>(`${API_PREFIX}/medicines/${id}`, {
    method: "PATCH",
    token,
    body: JSON.stringify(payload)
  });
}

export async function deleteMedicine(token: string, id: string): Promise<void> {
  return request<void>(`${API_PREFIX}/medicines/${id}`, {
    method: "DELETE",
    token
  });
}

export async function listCabinets(token: string): Promise<Page<Cabinet>> {
  return request<Page<Cabinet>>(`${API_PREFIX}/cabinets?limit=100`, { token });
}

export async function createCabinet(
  token: string,
  name: string
): Promise<Cabinet> {
  return request<Cabinet>(`${API_PREFIX}/cabinets`, {
    method: "POST",
    token,
    body: JSON.stringify({ name })
  });
}

export async function renameCabinet(
  token: string,
  id: string,
  name: string
): Promise<Cabinet> {
  return request<Cabinet>(`${API_PREFIX}/cabinets/${id}`, {
    method: "PATCH",
    token,
    body: JSON.stringify({ name })
  });
}

export async function deleteCabinet(token: string, id: string): Promise<void> {
  return request<void>(`${API_PREFIX}/cabinets/${id}`, {
    method: "DELETE",
    token
  });
}

export async function listCabinetItems(
  token: string,
  cabinetId: string,
  filters: { low_stock?: boolean; expired?: boolean } = {}
): Promise<Page<CabinetItem>> {
  const path = `${API_PREFIX}/cabinets/${cabinetId}/items`;
  const params = query({
    limit: 100,
    low_stock: filters.low_stock ? "true" : undefined,
    expired: filters.expired ? "true" : undefined
  });
  return request<Page<CabinetItem>>(`${path}${params}`, { token });
}

export async function createCabinetItem(
  token: string,
  cabinetId: string,
  payload: Pick<
    CabinetItem,
    "medicine_id" | "date_of_manufacture" | "quantity" | "minimum_quantity"
  >
): Promise<CabinetItem> {
  return request<CabinetItem>(`${API_PREFIX}/cabinets/${cabinetId}/items`, {
    method: "POST",
    token,
    body: JSON.stringify(payload)
  });
}

export async function updateCabinetItem(
  token: string,
  cabinetId: string,
  id: string,
  payload: Partial<
    Pick<CabinetItem, "quantity" | "minimum_quantity" | "date_of_manufacture">
  >
): Promise<CabinetItem> {
  return request<CabinetItem>(
    `${API_PREFIX}/cabinets/${cabinetId}/items/${id}`,
    {
      method: "PATCH",
      token,
      body: JSON.stringify(payload)
    }
  );
}

export async function deleteCabinetItem(
  token: string,
  cabinetId: string,
  id: string
): Promise<void> {
  return request<void>(`${API_PREFIX}/cabinets/${cabinetId}/items/${id}`, {
    method: "DELETE",
    token
  });
}

export async function listRecommendations(
  token: string,
  illnessId: string,
  cabinetId?: string
): Promise<{ items: Recommendation[] }> {
  return request<{ items: Recommendation[] }>(
    `${API_PREFIX}/recommendations${query({ illness_id: illnessId, cabinet_id: cabinetId })}`,
    { token }
  );
}

export async function listRestock(token: string): Promise<Page<RestockItem>> {
  return request<Page<RestockItem>>(`${API_PREFIX}/restock-list?limit=100`, {
    token
  });
}

export async function createRestock(
  token: string,
  medicineId: string,
  desiredQuantity: number
): Promise<RestockItem> {
  return request<RestockItem>(`${API_PREFIX}/restock-list`, {
    method: "POST",
    token,
    body: JSON.stringify({
      medicine_id: medicineId,
      desired_quantity: desiredQuantity
    })
  });
}

export async function updateRestock(
  token: string,
  id: string,
  desiredQuantity: number
): Promise<RestockItem> {
  return request<RestockItem>(`${API_PREFIX}/restock-list/${id}`, {
    method: "PATCH",
    token,
    body: JSON.stringify({ desired_quantity: desiredQuantity })
  });
}

export async function deleteRestock(token: string, id: string): Promise<void> {
  return request<void>(`${API_PREFIX}/restock-list/${id}`, {
    method: "DELETE",
    token
  });
}

export function normalizeError(error: unknown): AppError {
  if (
    typeof error === "object" &&
    error &&
    "code" in error &&
    "message" in error
  ) {
    return error as AppError;
  }
  return {
    code: "UNKNOWN",
    message: error instanceof Error ? error.message : "Неожиданная ошибка"
  };
}
