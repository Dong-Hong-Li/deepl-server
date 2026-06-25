/** 界面中文文案与选项标签（API 仍传英文枚举值） */

/** DeepL 专用语言码，优先于 Intl 解析结果 */
const DEEPL_LANGUAGE_OVERRIDES: Record<string, string> = {
  en: "英语",
  "en-gb": "英语（英国）",
  "en-us": "英语（美国）",
  es: "西班牙语",
  "es-419": "西班牙语（拉丁美洲）",
  nb: "挪威语（书面）",
  pt: "葡萄牙语",
  "pt-br": "葡萄牙语（巴西）",
  "pt-pt": "葡萄牙语（葡萄牙）",
  zh: "中文",
  "zh-hans": "中文（简体）",
  "zh-hant": "中文（繁体）",
};

const languageDisplayNames = new Intl.DisplayNames(["zh-CN"], {
  type: "language",
});

function normalizeLanguageCode(code: string): string[] {
  const lower = code.toLowerCase();
  const [base, region] = lower.split("-");
  const variants = new Set<string>([lower, base]);

  if (region) {
    variants.add(`${base}-${region.toUpperCase()}`);
    if (region.length === 4) {
      variants.add(`${base}-${region[0].toUpperCase()}${region.slice(1)}`);
    }
  }

  return [...variants];
}

function resolveIntlLanguageName(code: string): string | undefined {
  for (const variant of normalizeLanguageCode(code)) {
    try {
      const name = languageDisplayNames.of(variant);
      if (name && /[\u4e00-\u9fff]/.test(name)) {
        return name;
      }
    } catch {
      // ignore unsupported codes
    }
  }
  return undefined;
}

export function getLanguageLabel(code: string): string {
  const key = code.toLowerCase();
  if (DEEPL_LANGUAGE_OVERRIDES[key]) {
    return DEEPL_LANGUAGE_OVERRIDES[key];
  }
  return resolveIntlLanguageName(code) ?? key;
}

export function sortLanguagesByChineseLabel<T extends { code: string }>(
  languages: T[],
): T[] {
  return [...languages].sort((a, b) =>
    getLanguageLabel(a.code).localeCompare(getLanguageLabel(b.code), "zh-CN"),
  );
}

/** 搜索匹配：中文名、语言码、API 英文名 */
export function matchesLanguageSearch(
  lang: { code: string; name: string },
  query: string,
): boolean {
  const q = query.trim().toLowerCase();
  if (!q) return true;

  const zh = getLanguageLabel(lang.code);
  if (zh.includes(query.trim())) return true;
  if (lang.code.toLowerCase().includes(q)) return true;
  if (lang.name.toLowerCase().includes(q)) return true;
  return false;
}

export interface SelectOption {
  value: string;
  label: string;
}

export const FORMALITY_OPTIONS: SelectOption[] = [
  { value: "", label: "不指定" },
  { value: "default", label: "默认" },
  { value: "less", label: "较非正式" },
  { value: "more", label: "较正式" },
  { value: "prefer_less", label: "倾向非正式" },
  { value: "prefer_more", label: "倾向正式" },
];

export const WRITING_STYLE_OPTIONS: SelectOption[] = [
  { value: "", label: "不指定" },
  { value: "academic", label: "学术" },
  { value: "business", label: "商务" },
  { value: "casual", label: "休闲" },
  { value: "simple", label: "简洁" },
  { value: "default", label: "默认" },
  { value: "prefer_academic", label: "倾向学术" },
  { value: "prefer_business", label: "倾向商务" },
  { value: "prefer_casual", label: "倾向休闲" },
  { value: "prefer_simple", label: "倾向简洁" },
];

export const WRITING_TONE_OPTIONS: SelectOption[] = [
  { value: "", label: "不指定" },
  { value: "confident", label: "自信" },
  { value: "diplomatic", label: "得体" },
  { value: "enthusiastic", label: "热情" },
  { value: "friendly", label: "友好" },
  { value: "default", label: "默认" },
  { value: "prefer_confident", label: "倾向自信" },
  { value: "prefer_diplomatic", label: "倾向得体" },
  { value: "prefer_enthusiastic", label: "倾向热情" },
  { value: "prefer_friendly", label: "倾向友好" },
];

/** 将常见英文 API 错误信息转为中文提示 */
export function translateApiError(message: string): string {
  const map: Record<string, string> = {
    "text is required": "请输入文本",
    "targetLangCode is required": "请选择目标语言",
    "sourceLangCode is required when glossaryId is set":
      "使用术语表时必须指定源语言",
    "invalid request body": "请求格式无效",
    unauthorized: "未授权，请检查访问令牌",
  };

  const lower = message.toLowerCase();
  for (const [key, zh] of Object.entries(map)) {
    if (lower.includes(key.toLowerCase())) {
      return zh;
    }
  }
  if (lower.includes("403") || lower.includes("forbidden")) {
    if (lower.includes("write") || lower.includes("rephrase") || lower.includes("pro")) {
      return "文本改写需要 DeepL API Pro 订阅，Free 版 Key 不支持此功能";
    }
    return "无权访问该接口，请检查 API Key 权限";
  }
  if (lower.includes("rephrasing failed")) {
    if (lower.includes("403") || lower.includes("pro")) {
      return "文本改写需要 DeepL API Pro 订阅，Free 版 Key 不支持此功能";
    }
    return message;
  }
  if (lower.includes("translation failed")) return message;
  if (lower.includes("failed to get")) return "获取数据失败，请检查服务是否启动";
  return message;
}
