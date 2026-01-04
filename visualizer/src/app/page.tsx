'use client';

import { useState } from 'react';
import Header from '@/components/Layout/Header';
import TabBar, { TabId } from '@/components/Layout/TabBar';
import WalkthroughContainer from '@/components/Walkthrough/WalkthroughContainer';
import SandboxContainer from '@/components/Sandbox/SandboxContainer';
import DeepDiveContainer from '@/components/DeepDive/DeepDiveContainer';

export default function Home() {
  const [activeTab, setActiveTab] = useState<TabId>('walkthrough');

  return (
    <div className="h-screen flex flex-col bg-[#0a0a0f] overflow-hidden">
      <Header />
      <TabBar activeTab={activeTab} onTabChange={setActiveTab} />

      <main className="flex-1 flex flex-col overflow-hidden">
        {activeTab === 'walkthrough' && <WalkthroughContainer />}
        {activeTab === 'sandbox' && <SandboxContainer />}
        {activeTab === 'deepdive' && <DeepDiveContainer />}
      </main>
    </div>
  );
}
