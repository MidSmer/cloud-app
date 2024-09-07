import { createContext, useContext, useReducer, Dispatch, useState } from 'react';
import { Player } from './player';
import { GameStatus } from './game';

export const GomokuDispatchContext = createContext(null as [GomokuState, Dispatch<GomokuAction>] | null);

export const useGomokuDispatch = () => useContext(GomokuDispatchContext);

export type PieceState = {
    rowNo: number,
    colNo: number,
    piece: number,
}

export type GomokuState = PieceState[];

export const TABLE_SIZE = 20

type gomokuResult = 'continue' | 'switch' | 'win' | 'lose';
export type RelatedAction = (result: gomokuResult) => void;

type GomokuAction = {
    type: string,
    payload: {
        rowNo: number,
        colNo: number,
        player: Player,
    }
};

export type GomokuHistory = {
    type: string,
    data: GomokuRecord | GomokuFall,
}[];

type GomokuRecord = PieceState[]

type GomokuFall = {
    player: Player,
    pieceState: PieceState,
}

export const GomokuHistoryContext = createContext(null as [GomokuHistory, Dispatch<GomokuHistory>] | null);

export const useGomokuHistoryContext = () => useContext(GomokuHistoryContext);

export function CheckGoalAchievement(pieceState: PieceState, gomokuState: GomokuState) {
    let searchSamePiece = (piece: number) => {
        let checkBeside = (position: PieceState, action: { row: number, col: number }): number => {
            let row = position.rowNo + action.row;
            let col = position.colNo + action.col;
            if (row < 0 || row >= TABLE_SIZE || col < 0 || col >= TABLE_SIZE) {
                return 0;
            }

            let cell = gomokuState.find((cell) => cell.rowNo === row && cell.colNo === col);
            if (cell && cell.piece === piece) {
                return 1 + checkBeside(cell, action);
            }

            return 0;
        }

        return Math.max(
            checkBeside(pieceState!, { row: 0, col: 1 }) + checkBeside(pieceState!, { row: 0, col: -1 }) + 1,
            checkBeside(pieceState!, { row: 1, col: 0 }) + checkBeside(pieceState!, { row: -1, col: 0 }) + 1,
            checkBeside(pieceState!, { row: 1, col: 1 }) + checkBeside(pieceState!, { row: -1, col: -1 }) + 1,
            checkBeside(pieceState!, { row: 1, col: -1 }) + checkBeside(pieceState!, { row: -1, col: 1 }) + 1
        )
    }

    let count = searchSamePiece(pieceState!.piece);
    console.log('数量：', count)

    return count >= 5;
}

const gomokuReducer = (state: GomokuState, action: GomokuAction) => {
    switch (action.type) {
        case 'fall': {
            const { rowNo, colNo, player } = action.payload;

            const newGomokuStateData = state.map((cell) => {
                if (cell.rowNo === rowNo && cell.colNo === colNo) {
                    return { ...cell, piece: player.piece }
                }
                return cell
            })

            return newGomokuStateData;
        }
        default: {
            throw Error('Unknown action: ' + action.type);
        }
    }
}


export const useGomokuReducer = (initState: GomokuState) => {
    return useReducer(gomokuReducer, initState);
}




