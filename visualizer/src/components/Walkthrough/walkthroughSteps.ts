// Walkthrough steps - each step has narrative text and a diagram state
export interface DiagramState {
  users: { name: string; isOwner?: boolean }[];
  files: { name: string; owner: string }[];
  shares: { from: string; to: string }[];
  showIntNodes: boolean;
  highlightNodes?: string[]; // node ids to highlight
  revokedUser?: string; // user being revoked
  showBrokenState?: boolean; // show the broken design issue
}

export interface WalkthroughStep {
  id: number;
  title: string;
  narrative: string;
  diagram: DiagramState;
  note?: string; // optional side note
}

export const walkthroughSteps: WalkthroughStep[] = [
  // Intro: Context and Threat Model
  {
    id: 1,
    title: "Welcome",
    narrative: "This visualizer explores the architecture of a secure, end-to-end encrypted file sharing system—built as part of UC Berkeley's CS 161 Computer Security course.",
    diagram: {
      users: [],
      files: [],
      shares: [],
      showIntNodes: false,
    },
  },
  {
    id: 2,
    title: "Two Adversaries",
    narrative: "The system defends against two distinct threats that don't collude. Understanding both is key to understanding why the architecture looks the way it does.",
    diagram: {
      users: [],
      files: [],
      shares: [],
      showIntNodes: false,
    },
  },
  {
    id: 3,
    title: "Adversary #1: The Server",
    narrative: "The Datastore Adversary controls the storage server completely. It can read, modify, and delete ANY data. It observes all API calls and knows your source code. But it cannot enumerate what it doesn't know exists.",
    diagram: {
      users: [],
      files: [],
      shares: [],
      showIntNodes: false,
    },
    note: "Per Kerckhoff's principle—security through obscurity is not security.",
  },
  {
    id: 4,
    title: "Adversary #2: The Revoked User",
    narrative: "A previously authorized user whose access was revoked. Before revocation, they recorded UUIDs, keys, and values. After revocation, they can still directly access the server—but only at locations they already know.",
    diagram: {
      users: [],
      files: [],
      shares: [],
      showIntNodes: false,
    },
    note: "Unlike the server, revoked users have limited visibility but active attack capability.",
  },
  {
    id: 5,
    title: "What We Must Protect",
    narrative: "Confidentiality of file contents and filenames (IND-CPA security). Post-revocation secrecy—revoked users learn nothing about future updates. Integrity—detect any tampering. Authorization—only valid users access files.",
    diagram: {
      users: [],
      files: [],
      shares: [],
      showIntNodes: false,
    },
  },
  {
    id: 6,
    title: "The Constraints",
    narrative: "Users authenticate with username/password only. No PKI, no trusted third party. The server provides storage but zero trust. Users must be able to share files and revoke access at will.",
    diagram: {
      users: [],
      files: [],
      shares: [],
      showIntNodes: false,
    },
  },
  {
    id: 7,
    title: "The Challenge",
    narrative: "Build a Dropbox-like system where the server never sees your files. End-to-end encryption. Zero knowledge. And efficient append operations—adding 100 bytes to a 10TB file shouldn't re-upload the whole thing.",
    diagram: {
      users: [],
      files: [],
      shares: [],
      showIntNodes: false,
    },
  },
  {
    id: 8,
    title: "The Key Insight",
    narrative: "Encryption isn't the hard part—it's key management. Anyone can encrypt a file. But sharing encrypted files while supporting revocation? That's where things get interesting.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [],
      showIntNodes: false,
    },
  },
  {
    id: 9,
    title: "Alice Creates a File",
    narrative: "Alice wants to share a file. She encrypts it with a symmetric key, storing it on the server. She keeps her key in her personal mailbox—encrypted with her own public key.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [],
      showIntNodes: false,
      highlightNodes: ["user-alice", "file-secret"],
    },
  },
  {
    id: 10,
    title: "Sharing with Bob",
    narrative: "Alice wants to share with Bob. She encrypts the file key with Bob's public key and places it in Bob's mailbox. Now Bob can decrypt too.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }],
      showIntNodes: false,
    },
  },
  {
    id: 11,
    title: "Bob Shares with Carol",
    narrative: "Bob decides Carol should have access too. He copies his key to Carol's mailbox. Simple, right?",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }, { from: "Bob", to: "Carol" }],
      showIntNodes: false,
    },
  },
  // Act 2: The Problem
  {
    id: 12,
    title: "The Revocation Problem",
    narrative: "Now Alice wants to revoke Bob's access. She re-encrypts the file with a NEW key, updates her mailbox... but what about Carol?",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }, { from: "Bob", to: "Carol" }],
      showIntNodes: false,
      revokedUser: "Bob",
    },
  },
  {
    id: 13,
    title: "Alice Can't See Carol",
    narrative: "Alice doesn't even know Carol exists! Bob shared with her directly. Carol's mailbox still points to the OLD file location with the OLD key.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }, { from: "Bob", to: "Carol" }],
      showIntNodes: false,
      revokedUser: "Bob",
      showBrokenState: true,
      highlightNodes: ["pm-carol"],
    },
    note: "Carol's mailbox is now a dangling pointer!",
  },
  {
    id: 14,
    title: "The Dilemma",
    narrative: "Either Carol keeps access (defeating revocation) or Carol loses access even though Alice only revoked Bob. Both outcomes are wrong.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }, { from: "Bob", to: "Carol" }],
      showIntNodes: false,
      showBrokenState: true,
    },
  },
  // Act 3: The Solution
  {
    id: 15,
    title: "Enter the IntNode",
    narrative: "The solution: add an intermediate layer. Instead of mailboxes pointing directly to the file, they point to an IntNode (intermediate node) which then points to the file.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [],
      showIntNodes: true,
      highlightNodes: ["sm-alice"],
    },
  },
  {
    id: 16,
    title: "Owner's IntNode",
    narrative: "When Alice creates a file, she gets her own IntNode. This IntNode contains the file's symmetric key and points to the actual FileStruct.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [],
      showIntNodes: true,
    },
  },
  {
    id: 17,
    title: "Sharing Creates New IntNode",
    narrative: "When Alice (the owner) shares with Bob, she creates a NEW IntNode for Bob. Bob's mailbox points to his IntNode, which points to the file.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }],
      showIntNodes: true,
      highlightNodes: ["sm-bob"],
    },
    note: "Owner sharing = NEW IntNode",
  },
  {
    id: 18,
    title: "Non-Owner Shares Same IntNode",
    narrative: "When Bob shares with Carol, something different happens. Carol joins Bob's EXISTING IntNode—she becomes part of Bob's \"cohort.\"",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }, { from: "Bob", to: "Carol" }],
      showIntNodes: true,
      highlightNodes: ["sm-bob"],
    },
    note: "Non-owner sharing = SAME IntNode (joins cohort)",
  },
  {
    id: 19,
    title: "The Cohort Structure",
    narrative: "Now we have a clear structure: Alice has her IntNode, Bob and Carol share another. Each IntNode represents a \"cohort\" of users who got access through the same path.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }, { from: "Bob", to: "Carol" }],
      showIntNodes: true,
    },
  },
  // Act 4: Revocation Works
  {
    id: 20,
    title: "Revocation with IntNodes",
    narrative: "Now when Alice revokes Bob, she simply DELETES Bob's IntNode. One operation severs access for Bob AND everyone in his cohort.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }, { from: "Bob", to: "Carol" }],
      showIntNodes: true,
      revokedUser: "Bob",
      highlightNodes: ["sm-bob"],
    },
  },
  {
    id: 21,
    title: "Clean Severance",
    narrative: "Bob and Carol both lose access. But Alice's IntNode is untouched—she still has access. And if she had shared with Dave directly, his IntNode would be untouched too.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Dave" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Dave" }],
      showIntNodes: true,
    },
    note: "Dave's cohort is independent of Bob's",
  },
  // Act 5: The Deeper Insight
  {
    id: 22,
    title: "Why This Works",
    narrative: "The magic is in the two-layer indirection. Alice can't track every downstream user, but she CAN track the IntNodes she created. Each IntNode is a revocation handle.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }, { name: "Dave" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }, { from: "Bob", to: "Carol" }, { from: "Alice", to: "Dave" }],
      showIntNodes: true,
    },
  },
  {
    id: 23,
    title: "After Revocation Update",
    narrative: "After deleting Bob's IntNode, Alice re-encrypts the file, moves it to a new location, and updates her remaining IntNodes (her own and Dave's) to point to the new location.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Dave" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Dave" }],
      showIntNodes: true,
      highlightNodes: ["sm-alice", "sm-dave"],
    },
    note: "Alice only needs to update IntNodes she controls",
  },
  {
    id: 24,
    title: "The Principle",
    narrative: "Revoke the person, sever the cohort. It's elegant because it accepts reality: Alice can't control what Bob does after sharing, but she CAN control the access point she gave him.",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }, { name: "Dave" }, { name: "Eve" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [
        { from: "Alice", to: "Bob" },
        { from: "Bob", to: "Carol" },
        { from: "Alice", to: "Dave" },
        { from: "Dave", to: "Eve" },
      ],
      showIntNodes: true,
    },
  },
  {
    id: 25,
    title: "Conclusion",
    narrative: "This is the IntNode pattern: a layer of indirection that transforms an intractable key management problem into a simple graph operation. Try it yourself in the Sandbox!",
    diagram: {
      users: [{ name: "Alice", isOwner: true }, { name: "Bob" }, { name: "Carol" }, { name: "Dave" }],
      files: [{ name: "secret.txt", owner: "Alice" }],
      shares: [{ from: "Alice", to: "Bob" }, { from: "Bob", to: "Carol" }, { from: "Alice", to: "Dave" }],
      showIntNodes: true,
    },
  },
];
