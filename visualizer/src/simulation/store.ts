import { create } from 'zustand';
import { Node, Edge } from 'reactflow';
import {
  User,
  File,
  PersonalMailbox,
  SharedMailbox,
  Share,
  SimulationEvent,
  DesignMode,
  RevokeResult,
} from './types';

interface SimulationState {
  // Design mode
  designMode: DesignMode;

  // Entities
  users: User[];
  files: File[];
  personalMailboxes: PersonalMailbox[];
  sharedMailboxes: SharedMailbox[]; // These are "IntNodes" in the blog terminology
  shares: Share[];

  // Event log
  events: SimulationEvent[];

  // Selected node for inspection
  selectedNodeId: string | null;

  // Actions
  addUser: (name: string) => void;
  createFile: (ownerName: string, filename: string) => void;
  shareFile: (fromName: string, toName: string, filename: string) => void;
  revokeAccess: (ownerName: string, targetName: string, filename: string) => RevokeResult;
  setDesignMode: (mode: DesignMode) => void;
  setSelectedNode: (nodeId: string | null) => void;
  reset: () => void;

  // Derived data for React Flow
  getNodes: () => Node[];
  getEdges: () => Edge[];
}

const generateId = () => Math.random().toString(36).substring(2, 9);

const createEvent = (
  type: SimulationEvent['type'],
  description: string,
  details?: SimulationEvent['details']
): SimulationEvent => ({
  id: generateId(),
  timestamp: Date.now(),
  type,
  description,
  details,
});

