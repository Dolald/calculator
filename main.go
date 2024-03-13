package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isOperator(char string) bool {
	operators := map[string]bool{"+": true, "-": true, "*": true, "/": true}
	return operators[char]
}

func getPriority(char string) int {
	priority := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}
	return priority[char]
}

func changeAllStacks(stack []string, outPut []float64) ([]string, []float64, error) {

	leftNum, rightNum := outPut[len(outPut)-2], (outPut)[len(outPut)-1]
	top := stack[len(stack)-1]

	switch {
	case top == "+":
		if stack[len(stack)-2] == "-" {
			outPut[len(outPut)-2] = -leftNum + rightNum
			stack[len(stack)-2] = "+"
		} else {
			outPut[len(outPut)-2] = leftNum + rightNum
		}
	case top == "-":
		if stack[len(stack)-2] == "-" {
			outPut[len(outPut)-2] = -leftNum - rightNum
			stack[len(stack)-2] = "+"
		} else {
			outPut[len(outPut)-2] = leftNum - rightNum
		}
	case top == "*":
		outPut[len(outPut)-2] = leftNum * rightNum
	case top == "/":
		if rightNum == 0 {
			return stack, outPut, fmt.Errorf("не дели на 0, а то будет грустно")
		}
		outPut[len(outPut)-2] = leftNum / rightNum
	}

	outPut = outPut[:len(outPut)-1]
	stack = stack[:len(stack)-1]

	return stack, outPut, nil
}

func ifEnteredOneSymbol(stringArray []string) (float64, error) {
	oneNumber, _ := strconv.ParseFloat(stringArray[0], 64) // проверка на наличие только одного элемента в строке
	if !(stringArray[0] >= "-1" && stringArray[0] <= "9") {
		return oneNumber, fmt.Errorf("ошибка, давай без этих, там, непонятных знаков")
	}
	return oneNumber, nil
}

func infixToRPN(stringArray []string) (float64, error) {
	var err error
	stack := []string{}
	outPut := []float64{}

	if len(stringArray) == 1 {
		oneNumber, err := ifEnteredOneSymbol(stringArray)
		if err != nil {
			return 0, err
		}
		return oneNumber, nil
	}

	stack, outPut, err = OperateBeforeRemovingAllParenthesises(stringArray, stack, outPut, err)
	if err != nil {
		return 0, err
	}

	for len(outPut) != 2 { // проводим операции до тех пор, пока в outPut не останется 2 операнда
		stack, outPut, err = changeAllStacks(stack, outPut)
		if err != nil {
			return 0, err
		}
	}

	leftNum, rightNum := outPut[len(outPut)-2], outPut[len(outPut)-1] // проводим операцию с последним арифметическим действием
	top := stack[len(stack)-1]

	switch {
	case top == "+":
		outPut[len(outPut)-2] = leftNum + rightNum
	case top == "-":
		outPut[len(outPut)-2] = leftNum - rightNum
	case top == "*":
		outPut[len(outPut)-2] = leftNum * rightNum
	case top == "/":
		if rightNum == 0 {
			return outPut[0], fmt.Errorf("ошибка, нельзя делить на 0")
		}
		outPut[len(outPut)-2] = leftNum / rightNum
	}

	outPut = outPut[:len(outPut)-1]

	return outPut[0], nil
}

func OperateBeforeRemovingAllParenthesises(stringArray []string, stack []string, outPut []float64, err error) ([]string, []float64, error) {
	for i, token := range stringArray {

		if token == ")" {
			stack, outPut, err = operationWithOpenParenthesis(stack, outPut)
			if err != nil {
				return stack, outPut, err
			}
		}

		if token == "(" {
			stack = append(stack, "(")
		}

		if isOperator(token) {
			if ifLastOperationIsGreaterThanCurrentOrEqual(stack, token) {
				stack, outPut, err = changeAllStacks(stack, outPut)
				if err != nil {
					return stack, outPut, err
				}
			}

			stack = append(stack, token)

			if token == "-" && i == 0 { // исключительный случай, когда первый элемент стека "-" и outPut пустой, значит следущий элемент - "("
				outPut = append(outPut, 0)
			}

		} else if !isOperator(token) && token != ")" && token != "(" {
			outPut, err = appendOutPutStack(outPut, token)
			if err != nil {
				return stack, outPut, err
			}
		}
	}
	return stack, outPut, nil
}

func appendOutPutStack(outPut []float64, token string) ([]float64, error) {
	newToken, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return outPut, fmt.Errorf("ошибка")
	}
	outPut = append(outPut, newToken)
	return outPut, err
}

