package config

import (
	"fmt"
	"strconv"
	"strings"

	"creeper/internal/common"
)

// Expression 表达式接口
type Expression interface {
	Interpret(context *ConfigContext) (interface{}, error)
}

// ConfigContext 配置上下文
type ConfigContext struct {
	variables map[string]interface{}
	logger    *common.Logger
}

// NewConfigContext 创建配置上下文
func NewConfigContext() *ConfigContext {
	return &ConfigContext{
		variables: make(map[string]interface{}),
		logger:    common.GetLogger(),
	}
}

// SetVariable 设置变量
func (cc *ConfigContext) SetVariable(name string, value interface{}) {
	cc.variables[name] = value
}

// GetVariable 获取变量
func (cc *ConfigContext) GetVariable(name string) (interface{}, bool) {
	value, exists := cc.variables[name]
	return value, exists
}

// LiteralExpression 字面量表达式
type LiteralExpression struct {
	value interface{}
}

// NewLiteralExpression 创建字面量表达式
func NewLiteralExpression(value interface{}) *LiteralExpression {
	return &LiteralExpression{value: value}
}

func (le *LiteralExpression) Interpret(context *ConfigContext) (interface{}, error) {
	return le.value, nil
}

// VariableExpression 变量表达式
type VariableExpression struct {
	name string
}

// NewVariableExpression 创建变量表达式
func NewVariableExpression(name string) *VariableExpression {
	return &VariableExpression{name: name}
}

func (ve *VariableExpression) Interpret(context *ConfigContext) (interface{}, error) {
	value, exists := context.GetVariable(ve.name)
	if !exists {
		return nil, fmt.Errorf("变量未定义: %s", ve.name)
	}
	return value, nil
}

// BinaryExpression 二元表达式
type BinaryExpression struct {
	left     Expression
	operator string
	right    Expression
}

// NewBinaryExpression 创建二元表达式
func NewBinaryExpression(left Expression, operator string, right Expression) *BinaryExpression {
	return &BinaryExpression{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (be *BinaryExpression) Interpret(context *ConfigContext) (interface{}, error) {
	leftValue, err := be.left.Interpret(context)
	if err != nil {
		return nil, err
	}

	rightValue, err := be.right.Interpret(context)
	if err != nil {
		return nil, err
	}

	switch be.operator {
	case "+":
		return be.add(leftValue, rightValue)
	case "-":
		return be.subtract(leftValue, rightValue)
	case "*":
		return be.multiply(leftValue, rightValue)
	case "/":
		return be.divide(leftValue, rightValue)
	case "==":
		return be.equal(leftValue, rightValue)
	case "!=":
		return be.notEqual(leftValue, rightValue)
	case "&&":
		return be.and(leftValue, rightValue)
	case "||":
		return be.or(leftValue, rightValue)
	default:
		return nil, fmt.Errorf("不支持的运算符: %s", be.operator)
	}
}

// add 加法运算
func (be *BinaryExpression) add(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		if r, ok := right.(int); ok {
			return l + r, nil
		}
	case float64:
		if r, ok := right.(float64); ok {
			return l + r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l + r, nil
		}
	}
	return nil, fmt.Errorf("不支持的加法运算类型")
}

// subtract 减法运算
func (be *BinaryExpression) subtract(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		if r, ok := right.(int); ok {
			return l - r, nil
		}
	case float64:
		if r, ok := right.(float64); ok {
			return l - r, nil
		}
	}
	return nil, fmt.Errorf("不支持的减法运算类型")
}

// multiply 乘法运算
func (be *BinaryExpression) multiply(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		if r, ok := right.(int); ok {
			return l * r, nil
		}
	case float64:
		if r, ok := right.(float64); ok {
			return l * r, nil
		}
	}
	return nil, fmt.Errorf("不支持的乘法运算类型")
}

// divide 除法运算
func (be *BinaryExpression) divide(left, right interface{}) (interface{}, error) {
	switch l := left.(type) {
	case int:
		if r, ok := right.(int); ok {
			if r == 0 {
				return nil, fmt.Errorf("除数不能为零")
			}
			return l / r, nil
		}
	case float64:
		if r, ok := right.(float64); ok {
			if r == 0 {
				return nil, fmt.Errorf("除数不能为零")
			}
			return l / r, nil
		}
	}
	return nil, fmt.Errorf("不支持的除法运算类型")
}

// equal 相等比较
func (be *BinaryExpression) equal(left, right interface{}) (interface{}, error) {
	return left == right, nil
}

// notEqual 不等比较
func (be *BinaryExpression) notEqual(left, right interface{}) (interface{}, error) {
	return left != right, nil
}

