package evaluator

import (
	"Interpreter/ast"
	"Interpreter/object"
	"fmt"
)

// 每一个AST 节点都实现了 ast.Node
func Eval(node ast.Node, ctx *object.Context) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, ctx)
		//return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, ctx)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.BlockStatement:
		return evalBlockStatement(node, ctx)
		//return evalStatements(node.Statements)
	case *ast.PrefixExpression:
		right := Eval(node.Right, ctx)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, ctx)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, ctx)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, ctx)

	case *ast.ReturnStatement:
		// 如果是返回语句，则继续计算返回表达式
		val := Eval(node.ReturnValue, ctx)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, ctx) // LetStatement 语句表达式存在 Value 中
		if isError(val) {
			return val
		}
		ctx.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, ctx)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		// 给每一个函数增加一个全局上下文环境，然后再给一个函数增加一个自有的函数栈，上下文环境是一个指针，因此，每个函数都能修改全局上下文
		return &object.Function{Parameters: params, Body: body, Ctx: ctx}

	case *ast.CallExpression:
		function := Eval(node.Function, ctx) // 计算函数参数中的表达式
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, ctx) // 计算每一个参数
		if len(args) == 1 && isError(args[0]) {
			// 如果只有一个参数 且是一个错误，则返回该错误
			return args[0]
		}
		callFunction(function, args)
	}
	return nil
}

func callFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	fnCtx := callFunctionCtx(function, args)

	// function.Body 是一个代码块
	evaluated := Eval(function.Body, fnCtx) // 计算函数体，然后用自身的上下文环境
	// 将一个ReturnValue 对象，转换成一个 Integer 或者 Boolean object
	return unwarpReturnValue(evaluated)
}

func callFunctionCtx(fn *object.Function, args []object.Object) *object.Context {
	// fn.ctx 将全局的上下文环境保留
	ctx := object.NewFunctionStackContext(fn.Ctx)
	for paramIdx, param := range fn.Parameters {
		// param.Value 代表原始变量名，args[paramIdx] 代表该变量被计算的值
		ctx.Set(param.Value, args[paramIdx])
		// 加入到自身函数栈中
	}
	return ctx
}
func unwarpReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		// 将ReturnValue 中的值取出来
		return returnValue.Value
	}
	return obj
}
func evalExpressions(exps []ast.Expression, ctx *object.Context) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, ctx)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIdentifier(node *ast.Identifier, ctx *object.Context) object.Object {
	val, ok := ctx.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}

func evalProgram(program *ast.Program, ctx *object.Context) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, ctx)
		// 如果 遇到return 语句 则，结束剩下语句的解析
		//if returnValue, ok := result.(*object.ReturnValue); ok {
		//	return returnValue.Value
		//}
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		// 这里直接进的了指标比较因为 BooleanObject 类型的
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type(): // 二元操作符两边的类型不相等，不能计算
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
		//return NULL
	}
}
func evalIfExpression(ie *ast.IfExpression, ctx *object.Context) object.Object {
	condition := Eval(ie.Condition, ctx) // 先计算if 的条件
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, ctx)
	} else if ie.Alternative != nil { // 如果有else 节点 则 计算else
		return Eval(ie.Alternative, ctx)
	} else {
		return NULL
	}
}
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true // 其他类型都为真
	}
}
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
		//return NULL
	}
}
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s %s %s", operator, right, right.Type())
		//return NULL
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
		//return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE

	}
}
func evalStatements(stmts []ast.Statement, ctx *object.Context) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, ctx)

		if resultValue, ok := result.(*object.ReturnValue); ok {
			return resultValue.Value
		}
	}
	return result
}
func evalBlockStatement(block *ast.BlockStatement, ctx *object.Context) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, ctx)
		// 如果是返回值类型的对象，则返回具体类型，并且遇到return 语句则退出当前解析，不再计算剩下的表达式.
		//if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
		//	return result
		//}
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// nativeBoolToBooleanObject 减少内存申请
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return object.ERROR_OBJ == obj.Type()
	}
	return false
}
