import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { I18nextProvider } from "react-i18next";
import "./index.css";
import i18n from "./i18n";
import App from "./App";
import { ToastProvider } from "./context/ToastContext";
import { ThemeProvider } from "./theme/ThemeContext";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <I18nextProvider i18n={i18n}>
      <ThemeProvider>
        <ToastProvider>
          <App />
        </ToastProvider>
      </ThemeProvider>
    </I18nextProvider>
  </StrictMode>
);
