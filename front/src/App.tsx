import { useCallback, useEffect, useState } from "react";
import { Layout } from "./components/Layout";
import { Store } from "./components/Store";
import { Forge } from "./components/Forge";
import { Exchange } from "./components/Exchange";
import { DailyLogin } from "./components/DailyLogin";
import { ProfilePage } from "./components/ProfilePage";
import { InventoryPage } from "./components/InventoryPage";
import { TransactionsPage } from "./components/TransactionsPage";
import { AnimatePresence, motion } from "framer-motion";
import { useAuthStore } from "./store/authStore";
import type { AppTab } from "./types/app";

export default function App() {
  const [activeTab, setActiveTab] = useState<AppTab>("store");
  const initAuth = useAuthStore((s) => s.initAuth);

  useEffect(() => {
    void initAuth();
  }, [initAuth]);

  const onAuth401 = useCallback(() => setActiveTab("store"), []);

  const variants = {
    initial: { opacity: 0, x: 20 },
    animate: { opacity: 1, x: 0 },
    exit: { opacity: 0, x: -20 },
  };

  const renderMain = () => {
    switch (activeTab) {
      case "store":
        return <Store />;
      case "forge":
        return <Forge />;
      case "exchange":
        return <Exchange />;
      case "daily":
        return <DailyLogin />;
      case "profile":
        return <ProfilePage onBack={() => setActiveTab("store")} />;
      case "inventory":
        return <InventoryPage onBack={() => setActiveTab("store")} />;
      case "transactions":
        return <TransactionsPage onBack={() => setActiveTab("store")} />;
      default:
        return <Store />;
    }
  };

  return (
    <Layout
      activeTab={activeTab}
      setActiveTab={setActiveTab}
      onAuth401={onAuth401}
    >
      <AnimatePresence mode="wait">
        <motion.div
          key={activeTab}
          initial="initial"
          animate="animate"
          exit="exit"
          variants={variants}
          transition={{ duration: 0.2 }}
        >
          {renderMain()}
        </motion.div>
      </AnimatePresence>
    </Layout>
  );
}
