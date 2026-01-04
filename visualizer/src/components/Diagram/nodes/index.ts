import UserNode from './UserNode';
import PersonalMailboxNode from './PersonalMailboxNode';
import SharedMailboxNode from './SharedMailboxNode';
import FileStructNode from './FileStructNode';
import ContentNode from './ContentNode';

export const nodeTypes = {
  user: UserNode,
  personalMailbox: PersonalMailboxNode,
  sharedMailbox: SharedMailboxNode,
  fileStruct: FileStructNode,
  content: ContentNode,
};

export { UserNode, PersonalMailboxNode, SharedMailboxNode, FileStructNode, ContentNode };
