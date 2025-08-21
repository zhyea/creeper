package parser

import (
	"fmt"
	"time"
)

// ParseCommand 解析命令接口
type ParseCommand interface {
	Execute() error
	Undo() error
	GetDescription() string
	GetResult() interface{}
}

// ParseNovelCommand 解析小说命令
type ParseNovelCommand struct {
	parser   *Parser
	path     string
	result   *Novel
	executed bool
}

// NewParseNovelCommand 创建解析小说命令
func NewParseNovelCommand(parser *Parser, path string) *ParseNovelCommand {
	return &ParseNovelCommand{
		parser: parser,
		path:   path,
	}
}

func (pnc *ParseNovelCommand) Execute() error {
	if pnc.executed {
		return fmt.Errorf("命令已执行")
	}
	
	novel, err := pnc.parser.ParseNovel(pnc.path)
	if err != nil {
		return err
	}
	
	pnc.result = novel
	pnc.executed = true
	return nil
}

func (pnc *ParseNovelCommand) Undo() error {
	if !pnc.executed {
		return fmt.Errorf("命令未执行，无法撤销")
	}
	
	pnc.result = nil
	pnc.executed = false
	return nil
}

func (pnc *ParseNovelCommand) GetDescription() string {
	return fmt.Sprintf("解析小说: %s", pnc.path)
}

func (pnc *ParseNovelCommand) GetResult() interface{} {
	return pnc.result
}

// ConvertChapterCommand 转换章节命令
type ConvertChapterCommand struct {
	adapter     *ChapterAdapter
	chapter     *Chapter
	sourceType  string
	originalHTML string
	executed    bool
}

// NewConvertChapterCommand 创建转换章节命令
func NewConvertChapterCommand(adapter *ChapterAdapter, chapter *Chapter, sourceType string) *ConvertChapterCommand {
	return &ConvertChapterCommand{
		adapter:    adapter,
		chapter:    chapter,
		sourceType: sourceType,
	}
}

func (ccc *ConvertChapterCommand) Execute() error {
	if ccc.executed {
		return fmt.Errorf("命令已执行")
	}
	
	// 保存原始 HTML 内容
	ccc.originalHTML = ccc.chapter.HTMLContent
	
	// 执行转换
	err := ccc.adapter.ConvertChapter(ccc.chapter, ccc.sourceType)
	if err != nil {
		return err
	}
	
	ccc.executed = true
	return nil
}

func (ccc *ConvertChapterCommand) Undo() error {
	if !ccc.executed {
		return fmt.Errorf("命令未执行，无法撤销")
	}
	
	// 恢复原始内容
	ccc.chapter.HTMLContent = ccc.originalHTML
	ccc.executed = false
	return nil
}

func (ccc *ConvertChapterCommand) GetDescription() string {
	return fmt.Sprintf("转换章节: %s (%s)", ccc.chapter.Title, ccc.sourceType)
}

func (ccc *ConvertChapterCommand) GetResult() interface{} {
	return ccc.chapter
}

// CommandInvoker 命令调用器
type CommandInvoker struct {
	commands []ParseCommand
	current  int
}

// NewCommandInvoker 创建命令调用器
func NewCommandInvoker() *CommandInvoker {
	return &CommandInvoker{
		commands: make([]ParseCommand, 0),
		current:  -1,
	}
}

// Execute 执行命令
func (ci *CommandInvoker) Execute(command ParseCommand) error {
	// 如果当前位置不在末尾，清除后面的命令
	if ci.current < len(ci.commands)-1 {
		ci.commands = ci.commands[:ci.current+1]
	}
	
	// 执行命令
	if err := command.Execute(); err != nil {
		return err
	}
	
	// 添加到命令历史
	ci.commands = append(ci.commands, command)
	ci.current++
	
	return nil
}

// Undo 撤销命令
func (ci *CommandInvoker) Undo() error {
	if ci.current < 0 {
		return fmt.Errorf("没有可撤销的命令")
	}
	
	command := ci.commands[ci.current]
	if err := command.Undo(); err != nil {
		return err
	}
	
	ci.current--
	return nil
}

// Redo 重做命令
func (ci *CommandInvoker) Redo() error {
	if ci.current >= len(ci.commands)-1 {
		return fmt.Errorf("没有可重做的命令")
	}
	
	ci.current++
	command := ci.commands[ci.current]
	return command.Execute()
}

