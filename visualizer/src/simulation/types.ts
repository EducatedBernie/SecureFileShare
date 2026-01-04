export interface User {
  id: string;
  name: string;
}

export interface File {
  id: string;
  filename: string;
  ownerId: string;
}

export interface PersonalMailbox {
  id: string;
  userId: string;
  pointsTo: string; // sharedMailboxId or fileStructId depending on design mode
}

export interface SharedMailbox {
  id: string;
  fileId: string;
  cohort: string[]; // userIds in this sharing cohort
  createdForUserId: string; // The user this IntNode was created for (used for revocation)
}

export interface Share {
  id: string;
  fromUserId: string;
  toUserId: string;
  fileId: string;
  sharedMailboxId?: string; // only in working design
}

export interface SimulationEvent {
  id: string;
  timestamp: number;
  type: 'user_added' | 'file_created' | 'file_shared' | 'access_revoked';
  description: string;
  details?: {
    incorrectlyRevoked?: string[];
    correctlyRevoked?: string[];
  };
}

export type DesignMode = 'broken' | 'working';

export interface RevokeResult {
  success: boolean;
  usersWhoLostAccess: string[];
  incorrectlyRevoked: string[]; // users who should NOT have lost access
}
