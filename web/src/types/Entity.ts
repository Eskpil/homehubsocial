import {State} from "./State";

export interface Entity {
    id: string
    name: string
    history: State[]
}