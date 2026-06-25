import { useCallback, useEffect, useId, useRef, useState } from "react";
import { getSourceLanguages, getTargetLanguages } from "../api/client";
import type { Language } from "../api/types";
import {
  getLanguageLabel,
  matchesLanguageSearch,
  sortLanguagesByChineseLabel,
  translateApiError,
} from "../i18n/zh";

interface LanguageSelectProps {
  id: string;
  label: string;
  languages: Language[];
  value: string;
  onChange: (code: string) => void;
  allowEmpty?: boolean;
  emptyLabel?: string;
  disabled?: boolean;
}

export function LanguageSelect({
  id,
  label,
  languages,
  value,
  onChange,
  allowEmpty = false,
  emptyLabel = "自动检测",
  disabled = false,
}: LanguageSelectProps) {
  const listId = useId();
  const rootRef = useRef<HTMLDivElement>(null);
  const searchRef = useRef<HTMLInputElement>(null);

  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");

  const selectedLabel = value
    ? getLanguageLabel(value)
    : allowEmpty
      ? emptyLabel
      : "请选择语言";

  const filtered = languages.filter((lang) => matchesLanguageSearch(lang, query));

  const close = useCallback(() => {
    setOpen(false);
    setQuery("");
  }, []);

  const selectCode = useCallback(
    (code: string) => {
      onChange(code);
      close();
    },
    [close, onChange],
  );

  useEffect(() => {
    if (!open) return;

    const onPointerDown = (event: MouseEvent) => {
      if (!rootRef.current?.contains(event.target as Node)) {
        close();
      }
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") close();
    };

    document.addEventListener("mousedown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    searchRef.current?.focus();

    return () => {
      document.removeEventListener("mousedown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [close, open]);

  return (
    <div className="field lang-combobox" ref={rootRef}>
      <span className="field-label" id={`${id}-label`}>
        {label}
      </span>
      <button
        type="button"
        id={id}
        className="combobox-trigger"
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-labelledby={`${id}-label`}
        disabled={disabled}
        onClick={() => setOpen((prev) => !prev)}
      >
        <span className="combobox-value">{selectedLabel}</span>
        <span className="combobox-chevron" aria-hidden>
          ▾
        </span>
      </button>

      {open && (
        <div className="combobox-panel">
          <input
            ref={searchRef}
            className="combobox-search input"
            type="search"
            value={query}
            placeholder="搜索语言（中文 / 英文 / 代码）"
            aria-controls={listId}
            onChange={(e) => setQuery(e.target.value)}
          />
          <ul className="combobox-list" id={listId} role="listbox">
            {allowEmpty && (
              <li>
                <button
                  type="button"
                  role="option"
                  aria-selected={value === ""}
                  className="combobox-option"
                  data-selected={value === ""}
                  onClick={() => selectCode("")}
                >
                  {emptyLabel}
                </button>
              </li>
            )}
            {filtered.length === 0 ? (
              <li className="combobox-empty">无匹配语言</li>
            ) : (
              filtered.map((lang) => (
                <li key={lang.code}>
                  <button
                    type="button"
                    role="option"
                    aria-selected={value === lang.code}
                    className="combobox-option"
                    data-selected={value === lang.code}
                    onClick={() => selectCode(lang.code)}
                  >
                    {getLanguageLabel(lang.code)}
                  </button>
                </li>
              ))
            )}
          </ul>
        </div>
      )}
    </div>
  );
}

interface UseLanguagesResult {
  sourceLanguages: Language[];
  targetLanguages: Language[];
  loading: boolean;
  error: string | null;
  reload: () => void;
}

export function useLanguages(): UseLanguagesResult {
  const [sourceLanguages, setSourceLanguages] = useState<Language[]>([]);
  const [targetLanguages, setTargetLanguages] = useState<Language[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [source, target] = await Promise.all([
        getSourceLanguages(),
        getTargetLanguages(),
      ]);
      setSourceLanguages(sortLanguagesByChineseLabel(source));
      setTargetLanguages(sortLanguagesByChineseLabel(target));
    } catch (err) {
      setError(
        err instanceof Error ? translateApiError(err.message) : "加载语言列表失败",
      );
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void reload();
  }, [reload]);

  return { sourceLanguages, targetLanguages, loading, error, reload };
}