// and 逻辑与
func (be *BinaryExpression) and(left, right interface{}) (interface{}, error) {
	if l, ok := left.(bool); ok {
		if r, ok := right.(bool); ok {
			return l && r, nil
		}
	}
	return nil, fmt.Errorf("逻辑与运算需要布尔类型")
}

// or 逻辑或
func (be *BinaryExpression) or(left, right interface{}) (interface{}, error) {
	if l, ok := left.(bool); ok {
		if r, ok := right.(bool); ok {
			return l || r, nil
		}
	}
	return nil, fmt.Errorf("逻辑或运算需要布尔类型")
}

// FunctionExpression 函数表达式
type FunctionExpression struct {
	name      string
	arguments []Expression
}

// NewFunctionExpression 创建函数表达式
func NewFunctionExpression(name string, arguments []Expression) *FunctionExpression {
	return &FunctionExpression{
		name:      name,
		arguments: arguments,
	}
}

func (fe *FunctionExpression) Interpret(context *ConfigContext) (interface{}, error) {
	// 解析参数
	args := make([]interface{}, len(fe.arguments))
	for i, arg := range fe.arguments {
		value, err := arg.Interpret(context)
		if err != nil {
			return nil, err
		}
		args[i] = value
	}

	// 执行函数
	switch fe.name {
	case "len":
		return fe.len(args)
	case "concat":
		return fe.concat(args)
	case "toInt":
		return fe.toInt(args)
	case "toString":
		return fe.toString(args)
	case "toBool":
		return fe.toBool(args)
	default:
		return nil, fmt.Errorf("未知函数: %s", fe.name)
	}
}

// len 获取长度
func (fe *FunctionExpression) len(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len函数需要一个参数")
	}

	switch v := args[0].(type) {
	case string:
		return len(v), nil
	case []interface{}:
		return len(v), nil
	default:
		return nil, fmt.Errorf("len函数不支持的类型")
	}
}

// concat 字符串连接
func (fe *FunctionExpression) concat(args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("concat函数至少需要两个参数")
	}

	var result strings.Builder
	for _, arg := range args {
		if str, ok := arg.(string); ok {
			result.WriteString(str)
		} else {
			return nil, fmt.Errorf("concat函数参数必须是字符串")
		}
	}

	return result.String(), nil
}

// toInt 转换为整数
func (fe *FunctionExpression) toInt(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("toInt函数需要一个参数")
	}

	switch v := args[0].(type) {
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return nil, fmt.Errorf("toInt函数不支持的类型")
	}
}

// toString 转换为字符串
func (fe *FunctionExpression) toString(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("toString函数需要一个参数")
	}

	return fmt.Sprintf("%v", args[0]), nil
}

// toBool 转换为布尔值
func (fe *FunctionExpression) toBool(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("toBool函数需要一个参数")
	}

	switch v := args[0].(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case int:
		return v != 0, nil
	default:
		return nil, fmt.Errorf("toBool函数不支持的类型")
	}
}

// ConfigInterpreter 配置解释器
type ConfigInterpreter struct {
	parser  *ConfigParser
	logger  *common.Logger
}

// NewConfigInterpreter 创建配置解释器
func NewConfigInterpreter() *ConfigInterpreter {
	return &ConfigInterpreter{
		parser: NewConfigParser(),
		logger: common.GetLogger(),
	}
}

// Interpret 解释配置表达式
func (ci *ConfigInterpreter) Interpret(expression string, context *ConfigContext) (interface{}, error) {
	// 解析表达式
	expr, err := ci.parser.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("解析表达式失败: %w", err)
	}

	// 解释表达式
	return expr.Interpret(context)
}

// ConfigParser 配置解析器
type ConfigParser struct {
	logger *common.Logger
}

// NewConfigParser 创建配置解析器
func NewConfigParser() *ConfigParser {
	return &ConfigParser{
		logger: common.GetLogger(),
	}
}

// Parse 解析表达式
func (cp *ConfigParser) Parse(expression string) (Expression, error) {
	// 简化解析器实现
	// 实际项目中可以使用更复杂的解析器
	
	expression = strings.TrimSpace(expression)
	
	// 检查是否是字面量
	if cp.isLiteral(expression) {
		return cp.parseLiteral(expression)
	}
	
	// 检查是否是变量
	if cp.isVariable(expression) {
		return cp.parseVariable(expression)
	}
	
	// 检查是否是函数调用
	if cp.isFunctionCall(expression) {
		return cp.parseFunctionCall(expression)
	}
	
	// 检查是否是二元表达式
	if cp.isBinaryExpression(expression) {
		return cp.parseBinaryExpression(expression)
	}
	
	return nil, fmt.Errorf("无法解析表达式: %s", expression)
}

