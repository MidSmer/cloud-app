import { useState } from "react"
import {
    GomokuDispatchContext, useGomokuDispatch, useGomokuReducer, GomokuState, 
    GomokuHistoryContext, GomokuHistory,
    useGomokuHistoryContext, TABLE_SIZE,
    CheckGoalAchievement
} from "@/hooks/gomoku"
import { usePlayerContext } from "@/hooks/player"
import { useGameStatus } from '@/hooks/game'


export default function Tabletop() {
    return (
        <GomokuTable />
    )
}

function GomokuTable() {
    const initGomokuData: GomokuState = []
    Array.from(Array(TABLE_SIZE).keys()).forEach((rowNo: number) => {
        Array.from(Array(TABLE_SIZE).keys()).forEach((colNo) => {
            initGomokuData.push({
                rowNo,
                colNo,
                piece: 0,
            })
        })
    })

    const [gomokuHistory, setGomokuHistory] = useState([{
        'type': 'record',
        'data': initGomokuData
    }] as GomokuHistory)
    const [gomokuState, dispatch] = useGomokuReducer(initGomokuData)

    const groupedCells = gomokuState.reduce((acc: { [no: number]: JSX.Element[] }, { rowNo, colNo, piece }) => {
        if (!acc[rowNo]) {
            acc[rowNo] = [];
        }
        acc[rowNo].push(<CrossCell key={`${rowNo}-${colNo}`} rowNo={rowNo} colNo={colNo} piece={piece} />);
        return acc;
    }, {});

    return (
        <GomokuDispatchContext.Provider value={[gomokuState, dispatch]}>
            <GomokuHistoryContext.Provider value={[gomokuHistory, setGomokuHistory]}>
                <div className="flex flex-col bg-orange-200 w-min">
                    {Object.values(groupedCells).map((cells, index) => (
                        <div key={index} className="flex flex-row">
                            {cells}
                        </div>
                    ))}
                </div>
            </GomokuHistoryContext.Provider>
        </GomokuDispatchContext.Provider>
    )
}

// 十字格单元
function CrossCell({ rowNo, colNo, piece }: {
    rowNo: number,
    colNo: number,
    piece: number,
}) {
    const { player, dispatch: playerDispatch } = usePlayerContext()
    const gomokuHistoryContext = useGomokuHistoryContext()
    const gomokuStateContext = useGomokuDispatch()
    const { gameStatus, dispatch: setGameStatus } = useGameStatus()

    if (!gomokuHistoryContext || !gomokuStateContext) {
        return <></>
    }

    const [gomokuHistory, setGomokuHistory] = gomokuHistoryContext
    const [gomokuState, gomokuDispatch] = gomokuStateContext

    const onClick = () => {
        if (piece !== 0) {
            return
        }

        setGomokuHistory([...gomokuHistory, {
            type: 'fall',
            data: {
                player,
                pieceState: {
                    rowNo,
                    colNo,
                    piece: player.piece,
                }
            }
        }])

        gomokuDispatch?.({
            type: 'fall',
            payload: {
                rowNo,
                colNo,
                player,
            }
        })

        if (CheckGoalAchievement({rowNo, colNo, piece: player.piece}, gomokuState)) {
            if (gameStatus === 'a-playing') {
                setGameStatus?.('b-playing')
            } else if (gameStatus === 'b-playing') {
                setGameStatus?.('win')
            }
        } else {
            playerDispatch?.({
                type: 'switch',
                payload: {
                    player,
                }
            })
        }

    }

    let pieceShow = "invisible"
    if (piece == 1) {
        pieceShow = "fill-black"
    } else if (piece == 2) {
        pieceShow = "fill-white"
    }

    return (
        <div className="relative w-10 h-10">
            <div className="absolute top-0 bottom-0 left-1/2 border-x border-solid border-black">
            </div>

            <div className="absolute top-1/2 left-0 right-0 border-y border-solid border-black">
            </div>

            <div className="absolute flex justify-center items-center w-full h-full">
                <div className="absolute  w-3/4 h-3/4 cursor-pointer" onClick={onClick}>
                    <div className={pieceShow}>
                        <svg viewBox="0 0 100 100" preserveAspectRatio="none">
                            <circle cx="50" cy="50" r="50" />
                        </svg>
                    </div>
                </div>
            </div>
        </div>
    )
}