package parser

import (
	"strings"

	"kv_db/internal/database/compute"
)

type state interface {
	appendLetter(letter byte)
	skipLetter()
}

type StateMachine struct {
	currentState state
	tokens       []string
	sb           strings.Builder

	initialState    state
	wordState       state
	whiteSpaceState state
}

func NewStateMachine() *StateMachine {
	sm := &StateMachine{}
	sm.initialState = &initialState{sm}
	sm.wordState = &wordState{sm}
	sm.whiteSpaceState = &whiteSpaceState{sm}
	sm.currentState = sm.initialState
	return sm
}

func (sm *StateMachine) Parse(query string) ([]string, error) {
	for i := 0; i < len(query); i++ {
		symbol := query[i]
		switch {
		case isWhiteSpace(symbol):
			sm.currentState.skipLetter()
		case isLetter(symbol):
			sm.currentState.appendLetter(symbol)
		default:
			return nil, compute.ErrInvalidSymbol
		}
	}

	sm.currentState.skipLetter()
	return sm.tokens, nil
}

func (sm *StateMachine) setState(st state) {
	sm.currentState = st
}

func (sm *StateMachine) writeLetter(letter byte) {
	sm.sb.WriteByte(letter)
}

func (sm *StateMachine) saveToken() {
	sm.tokens = append(sm.tokens, sm.sb.String())
	sm.sb.Reset()
}

func (sm *StateMachine) getWordState() state {
	return sm.wordState
}

func (sm *StateMachine) getWhiteSpaceState() state {
	return sm.whiteSpaceState
}

type initialState struct {
	sm *StateMachine
}

func (ins *initialState) appendLetter(letter byte) {
	ins.sm.writeLetter(letter)
	ins.sm.setState(ins.sm.getWordState())
}

func (ins *initialState) skipLetter() {
	ins.sm.setState(ins.sm.getWhiteSpaceState())
}

type wordState struct {
	sm *StateMachine
}

func (ws *wordState) appendLetter(letter byte) {
	ws.sm.writeLetter(letter)
}

func (ws *wordState) skipLetter() {
	ws.sm.saveToken()
	ws.sm.setState(ws.sm.getWhiteSpaceState())
}

type whiteSpaceState struct {
	sm *StateMachine
}

func (ws *whiteSpaceState) appendLetter(letter byte) {
	ws.sm.writeLetter(letter)
	ws.sm.setState(ws.sm.getWordState())
}

func (ws *whiteSpaceState) skipLetter() {
}

func isWhiteSpace(symbol byte) bool {
	return symbol == '\t' || symbol == '\n' || symbol == ' '
}

func isLetter(symbol byte) bool {
	return (symbol >= 'a' && symbol <= 'z') ||
		(symbol >= 'A' && symbol <= 'Z') ||
		(symbol >= '0' && symbol <= '9') ||
		(symbol == '_')
}
