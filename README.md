# End-to-End Encrypted File Sharing System

A secure, Dropbox-like file sharing system designed to protect user data against both a malicious storage server and revoked users. Built as part of UC Berkeley's CS 161 (Computer Security) curriculum.

**Authors:** Bernie Miao & Lyna Jiang  
**Course:** CS 161 – Computer Security, Summer 2022  
**Language:** Go

---

## Blog Post

For a detailed walkthrough of the design evolution, including the failures that led to the final architecture, see: [I Built a Dropbox Clone That Even the Server Can't Read](#) *([link to your blog post](https://medium.com/@berniem4483/i-built-a-dropbox-clone-that-even-the-server-cant-read-8599f15647c6))*


## Overview

This project implements a client-side application for secure file storage and sharing. The system guarantees confidentiality and integrity of user data even when the storage server is completely untrusted—it can read, modify, or delete any stored data.

### Threat Model

The system defends against two adversaries:

| Adversary | Capabilities | Limitations |
|-----------|--------------|-------------|
| **Datastore Adversary** | Read/modify/delete any stored data, observe all API calls, knows source code | Cannot perform rollback attacks |
| **Revoked User Adversary** | Remembers all UUIDs and keys accessed before revocation, can directly call Datastore APIs | Cannot enumerate unknown UUIDs |

### Security Guarantees

- **Confidentiality**: File contents and filenames are hidden from the server (IND-CPA secure)
- **Integrity**: All tampering is detected; operations fail safely on corrupted data
- **Access Control**: Only the owner and users with valid, unrevoked invitations can access files
- **Forward Secrecy**: Revoked users cannot access future file updates

---

## Architecture

The system uses a layered architecture where each layer provides a specific security property:

```
User Struct → Personal Mailbox → Shared Mailbox (IntNode) → FileStruct → Content Chunks
                    ↑                     ↑                       ↑
              (per-user,            (per-cohort,            (integrity
               for lookup)          for access control)      wrapper)
```

### Key Design Decisions

**Two-Layer Mailbox Structure**  
Personal mailboxes (deterministic UUIDs) handle file lookup. Shared mailboxes (random UUIDs, owner-controlled) handle access control. This separation enables O(direct shares) revocation instead of O(total users).

**Linked-List File Storage**  
File content is stored as a linked list of encrypted chunks. Appending data only requires uploading the new chunk and updating pointers—not re-encrypting the entire file. This satisfies the efficiency requirement that appending to a 10TB file uses bandwidth proportional to the append size, not the file size.

**Cohort-Based Sharing**  
When Alice shares with Bob, she creates a shared mailbox for Bob's "cohort." If Bob shares with Carol, Carol points to the same shared mailbox. Revoking Bob deletes one shared mailbox and severs access for his entire downstream tree.

### Struct Summary

| Struct | Purpose | UUID | Encryption |
|--------|---------|------|------------|
| User | Store user credentials | `Hash(username)` | Argon2-derived key |
| Personal Mailbox | Per-user file namespace | `Hash(filename—recipient+sender)` | Owner's public key |
| Shared Mailbox (IntNode) | Cohort access point | Random | Random symmetric key |
| FileStruct | File metadata & pointers | Random | Derived from IntNode key |
| Content | Actual file data | Random | Derived key per chunk |
| Invitation | Pending share | `Hash(filename+recipient+sender)` | Recipient's public key |

---

## Cryptographic Primitives

- **Symmetric Encryption**: AES-128-CTR
- **Asymmetric Encryption**: RSA-2048 OAEP (hybrid encryption for large data)
- **Digital Signatures**: RSA-2048 PSS
- **Key Derivation**: Argon2 (password-based), HashKDF (key derivation)
- **Integrity**: HMAC-SHA512 (encrypt-then-MAC pattern)
- **Hashing**: SHA-512

---

## API

```go
// User Authentication
func InitUser(username string, password string) (*User, error)
func GetUser(username string, password string) (*User, error)

// File Operations
func (user *User) StoreFile(filename string, content []byte) error
func (user *User) LoadFile(filename string) ([]byte, error)
func (user *User) AppendToFile(filename string, content []byte) error

// Sharing
func (user *User) CreateInvitation(filename string, recipientUsername string) (UUID, error)
func (user *User) AcceptInvitation(senderUsername string, invitationPtr UUID, filename string) error
func (user *User) RevokeAccess(filename string, recipientUsername string) error
```

---

## Design Documents

- [`CS_161_Final_Design_Document.pdf`](./CS_161_Final_Design_Document.pdf) – Complete system design with struct definitions, key derivations, and method descriptions

---

## Key Challenges Solved

### The Revocation Problem
In early designs, revoking a user required updating every downstream user's access—but the owner doesn't know who downstream users are. 

**Solution**: Separate personal lookup (deterministic, per-user) from shared access (random UUID, owner-controlled). Non-owners share their existing shared mailbox rather than creating new ones. This makes revocation O(direct shares) instead of O(total users).

### Append Efficiency
Naively, appending to a file would require re-encrypting and re-uploading the entire file.

**Solution**: Store content as a linked list of independently encrypted chunks. Each chunk uses a derived key based on the append counter. Appending only touches new data and pointer updates.

### Integrity Without Trusted Storage
The server can modify any stored data, so every retrieval must verify integrity.

**Solution**: Encrypt-then-MAC pattern throughout. HMAC verification happens before decryption to prevent chosen-ciphertext attacks. Dedicated HMAC struct wrapper avoids concatenation ambiguity.

---

## Academic Integrity Notice

This repository contains design documentation only. Per UC Berkeley's academic integrity policy, source code is not included. If you're a current CS 161 student, please develop your own solution—the learning comes from the struggle.

---

## Acknowledgments

- **Lyna Jiang** – Project partner and co-designer
- **UC Berkeley CS 161 Course Staff** – For designing a project that teaches security through pain.
