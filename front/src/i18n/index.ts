import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import en from "./locales/en.json";
import ru from "./locales/ru.json";

const STORAGE_KEY = "lootforge_lang";

const saved =
  typeof localStorage !== "undefined"
    ? localStorage.getItem(STORAGE_KEY)
    : null;

void i18n.use(initReactI18next).init({
  resources: {
    ru: { translation: ru },
    en: { translation: en },
  },
  lng: saved === "en" ? "en" : "ru",
  fallbackLng: "ru",
  interpolation: { escapeValue: false },
});

i18n.on("languageChanged", (lng) => {
  document.documentElement.lang = lng === "en" ? "en" : "ru";
});
document.documentElement.lang = i18n.language === "en" ? "en" : "ru";

export function setAppLanguage(lng: "ru" | "en") {
  void i18n.changeLanguage(lng);
  localStorage.setItem(STORAGE_KEY, lng);
}

export default i18n;
