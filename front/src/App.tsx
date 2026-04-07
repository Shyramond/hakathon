import { useState } from 'react';
import { Layout } from './components/Layout';
import { Store } from './components/Store';
import { Forge } from './components/Forge';
import { Exchange } from './components/Exchange';
import { DailyLogin } from './components/DailyLogin';
import { AnimatePresence, motion } from 'framer-motion';

export default function App() {
  const [activeTab, setActiveTab] = useState('store');

  const variants = {
    initial: { opacity: 0, x: 20 },
    animate: { opacity: 1, x: 0 },
    exit: { opacity: 0, x: -20 },
  };

  return (
    <Layout activeTab={activeTab} setActiveTab={setActiveTab}>
      <AnimatePresence mode="wait">
        <motion.div
          key={activeTab}
          initial="initial"
          animate="animate"
          exit="exit"
          variants={variants}
          transition={{ duration: 0.2 }}
        >
          {activeTab === 'store' && <Store />}
          {activeTab === 'forge' && <Forge />}
          {activeTab === 'exchange' && <Exchange />}
          {activeTab === 'daily' && <DailyLogin />}
        </motion.div>
      </AnimatePresence>
    </Layout>
  );
}
