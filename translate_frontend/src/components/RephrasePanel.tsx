import { useState } from "react";
import { rephraseText } from "../api/client";
import {
  translateApiError,
  WRITING_STYLE_OPTIONS,
  WRITING_TONE_OPTIONS,
} from "../i18n/zh";

export function RephrasePanel() {
  const [inputText, setInputText] = useState("");
  const [outputText, setOutputText] = useState("");
  const [style, setStyle] = useState("");
  const [tone, setTone] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleRephrase = async () => {
    if (!inputText.trim()) {
      setError("请输入要改写的文本");
      return;
    }

    setSubmitting(true);
    setError(null);
    try {
      const result = await rephraseText({
        text: inputText,
        ...(style ? { style } : {}),
        ...(tone ? { tone } : {}),
      });
      setOutputText(result.text);
    } catch (err) {
      setError(
        err instanceof Error ? translateApiError(err.message) : "改写失败",
      );
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="rephrase-panel">
      <div className="controls-row controls-row--secondary">
        <label className="field">
          <span className="field-label">写作风格</span>
          <select
            className="select"
            value={style}
            onChange={(e) => setStyle(e.target.value)}
          >
            {WRITING_STYLE_OPTIONS.map((opt) => (
              <option key={opt.value || "style-empty"} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </label>
        <label className="field">
          <span className="field-label">写作语气</span>
          <select
            className="select"
            value={tone}
            onChange={(e) => setTone(e.target.value)}
          >
            {WRITING_TONE_OPTIONS.map((opt) => (
              <option key={opt.value || "tone-empty"} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
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
            placeholder="输入要改写的文本…"
            rows={10}
          />
        </div>
        <div className="editor-box">
          <div className="editor-header">
            <span>改写结果</span>
          </div>
          <textarea
            className="textarea textarea--output"
            value={outputText}
            readOnly
            placeholder="改写结果将显示在这里…"
            rows={10}
          />
        </div>
      </div>

      {error && <p className="form-error">{error}</p>}

      <div className="actions">
        <button
          type="button"
          className="btn btn-primary"
          onClick={handleRephrase}
          disabled={submitting || !inputText.trim()}
        >
          {submitting ? "改写中…" : "改写"}
        </button>
      </div>
    </div>
  );
}
