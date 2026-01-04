'use client';

import Link from 'next/link';

export default function Header() {
  return (
    <header className="border-b border-zinc-800 bg-zinc-900/80 backdrop-blur-sm">
      <div className="max-w-7xl mx-auto px-4 py-4">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
          <div>
            <h1 className="text-xl font-bold text-zinc-100 flex items-center gap-2">
              <span className="text-indigo-400">{"{"}</span>
              E2E Encrypted File Sharing
              <span className="text-indigo-400">{"}"}</span>
            </h1>
            <p className="text-sm text-zinc-500">
              An interactive visualization of the architecture
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Link
              href="https://medium.com/@berniem4483/8599f15647c6"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 px-4 py-2 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-white font-medium transition-all shadow-lg shadow-indigo-600/25 hover:shadow-indigo-500/40"
            >
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                <path d="M13 4v2.67l-1 1-1-1V4H7v7l5 5 5-5V4h-4zm-1 10.5L7 9.5V6h2v3.17l3 3 3-3V6h2v3.5l-5 5z" />
              </svg>
              Read the Threat Model
            </Link>
            <Link
              href="https://github.com/educatedbernie/securefileshare"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 px-4 py-2 rounded-lg border border-zinc-700 hover:border-zinc-500 text-zinc-300 hover:text-white font-medium transition-all hover:bg-zinc-800"
            >
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
              </svg>
              View Source
            </Link>
          </div>
        </div>
      </div>
    </header>
  );
}
