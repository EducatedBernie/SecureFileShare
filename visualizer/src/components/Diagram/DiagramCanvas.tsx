'use client';

import { useCallback } from 'react';
import ReactFlow, {
  Node,
  Edge,
  Background,
  Controls,
  MiniMap,
  useNodesState,
  useEdgesState,
  BackgroundVariant,
  ConnectionMode,
} from 'reactflow';
import 'reactflow/dist/style.css';

import { nodeTypes } from './nodes';

// Color mapping for minimap
const nodeColor = (node: Node) => {
  switch (node.type) {
    case 'user':
      return '#3B82F6';
    case 'personalMailbox':
      return '#8B5CF6';
    case 'sharedMailbox':
      return '#F59E0B';
    case 'fileStruct':
      return '#10B981';
    case 'content':
      return '#6B7280';
    default:
      return '#94A3B8';
  }
};

interface DiagramCanvasProps {
  initialNodes?: Node[];
  initialEdges?: Edge[];
  onNodeClick?: (node: Node) => void;
  fitView?: boolean;
}

export default function DiagramCanvas({
  initialNodes = [],
  initialEdges = [],
  onNodeClick,
  fitView = true,
}: DiagramCanvasProps) {
  const [nodes, , onNodesChange] = useNodesState(initialNodes);
  const [edges, , onEdgesChange] = useEdgesState(initialEdges);

  const handleNodeClick = useCallback(
    (_: React.MouseEvent, node: Node) => {
      onNodeClick?.(node);
    },
    [onNodeClick]
  );

  return (
    <div className="w-full h-full">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeClick={handleNodeClick}
        nodeTypes={nodeTypes}
        connectionMode={ConnectionMode.Loose}
        fitView={fitView}
        fitViewOptions={{ padding: 0.2 }}
        minZoom={0.5}
        maxZoom={2}
        attributionPosition="bottom-left"
      >
        <Background variant={BackgroundVariant.Dots} gap={16} size={1} />
        <Controls />
        <MiniMap
          nodeColor={nodeColor}
          nodeStrokeWidth={3}
          zoomable
          pannable
        />
      </ReactFlow>
    </div>
  );
}
