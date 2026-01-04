'use client';

import { useState, useMemo } from 'react';
import dynamic from 'next/dynamic';
import { Node, Edge } from 'reactflow';
import { walkthroughSteps, DiagramState } from './walkthroughSteps';

const DiagramCanvas = dynamic(
  () => import('@/components/Diagram/DiagramCanvas'),
  { ssr: false }
);

// Convert DiagramState to React Flow nodes and edges
function buildNodesAndEdges(state: DiagramState): { nodes: Node[]; edges: Edge[] } {
  const nodes: Node[] = [];
  const edges: Edge[] = [];
  let yOffset = 0;

  // User nodes
  state.users.forEach((user, index) => {
    const nodeId = `user-${user.name.toLowerCase()}`;
    nodes.push({
      id: nodeId,
      type: 'user',
      position: { x: index * 200, y: yOffset },
      data: {
        label: user.name,
        isOwner: user.isOwner,
        isRevoked: state.revokedUser === user.name,
      },
      className: state.highlightNodes?.includes(nodeId) ? 'ring-2 ring-yellow-400' : '',
    });
  });

  if (state.users.length === 0) {
    return { nodes, edges };
  }

  yOffset += 120;

  // Personal mailbox nodes (for users with file access)
  const usersWithAccess = new Set<string>();

  // Owner always has access
  const owner = state.users.find(u => u.isOwner);
  if (owner && state.files.length > 0) {
    usersWithAccess.add(owner.name);
  }

  // Build sharing graph
  state.shares.forEach(share => {
    if (usersWithAccess.has(share.from) || state.users.find(u => u.name === share.from)?.isOwner) {
      usersWithAccess.add(share.to);
    }
  });
  // Second pass to catch transitive shares
  state.shares.forEach(share => {
    if (usersWithAccess.has(share.from)) {
      usersWithAccess.add(share.to);
    }
  });

  // Remove revoked users and their downstream in broken state
  if (state.revokedUser && !state.showIntNodes) {
    usersWithAccess.delete(state.revokedUser);
    // In broken state, downstream users have dangling pointers
  }

  state.users.forEach((user, index) => {
    if (!usersWithAccess.has(user.name)) return;

    const nodeId = `pm-${user.name.toLowerCase()}`;
    const isBroken = state.showBrokenState && !state.users.find(u => u.name === user.name)?.isOwner;

    nodes.push({
      id: nodeId,
      type: 'personalMailbox',
      position: { x: index * 200, y: yOffset },
      data: {
        label: `${user.name}'s Mailbox`,
        owner: user.name,
        isBroken,
      },
      className: state.highlightNodes?.includes(nodeId) ? 'ring-2 ring-yellow-400' : '',
    });

    // Edge from user to mailbox
    edges.push({
      id: `edge-user-pm-${user.name}`,
      source: `user-${user.name.toLowerCase()}`,
      target: nodeId,
      animated: true,
      style: { stroke: state.revokedUser === user.name ? '#EF4444' : '#8B5CF6' },
    });
  });

  yOffset += 120;

  // IntNode layer (only in working design)
  if (state.showIntNodes && state.files.length > 0) {
    // Build cohorts based on sharing structure
    // Owner has their own IntNode
    // Each direct share from owner creates a new IntNode
    // Non-owner shares add to existing IntNode

    const cohorts: Map<string, string[]> = new Map();

    if (owner) {
      cohorts.set(owner.name, [owner.name]);
    }

    // First pass: owner's direct shares create new IntNodes
    state.shares.forEach(share => {
      const fromUser = state.users.find(u => u.name === share.from);
      if (fromUser?.isOwner) {
        cohorts.set(share.to, [share.to]);
      }
    });

    // Second pass: non-owner shares add to existing cohort
    state.shares.forEach(share => {
      const fromUser = state.users.find(u => u.name === share.from);
      if (!fromUser?.isOwner) {
        // Find which cohort the sharer is in
        const cohortKeys = Array.from(cohorts.keys());
        for (const key of cohortKeys) {
          const members = cohorts.get(key);
          if (members && members.includes(share.from)) {
            members.push(share.to);
            break;
          }
        }
      }
    });

    // Render IntNodes
    let intNodeIndex = 0;
    cohorts.forEach((members, createdFor) => {
      // Skip revoked cohorts
      if (state.revokedUser && (createdFor === state.revokedUser || members.includes(state.revokedUser))) {
        // Check if this is the revoked user's cohort
        const cohortOwner = state.users.find(u => u.name === createdFor);
        if (!cohortOwner?.isOwner) {
          return; // Skip this IntNode entirely if revoked
        }
      }

      const nodeId = `sm-${createdFor.toLowerCase()}`;
      const memberIndices = members
        .map(m => state.users.findIndex(u => u.name === m))
        .filter(i => i !== -1);
      const avgX = memberIndices.length > 0
        ? memberIndices.reduce((a, b) => a + b, 0) / memberIndices.length * 200
        : intNodeIndex * 220;

      nodes.push({
        id: nodeId,
        type: 'sharedMailbox',
        position: { x: avgX, y: yOffset },
        data: {
          label: `IntNode (${createdFor})`,
          cohort: members,
        },
        className: state.highlightNodes?.includes(nodeId) ? 'ring-2 ring-yellow-400' : '',
      });

      // Edges from mailboxes to IntNode
      members.forEach(member => {
        if (!usersWithAccess.has(member)) return;
        edges.push({
          id: `edge-pm-sm-${member}`,
          source: `pm-${member.toLowerCase()}`,
          target: nodeId,
          animated: true,
          style: { stroke: '#F59E0B' },
        });
      });

      intNodeIndex++;
    });

    yOffset += 120;

    // Edges from IntNodes to file
    if (state.files.length > 0) {
      cohorts.forEach((members, createdFor) => {
        if (state.revokedUser && createdFor === state.revokedUser) return;

        const cohortOwner = state.users.find(u => u.name === createdFor);
        if (state.revokedUser && !cohortOwner?.isOwner && members.includes(state.revokedUser)) return;

        edges.push({
          id: `edge-sm-file-${createdFor}`,
          source: `sm-${createdFor.toLowerCase()}`,
          target: `file-${state.files[0].name}`,
          animated: true,
          style: { stroke: '#10B981' },
        });
      });
    }
  } else if (state.files.length > 0) {
    // Broken design: mailboxes point directly to file
    state.users.forEach(user => {
      if (!usersWithAccess.has(user.name)) return;

      const isBroken = state.showBrokenState && state.revokedUser && user.name !== owner?.name;

      edges.push({
        id: `edge-pm-file-${user.name}`,
        source: `pm-${user.name.toLowerCase()}`,
        target: `file-${state.files[0].name}`,
        animated: !isBroken,
        style: {
          stroke: isBroken ? '#EF4444' : '#10B981',
          strokeDasharray: isBroken ? '5,5' : undefined,
        },
      });
    });
  }

  // File node
  if (state.files.length > 0) {
    const file = state.files[0];
    const centerX = state.users.length > 0 ? ((state.users.length - 1) * 200) / 2 : 100;

    nodes.push({
      id: `file-${file.name}`,
      type: 'fileStruct',
      position: { x: centerX, y: yOffset },
      data: {
        label: file.name,
        filename: file.name,
        owner: file.owner,
      },
      className: state.highlightNodes?.includes(`file-${file.name}`) ? 'ring-2 ring-yellow-400' : '',
    });

    yOffset += 120;

    // Content node
    nodes.push({
      id: `content-${file.name}`,
      type: 'content',
      position: { x: centerX, y: yOffset },
      data: {
        label: 'Content',
        chunkCount: 1,
      },
    });

    edges.push({
      id: `edge-file-content-${file.name}`,
      source: `file-${file.name}`,
      target: `content-${file.name}`,
      animated: true,
      style: { stroke: '#6B7280' },
    });
  }

  return { nodes, edges };
}

