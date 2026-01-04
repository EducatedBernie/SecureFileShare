'use client';

import { memo } from 'react';
import { Handle, Position, NodeProps } from 'reactflow';

export interface PersonalMailboxNodeData {
  label: string;
  owner: string;
  isRevoked?: boolean;
}

function PersonalMailboxNode({ data, selected }: NodeProps<PersonalMailboxNodeData>) {
  return (
    <div
      className={`
        px-4 py-3 rounded-lg border min-w-[140px]
        bg-gradient-to-br from-violet-600 to-violet-800 text-white border-violet-500/50
        shadow-lg shadow-violet-500/20
        ${selected ? 'ring-2 ring-violet-400 ring-offset-2 ring-offset-zinc-900' : ''}
        ${data.isRevoked ? 'opacity-40 grayscale' : ''}
        transition-all duration-200
      `}
    >
      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 !bg-violet-400 border-2 !border-zinc-900"
      />
      <div className="flex items-center gap-2">
        <span className="text-lg">📬</span>
        <div>
          <div className="font-semibold text-sm">{data.label}</div>
          <div className="text-xs text-violet-200/70">Personal Mailbox</div>
        </div>
      </div>
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 !bg-violet-400 border-2 !border-zinc-900"
      />
    </div>
  );
}

export default memo(PersonalMailboxNode);
