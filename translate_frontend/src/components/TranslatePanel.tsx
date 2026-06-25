import { useState } from "react";
import { translateText } from "../api/client";
import {
  FORMALITY_OPTIONS,
  getLanguageLabel,
  translateApiError,
} from "../i18n/zh";
import { LanguageSelect, useLanguages } from "./LanguageSelect";

export function TranslatePanel() {
  const { sourceLanguages, targetLanguages, loading, error, reload } =
    useLanguages();

  const [sourceLang, setSourceLang] = useState("");
  const [targetLang, setTargetLang] = useState("zh");
  const [formality, setFormality] = useState("");
  const [glossaryId, setGlossaryId] = useState("");
  const [inputText, setInputText] = useState("");
  const [outputText, setOutputText] = useState("");
  const [detectedLang, setDetectedLang] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const handleSwap = () => {
    if (!sourceLang || !targetLang) return;
    setSourceLang(targetLang);
    setTargetLang(sourceLang);
    setInputText(outputText);
    setOutputText("");
    setDetectedLang("");
  };

  const handleTranslate = async () => {
    if (!inputText.trim() || !targetLang) {
      setSubmitError("请输入文本并选择目标语言");
      return;
    }
    if (glossaryId && !sourceLang) {
      setSubmitError("使用术语表时必须指定源语言");
      return;
    }

    setSubmitting(true);
    setSubmitError(null);
    try {
      const result = await translateText({
        text: inputText,
        targetLangCode: targetLang,
        ...(sourceLang ? { sourceLangCode: sourceLang } : {}),
        ...(formality ? { formality } : {}),
        ...(glossaryId ? { glossaryId } : {}),
      });
      setOutputText(result.text);
      setDetectedLang(result.detectedSourceLang);
    } catch (err) {
      setSubmitError(
        err instanceof Error ? translateApiError(err.message) : "翻译失败",
      );
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <div className="panel-state">正在加载语言列表…</div>;
  }

  if (error) {
    return (
      <div className="panel-state panel-state--error">
        <p>{error}</p>
        <button type="button" className="btn btn-secondary" onClick={reload}>
          重试
        </button>
      </div>
    );
  }

  return (
    <div className="translate-panel">
      <div className="controls-row">
        <LanguageSelect
          id="source-lang"
          label="源语言"
          languages={sourceLanguages}
          value={sourceLang}
          onChange={setSourceLang}
          allowEmpty
        />
        <button
          type="button"
          className="swap-btn"
          onClick={handleSwap}
          title="交换语言"
          aria-label="交换语言"
        >
          ⇄
        </button>
        <LanguageSelect
          id="target-lang"
          label="目标语言"
          languages={targetLanguages}
          value={targetLang}
          onChange={setTargetLang}
        />
      </div>

      <div className="controls-row controls-row--secondary">
        <label className="field">
          <span className="field-label">正式程度</span>
          <select
            className="select"
            value={formality}
            onChange={(e) => setFormality(e.target.value)}
          >
            {FORMALITY_OPTIONS.map((opt) => (
              <option key={opt.value || "default-empty"} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </label>
        <label className="field field--grow">
          <span className="field-label">术语表 ID（可选）</span>
          <input
            className="input"
            type="text"
            value={glossaryId}
            onChange={(e) => setGlossaryId(e.target.value)}
            placeholder="请输入术语表 ID"
          />
        </label>
      </div>

      <div className="editor-grid">
        <div className="editor-box">
          <div className="editor-header">
            <span>原文</span>
            <span className="char-count">{inputText.length} 字符</span>
          </div>
          <textarea
            className="textarea"
            value={inputText}
            onChange={(e) => setInputText(e.target.value)}
            placeholder="输入要翻译的文本…"
            rows={10}
          />
        </div>
        <div className="editor-box">
          <div className="editor-header">
            <span>译文</span>
            {detectedLang && (
              <span className="detected-lang">
                检测到：{getLanguageLabel(detectedLang)}
              </span>
            )}
          </div>
          <textarea
            className="textarea textarea--output"
            value={outputText}
            readOnly
            placeholder="翻译结果将显示在这里…"
            rows={10}
          />
        </div>
      </div>

      {submitError && <p className="form-error">{submitError}</p>}

      <div className="actions">
        <button
          type="button"
          className="btn btn-primary"
          onClick={handleTranslate}
          disabled={submitting || !inputText.trim()}
        >
          {submitting ? "翻译中…" : "翻译"}
        </button>
      </div>
    </div>
  );
}
