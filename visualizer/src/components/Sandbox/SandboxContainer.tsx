'use client';

import { useState, useMemo, useCallback } from 'react';
import dynamic from 'next/dynamic';
import { useSimulationStore } from '@/simulation/store';
import { Node } from 'reactflow';

const DiagramCanvas = dynamic(
  () => import('@/components/Diagram/DiagramCanvas'),
  { ssr: false }
);

const DEFAULT_NAMES = ['Alice', 'Bob', 'Carol', 'Dave', 'Eve', 'Frank', 'Grace', 'Heidi', 'Ivan', 'Judy'];

export default function SandboxContainer() {
  const {
    designMode,
    users,
    files,
    events,
    personalMailboxes,
    sharedMailboxes,
    addUser,
    createFile,
    shareFile,
    revokeAccess,
    setDesignMode,
    reset,
    getNodes,
    getEdges,
  } = useSimulationStore();

  // Modal state for sharing
  const [shareModal, setShareModal] = useState<{ from: string; isOpen: boolean } | null>(null);
  const [revokeModal, setRevokeModal] = useState<{ isOpen: boolean } | null>(null);

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const baseNodes = useMemo(() => getNodes(), [users, files, personalMailboxes, sharedMailboxes, designMode, getNodes]);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  const edges = useMemo(() => getEdges(), [users, files, personalMailboxes, sharedMailboxes, designMode, getEdges]);

  // Add callbacks to nodes for hover actions
  const nodes = useMemo(() => {
    const hasFile = files.length > 0;
    return baseNodes.map(node => {
      if (node.type === 'user') {
        const userName = node.data.label;
        const userHasAccess = personalMailboxes.some(pm => {
          const user = users.find(u => u.name === userName);
          return user && pm.userId === user.id;
        });
        return {
          ...node,
          data: {
            ...node.data,
            hasFile: hasFile && userHasAccess,
            onShare: () => setShareModal({ from: userName, isOpen: true }),
          },
        };
      }
      if (node.type === 'fileStruct') {
        return {
          ...node,
          data: {
            ...node.data,
            onRevoke: () => setRevokeModal({ isOpen: true }),
          },
        };
      }
      return node;
    });
  }, [baseNodes, files.length, personalMailboxes, users]);

  // Get the next available default name
  const getNextName = useCallback(() => {
    const usedNames = new Set(users.map(u => u.name));
    return DEFAULT_NAMES.find(name => !usedNames.has(name)) || `User${users.length + 1}`;
  }, [users]);

  const handleAddUser = () => {
    const name = getNextName();
    addUser(name);
  };

  const handleCreateFile = () => {
    if (users.length === 0) {
      addUser(getNextName());
    }
    // Use the first user (owner) and a simple filename
    const owner = users[0]?.name || getNextName();
    if (!users[0]) {
      addUser(owner);
      setTimeout(() => {
        createFile(owner, 'secret.txt');
      }, 0);
    } else {
      createFile(owner, 'secret.txt');
    }
  };

  const handleNodeClick = (node: Node) => {
    // If clicking on a user node and there's a file, open share modal
    if (node.type === 'user' && files.length > 0) {
      const userName = node.data.label;
      setShareModal({ from: userName, isOpen: true });
    }
  };

  const handleShare = (to: string) => {
    if (shareModal && files.length > 0) {
      shareFile(shareModal.from, to, files[0].filename);
      setShareModal(null);
    }
  };

  const handleRevoke = (target: string) => {
    if (files.length > 0) {
      const owner = users.find(u => u.id === files[0].ownerId);
      if (owner) {
        revokeAccess(owner.name, target, files[0].filename);
      }
    }
    setRevokeModal(null);
  };

  // Get users who can be shared with (don't already have access from the sharer)
  const getShareableUsers = (fromName: string) => {
    return users.filter(u => u.name !== fromName);
  };

  // Get users who have access (for revocation) - only direct shares from owner
  const getUsersWithAccess = () => {
    if (files.length === 0) return [];
    const file = files[0];

    if (designMode === 'working') {
      // In working mode, we can only revoke users who have their own IntNode
      // (i.e., users the owner directly shared with)
      return sharedMailboxes
        .filter(sm => sm.fileId === file.id && sm.createdForUserId !== file.ownerId)
        .map(sm => users.find(u => u.id === sm.createdForUserId))
        .filter((u): u is NonNullable<typeof u> => u !== undefined);
    } else {
      // In broken mode, return all users except owner who have a mailbox
      return users.filter(u => {
        if (u.id === file.ownerId) return false;
        return personalMailboxes.some(pm => pm.userId === u.id);
      });
    }
  };

  const hasFile = files.length > 0;
  const owner = hasFile ? users.find(u => u.id === files[0].ownerId) : null;

  return (
    <div className="flex flex-col h-full">
      {/* Top toolbar */}
      <div className="bg-zinc-900 border-b border-zinc-800 px-4 py-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            {/* Quick actions */}
            <button
              onClick={handleAddUser}
              disabled={users.length >= 10}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:bg-zinc-700 disabled:text-zinc-500 text-white rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
            >
              <span>+ Add User</span>
              {users.length > 0 && (
                <span className="text-blue-200 text-xs">({getNextName()})</span>
              )}
            </button>

            {!hasFile && users.length > 0 && (
              <button
                onClick={handleCreateFile}
                className="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg text-sm font-medium transition-colors"
              >
                + Create File
              </button>
            )}

            {hasFile && (
              <button
                onClick={() => setRevokeModal({ isOpen: true })}
                disabled={getUsersWithAccess().length === 0}
                className="px-4 py-2 bg-red-600 hover:bg-red-500 disabled:bg-zinc-700 disabled:text-zinc-500 text-white rounded-lg text-sm font-medium transition-colors"
              >
                Revoke Access
              </button>
            )}

            <button
              onClick={reset}
              className="px-4 py-2 border border-zinc-700 text-zinc-400 hover:text-zinc-200 hover:bg-zinc-800 rounded-lg text-sm transition-colors"
            >
              Reset
            </button>
          </div>

          {/* Design mode toggle */}
          <div className="flex items-center gap-4">
            <span className="text-zinc-500 text-sm">Design:</span>
            <div className="flex bg-zinc-800 rounded-lg p-1">
              <button
                onClick={() => setDesignMode('broken')}
                className={`px-3 py-1.5 rounded text-sm font-medium transition-colors ${
                  designMode === 'broken'
                    ? 'bg-red-600 text-white'
                    : 'text-zinc-400 hover:text-zinc-200'
                }`}
              >
                Broken
              </button>
              <button
                onClick={() => setDesignMode('working')}
                className={`px-3 py-1.5 rounded text-sm font-medium transition-colors ${
                  designMode === 'working'
                    ? 'bg-emerald-600 text-white'
                    : 'text-zinc-400 hover:text-zinc-200'
                }`}
              >
                Working
              </button>
            </div>
          </div>
        </div>

        {/* Instructions */}
        {users.length === 0 && (
          <p className="text-zinc-500 text-sm mt-2">
            Click &ldquo;Add User&rdquo; to start, then create a file and share it.
          </p>
        )}
        {users.length > 0 && !hasFile && (
          <p className="text-zinc-500 text-sm mt-2">
            Click &ldquo;Create File&rdquo; to create a file owned by {users[0].name}.
          </p>
        )}
        {hasFile && (
          <p className="text-zinc-500 text-sm mt-2">
            Click on a user node to share the file from them. Owner: <span className="text-indigo-400">{owner?.name}</span>
          </p>
        )}
      </div>

      {/* Main content */}
      <div className="flex-1 flex">
        {/* Diagram */}
        <div className="flex-1 bg-[#0a0a0f]">
          {nodes.length > 0 ? (
            <DiagramCanvas
              key={`${nodes.length}-${edges.length}-${designMode}`}
              initialNodes={nodes}
              initialEdges={edges}
              onNodeClick={handleNodeClick}
            />
          ) : (
            <div className="h-full flex items-center justify-center">
              <div className="text-center">
                <div className="text-6xl mb-4">🔐</div>
                <div className="text-2xl text-zinc-600 font-light mb-2">E2E File Sharing Visualizer</div>
                <p className="text-zinc-500">Add users and create a file to see the architecture</p>
              </div>
            </div>
          )}
        </div>

        {/* Event log sidebar - simplified */}
        <div className="w-64 border-l border-zinc-800 bg-zinc-900 p-4 overflow-y-auto">
          <h3 className="font-semibold text-zinc-100 mb-3 text-sm">Event Log</h3>
          <div className="space-y-2 text-xs">
            {events.length === 0 ? (
              <p className="text-zinc-500 italic">No events yet...</p>
            ) : (
              events
                .slice()
                .reverse()
                .slice(0, 10)
                .map((event) => (
                  <div
                    key={event.id}
                    className={`p-2 rounded border ${
                      event.type === 'access_revoked'
                        ? 'bg-red-900/30 text-red-400 border-red-800/50'
                        : event.type === 'file_shared'
                        ? 'bg-violet-900/30 text-violet-400 border-violet-800/50'
                        : event.type === 'file_created'
                        ? 'bg-emerald-900/30 text-emerald-400 border-emerald-800/50'
                        : 'bg-indigo-900/30 text-indigo-400 border-indigo-800/50'
                    }`}
                  >
                    {event.description}
                  </div>
                ))
            )}
          </div>

          {/* Mode explanation */}
          <div className="mt-6 pt-4 border-t border-zinc-800">
            <div
              className={`p-3 rounded-lg text-xs border ${
                designMode === 'broken'
                  ? 'bg-red-900/20 text-red-400 border-red-800/50'
                  : 'bg-emerald-900/20 text-emerald-400 border-emerald-800/50'
              }`}
            >
              {designMode === 'broken' ? (
                <>
                  <strong className="block mb-1">Broken Design</strong>
                  <p className="opacity-80">
                    No IntNode layer. Revocation breaks sharing chains incorrectly.
                  </p>
                </>
              ) : (
                <>
                  <strong className="block mb-1">Working Design</strong>
                  <p className="opacity-80">
                    IntNode layer enables proper cohort-based revocation.
                  </p>
                </>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Share Modal */}
      {shareModal?.isOpen && (
        <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
          <div className="bg-zinc-900 border border-zinc-700 rounded-xl p-6 w-80 shadow-2xl">
            <h3 className="text-lg font-semibold text-zinc-100 mb-4">
              Share from {shareModal.from}
            </h3>
            <p className="text-zinc-400 text-sm mb-4">
              {shareModal.from === owner?.name
                ? 'This will create a NEW IntNode for the recipient.'
                : 'This will share the SAME IntNode (recipient joins the cohort).'}
            </p>
            <div className="space-y-2">
              {getShareableUsers(shareModal.from).map((user) => (
                <button
                  key={user.id}
                  onClick={() => handleShare(user.name)}
                  className="w-full px-4 py-2 bg-zinc-800 hover:bg-violet-600 text-zinc-200 rounded-lg text-sm transition-colors text-left"
                >
                  Share with <span className="font-semibold">{user.name}</span>
                </button>
              ))}
            </div>
            <button
              onClick={() => setShareModal(null)}
              className="w-full mt-4 px-4 py-2 border border-zinc-700 text-zinc-400 hover:text-zinc-200 rounded-lg text-sm transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {/* Revoke Modal */}
      {revokeModal?.isOpen && (
        <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
          <div className="bg-zinc-900 border border-zinc-700 rounded-xl p-6 w-80 shadow-2xl">
            <h3 className="text-lg font-semibold text-zinc-100 mb-4">
              Revoke Access
            </h3>
            <p className="text-zinc-400 text-sm mb-4">
              {designMode === 'working'
                ? 'Select a user to revoke. This deletes their IntNode and severs their entire cohort.'
                : 'Select a user to revoke. In the broken design, this may incorrectly affect other users.'}
            </p>
            <div className="space-y-2">
              {getUsersWithAccess().map((user) => (
                <button
                  key={user.id}
                  onClick={() => handleRevoke(user.name)}
                  className="w-full px-4 py-2 bg-zinc-800 hover:bg-red-600 text-zinc-200 rounded-lg text-sm transition-colors text-left"
                >
                  Revoke <span className="font-semibold">{user.name}</span>
                </button>
              ))}
            </div>
            <button
              onClick={() => setRevokeModal(null)}
              className="w-full mt-4 px-4 py-2 border border-zinc-700 text-zinc-400 hover:text-zinc-200 rounded-lg text-sm transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
