'use client';

import { motion } from 'framer-motion';

// Reusable animated components
const FloatingIcon = ({ children, delay = 0, className = '' }: { children: React.ReactNode; delay?: number; className?: string }) => (
  <motion.div
    initial={{ opacity: 0, y: 20 }}
    animate={{ opacity: 1, y: 0 }}
    transition={{ duration: 0.5, delay }}
    className={className}
  >
    {children}
  </motion.div>
);

const PulsingRing = ({ color, size = 'w-32 h-32', delay = 0 }: { color: string; size?: string; delay?: number }) => (
  <motion.div
    className={`absolute ${size} rounded-full border-2 ${color}`}
    initial={{ scale: 0.8, opacity: 0 }}
    animate={{ scale: [1, 1.2, 1], opacity: [0.5, 0.2, 0.5] }}
    transition={{ duration: 3, repeat: Infinity, delay }}
  />
);

// Welcome slide - lock icon with encryption visualization
export function WelcomeVisual() {
  return (
    <div className="h-full flex items-center justify-center relative overflow-hidden">
      {/* Background grid */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute inset-0" style={{
          backgroundImage: 'linear-gradient(rgba(99, 102, 241, 0.3) 1px, transparent 1px), linear-gradient(90deg, rgba(99, 102, 241, 0.3) 1px, transparent 1px)',
          backgroundSize: '40px 40px'
        }} />
      </div>

      {/* Central lock icon */}
      <FloatingIcon className="relative">
        <div className="text-8xl mb-4">🔐</div>
        <motion.div
          className="absolute -inset-8 border-2 border-indigo-500/30 rounded-full"
          animate={{ rotate: 360 }}
          transition={{ duration: 20, repeat: Infinity, ease: "linear" }}
        />
        <motion.div
          className="absolute -inset-16 border border-indigo-500/20 rounded-full"
          animate={{ rotate: -360 }}
          transition={{ duration: 30, repeat: Infinity, ease: "linear" }}
        />
      </FloatingIcon>

      {/* Floating encryption symbols */}
      <motion.div
        className="absolute top-1/4 left-1/4 text-2xl text-indigo-400/40 font-mono"
        animate={{ opacity: [0.2, 0.5, 0.2], y: [-5, 5, -5] }}
        transition={{ duration: 4, repeat: Infinity }}
      >
        AES-256
      </motion.div>
      <motion.div
        className="absolute bottom-1/3 right-1/4 text-2xl text-violet-400/40 font-mono"
        animate={{ opacity: [0.2, 0.5, 0.2], y: [5, -5, 5] }}
        transition={{ duration: 4, repeat: Infinity, delay: 1 }}
      >
        RSA-2048
      </motion.div>
    </div>
  );
}

// Two Adversaries - split screen with two threat actors
export function TwoAdversariesVisual() {
  return (
    <div className="h-full flex items-center justify-center gap-16 relative">
      {/* Server Adversary */}
      <FloatingIcon delay={0.2} className="flex flex-col items-center">
        <div className="relative">
          <div className="w-24 h-32 bg-gradient-to-b from-red-900/50 to-red-950/50 border border-red-500/30 rounded-lg flex items-center justify-center">
            <div className="text-4xl">🖥️</div>
          </div>
          <motion.div
            className="absolute -top-2 -right-2 w-6 h-6 bg-red-500 rounded-full flex items-center justify-center text-xs"
            animate={{ scale: [1, 1.2, 1] }}
            transition={{ duration: 2, repeat: Infinity }}
          >
            👁️
          </motion.div>
        </div>
        <span className="mt-3 text-red-400 font-medium text-sm">Server</span>
        <span className="text-zinc-500 text-xs">Sees everything</span>
      </FloatingIcon>

      {/* VS divider */}
      <div className="text-zinc-600 text-2xl font-bold">VS</div>

      {/* Revoked User Adversary */}
      <FloatingIcon delay={0.4} className="flex flex-col items-center">
        <div className="relative">
          <div className="w-24 h-32 bg-gradient-to-b from-amber-900/50 to-amber-950/50 border border-amber-500/30 rounded-lg flex items-center justify-center">
            <div className="text-4xl">👤</div>
          </div>
          <motion.div
            className="absolute -top-2 -right-2 w-6 h-6 bg-amber-500 rounded-full flex items-center justify-center"
            animate={{ rotate: [0, 10, -10, 0] }}
            transition={{ duration: 2, repeat: Infinity }}
          >
            <span className="text-xs">🔑</span>
          </motion.div>
        </div>
        <span className="mt-3 text-amber-400 font-medium text-sm">Revoked User</span>
        <span className="text-zinc-500 text-xs">Remembers secrets</span>
      </FloatingIcon>
    </div>
  );
}

