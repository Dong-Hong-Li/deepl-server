import type {
  ApiError,
  Language,
  RephraseTextRequest,
  RephraseTextResponse,
  TranslateTextRequest,
  TranslateTextResponse,
} from "./types";
import { translateApiError } from "../i18n/zh";

const baseURL = import.meta.env.VITE_API_BASE_URL ?? "";
const authToken = import.meta.env.VITE_AUTH_TOKEN ?? "";

function normalizeLanguage(raw: Record<string, string>): Language {
  return {
    code: raw.code ?? raw.Code ?? "",
    name: raw.name ?? raw.Name ?? raw.code ?? raw.Code ?? "",
  };
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers);
  if (init?.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (authToken) {
    headers.set("Authorization", `Bearer ${authToken}`);
  }

  const response = await fetch(`${baseURL}${path}`, {
    ...init,
    headers,
  });

  const data = (await response.json()) as T | ApiError;

  if (!response.ok) {
    const message =
      typeof data === "object" && data !== null && "error" in data
        ? translateApiError((data as ApiError).error)
        : `请求失败（${response.status}）`;
    throw new Error(message);
  }

  return data as T;
}

export async function checkHealth(): Promise<boolean> {
  try {
    const data = await request<{ status: string }>("/health");
    return data.status === "ok";
  } catch {
    return false;
  }
}

export async function getSourceLanguages(): Promise<Language[]> {
  const data = await request<{ languages: Record<string, string>[] }>(
    "/api/languages/source",
  );
  return (data.languages ?? []).map(normalizeLanguage);
}

export async function getTargetLanguages(): Promise<Language[]> {
  const data = await request<{ languages: Record<string, string>[] }>(
    "/api/languages/target",
  );
  return (data.languages ?? []).map(normalizeLanguage);
}

export async function translateText(
  body: TranslateTextRequest,
): Promise<TranslateTextResponse> {
  return request<TranslateTextResponse>("/api/translate", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function rephraseText(
  body: RephraseTextRequest,
): Promise<RephraseTextResponse> {
  return request<RephraseTextResponse>("/api/rephrase", {
    method: "POST",
    body: JSON.stringify(body),
  });
}
