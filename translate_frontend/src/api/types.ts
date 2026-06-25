export interface Language {
  code: string;
  name: string;
}

export interface TranslateTextRequest {
  text: string;
  targetLangCode: string;
  sourceLangCode?: string;
  formality?: string;
  glossaryId?: string;
}

export interface TranslateTextResponse {
  text: string;
  detectedSourceLang: string;
  targetLangCode: string;
}

export interface RephraseTextRequest {
  text: string;
  style?: string;
  tone?: string;
}

export interface RephraseTextResponse {
  text: string;
}

export interface ApiError {
  error: string;
}