// isLiteral 检查是否是字面量
func (cp *ConfigParser) isLiteral(expression string) bool {
	// 检查数字
	if _, err := strconv.Atoi(expression); err == nil {
		return true
	}
	
	// 检查浮点数
	if _, err := strconv.ParseFloat(expression, 64); err == nil {
		return true
	}
	
	// 检查布尔值
	if expression == "true" || expression == "false" {
		return true
	}
	
	// 检查字符串（带引号）
	if (strings.HasPrefix(expression, `"`) && strings.HasSuffix(expression, `"`)) ||
		(strings.HasPrefix(expression, `'`) && strings.HasSuffix(expression, `'`)) {
		return true
	}
	
	return false
}

// parseLiteral 解析字面量
func (cp *ConfigParser) parseLiteral(expression string) (Expression, error) {
	// 解析数字
	if num, err := strconv.Atoi(expression); err == nil {
		return NewLiteralExpression(num), nil
	}
	
	// 解析浮点数
	if num, err := strconv.ParseFloat(expression, 64); err == nil {
		return NewLiteralExpression(num), nil
	}
	
	// 解析布尔值
	if expression == "true" {
		return NewLiteralExpression(true), nil
	}
	if expression == "false" {
		return NewLiteralExpression(false), nil
	}
	
	// 解析字符串
	if strings.HasPrefix(expression, `"`) && strings.HasSuffix(expression, `"`) {
		return NewLiteralExpression(strings.Trim(expression, `"`)), nil
	}
	if strings.HasPrefix(expression, `'`) && strings.HasSuffix(expression, `'`) {
		return NewLiteralExpression(strings.Trim(expression, `'`)), nil
	}
	
	return nil, fmt.Errorf("无法解析字面量: %s", expression)
}

// isVariable 检查是否是变量
func (cp *ConfigParser) isVariable(expression string) bool {
	// 变量以字母或下划线开头，包含字母、数字、下划线
	if len(expression) == 0 {
		return false
	}
	
	first := expression[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_') {
		return false
	}
	
	for _, char := range expression[1:] {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || 
			(char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	
	return true
}

// parseVariable 解析变量
func (cp *ConfigParser) parseVariable(expression string) (Expression, error) {
	return NewVariableExpression(expression), nil
}

// isFunctionCall 检查是否是函数调用
func (cp *ConfigParser) isFunctionCall(expression string) bool {
	return strings.Contains(expression, "(") && strings.Contains(expression, ")")
}

// parseFunctionCall 解析函数调用
func (cp *ConfigParser) parseFunctionCall(expression string) (Expression, error) {
	// 简化实现
	// 实际项目中需要更复杂的解析逻辑
	
	// 提取函数名
	openParen := strings.Index(expression, "(")
	if openParen == -1 {
		return nil, fmt.Errorf("无效的函数调用: %s", expression)
	}
	
	funcName := strings.TrimSpace(expression[:openParen])
	
	// 提取参数
	closeParen := strings.LastIndex(expression, ")")
	if closeParen == -1 {
		return nil, fmt.Errorf("无效的函数调用: %s", expression)
	}
	
	argsStr := strings.TrimSpace(expression[openParen+1 : closeParen])
	
	// 解析参数（简化实现）
	var args []Expression
	if argsStr != "" {
		argStrs := strings.Split(argsStr, ",")
		for _, argStr := range argStrs {
			argStr = strings.TrimSpace(argStr)
			arg, err := cp.Parse(argStr)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
	}
	
	return NewFunctionExpression(funcName, args), nil
}

// isBinaryExpression 检查是否是二元表达式
func (cp *ConfigParser) isBinaryExpression(expression string) bool {
	operators := []string{"+", "-", "*", "/", "==", "!=", "&&", "||"}
	for _, op := range operators {
		if strings.Contains(expression, op) {
			return true
		}
	}
	return false
}

// parseBinaryExpression 解析二元表达式
func (cp *ConfigParser) parseBinaryExpression(expression string) (Expression, error) {
	// 简化实现
	// 实际项目中需要更复杂的解析逻辑
	
	operators := []string{"+", "-", "*", "/", "==", "!=", "&&", "||"}
	
	for _, op := range operators {
		if strings.Contains(expression, op) {
			parts := strings.Split(expression, op)
			if len(parts) == 2 {
				left, err := cp.Parse(strings.TrimSpace(parts[0]))
				if err != nil {
					return nil, err
				}
				
				right, err := cp.Parse(strings.TrimSpace(parts[1]))
				if err != nil {
					return nil, err
				}
				
				return NewBinaryExpression(left, op, right), nil
			}
		}
	}
	
	return nil, fmt.Errorf("无法解析二元表达式: %s", expression)
}
