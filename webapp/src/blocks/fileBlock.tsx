// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import {Block, createBlock} from './block'

type AttachmentBlockFields = {
    attachmentId: string
}

type AttachmentBlock = Block & {
    type: 'attachment'
    fields: AttachmentBlockFields
}

function createAttachmentBlock(block?: Block): AttachmentBlock {
    return {
        ...createBlock(block),
        type: 'attachment',
        fields: {
            attachmentId: block?.fields.attachmentId || '',
        },
    }
}

export {AttachmentBlock, createAttachmentBlock}
