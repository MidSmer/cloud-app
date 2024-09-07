import { createContext, useContext, useReducer, Dispatch } from 'react'

export type Player = {
    name: string,
    piece: 1 | 2,
}

type PlayerAction = {
    type: string,
    payload: {
        player: Player,
    }
};

export const PlayerContext = createContext({} as { player: Player, dispatch: Dispatch<PlayerAction> | null });
export const usePlayerContext = () => useContext(PlayerContext);

const playerReducer = (state: Player, action: PlayerAction) => {
    switch (action.type) {
        case 'switch': {
            const { player } = action.payload;

            if (player.piece === 1) {
                return { name: "白棋", piece: 2 } as Player
            } else {
                return { name: "黑棋", piece: 1 } as Player
            }
        }
        default: {
            throw Error('Unknown action: ' + action.type);
        }
    }
}

export const usePlayerReducer = (initState: Player) => {
    return useReducer(playerReducer, initState);
} 