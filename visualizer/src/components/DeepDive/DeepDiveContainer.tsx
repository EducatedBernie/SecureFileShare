'use client';

import { useState } from 'react';

interface Section {
  id: string;
  title: string;
  color: string;
  borderColor: string;
  description: string;
  details: {
    purpose: string;
    fields?: { name: string; type: string; description: string }[];
    operations?: { name: string; description: string }[];
    code?: string;
    notes?: string[];
  };
}

const sections: Section[] = [
  {
    id: 'user',
    title: 'User Struct',
    color: 'from-blue-600 to-blue-800',
    borderColor: 'border-blue-500/50',
    description: 'Entry point containing user identity and cryptographic keys',
    details: {
      purpose: 'The User struct is the root of all user operations. It stores the username and private keys needed to decrypt personal data and sign operations.',
      fields: [
        { name: 'Username', type: 'string', description: 'Unique identifier for the user' },
        { name: 'PrivateDecKey', type: 'PKEDecKey', description: 'RSA private key for decrypting data encrypted with user\'s public key' },
        { name: 'PrivateDigitalSig', type: 'DSSignKey', description: 'RSA private key for signing data to prove authenticity' },
      ],
      code: `type User struct {
    Username          string
    PrivateDecKey     userlib.PrivateKeyType
    PrivateDigitalSig userlib.DSSignKey
}`,
      notes: [
        'Public keys are stored in KeyStore, not in the User struct',
        'User struct is encrypted with Argon2-derived key from password',
        'Stored at UUID = Hash(username)[:16]',
      ],
    },
  },
  {
    id: 'mailbox',
    title: 'Personal Mailbox',
    color: 'from-violet-600 to-violet-800',
    borderColor: 'border-violet-500/50',
    description: 'Per-user encrypted pointer to shared resources',
    details: {
      purpose: 'A personal mailbox links a user to a file. It\'s encrypted with the user\'s public key and signed with their private key, ensuring only they can read it.',
      fields: [
        { name: 'FileStructID', type: 'uuid.UUID', description: 'Points to the IntNode (intermediate node)' },
        { name: 'FileKey', type: '[]byte', description: 'Symmetric key to decrypt the IntNode' },
      ],
      code: `type Mailbox struct {
    FileStructID uuid.UUID
    FileKey      []byte
}`,
      operations: [
        { name: 'Location', description: 'UUID = Hash(filename + "-" + username)[:16]' },
        { name: 'Encryption', description: 'PKE encrypted with user\'s public key' },
        { name: 'Authentication', description: 'Signed with user\'s private signing key' },
      ],
      notes: [
        'Each user has one mailbox per file they can access',
        'Mailbox is re-encrypted when sharing (different user\'s public key)',
        'Signature prevents tampering by malicious server',
      ],
    },
  },
  {
    id: 'intnode',
    title: 'IntNode (Intermediate Node)',
    color: 'from-amber-600 to-amber-800',
    borderColor: 'border-amber-500/50',
    description: 'The key innovation: enables cohort-based revocation',
    details: {
      purpose: 'IntNodes create an indirection layer between mailboxes and files. When an owner shares, they create a NEW IntNode for the recipient. When a non-owner shares, the recipient joins the SAME IntNode. This groups users into "cohorts" that can be revoked together.',
      fields: [
        { name: 'FileStructID', type: 'uuid.UUID', description: 'Points to the HMAC-wrapped FileStruct' },
        { name: 'FileKey', type: '[]byte', description: 'Symmetric key to decrypt the FileStruct' },
      ],
      code: `// IntNode uses the same Mailbox struct
type Mailbox struct {
    FileStructID uuid.UUID  // -> HMACStruct -> FileStruct
    FileKey      []byte     // Key to decrypt FileStruct
}`,
      operations: [
        { name: 'Owner shares', description: 'Creates NEW IntNode for recipient with derived key' },
        { name: 'Non-owner shares', description: 'Shares SAME IntNode (recipient joins cohort)' },
        { name: 'Revocation', description: 'Delete IntNode → entire cohort loses access' },
      ],
      notes: [
        'This is THE key insight of the design',
        'Owner tracks IntNodes they created (one per direct share)',
        'Non-owners cannot create new IntNodes',
        'Stored with HMAC for integrity protection',
      ],
    },
  },
  {
    id: 'hmac',
    title: 'HMAC Struct',
    color: 'from-teal-600 to-teal-800',
    borderColor: 'border-teal-500/50',
    description: 'Integrity verification wrapper for encrypted data',
    details: {
      purpose: 'Wraps encrypted data with an HMAC to detect tampering. Used to protect FileStructs and content from modification by a malicious server.',
      fields: [
        { name: 'Hmac', type: '[]byte', description: 'HMAC-SHA256 of the encrypted data' },
        { name: 'Encryption', type: '[]byte', description: 'The encrypted (ciphertext) payload' },
      ],
      code: `type HMAC struct {
    Hmac       []byte
    Encryption []byte
}

func checkHmacStruct(hmac HMAC, key []byte) bool {
    checkHmac, err := userlib.HMACEval(key, hmac.Encryption)
    if err != nil {
        return false
    }
    return userlib.HMACEqual(checkHmac, hmac.Hmac)
}`,
      notes: [
        'HMAC is computed over ciphertext (encrypt-then-MAC)',
        'Same key used for both encryption and HMAC derivation',
        'Provides integrity but not authenticity (no signatures)',
      ],
    },
  },
  {
    id: 'filestruct',
    title: 'FileStruct',
    color: 'from-emerald-600 to-emerald-800',
    borderColor: 'border-emerald-500/50',
    description: 'File metadata, ownership, and content pointers',
    details: {
      purpose: 'Contains all metadata about a file: who owns it, where the content is, who has been shared access, and pointers for efficient append operations.',
      fields: [
        { name: 'Owner', type: '[]byte', description: 'Hash of owner\'s username' },
        { name: 'UUIDContents', type: 'uuid.UUID', description: 'Location of encrypted file content' },
        { name: 'SharedUsersUUID', type: 'uuid.UUID', description: 'List of directly shared users (hashed)' },
        { name: 'InvitedUsersUUID', type: 'uuid.UUID', description: 'Users invited but not yet accepted' },
        { name: 'NextNode', type: 'uuid.UUID', description: 'For append: pointer to next chunk' },
        { name: 'LastNode', type: 'uuid.UUID', description: 'For append: pointer to last chunk (O(1) append)' },
        { name: 'NumAppends', type: 'int', description: 'Count of appended chunks' },
        { name: 'IsDummy', type: 'bool', description: 'True for head node (sentinel)' },
      ],
      code: `type File struct {
    Owner            []byte
    UUIDContents     uuid.UUID
    SharedUsersUUID  uuid.UUID
    InvitedUsersUUID uuid.UUID
    NextNode         uuid.UUID
    LastNode         uuid.UUID
    NumAppends       int
    IsDummy          bool
}`,
      notes: [
        'Owner is stored as hash to prevent leaking username',
        'SharedUsers list enables owner to update all IntNodes on revocation',
        'Linked list structure for efficient appends',
      ],
    },
  },
  {
    id: 'content',
    title: 'Content Chunks',
    color: 'from-zinc-600 to-zinc-800',
    borderColor: 'border-zinc-500/50',
    description: 'Encrypted file content with integrity protection',
    details: {
      purpose: 'Actual file content, encrypted with the file\'s symmetric key. Stored separately from metadata to allow efficient content updates without touching the FileStruct.',
      operations: [
        { name: 'Storage format', description: 'HMAC(64 bytes) || Ciphertext' },
        { name: 'Encryption', description: 'AES-CTR with random IV prepended' },
        { name: 'Key derivation', description: 'Append chunks use HashKDF(fileKey, "content" + index)' },
      ],
      code: `// Store content
iv := userlib.RandomBytes(16)
encryptedContent := userlib.SymEnc(fileKey, iv, content)
contentHmac, _ := calculateHMAC(fileKey, encryptedContent)
userlib.DatastoreSet(uuid, append(contentHmac, encryptedContent...))

// Load content
func fetchFileContents(uuid uuid.UUID, key []byte) ([]byte, error) {
    userObj, ok := userlib.DatastoreGet(uuid)
    hmac, cipher := userObj[:64], userObj[64:]
    // Verify HMAC, then decrypt
    plaintext := userlib.SymDec(key, cipher)
    return plaintext, nil
}`,
      notes: [
        'Each append creates a new chunk with derived key',
        'Load operation walks the linked list concatenating chunks',
        'StoreFile overwrites and resets the linked list',
      ],
    },
  },
  {
    id: 'invitation',
    title: 'Invitation',
    color: 'from-pink-600 to-pink-800',
    borderColor: 'border-pink-500/50',
    description: 'Secure key transfer mechanism for sharing',
    details: {
      purpose: 'When sharing a file, an Invitation struct securely transfers the IntNode location and decryption key to the recipient. It\'s encrypted with the recipient\'s public key.',
      fields: [
        { name: 'IntNodeUUID', type: 'uuid.UUID', description: 'Location of IntNode (new or shared)' },
        { name: 'Key', type: '[]byte', description: 'Key to decrypt the IntNode' },
      ],
      code: `type Invitation struct {
    IntNodeUUID uuid.UUID
    Key         []byte
}

// Owner sharing: creates NEW IntNode
if yesOwner {
    recipientFileKey := HashKDF(mailbox.FileKey, recipientUsernameHash)
    // ... create new IntNode for recipient
    invite = createInvitationStruct(newIntNodePtr, recipientFileKey)
} else {
    // Non-owner: shares SAME IntNode
    invite = createInvitationStruct(mailbox.FileStructID, mailbox.FileKey)
}`,
      operations: [
        { name: 'Location', description: 'UUID = Hash(filename + "+" + recipient + "+" + sender)[:16]' },
        { name: 'Encryption', description: 'PKE encrypted with recipient\'s public key' },
        { name: 'Authentication', description: 'Signed by sender to prove origin' },
      ],
      notes: [
        'Invitation is deleted after acceptance',
        'If revoked before acceptance, invitation is also deleted',
        'Recipient creates their own Mailbox pointing to the IntNode',
      ],
    },
  },
];

