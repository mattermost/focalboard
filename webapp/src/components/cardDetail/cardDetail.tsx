// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React, {useState, useRef, useEffect} from 'react'
import {FormattedMessage} from 'react-intl'

import {BlockIcons} from '../../blockIcons'
import mutator from '../../mutator'
import {BoardTree} from '../../viewModel/boardTree'
import {CardTree} from '../../viewModel/cardTree'
import Button from '../../widgets/buttons/button'
import {Focusable} from '../../widgets/editable'
import EditableArea from '../../widgets/editableArea'
import EmojiIcon from '../../widgets/icons/emoji'

import BlockIconSelector from '../blockIconSelector'

import CommentsList from './commentsList'
import CardDetailContents from './cardDetailContents'
import CardDetailContentsMenu from './cardDetailContentsMenu'
import CardDetailProperties from './cardDetailProperties'

import './cardDetail.scss'

type Props = {
    boardTree: BoardTree
    cardTree: CardTree
    readonly: boolean
}

const CardDetail = (props: Props): JSX.Element|null => {
    const {cardTree} = props
    const {card, comments} = cardTree
    const [title, setTitle] = useState(cardTree.card.title)
    const titleRef = useRef<Focusable>(null)
    const titleValueRef = useRef(title)
    titleValueRef.current = title

    useEffect(() => {
        if (!title) {
            titleRef.current?.focus()
        }
    }, [])

    useEffect(() => {
        return () => {
            if (titleValueRef.current !== cardTree?.card.title) {
                mutator.changeTitle(card, titleValueRef.current)
            }
        }
    }, [cardTree])

    if (!cardTree) {
        return null
    }

    return (
        <>
            <div className='CardDetail content'>
                <BlockIconSelector
                    block={card}
                    size='l'
                    readonly={props.readonly}
                />
                {!props.readonly && !card.icon &&
                    <div className='add-buttons'>
                        <Button
                            onClick={() => {
                                const newIcon = BlockIcons.shared.randomIcon()
                                mutator.changeIcon(card, newIcon)
                            }}
                            icon={<EmojiIcon/>}
                        >
                            <FormattedMessage
                                id='CardDetail.add-icon'
                                defaultMessage='Add icon'
                            />
                        </Button>
                    </div>}

                <EditableArea
                    ref={titleRef}
                    className='title'
                    value={title}
                    placeholderText='Untitled'
                    onChange={(newTitle: string) => setTitle(newTitle)}
                    saveOnEsc={true}
                    onSave={() => {
                        if (title !== props.cardTree.card.title) {
                            mutator.changeTitle(card, title)
                        }
                    }}
                    onCancel={() => setTitle(props.cardTree.card.title)}
                    readonly={props.readonly}
                    spellCheck={true}
                />

                {/* Property list */}

                <CardDetailProperties
                    boardTree={props.boardTree}
                    cardTree={props.cardTree}
                    readonly={props.readonly}
                />

                {/* Comments */}

                <hr/>
                <CommentsList
                    comments={comments}
                    rootId={card.rootId}
                    cardId={card.id}
                    readonly={props.readonly}
                />
            </div>

            {/* Content blocks */}

            <div className='CardDetail content fullwidth'>
                <CardDetailContents
                    cardTree={props.cardTree}
                    readonly={props.readonly}
                />
            </div>

            {!props.readonly &&
                <CardDetailContentsMenu card={props.cardTree.card}/>
            }
        </>
    )
}

export default CardDetail