// CanUndo 是否可以撤销
func (ci *CommandInvoker) CanUndo() bool {
	return ci.current >= 0
}

// CanRedo 是否可以重做
func (ci *CommandInvoker) CanRedo() bool {
	return ci.current < len(ci.commands)-1
}

// GetHistory 获取命令历史
func (ci *CommandInvoker) GetHistory() []string {
	history := make([]string, len(ci.commands))
	for i, command := range ci.commands {
		status := "✓"
		if i > ci.current {
			status = "○"
		}
		history[i] = fmt.Sprintf("%s %s", status, command.GetDescription())
	}
	return history
}

// Clear 清空命令历史
func (ci *CommandInvoker) Clear() {
	ci.commands = make([]ParseCommand, 0)
	ci.current = -1
}

// BatchParseCommand 批量解析命令
type BatchParseCommand struct {
	commands    []ParseCommand
	results     []interface{}
	executed    bool
	failOnError bool
}

// NewBatchParseCommand 创建批量解析命令
func NewBatchParseCommand(failOnError bool) *BatchParseCommand {
	return &BatchParseCommand{
		commands:    make([]ParseCommand, 0),
		results:     make([]interface{}, 0),
		failOnError: failOnError,
	}
}

// AddCommand 添加命令
func (bpc *BatchParseCommand) AddCommand(command ParseCommand) {
	bpc.commands = append(bpc.commands, command)
}

func (bpc *BatchParseCommand) Execute() error {
	if bpc.executed {
		return fmt.Errorf("批量命令已执行")
	}
	
	bpc.results = make([]interface{}, 0)
	
	for i, command := range bpc.commands {
		if err := command.Execute(); err != nil {
			if bpc.failOnError {
				// 撤销已执行的命令
				for j := i - 1; j >= 0; j-- {
					bpc.commands[j].Undo()
				}
				return fmt.Errorf("批量命令执行失败: %w", err)
			}
			// 继续执行，但记录错误
			bpc.results = append(bpc.results, err)
		} else {
			bpc.results = append(bpc.results, command.GetResult())
		}
	}
	
	bpc.executed = true
	return nil
}

func (bpc *BatchParseCommand) Undo() error {
	if !bpc.executed {
		return fmt.Errorf("批量命令未执行，无法撤销")
	}
	
	// 逆序撤销所有命令
	for i := len(bpc.commands) - 1; i >= 0; i-- {
		if err := bpc.commands[i].Undo(); err != nil {
			return fmt.Errorf("撤销命令失败: %w", err)
		}
	}
	
	bpc.executed = false
	bpc.results = make([]interface{}, 0)
	return nil
}

func (bpc *BatchParseCommand) GetDescription() string {
	return fmt.Sprintf("批量解析命令 (%d个子命令)", len(bpc.commands))
}

func (bpc *BatchParseCommand) GetResult() interface{} {
	return bpc.results
}

// TimedParseCommand 计时解析命令
type TimedParseCommand struct {
	command   ParseCommand
	startTime time.Time
	duration  time.Duration
	executed  bool
}

// NewTimedParseCommand 创建计时解析命令
func NewTimedParseCommand(command ParseCommand) *TimedParseCommand {
	return &TimedParseCommand{
		command: command,
	}
}

func (tpc *TimedParseCommand) Execute() error {
	if tpc.executed {
		return fmt.Errorf("计时命令已执行")
	}
	
	tpc.startTime = time.Now()
	err := tpc.command.Execute()
	tpc.duration = time.Since(tpc.startTime)
	tpc.executed = true
	
	return err
}

func (tpc *TimedParseCommand) Undo() error {
	if !tpc.executed {
		return fmt.Errorf("计时命令未执行，无法撤销")
	}
	
	err := tpc.command.Undo()
	if err == nil {
		tpc.executed = false
		tpc.duration = 0
	}
	
	return err
}

func (tpc *TimedParseCommand) GetDescription() string {
	desc := tpc.command.GetDescription()
	if tpc.executed {
		desc += fmt.Sprintf(" (耗时: %v)", tpc.duration)
	}
	return desc
}

func (tpc *TimedParseCommand) GetResult() interface{} {
	return map[string]interface{}{
		"result":   tpc.command.GetResult(),
		"duration": tpc.duration,
	}
}

// GetDuration 获取执行时间
func (tpc *TimedParseCommand) GetDuration() time.Duration {
	return tpc.duration
}
