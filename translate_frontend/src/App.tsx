import { useEffect, useState } from "react";
import { checkHealth } from "./api/client";
import { RephrasePanel } from "./components/RephrasePanel";
import { TranslatePanel } from "./components/TranslatePanel";

type Tab = "translate" | "rephrase";

export default function App() {
  const [tab, setTab] = useState<Tab>("translate");
  const [online, setOnline] = useState<boolean | null>(null);

  useEffect(() => {
    void checkHealth().then(setOnline);
  }, []);

  return (
    <div className="app">
      <header className="header">
        <div className="header-inner">
          <div className="brand">
            <span className="brand-icon">译</span>
            <div>
              <h1 className="brand-title">DeepL 翻译工作台</h1>
              <p className="brand-subtitle">文本翻译 · 智能改写</p>
            </div>
          </div>
          <div className="status-badge" data-online={online}>
            <span className="status-dot" />
            {online === null
              ? "检测连接…"
              : online
                ? "服务已连接"
                : "服务未连接"}
          </div>
        </div>
      </header>

      <main className="main">
        <nav className="tabs" role="tablist">
          <button
            type="button"
            role="tab"
            className="tab"
            data-active={tab === "translate"}
            onClick={() => setTab("translate")}
            aria-selected={tab === "translate"}
          >
            文本翻译
          </button>
          <button
            type="button"
            role="tab"
            className="tab"
            data-active={tab === "rephrase"}
            onClick={() => setTab("rephrase")}
            aria-selected={tab === "rephrase"}
          >
            文本改写
          </button>
        </nav>

        <section className="card">
          {tab === "translate" ? <TranslatePanel /> : <RephrasePanel />}
        </section>
      </main>

      <footer className="footer">
        已连接接口：语言列表 · 文本翻译 · 文本改写
      </footer>
    </div>
  );
}