export const useSimulationStore = create<SimulationState>((set, get) => ({
  designMode: 'working',
  users: [],
  files: [],
  personalMailboxes: [],
  sharedMailboxes: [],
  shares: [],
  events: [],
  selectedNodeId: null,

  addUser: (name: string) => {
    const id = generateId();
    set((state) => ({
      users: [...state.users, { id, name }],
      events: [...state.events, createEvent('user_added', `${name} joined`)],
    }));
  },

  createFile: (ownerName: string, filename: string) => {
    const state = get();
    const owner = state.users.find((u) => u.name === ownerName);
    if (!owner) return;

    const fileId = generateId();
    const personalMailboxId = generateId();

    if (state.designMode === 'working') {
      // Working design: Owner -> Owner's Mailbox -> Owner's IntNode -> FileStruct
      // Owner gets their own IntNode (which only they use)
      const ownerIntNodeId = generateId();
      set((state) => ({
        files: [...state.files, { id: fileId, filename, ownerId: owner.id }],
        // Owner's IntNode - this one is special, only owner uses it
        sharedMailboxes: [
          ...state.sharedMailboxes,
          { id: ownerIntNodeId, fileId, cohort: [owner.id], createdForUserId: owner.id },
        ],
        personalMailboxes: [
          ...state.personalMailboxes,
          { id: personalMailboxId, userId: owner.id, pointsTo: ownerIntNodeId },
        ],
        events: [
          ...state.events,
          createEvent('file_created', `${ownerName} created "${filename}"`),
        ],
      }));
    } else {
      // Broken design: Owner -> Owner's Mailbox -> FileStruct directly (no IntNode layer)
      set((state) => ({
        files: [...state.files, { id: fileId, filename, ownerId: owner.id }],
        personalMailboxes: [
          ...state.personalMailboxes,
          { id: personalMailboxId, userId: owner.id, pointsTo: fileId },
        ],
        events: [
          ...state.events,
          createEvent('file_created', `${ownerName} created "${filename}"`),
        ],
      }));
    }
  },

  shareFile: (fromName: string, toName: string, filename: string) => {
    const state = get();
    const fromUser = state.users.find((u) => u.name === fromName);
    const toUser = state.users.find((u) => u.name === toName);
    const file = state.files.find((f) => f.filename === filename);
    if (!fromUser || !toUser || !file) return;

    const shareId = generateId();
    const newPersonalMailboxId = generateId();
    const isOwnerSharing = file.ownerId === fromUser.id;

    if (state.designMode === 'working') {
      // WORKING DESIGN (with IntNodes)
      // Key insight from blog:
      // - When OWNER shares: create a NEW IntNode for that recipient's cohort
      // - When NON-OWNER shares: recipient joins the SAME IntNode as the sharer

      if (isOwnerSharing) {
        // Owner is sharing -> create a NEW IntNode for this recipient
        // This IntNode will be shared by the recipient and anyone they share with
        const newIntNodeId = generateId();
        set((state) => ({
          shares: [
            ...state.shares,
            { id: shareId, fromUserId: fromUser.id, toUserId: toUser.id, fileId: file.id, sharedMailboxId: newIntNodeId },
          ],
          // Create new IntNode for this sharing cohort
          sharedMailboxes: [
            ...state.sharedMailboxes,
            { id: newIntNodeId, fileId: file.id, cohort: [toUser.id], createdForUserId: toUser.id },
          ],
          // Recipient's mailbox points to their new IntNode
          personalMailboxes: [
            ...state.personalMailboxes,
            { id: newPersonalMailboxId, userId: toUser.id, pointsTo: newIntNodeId },
          ],
          events: [
            ...state.events,
            createEvent('file_shared', `${fromName} shared "${filename}" with ${toName} (new IntNode created)`),
          ],
        }));
      } else {
        // Non-owner is sharing -> recipient joins the SAME IntNode as the sharer
        // Find which IntNode the sharer uses
        const sharerMailbox = state.personalMailboxes.find(
          (pm) => pm.userId === fromUser.id && state.sharedMailboxes.some(sm => sm.id === pm.pointsTo && sm.fileId === file.id)
        );
        if (!sharerMailbox) return;

        const sharerIntNodeId = sharerMailbox.pointsTo;

        set((state) => ({
          shares: [
            ...state.shares,
            { id: shareId, fromUserId: fromUser.id, toUserId: toUser.id, fileId: file.id, sharedMailboxId: sharerIntNodeId },
          ],
          // Recipient's mailbox points to the SAME IntNode as the sharer
          personalMailboxes: [
            ...state.personalMailboxes,
            { id: newPersonalMailboxId, userId: toUser.id, pointsTo: sharerIntNodeId },
          ],
          // Update the cohort to include the new user
          sharedMailboxes: state.sharedMailboxes.map((sm) =>
            sm.id === sharerIntNodeId
              ? { ...sm, cohort: [...sm.cohort, toUser.id] }
              : sm
          ),
          events: [
            ...state.events,
            createEvent('file_shared', `${fromName} shared "${filename}" with ${toName} (same IntNode)`),
          ],
        }));
      }
    } else {
      // BROKEN DESIGN (checkpoint design - no IntNode layer)
      // Everyone's mailbox points DIRECTLY to the FileStruct
      // This is why revocation breaks - Alice can't update mailboxes she doesn't know about
      set((state) => ({
        shares: [
          ...state.shares,
          { id: shareId, fromUserId: fromUser.id, toUserId: toUser.id, fileId: file.id },
        ],
        personalMailboxes: [
          ...state.personalMailboxes,
          { id: newPersonalMailboxId, userId: toUser.id, pointsTo: file.id },
        ],
        events: [
          ...state.events,
          createEvent('file_shared', `${fromName} shared "${filename}" with ${toName}`),
        ],
      }));
    }
  },

  revokeAccess: (ownerName: string, targetName: string, filename: string): RevokeResult => {
    const state = get();
    const owner = state.users.find((u) => u.name === ownerName);
    const target = state.users.find((u) => u.name === targetName);
    const file = state.files.find((f) => f.filename === filename);
    if (!owner || !target || !file || file.ownerId !== owner.id) {
      return { success: false, usersWhoLostAccess: [], incorrectlyRevoked: [] };
    }

    if (state.designMode === 'working') {
      // WORKING DESIGN: Delete the target's IntNode
      // This severs access for the target AND everyone who got access through them
      // But other IntNodes (other cohorts) remain intact

      // Find the IntNode that was created when owner shared with target
      const targetIntNode = state.sharedMailboxes.find(
        (sm) => sm.fileId === file.id && sm.createdForUserId === target.id
      );

      if (!targetIntNode) {
        return { success: false, usersWhoLostAccess: [], incorrectlyRevoked: [] };
      }

      // Find all users whose mailboxes point to this IntNode (the whole cohort)
      const affectedMailboxes = state.personalMailboxes.filter(
        (pm) => pm.pointsTo === targetIntNode.id
      );
      const usersWhoLostAccess = affectedMailboxes.map((pm) => {
        const user = state.users.find((u) => u.id === pm.userId);
        return user?.name || pm.userId;
      });

      set((state) => ({
        // Delete the IntNode
        sharedMailboxes: state.sharedMailboxes.filter((sm) => sm.id !== targetIntNode.id),
        // Delete all mailboxes pointing to that IntNode
        personalMailboxes: state.personalMailboxes.filter((pm) => pm.pointsTo !== targetIntNode.id),
        // Remove related shares
        shares: state.shares.filter((s) => s.sharedMailboxId !== targetIntNode.id),
        events: [
          ...state.events,
          createEvent('access_revoked', `${ownerName} revoked ${targetName}'s access to "${filename}"`, {
            correctlyRevoked: usersWhoLostAccess,
          }),
        ],
      }));

      return { success: true, usersWhoLostAccess, incorrectlyRevoked: [] };
    } else {
      // BROKEN DESIGN: This is the bug we're demonstrating!
      // Alice has to re-encrypt and move the file to a new UUID
      // But she can only update mailboxes SHE created
      // Mailboxes created by Bob/Carol/etc. still point to the old location

      // Find all shares downstream of target
      const findDownstreamUsers = (userId: string): string[] => {
        const downstream: string[] = [];
        const directShares = state.shares.filter(
          (s) => s.fromUserId === userId && s.fileId === file.id
        );
        for (const share of directShares) {
          downstream.push(share.toUserId);
          downstream.push(...findDownstreamUsers(share.toUserId));
        }
        return downstream;
      };

      const downstreamUserIds = findDownstreamUsers(target.id);
      const allRevokedIds = [target.id, ...downstreamUserIds];

      // In the broken design, Alice can only see shares SHE made
      // She'll update those, but downstream shares she doesn't know about are broken
      const sharesFromOwner = state.shares.filter(
        (s) => s.fromUserId === owner.id && s.fileId === file.id
      );
      const usersOwnerDirectlySharedWith = sharesFromOwner.map((s) => s.toUserId);

      // Users who should keep access: those owner shared with directly (except target)
      const shouldKeepAccess = usersOwnerDirectlySharedWith.filter((id) => id !== target.id);

      // Users who incorrectly lost access: downstream of someone who should keep access
      // but their mailbox pointed to the old FileStruct location
      const incorrectlyRevoked: string[] = [];
      for (const userId of shouldKeepAccess) {
        const theirDownstream = findDownstreamUsers(userId);
        incorrectlyRevoked.push(...theirDownstream);
      }
      const incorrectlyRevokedNames = incorrectlyRevoked.map((id) => {
        const user = state.users.find((u) => u.id === id);
        return user?.name || id;
      });

      const usersWhoLostAccess = allRevokedIds.map((id) => {
        const user = state.users.find((u) => u.id === id);
        return user?.name || id;
      });

      set((state) => ({
        personalMailboxes: state.personalMailboxes.filter(
          (pm) => !allRevokedIds.includes(pm.userId) || !state.files.some(f => f.id === pm.pointsTo && f.id === file.id)
        ),
        shares: state.shares.filter(
          (s) => !(allRevokedIds.includes(s.toUserId) && s.fileId === file.id)
        ),
        events: [
          ...state.events,
          createEvent('access_revoked', `${ownerName} revoked ${targetName}'s access to "${filename}"`, {
            correctlyRevoked: usersWhoLostAccess,
            incorrectlyRevoked: incorrectlyRevokedNames,
          }),
        ],
      }));

      return { success: true, usersWhoLostAccess, incorrectlyRevoked: incorrectlyRevokedNames };
    }
  },

  setDesignMode: (mode: DesignMode) => {
    set({ designMode: mode });
    get().reset();
  },

  setSelectedNode: (nodeId: string | null) => {
    set({ selectedNodeId: nodeId });
  },

  reset: () => {
    set({
      users: [],
      files: [],
      personalMailboxes: [],
      sharedMailboxes: [],
      shares: [],
      events: [],
      selectedNodeId: null,
    });
  },

  getNodes: (): Node[] => {
    const state = get();
    const nodes: Node[] = [];
    let yOffset = 0;

    // User nodes (top row)
    state.users.forEach((user, index) => {
      const file = state.files[0];
      nodes.push({
        id: `user-${user.id}`,
        type: 'user',
        position: { x: index * 200, y: yOffset },
        data: {
          label: user.name,
          isOwner: file?.ownerId === user.id,
        },
      });
    });
    yOffset += 120;

    // Personal mailbox nodes
    state.personalMailboxes.forEach((pm) => {
      const user = state.users.find((u) => u.id === pm.userId);
      const userIndex = state.users.findIndex((u) => u.id === pm.userId);
      nodes.push({
        id: `pm-${pm.id}`,
        type: 'personalMailbox',
        position: { x: userIndex * 200, y: yOffset },
        data: {
          label: `${user?.name}'s Mailbox`,
          owner: user?.name || '',
        },
      });
    });
    yOffset += 120;

    // IntNode/Shared mailbox nodes (only in working design)
    if (state.designMode === 'working') {
      // Position IntNodes based on their cohort - try to center them under their users
      state.sharedMailboxes.forEach((sm, index) => {
        const cohortNames = sm.cohort.map((userId) => {
          const user = state.users.find((u) => u.id === userId);
          return user?.name || userId;
        });

        // Find the user this IntNode was created for
        const createdForUser = state.users.find((u) => u.id === sm.createdForUserId);
        const createdForName = createdForUser?.name || 'Owner';

        // Calculate x position as average of cohort members' positions
        const cohortIndices = sm.cohort.map((userId) =>
          state.users.findIndex((u) => u.id === userId)
        ).filter((i) => i !== -1);
        const avgX = cohortIndices.length > 0
          ? cohortIndices.reduce((a, b) => a + b, 0) / cohortIndices.length * 200
          : index * 220;

        nodes.push({
          id: `sm-${sm.id}`,
          type: 'sharedMailbox',
          position: { x: avgX, y: yOffset },
          data: {
            label: `IntNode (${createdForName})`,
            cohort: cohortNames,
          },
        });
      });
      yOffset += 120;
    }

    // File struct nodes
    state.files.forEach((file) => {
      const owner = state.users.find((u) => u.id === file.ownerId);
      // Center the file under all users
      const centerX = state.users.length > 0 ? ((state.users.length - 1) * 200) / 2 : 100;
      nodes.push({
        id: `file-${file.id}`,
        type: 'fileStruct',
        position: { x: centerX, y: yOffset },
        data: {
          label: file.filename,
          filename: file.filename,
          owner: owner?.name || '',
        },
      });
    });
    yOffset += 120;

    // Content nodes
    state.files.forEach((file) => {
      const centerX = state.users.length > 0 ? ((state.users.length - 1) * 200) / 2 : 100;
      nodes.push({
        id: `content-${file.id}`,
        type: 'content',
        position: { x: centerX, y: yOffset },
        data: {
          label: 'Content',
          chunkCount: 1,
        },
      });
    });

    return nodes;
  },

  getEdges: (): Edge[] => {
    const state = get();
    const edges: Edge[] = [];

    // User -> Personal Mailbox edges
    state.personalMailboxes.forEach((pm) => {
      edges.push({
        id: `edge-user-pm-${pm.id}`,
        source: `user-${pm.userId}`,
        target: `pm-${pm.id}`,
        animated: true,
        style: { stroke: '#8B5CF6' },
      });
    });

    if (state.designMode === 'working') {
      // Personal Mailbox -> IntNode edges
      state.personalMailboxes.forEach((pm) => {
        edges.push({
          id: `edge-pm-sm-${pm.id}`,
          source: `pm-${pm.id}`,
          target: `sm-${pm.pointsTo}`,
          animated: true,
          style: { stroke: '#F59E0B' },
        });
      });

      // IntNode -> FileStruct edges
      state.sharedMailboxes.forEach((sm) => {
        edges.push({
          id: `edge-sm-file-${sm.id}`,
          source: `sm-${sm.id}`,
          target: `file-${sm.fileId}`,
          animated: true,
          style: { stroke: '#10B981' },
        });
      });
    } else {
      // Broken design: Personal Mailbox -> FileStruct directly (no IntNode layer)
      state.personalMailboxes.forEach((pm) => {
        edges.push({
          id: `edge-pm-file-${pm.id}`,
          source: `pm-${pm.id}`,
          target: `file-${pm.pointsTo}`,
          animated: true,
          style: { stroke: '#10B981' },
        });
      });
    }

    // FileStruct -> Content edges
    state.files.forEach((file) => {
      edges.push({
        id: `edge-file-content-${file.id}`,
        source: `file-${file.id}`,
        target: `content-${file.id}`,
        animated: true,
        style: { stroke: '#6B7280' },
      });
    });

    return edges;
  },
}));
