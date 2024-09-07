import { createContext, useContext, useReducer, Dispatch, useState } from 'react';

// a-playing: 黑子先下 b-playing: 白子先下，时间倒计
export type GameStatus = 'ready' | 'a-playing' | 'b-playing' | 'win' | 'lose' ;

export const GameStatusContext = createContext({} as {gameStatus: GameStatus, dispatch: Dispatch<GameStatus> | null});

export const useGameStatus = () => useContext(GameStatusContext);
