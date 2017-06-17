package vm

import (
	"testing"

	"github.com/goby-lang/goby/bytecode"
	"github.com/goby-lang/goby/lexer"
	"github.com/goby-lang/goby/parser"
)

func testEval(t *testing.T, input string) Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	g := bytecode.NewGenerator(program)
	bytecodes := g.GenerateByteCode(program)
	return testExec(bytecodes)
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testExec(bytecodes string) Object {
	v := New("./", []string{})
	v.ExecBytecodes(bytecodes, "./")

	return v.mainThread.stack.top().Target
}

func testIntegerObject(t *testing.T, obj Object, expected int) bool {
	switch result := obj.(type) {
	case *IntegerObject:
		if result.Value != expected {
			t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Error(result.Message)
		return false
	default:
		t.Errorf("object is not Integer. got=%T (%+v).", obj, obj)
		return false
	}
}

func testNullObject(t *testing.T, obj Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj Object, expected string) bool {
	switch result := obj.(type) {
	case *StringObject:
		if result.Value != expected {
			t.Errorf("object has wrong value. expect=%s, got=%s", expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Error(result.Message)
		return false
	default:
		t.Errorf("object is not String. got=%T (%+v).", obj, obj)
		return false
	}
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	switch result := obj.(type) {
	case *BooleanObject:
		if result.Value != expected {
			t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Error(result.Message)
		return false
	default:
		t.Errorf("object is not Boolean. got=%T (%+v).", obj, obj)
		return false
	}
}

func testArrayObject(t *testing.T, obj Object, expected *ArrayObject) bool {
	result, ok := obj.(*ArrayObject)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}

	if len(result.Elements) != len(expected.Elements) {
		t.Fatalf("Don't equals length of array. Expect %d, got=%d", len(expected.Elements), len(result.Elements))
	}

	for i := 0; i < len(result.Elements); i++ {
		intObj, ok := expected.Elements[i].(*IntegerObject)
		if ok {
			testIntegerObject(t, result.Elements[i], intObj.Value)
			continue
		}
		str, ok := expected.Elements[i].(*StringObject)
		if ok {
			testStringObject(t, result.Elements[i], str.Value)
			continue
		}

		b, ok := expected.Elements[i].(*BooleanObject)
		if ok {
			testBooleanObject(t, result.Elements[i], b.Value)
			continue
		}

		t.Fatalf("object is wrong type %T", expected.Elements[i])
	}

	return true
}

func checkExpected(t *testing.T, evaluated Object, expected interface{}) {
	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, evaluated, expected)
	case string:
		testStringObject(t, evaluated, expected)
	case bool:
		testBooleanObject(t, evaluated, expected)
	case nil:
		_, ok := evaluated.(*NullObject)

		if !ok {
			t.Fatalf("expect result should be Null. got=%T", evaluated)
		}
	}
}

func isError(obj Object) bool {
	if obj != nil {
		_, ok := obj.(*Error)
		return ok
	}
	return false
}
