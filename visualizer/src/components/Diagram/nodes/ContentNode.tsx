'use client';

import { memo } from 'react';
import { Handle, Position, NodeProps } from 'reactflow';

export interface ContentNodeData {
  label: string;
  chunkCount?: number;
}

function ContentNode({ data, selected }: NodeProps<ContentNodeData>) {
  return (
    <div
      className={`
        px-4 py-3 rounded-lg border min-w-[120px]
        bg-gradient-to-br from-zinc-700 to-zinc-900 text-white border-zinc-600/50
        shadow-lg shadow-zinc-500/10
        ${selected ? 'ring-2 ring-zinc-400 ring-offset-2 ring-offset-zinc-900' : ''}
        transition-all duration-200
      `}
    >
      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 !bg-zinc-400 border-2 !border-zinc-900"
      />
      <div className="flex items-center gap-2">
        <span className="text-lg">🔒</span>
        <div>
          <div className="font-semibold text-sm">{data.label}</div>
          <div className="text-xs text-zinc-400">Encrypted Content</div>
        </div>
      </div>
      {data.chunkCount && (
        <div className="text-xs mt-1 text-zinc-500">
          {data.chunkCount} chunk{data.chunkCount > 1 ? 's' : ''}
        </div>
      )}
    </div>
  );
}

export default memo(ContentNode);