// Server Adversary - menacing server with capabilities
export function ServerAdversaryVisual() {
  return (
    <div className="h-full flex items-center justify-center relative">
      {/* Central server */}
      <div className="relative">
        <FloatingIcon>
          <div className="w-40 h-48 bg-gradient-to-b from-zinc-800 to-zinc-900 border border-red-500/40 rounded-xl flex flex-col items-center justify-center shadow-lg shadow-red-500/10">
            <div className="text-6xl mb-2">🖥️</div>
            <div className="text-red-400 font-mono text-xs">UNTRUSTED</div>
          </div>
        </FloatingIcon>

        {/* Capability rays */}
        {['Read', 'Modify', 'Delete', 'Observe'].map((cap, i) => (
          <motion.div
            key={cap}
            className="absolute text-xs font-mono text-red-400/70 whitespace-nowrap"
            style={{
              top: '50%',
              left: '50%',
              transform: `rotate(${i * 90 - 45}deg) translateX(120px)`,
            }}
            initial={{ opacity: 0 }}
            animate={{ opacity: [0.3, 0.7, 0.3] }}
            transition={{ duration: 2, repeat: Infinity, delay: i * 0.3 }}
          >
            {cap}
          </motion.div>
        ))}

        {/* Pulsing danger rings */}
        <PulsingRing color="border-red-500/30" size="w-48 h-56" />
        <PulsingRing color="border-red-500/20" size="w-64 h-72" delay={0.5} />
      </div>
    </div>
  );
}

// Revoked User - user with old keys trying to access
export function RevokedUserVisual() {
  return (
    <div className="h-full flex items-center justify-center gap-8 relative">
      {/* Revoked user with keys */}
      <FloatingIcon className="flex flex-col items-center">
        <div className="relative">
          <div className="w-28 h-36 bg-gradient-to-b from-amber-900/30 to-amber-950/30 border border-amber-500/40 rounded-xl flex flex-col items-center justify-center">
            <div className="text-5xl">👤</div>
            <div className="absolute -bottom-2 left-1/2 -translate-x-1/2 bg-red-500/80 text-[10px] px-2 py-0.5 rounded text-white font-medium">
              REVOKED
            </div>
          </div>
          {/* Old keys floating around */}
          <motion.div
            className="absolute -right-4 top-4 text-2xl"
            animate={{ rotate: [0, 15, 0], y: [-2, 2, -2] }}
            transition={{ duration: 3, repeat: Infinity }}
          >
            🔑
          </motion.div>
          <motion.div
            className="absolute -left-4 top-8 text-xl opacity-60"
            animate={{ rotate: [0, -10, 0], y: [2, -2, 2] }}
            transition={{ duration: 2.5, repeat: Infinity }}
          >
            🔑
          </motion.div>
        </div>
        <span className="mt-4 text-amber-400 font-medium text-sm">Has old keys</span>
      </FloatingIcon>

      {/* Arrow */}
      <motion.div
        className="text-3xl text-zinc-600"
        animate={{ x: [0, 10, 0] }}
        transition={{ duration: 1.5, repeat: Infinity }}
      >
        →
      </motion.div>

      {/* Database with question mark */}
      <FloatingIcon delay={0.3} className="flex flex-col items-center">
        <div className="relative">
          <div className="w-28 h-36 bg-gradient-to-b from-zinc-800 to-zinc-900 border border-zinc-600 rounded-xl flex flex-col items-center justify-center">
            <div className="text-4xl">📦</div>
            <div className="text-zinc-500 text-xs mt-2 font-mono">UUID: ???</div>
          </div>
          <motion.div
            className="absolute -top-3 -right-3 text-2xl"
            animate={{ scale: [1, 1.2, 1], rotate: [0, 10, -10, 0] }}
            transition={{ duration: 2, repeat: Infinity }}
          >
            ❓
          </motion.div>
        </div>
        <span className="mt-4 text-zinc-400 font-medium text-sm">Limited visibility</span>
      </FloatingIcon>
    </div>
  );
}