const cryptoPrimitives = [
  { name: 'Symmetric Encryption', algorithm: 'AES-CTR', usage: 'File content, structs' },
  { name: 'Asymmetric Encryption', algorithm: 'RSA-OAEP', usage: 'Mailboxes, invitations' },
  { name: 'Digital Signatures', algorithm: 'RSA-PSS', usage: 'Mailbox authentication' },
  { name: 'Password KDF', algorithm: 'Argon2', usage: 'User struct encryption key' },
  { name: 'Key Derivation', algorithm: 'HKDF-SHA512', usage: 'Per-recipient keys, append keys' },
  { name: 'MAC', algorithm: 'HMAC-SHA256', usage: 'Integrity of all stored data' },
  { name: 'Hashing', algorithm: 'SHA-512', usage: 'UUID generation, user identification' },
];

export default function DeepDiveContainer() {
  const [expandedSection, setExpandedSection] = useState<string | null>(null);

  const toggleSection = (id: string) => {
    setExpandedSection(expandedSection === id ? null : id);
  };

  return (
    <div className="h-full overflow-y-auto bg-[#0a0a0f]">
      <div className="max-w-4xl mx-auto p-6 md:p-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-2xl md:text-3xl font-bold text-zinc-100 mb-2">
            Architecture Deep Dive
          </h1>
          <p className="text-zinc-400">
            Explore each layer of the E2E encrypted file sharing system. Click any component to see implementation details.
          </p>
        </div>

        {/* Data Flow Overview */}
        <div className="mb-8 p-4 bg-zinc-900/50 border border-zinc-800 rounded-lg">
          <h2 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-3">Data Flow</h2>
          <div className="flex flex-wrap items-center gap-2 text-sm">
            <span className="px-2 py-1 bg-blue-900/50 text-blue-300 rounded">User</span>
            <span className="text-zinc-600">→</span>
            <span className="px-2 py-1 bg-violet-900/50 text-violet-300 rounded">Mailbox</span>
            <span className="text-zinc-600">→</span>
            <span className="px-2 py-1 bg-amber-900/50 text-amber-300 rounded">IntNode</span>
            <span className="text-zinc-600">→</span>
            <span className="px-2 py-1 bg-teal-900/50 text-teal-300 rounded">HMAC</span>
            <span className="text-zinc-600">→</span>
            <span className="px-2 py-1 bg-emerald-900/50 text-emerald-300 rounded">FileStruct</span>
            <span className="text-zinc-600">→</span>
            <span className="px-2 py-1 bg-zinc-700 text-zinc-300 rounded">Content</span>
          </div>
        </div>

        {/* Layer Stack */}
        <div className="space-y-3 mb-10">
          {sections.map((section) => (
            <div key={section.id}>
              <button
                onClick={() => toggleSection(section.id)}
                className="w-full text-left group"
              >
                <div
                  className={`
                    bg-gradient-to-r ${section.color} text-white p-4 rounded-lg
                    border ${section.borderColor}
                    transform transition-all duration-200
                    hover:scale-[1.01] hover:shadow-lg
                    flex items-center justify-between
                  `}
                >
                  <div>
                    <div className="font-semibold text-lg">{section.title}</div>
                    <div className="text-sm opacity-80">{section.description}</div>
                  </div>
                  <svg
                    className={`w-5 h-5 transition-transform duration-200 ${
                      expandedSection === section.id ? 'rotate-90' : ''
                    }`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </div>
              </button>

              {/* Expanded Details */}
              {expandedSection === section.id && (
                <div className="mt-2 ml-4 border-l-2 border-zinc-700 pl-4 py-4 space-y-4 animate-in slide-in-from-top-2 duration-200">
                  {/* Purpose */}
                  <div>
                    <h4 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-2">Purpose</h4>
                    <p className="text-zinc-300 leading-relaxed">{section.details.purpose}</p>
                  </div>

                  {/* Fields */}
                  {section.details.fields && (
                    <div>
                      <h4 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-2">Fields</h4>
                      <div className="space-y-2">
                        {section.details.fields.map((field) => (
                          <div key={field.name} className="flex items-start gap-3 text-sm">
                            <code className="px-2 py-0.5 bg-zinc-800 text-indigo-400 rounded font-mono shrink-0">
                              {field.name}
                            </code>
                            <span className="text-zinc-500 shrink-0">({field.type})</span>
                            <span className="text-zinc-400">{field.description}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Operations */}
                  {section.details.operations && (
                    <div>
                      <h4 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-2">Operations</h4>
                      <div className="space-y-2">
                        {section.details.operations.map((op) => (
                          <div key={op.name} className="flex items-start gap-3 text-sm">
                            <span className="text-emerald-400 shrink-0 font-medium">{op.name}:</span>
                            <span className="text-zinc-400">{op.description}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Code */}
                  {section.details.code && (
                    <div>
                      <h4 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-2">Implementation</h4>
                      <pre className="bg-zinc-900 border border-zinc-800 rounded-lg p-4 overflow-x-auto text-sm">
                        <code className="text-zinc-300 font-mono whitespace-pre">{section.details.code}</code>
                      </pre>
                    </div>
                  )}

                  {/* Notes */}
                  {section.details.notes && (
                    <div>
                      <h4 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-2">Key Points</h4>
                      <ul className="space-y-1">
                        {section.details.notes.map((note, i) => (
                          <li key={i} className="flex items-start gap-2 text-sm text-zinc-400">
                            <span className="text-amber-500 mt-1">•</span>
                            <span>{note}</span>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Crypto Primitives */}
        <div className="mb-10">
          <h2 className="text-xl font-bold text-zinc-100 mb-4">Cryptographic Primitives</h2>
          <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-800">
                  <th className="text-left p-3 text-zinc-400 font-medium">Purpose</th>
                  <th className="text-left p-3 text-zinc-400 font-medium">Algorithm</th>
                  <th className="text-left p-3 text-zinc-400 font-medium hidden md:table-cell">Used For</th>
                </tr>
              </thead>
              <tbody>
                {cryptoPrimitives.map((prim, i) => (
                  <tr key={prim.name} className={i % 2 === 0 ? 'bg-zinc-900/30' : ''}>
                    <td className="p-3 text-zinc-300">{prim.name}</td>
                    <td className="p-3">
                      <code className="text-indigo-400 font-mono">{prim.algorithm}</code>
                    </td>
                    <td className="p-3 text-zinc-500 hidden md:table-cell">{prim.usage}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Key Insight Box */}
        <div className="bg-gradient-to-r from-amber-900/20 to-amber-800/10 border border-amber-700/30 rounded-lg p-6">
          <h2 className="text-lg font-bold text-amber-400 mb-3 flex items-center gap-2">
            <span>💡</span> The Key Insight: IntNode Pattern
          </h2>
          <div className="text-zinc-300 space-y-3">
            <p>
              The IntNode creates a layer of indirection that transforms an intractable problem into a simple graph operation:
            </p>
            <ul className="space-y-2 ml-4">
              <li className="flex items-start gap-2">
                <span className="text-emerald-400">✓</span>
                <span><strong>Owner sharing</strong> creates a NEW IntNode → recipient gets their own revocation handle</span>
              </li>
              <li className="flex items-start gap-2">
                <span className="text-emerald-400">✓</span>
                <span><strong>Non-owner sharing</strong> shares the SAME IntNode → recipient joins the sharer&apos;s cohort</span>
              </li>
              <li className="flex items-start gap-2">
                <span className="text-emerald-400">✓</span>
                <span><strong>Revocation</strong> deletes the IntNode → entire cohort loses access in one operation</span>
              </li>
            </ul>
            <p className="text-zinc-400 text-sm mt-4">
              The owner doesn&apos;t need to track every downstream user—only the IntNodes they created. Each IntNode is a single point of revocation for an entire sharing subtree.
            </p>
          </div>
        </div>

        {/* Footer */}
        <div className="mt-10 pt-6 border-t border-zinc-800 text-center text-zinc-500 text-sm">
          <p>Based on the CS 161 secure file sharing project design</p>
        </div>
      </div>
    </div>
  );
}