func ifLastOperationIsGreaterThanCurrentOrEqual(stack []string, token string) bool {
	return len(stack) > 0 && ((getPriority(stack[len(stack)-1]) > getPriority(token)) || (getPriority(stack[len(stack)-1]) == 2 && getPriority(token) == 2))
}

func operationWithOpenParenthesis(stack []string, outPut []float64) ([]string, []float64, error) {
	for i := range stack {
		if doesExistOperationBeforeOpenParenthesis(stack, i) { // если после "(" есть хотя бы одна операция, продолжаем
			for stack[len(stack)-1] != "(" {
				stack, outPut, _ = changeAllStacks(stack, outPut)
			}
			break
		} else {
			return stack, outPut, fmt.Errorf("Ошибка")
		}
	}
	stack = stack[:len(stack)-1] // убираем "(" из стека и продолжаем дальше
	return stack, outPut, nil
}

func doesExistOperationBeforeOpenParenthesis(stack []string, i int) bool {
	// if stack[i] == "(" && i != len(stack)-1 { // если после "(" есть хотя бы одна операция, продолжаем
	// 	return true
	if stack[i] == "(" && i == len(stack)-1 {
		return false
	}
	return true
}

func checkForCorrectString(enterString string) error {
	bracketStack := make([]rune, 0)
	bracketStackForCheck := make([]rune, 0)

	var countOperations, countLeftTypeSampethizes, countRightTypeSampethizes int

	countSpaces := strings.Count(enterString, " ")
	enterString = strings.ReplaceAll(enterString, " ", "")

	for i := 0; i < len(enterString); i++ { // количество операций

		switch enterString[i] {
		case '+', '-', '*', '/', '(', ')': // подсчёт количества операций, исключая случай, когда самый первый элемент "-"
			countOperations++
			if i == 0 && enterString[i] == '-' {
				countOperations--
			}
			if i == 0 && enterString[i] == '-' && enterString[i+1] == '(' { // исключительный случай
				countOperations++
				countSpaces++
			}
		}

		if enterString[len(enterString)-1] == byte(enterString[i]) { // проверка на последний символ
			if !(enterString[i] >= '0' && enterString[i] <= '9') && enterString[i] != ')' {
				return fmt.Errorf("ошибка, вы ввели не верный последний символ")
			}
		}

		switch enterString[i] {
		case '+', '-', '*', '/', '(':
			if enterString[i+1] == '-' && (enterString[i+2] >= '0' && enterString[i+2] <= '9') {
				countOperations--
			}
		}

		if enterString[i] == ')' {
			countLeftTypeSampethizes++
			bracketStackForCheck = append(bracketStackForCheck, ')')
		} else if enterString[i] == '(' {
			countRightTypeSampethizes++
			bracketStackForCheck = append(bracketStackForCheck, '(')
		}
	}

	if len(bracketStackForCheck) >= 1 {

		if bracketStackForCheck[0] != '(' {
			return fmt.Errorf("ошибка, не верная первая скобка")
		}

		if countLeftTypeSampethizes != countRightTypeSampethizes { // проверка количества скобок
			return fmt.Errorf("ошибка, неверное количество скобок")
		}
	}

	if countOperations*2-(countLeftTypeSampethizes+countRightTypeSampethizes) != countSpaces {
		return fmt.Errorf("ошибка, пиши нормально, с пробелами строку")
	}

	for i, char := range enterString {
		switch char {
		case '+', '-', '*', '/':
			if enterString[i+1] == '*' || enterString[i+1] == '+' || enterString[i+1] == '/' { //проверка на верность следующей операции сразу после текущей
				return fmt.Errorf("ошибка, вы ввели не правильную операцию")
			}
		}

		if !(char >= '0' && char <= '9') && !(char == '-' || char == '/' || char == '*' || char == '+' || char == '(' || char == ')') { // проверка на верность самих символов
			return fmt.Errorf("ошибка, вы ввели не правильную операцию, наверно")
		}

		if char == ')' {
			bracketStack = bracketStack[:len(bracketStack)-1]
		}

		if char == '(' {
			bracketStack = append(bracketStack, char)
		}
	}

	if len(bracketStack) != 0 {
		return fmt.Errorf("ошибка, не правильно расставлены скобки")
	}

	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите строку: ")
	enterString, _ := reader.ReadString('\n')
	enterString = enterString[:len(enterString)-1] // Удаление символа новой строки из строки

	err := checkForCorrectString(enterString)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	stringArray := strings.Split(enterString, " ")

	result, err := infixToRPN(stringArray)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
