import { useContext, useEffect, useState } from "react"
import { PlayerContext } from "@/hooks/player"
import { useGameStatus } from "@/hooks/game"

export default function Kanban() {
    return (
        <div className="w-24">
            <Timer />
            <PlayerInfo />
        </div>
    )
}

function Timer() {
    const [ totalSecond, setTotalSecond ] = useState(0)
    const { gameStatus, dispatch: setGameStatus } = useGameStatus()

    useEffect(() => {
        const id = setInterval(() => {
            if (gameStatus === 'a-playing') {
                setTotalSecond(second => second + 1)
            } else if (gameStatus === 'b-playing') {
                if (totalSecond === 0) {
                    setGameStatus?.('win')
                } else {
                    setTotalSecond(second => second - 1)
                }
            } else if (gameStatus === 'win') {
                setTotalSecond(0)
            } else if (gameStatus === 'lose') {
                setTotalSecond(0)
            }
        }, 1000);
        return () => clearInterval(id);
    }, [gameStatus, setGameStatus, totalSecond]);

    const minute = ('00' + Math.floor(totalSecond / 60)).slice(-2)
    const second = ('00' + (totalSecond % 60)).slice(-2)

    return (
        <div className="flex flex-row justify-center items-center">
            <div className="flex flex-col justify-center items-center">
                <div className="text-4xl">{minute}</div>
            </div>
            <div className="text-4xl">:</div>
            <div className="flex flex-col justify-center items-center">
                <div className="text-4xl">{second}</div>
            </div>
        </div>
    )
}

function PlayerInfo() {
    const playerContext = useContext(PlayerContext)
    const player = playerContext.player

    return (
        <div className="flex flex-col justify-center items-center">
            <div className="text-2xl">{player.name}</div>
        </div>
    )
}