'use client';

import { memo } from 'react';
import { Handle, Position, NodeProps } from 'reactflow';

export interface SharedMailboxNodeData {
  label: string;
  cohort?: string[];
  isRevoked?: boolean;
}

function SharedMailboxNode({ data, selected }: NodeProps<SharedMailboxNodeData>) {
  return (
    <div
      className={`
        px-4 py-3 rounded-lg border min-w-[160px]
        bg-gradient-to-br from-amber-600 to-amber-800 text-white border-amber-500/50
        shadow-lg shadow-amber-500/20
        ${selected ? 'ring-2 ring-amber-400 ring-offset-2 ring-offset-zinc-900' : ''}
        ${data.isRevoked ? 'opacity-40 grayscale border-red-500 border-dashed' : ''}
        transition-all duration-200
      `}
    >
      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 !bg-amber-400 border-2 !border-zinc-900"
      />
      <div className="flex items-center gap-2">
        <span className="text-lg">📮</span>
        <div>
          <div className="font-semibold text-sm">{data.label}</div>
          <div className="text-xs text-amber-200/70">Shared Mailbox</div>
        </div>
      </div>
      {data.cohort && data.cohort.length > 0 && (
        <div className="text-xs mt-2 bg-amber-500/30 border border-amber-400/30 rounded px-2 py-1">
          Cohort: {data.cohort.join(', ')}
        </div>
      )}
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 !bg-amber-400 border-2 !border-zinc-900"
      />
    </div>
  );
}

export default memo(SharedMailboxNode);