// Security Properties - shield with checkmarks
export function SecurityPropertiesVisual() {
  const properties = [
    { icon: '🔒', label: 'Confidentiality', color: 'text-blue-400' },
    { icon: '✓', label: 'Integrity', color: 'text-emerald-400' },
    { icon: '🛡️', label: 'Authorization', color: 'text-violet-400' },
  ];

  return (
    <div className="h-full flex items-center justify-center">
      <div className="relative">
        {/* Central shield */}
        <FloatingIcon className="relative z-10">
          <div className="w-32 h-40 flex items-center justify-center">
            <svg viewBox="0 0 100 120" className="w-full h-full">
              <defs>
                <linearGradient id="shieldGrad" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stopColor="#10B981" stopOpacity="0.3" />
                  <stop offset="100%" stopColor="#10B981" stopOpacity="0.1" />
                </linearGradient>
              </defs>
              <path
                d="M50 5 L95 25 L95 60 Q95 100 50 115 Q5 100 5 60 L5 25 Z"
                fill="url(#shieldGrad)"
                stroke="#10B981"
                strokeWidth="2"
              />
              <motion.path
                d="M30 60 L45 75 L70 45"
                fill="none"
                stroke="#10B981"
                strokeWidth="4"
                strokeLinecap="round"
                initial={{ pathLength: 0 }}
                animate={{ pathLength: 1 }}
                transition={{ duration: 1, delay: 0.5 }}
              />
            </svg>
          </div>
        </FloatingIcon>

        {/* Orbiting properties */}
        {properties.map((prop, i) => (
          <motion.div
            key={prop.label}
            className="absolute flex items-center gap-2 bg-zinc-900/80 px-3 py-1.5 rounded-full border border-zinc-700"
            style={{
              top: `${30 + i * 30}%`,
              left: i % 2 === 0 ? '-80px' : 'auto',
              right: i % 2 === 1 ? '-80px' : 'auto',
            }}
            initial={{ opacity: 0, x: i % 2 === 0 ? -20 : 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.3 + i * 0.2 }}
          >
            <span>{prop.icon}</span>
            <span className={`text-xs font-medium ${prop.color}`}>{prop.label}</span>
          </motion.div>
        ))}
      </div>
    </div>
  );
}

// Constraints - locked down environment
export function ConstraintsVisual() {
  return (
    <div className="h-full flex items-center justify-center gap-12">
      {/* No local storage */}
      <FloatingIcon delay={0} className="flex flex-col items-center">
        <div className="relative w-20 h-24 bg-zinc-800/50 border border-zinc-600 rounded-lg flex items-center justify-center">
          <div className="text-3xl">💾</div>
          <motion.div
            className="absolute inset-0 flex items-center justify-center"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.5 }}
          >
            <div className="w-16 h-0.5 bg-red-500 rotate-45" />
          </motion.div>
        </div>
        <span className="mt-2 text-zinc-400 text-xs text-center">No local<br/>storage</span>
      </FloatingIcon>

      {/* Password only */}
      <FloatingIcon delay={0.2} className="flex flex-col items-center">
        <div className="w-20 h-24 bg-zinc-800/50 border border-zinc-600 rounded-lg flex flex-col items-center justify-center gap-1">
          <div className="text-2xl">👤</div>
          <div className="flex gap-0.5">
            {[...Array(6)].map((_, i) => (
              <motion.div
                key={i}
                className="w-2 h-2 bg-indigo-400 rounded-full"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.5 + i * 0.1 }}
              />
            ))}
          </div>
        </div>
        <span className="mt-2 text-zinc-400 text-xs text-center">Password<br/>only</span>
      </FloatingIcon>

      {/* Zero trust */}
      <FloatingIcon delay={0.4} className="flex flex-col items-center">
        <div className="relative w-20 h-24 bg-zinc-800/50 border border-red-500/30 rounded-lg flex items-center justify-center">
          <div className="text-3xl">🖥️</div>
          <motion.div
            className="absolute -bottom-1 left-1/2 -translate-x-1/2 bg-red-500/80 text-[8px] px-1.5 py-0.5 rounded text-white font-bold"
            animate={{ opacity: [0.7, 1, 0.7] }}
            transition={{ duration: 2, repeat: Infinity }}
          >
            ZERO TRUST
          </motion.div>
        </div>
        <span className="mt-2 text-zinc-400 text-xs text-center">Untrusted<br/>server</span>
      </FloatingIcon>
    </div>
  );
}

