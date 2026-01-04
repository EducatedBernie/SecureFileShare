'use client';

export type TabId = 'walkthrough' | 'sandbox' | 'deepdive';

interface Tab {
  id: TabId;
  label: string;
}

const tabs: Tab[] = [
  { id: 'walkthrough', label: 'Guided Walkthrough' },
  { id: 'sandbox', label: 'Sandbox' },
  { id: 'deepdive', label: 'Architecture Deep Dive' },
];

interface TabBarProps {
  activeTab: TabId;
  onTabChange: (tab: TabId) => void;
}

export default function TabBar({ activeTab, onTabChange }: TabBarProps) {
  return (
    <div className="border-b border-zinc-800 bg-zinc-900">
      <div className="max-w-7xl mx-auto px-4">
        <nav className="flex gap-1" aria-label="Tabs">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => onTabChange(tab.id)}
              className={`
                px-4 py-3 text-sm font-medium border-b-2 transition-all
                ${
                  activeTab === tab.id
                    ? 'border-indigo-500 text-indigo-400 bg-zinc-800/50'
                    : 'border-transparent text-zinc-500 hover:text-zinc-300 hover:border-zinc-600'
                }
              `}
            >
              {tab.label}
            </button>
          ))}
        </nav>
      </div>
    </div>
  );
}
