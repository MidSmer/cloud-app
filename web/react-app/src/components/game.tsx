import { useState } from "react"
import { GameStatusContext, useGameStatus, GameStatus } from "@/hooks/game"
import {  PlayerContext, Player, usePlayerReducer } from "@/hooks/player"
import Kanban from "@/components/kanban"
import  Tabletop  from "./tabletop"

export default function Game() {
    const [ gameStatus, setGamestatus ] = useState("ready" as GameStatus);
    let player = { name: "黑棋", piece: 1 } as Player
    let tabletopKey = 'a'

    if (gameStatus === "b-playing") {
        player = { name: "白棋", piece: 2 } as Player
        tabletopKey = 'b'
    }

    const [ playerContext, playerDispatch ] = usePlayerReducer(player)

    return (
        <GameStatusContext.Provider value={{ gameStatus, dispatch: setGamestatus }}>
            <PlayerContext.Provider value={{ player: playerContext, dispatch: playerDispatch }}>
                {(gameStatus !== "a-playing" && gameStatus !== "b-playing") && <GameOverlay />}
                {(gameStatus === "win" || gameStatus === "lose") && <GameOverlay />}
                <div className="w-full h-screen flex justify-center items-center">
                    <div className="relative ">
                        <Tabletop key={tabletopKey}/>
                        <div className="absolute -right-24 top-1/2">
                            <Kanban />
                        </div>
                    </div>
                </div>
            </PlayerContext.Provider>
        </GameStatusContext.Provider>
    )
}

function GameOverlay() {
    const { gameStatus, dispatch: setGamestatus } = useGameStatus();

    const startGame = () => {
        setGamestatus?.("a-playing" as GameStatus)
    }

    const restartGame = () => {
        setGamestatus?.("a-playing" as GameStatus)
    }

    return (
        <>
            <div className="z-50 w-full h-screen absolute flex justify-center items-center bg-white/80">
                {gameStatus === "ready" && (
                    <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                        onClick={startGame}>
                        Start Game
                    </button>
                )}
                {gameStatus === 'win' && (
                    <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                        onClick={restartGame}>
                        Restart Game
                    </button>
                )}
                            
            </div>
        </>
    )
}