export default function WalkthroughContainer() {
  const [currentStep, setCurrentStep] = useState(0);
  const step = walkthroughSteps[currentStep];

  const { nodes, edges } = useMemo(
    () => buildNodesAndEdges(step.diagram),
    [step.diagram]
  );

  const goNext = () => {
    if (currentStep < walkthroughSteps.length - 1) {
      setCurrentStep(currentStep + 1);
    }
  };

  const goBack = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const progress = ((currentStep + 1) / walkthroughSteps.length) * 100;

  return (
    <div className="flex flex-col h-full overflow-hidden">
      {/* Top narrative panel - front and center */}
      <div className="flex-shrink-0 bg-gradient-to-b from-zinc-900 to-zinc-900/95 border-b border-zinc-800">
        <div className="max-w-4xl mx-auto px-4 py-5">
          {/* Step counter and progress */}
          <div className="flex items-center gap-3 mb-3">
            <span className="text-xs font-mono text-zinc-500 bg-zinc-800 px-2 py-1 rounded">
              {currentStep + 1}/{walkthroughSteps.length}
            </span>
            <div className="flex-1 bg-zinc-800 rounded-full h-1">
              <div
                className="bg-indigo-500 h-1 rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
          </div>

          {/* Main content row: Back | Title + Narrative | Next */}
          <div className="flex items-center gap-4">
            {/* Back button */}
            <button
              onClick={goBack}
              disabled={currentStep === 0}
              className={`flex-shrink-0 w-12 h-12 rounded-full flex items-center justify-center transition-all ${
                currentStep === 0
                  ? 'text-zinc-700 cursor-not-allowed'
                  : 'text-zinc-400 hover:text-white hover:bg-zinc-800 border border-zinc-700 hover:border-zinc-600'
              }`}
              aria-label="Previous step"
            >
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>

            {/* Center content */}
            <div className="flex-1 min-w-0">
              <h2 className="text-xl md:text-2xl font-bold text-indigo-400 mb-2">{step.title}</h2>
              <p className="text-base md:text-lg text-zinc-300 leading-relaxed">{step.narrative}</p>

              {/* Note inline if present */}
              {step.note && (
                <div className="mt-3 flex items-start gap-2 text-amber-400 text-sm">
                  <span className="text-amber-500">💡</span>
                  <span>{step.note}</span>
                </div>
              )}
            </div>

            {/* Next button - extra prominent on first step */}
            {currentStep === 0 ? (
              <button
                onClick={goNext}
                className="flex-shrink-0 flex items-center gap-2 px-6 py-3 rounded-full text-white bg-indigo-600 hover:bg-indigo-500 shadow-lg shadow-indigo-600/40 transition-all animate-pulse hover:animate-none font-semibold"
                aria-label="Get started"
              >
                <span>Get Started</span>
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </button>
            ) : (
              <button
                onClick={goNext}
                disabled={currentStep === walkthroughSteps.length - 1}
                className={`flex-shrink-0 w-12 h-12 rounded-full flex items-center justify-center transition-all ${
                  currentStep === walkthroughSteps.length - 1
                    ? 'text-zinc-700 cursor-not-allowed'
                    : 'text-white bg-indigo-600 hover:bg-indigo-500 shadow-lg shadow-indigo-600/30'
                }`}
                aria-label="Next step"
              >
                <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </button>
            )}
          </div>

          {/* Step dots below */}
          <div className="flex items-center justify-center gap-1.5 mt-4">
            {walkthroughSteps.map((_, index) => (
              <button
                key={index}
                onClick={() => setCurrentStep(index)}
                className={`h-1.5 rounded-full transition-all ${
                  index === currentStep
                    ? 'bg-indigo-500 w-6'
                    : index < currentStep
                    ? 'bg-indigo-700 w-1.5 hover:bg-indigo-600'
                    : 'bg-zinc-700 w-1.5 hover:bg-zinc-600'
                }`}
                title={walkthroughSteps[index].title}
              />
            ))}
          </div>
        </div>
      </div>

      {/* Diagram area - takes remaining space */}
      <div className="flex-1 min-h-0 bg-[#0a0a0f] relative">
        {nodes.length > 0 ? (
          <div className="absolute inset-0">
            <DiagramCanvas
              key={`walkthrough-${currentStep}`}
              initialNodes={nodes}
              initialEdges={edges}
            />
          </div>
        ) : (
          <div className="h-full flex items-center justify-center">
            <div className="text-center max-w-md px-4">
              <div className="text-6xl mb-4">🔐</div>
              <h2 className="text-2xl font-light text-zinc-400 mb-2">
                End-to-End Encrypted File Sharing
              </h2>
              <p className="text-zinc-500">
                Click Next to begin the walkthrough
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
