// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

interface IUser {
    id: string,
    username: string,
    email: string,
    nickname: string,
    firstname: string,
    lastname: string,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    props: Record<string, any>,
    create_at: number,
    update_at: number,
    is_bot: boolean,
    roles: string,
}

interface UserWorkspace {
    id: string
    title: string
    boardCount: number
}

interface UserConfigPatch {
    updatedFields?: Record<string, string>
    deletedFields?: string[]
}

function parseUserProps(props: Array<UserPreference>): Record<string, UserPreference> {
    // const processedProps = props
    // const hiddenBoardIDs = props.hiddenBoardIDs ? JSON.parse(props.hiddenBoardIDs) : []
    // processedProps.hiddenBoardIDs = {}
    // hiddenBoardIDs.forEach((boardID: string) => processedProps.hiddenBoardIDs[boardID] = true)
    // return processedProps

    const processedProps: Record<string, UserPreference> = {}

    props.forEach((prop) => {
        const processedProp = prop
        if (prop.name === 'hiddenBoardIDs') {
            const hiddenBoardIDs = JSON.parse(processedProp.value)
            processedProp.value = {}
            hiddenBoardIDs.forEach((boardID: string) => processedProp.value[boardID] = true)
        }
        processedProps[processedProp.name] = processedProp
    })

    return processedProps
}

const UserPropPrefix = 'focalboard_'

interface UserPreference {
    user_id: string
    category: string
    name: string
    value: any
}

export {
    IUser,
    UserWorkspace,
    UserConfigPatch,
    UserPropPrefix,
    parseUserProps,
    UserPreference,
}
