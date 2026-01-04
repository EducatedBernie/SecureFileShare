'use client';

import { memo, useState } from 'react';
import { Handle, Position, NodeProps } from 'reactflow';

export interface UserNodeData {
  label: string;
  isOwner?: boolean;
  isRevoked?: boolean;
  hasFile?: boolean;
  onShare?: () => void;
}

function UserNode({ data, selected }: NodeProps<UserNodeData>) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <div
      className={`
        relative px-4 py-3 rounded-lg border min-w-[120px]
        bg-gradient-to-br from-blue-600 to-blue-800 text-white border-blue-500/50
        shadow-lg shadow-blue-500/20
        ${selected ? 'ring-2 ring-blue-400 ring-offset-2 ring-offset-zinc-900' : ''}
        ${data.isRevoked ? 'opacity-40 grayscale' : ''}
        transition-all duration-200 cursor-pointer
      `}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Hover action button */}
      {data.hasFile && data.onShare && (
        <div
          className={`
            absolute -top-3 -right-3 z-10
            transition-all duration-200 ease-out
            ${isHovered ? 'opacity-100 scale-100' : 'opacity-0 scale-75 pointer-events-none'}
          `}
        >
          <button
            onClick={(e) => {
              e.stopPropagation();
              data.onShare?.();
            }}
            className="
              group relative w-8 h-8 rounded-full
              bg-gradient-to-br from-indigo-500 to-purple-600
              border-2 border-indigo-400/50
              shadow-lg shadow-indigo-500/40
              hover:shadow-indigo-500/60 hover:scale-110
              transition-all duration-150
              flex items-center justify-center
            "
          >
            <span className="text-white text-lg leading-none">+</span>
            {/* Tooltip */}
            <div className="
              absolute bottom-full left-1/2 -translate-x-1/2 mb-2
              px-2 py-1 rounded bg-zinc-900 border border-zinc-700
              text-xs text-zinc-200 whitespace-nowrap
              opacity-0 group-hover:opacity-100 transition-opacity
              pointer-events-none
            ">
              Share file from {data.label}
            </div>
          </button>
        </div>
      )}

      <div className="flex items-center gap-2">
        <span className="text-lg">👤</span>
        <div>
          <div className="font-semibold">{data.label}</div>
          <div className="text-xs text-blue-200/70">User</div>
        </div>
      </div>
      {data.isOwner && (
        <div className="text-xs mt-1 bg-blue-500/50 border border-blue-400/30 rounded px-1.5 py-0.5 inline-block">
          Owner
        </div>
      )}
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 !bg-blue-400 border-2 !border-zinc-900"
      />
    </div>
  );
}

export default memo(UserNode);
