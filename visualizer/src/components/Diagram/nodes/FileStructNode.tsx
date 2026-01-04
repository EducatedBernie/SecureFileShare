'use client';

import { memo, useState } from 'react';
import { Handle, Position, NodeProps } from 'reactflow';

export interface FileStructNodeData {
  label: string;
  filename: string;
  owner: string;
  onRevoke?: () => void;
}

function FileStructNode({ data, selected }: NodeProps<FileStructNodeData>) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <div
      className={`
        relative px-4 py-3 rounded-lg border min-w-[140px]
        bg-gradient-to-br from-emerald-600 to-emerald-800 text-white border-emerald-500/50
        shadow-lg shadow-emerald-500/20
        ${selected ? 'ring-2 ring-emerald-400 ring-offset-2 ring-offset-zinc-900' : ''}
        transition-all duration-200 cursor-pointer
      `}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Hover action button - Revoke */}
      {data.onRevoke && (
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
              data.onRevoke?.();
            }}
            className="
              group relative w-8 h-8 rounded-full
              bg-gradient-to-br from-red-500 to-rose-600
              border-2 border-red-400/50
              shadow-lg shadow-red-500/40
              hover:shadow-red-500/60 hover:scale-110
              transition-all duration-150
              flex items-center justify-center
            "
          >
            <span className="text-white text-lg leading-none">×</span>
            {/* Tooltip */}
            <div className="
              absolute bottom-full left-1/2 -translate-x-1/2 mb-2
              px-2 py-1 rounded bg-zinc-900 border border-zinc-700
              text-xs text-zinc-200 whitespace-nowrap
              opacity-0 group-hover:opacity-100 transition-opacity
              pointer-events-none
            ">
              Revoke someone&apos;s access
            </div>
          </button>
        </div>
      )}

      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 !bg-emerald-400 border-2 !border-zinc-900"
      />
      <div className="flex items-center gap-2">
        <span className="text-lg">📄</span>
        <div>
          <div className="font-semibold text-sm">{data.filename}</div>
          <div className="text-xs text-emerald-200/70">FileStruct</div>
        </div>
      </div>
      <div className="text-xs mt-1 text-emerald-200/50">
        Owner: {data.owner}
      </div>
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 !bg-emerald-400 border-2 !border-zinc-900"
      />
    </div>
  );
}

export default memo(FileStructNode);
