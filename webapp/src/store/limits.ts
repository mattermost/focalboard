// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {createSlice, PayloadAction} from '@reduxjs/toolkit'

import {BoardsCloudLimits} from '../boardsCloudLimits'

import {initialLoad} from './initialLoad'

import {RootState} from './index'

type LimitsState = {
    limits: BoardsCloudLimits,
    hasLimits: boolean
}

const defaultLimits = {
    cards: 0,
    used_cards: 0,
    card_limit_timestamp: 0,
    views: 0,
}

const initialState = {
    limits: defaultLimits,
    hasLimits: false,
} as LimitsState

const limitsSlice = createSlice({
    name: 'limits',
    initialState,
    reducers: {
        setLimits: (state, action: PayloadAction<BoardsCloudLimits>) => {
            state.limits = action.payload
        },
        setCardLimitTimestamp: (state, action: PayloadAction<number>) => {
            state.limits.card_limit_timestamp = action.payload
        },
    },
    extraReducers: (builder) => {
        builder.addCase(initialLoad.fulfilled, (state, action) => {
            state.hasLimits = Boolean(action.payload.limits)
            state.limits = action.payload.limits || defaultLimits
        })
    },
})

export const {reducer} = limitsSlice
export const {setCardLimitTimestamp} = limitsSlice.actions

export const getLimits = (state: RootState): BoardsCloudLimits | undefined => state.limits.limits
export const getGCardLimitTimestamp = (state: RootState): number => state.limits.limits.card_limit_timestamp