// The Challenge - Dropbox-like with encryption overlay
export function ChallengeVisual() {
  return (
    <div className="h-full flex items-center justify-center relative">
      <div className="relative">
        {/* Cloud/Dropbox representation */}
        <FloatingIcon>
          <div className="relative">
            <div className="w-48 h-32 bg-gradient-to-b from-blue-900/30 to-blue-950/30 border border-blue-500/30 rounded-2xl flex items-center justify-center">
              <div className="text-5xl">☁️</div>
            </div>

            {/* Files inside */}
            <motion.div
              className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 flex gap-2"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.5 }}
            >
              {['📄', '📁', '🖼️'].map((icon, i) => (
                <motion.div
                  key={i}
                  className="text-2xl"
                  animate={{ y: [-2, 2, -2] }}
                  transition={{ duration: 2, repeat: Infinity, delay: i * 0.3 }}
                >
                  {icon}
                </motion.div>
              ))}
            </motion.div>

            {/* Encryption overlay */}
            <motion.div
              className="absolute inset-0 bg-gradient-to-r from-emerald-500/10 via-transparent to-emerald-500/10 rounded-2xl"
              animate={{ opacity: [0.3, 0.6, 0.3] }}
              transition={{ duration: 3, repeat: Infinity }}
            />
          </div>
        </FloatingIcon>

        {/* Lock icon */}
        <motion.div
          className="absolute -bottom-4 left-1/2 -translate-x-1/2 text-4xl"
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.8 }}
        >
          🔐
        </motion.div>

        {/* Efficiency indicator */}
        <motion.div
          className="absolute -right-24 top-1/2 -translate-y-1/2 text-xs text-zinc-500 flex items-center gap-2"
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 1 }}
        >
          <span className="text-emerald-400">⚡</span>
          <span>Efficient<br/>Appends</span>
        </motion.div>
      </div>
    </div>
  );
}

// Key Insight - key management visualization
export function KeyInsightVisual() {
  return (
    <div className="h-full flex items-center justify-center relative">
      {/* Simple encryption - easy */}
      <FloatingIcon className="absolute left-1/4 top-1/3 flex flex-col items-center">
        <div className="flex items-center gap-3">
          <div className="text-3xl">📄</div>
          <motion.div
            animate={{ x: [0, 5, 0] }}
            transition={{ duration: 1, repeat: Infinity }}
          >
            →
          </motion.div>
          <div className="text-3xl">🔒</div>
          <motion.div
            animate={{ x: [0, 5, 0] }}
            transition={{ duration: 1, repeat: Infinity, delay: 0.2 }}
          >
            →
          </motion.div>
          <div className="text-3xl">📦</div>
        </div>
        <span className="mt-2 text-emerald-400 text-sm flex items-center gap-1">
          <span>✓</span> Easy
        </span>
      </FloatingIcon>

      {/* Sharing + Revocation - hard */}
      <FloatingIcon delay={0.3} className="absolute right-1/4 bottom-1/3 flex flex-col items-center">
        <div className="relative">
          <div className="flex items-center gap-2">
            <div className="text-2xl">👤</div>
            <div className="text-2xl">👤</div>
            <div className="text-2xl">👤</div>
          </div>
          <motion.div
            className="absolute -top-2 left-1/2 -translate-x-1/2 flex gap-1"
            animate={{ rotate: [0, 10, -10, 0] }}
            transition={{ duration: 3, repeat: Infinity }}
          >
            <span className="text-xl">🔑</span>
            <span className="text-xl">🔑</span>
            <span className="text-xl">🔑</span>
          </motion.div>
        </div>
        <span className="mt-4 text-red-400 text-sm flex items-center gap-1">
          <span>✗</span> Key Management
        </span>
      </FloatingIcon>

      {/* Central question */}
      <motion.div
        className="text-6xl"
        animate={{ scale: [1, 1.1, 1], rotate: [0, 5, -5, 0] }}
        transition={{ duration: 4, repeat: Infinity }}
      >
        🤔
      </motion.div>
    </div>
  );
}

// Map step IDs to visual components
export const introVisuals: Record<number, React.ComponentType> = {
  1: WelcomeVisual,
  2: TwoAdversariesVisual,
  3: ServerAdversaryVisual,
  4: RevokedUserVisual,
  5: SecurityPropertiesVisual,
  6: ConstraintsVisual,
  7: ChallengeVisual,
  8: KeyInsightVisual,
};
