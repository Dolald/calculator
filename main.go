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

func changeAllStacks(stack []string, outPut []float64) (error, []string, []float64) {

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
			return fmt.Errorf("Не дели на 0, а то будет грустно"), stack, outPut
		}
		outPut[len(outPut)-2] = leftNum / rightNum
	}

	outPut = outPut[:len(outPut)-1]
	stack = stack[:len(stack)-1]

	return nil, stack, outPut
}

func infixToRPN(infix string) (float64, error) {
	stringArray := strings.Split(infix, " ")
	var err error

	stack := []string{}
	outPut := []float64{}

	if len(stringArray) == 1 {
		oneNumber, _ := strconv.ParseFloat(stringArray[0], 64) // проверка на наличие только одного элемента в строке
		if !(stringArray[0] >= "-1" && stringArray[0] <= "9") {
			return oneNumber, fmt.Errorf("Ошибкаы, давай без этих, там, непонятных знаков")
		}
		return oneNumber, nil
	}

	for i, token := range stringArray {

		if token == ")" {

			for i := range stack {
				fmt.Println(stack)

				if stack[i] == "(" && i != len(stack)-1 { // если после "(" есть хотя бы одна операция, продолжаем
					for stack[len(stack)-1] != "(" {
						_, stack, outPut = changeAllStacks(stack, outPut)
					}
					break
				} else if stack[i] == "(" && i == len(stack)-1 {
					return 0, fmt.Errorf("Ошибка")
				}
			}
			stack = stack[:len(stack)-1] // убираем "(" из стека и продолжаем дальше
		}
		fmt.Println(outPut)

		if token == "(" {
			stack = append(stack, "(")
		}

		if isOperator(token) { // если токен операция, проверяем дальше
			if len(stack) > 0 {
				// если прошлая операция больше, чем текущая или они равны 2, то есть * * или * / и тд - проводим предыдущую операцию
				if (getPriority(stack[len(stack)-1]) > getPriority(token)) || (getPriority(stack[len(stack)-1]) == 2 && getPriority(token) == 2) {
					err, stack, outPut = changeAllStacks(stack, outPut)
					if err != nil {
						return 0, err
					}
				}
			}
			stack = append(stack, token)

			if token == "-" && i == 0 { // исключительный случай, когда первый элемент стека "-" и outPut пустой, значит следущий элемент - "("
				outPut = append(outPut, 0)
			}

		} else if !isOperator(token) && token != ")" && token != "(" {

			newToken, err := strconv.ParseFloat(token, 64) // если не операция, а операнда
			if err != nil {
				return 0, fmt.Errorf("ошибка")
			}

			outPut = append(outPut, newToken)
		}
	}
	fmt.Println(stack, outPut)

	for len(outPut) != 2 {
		err, stack, outPut = changeAllStacks(stack, outPut)
		if err != nil {
			return 0, err
		}
		fmt.Println(stack, outPut)
	}

	leftNum, rightNum := (outPut)[len(outPut)-2], (outPut)[len(outPut)-1] // проводим операцию с последним арифметическим действием
	top := stack[len(stack)-1]

	switch {
	case top == "+":
		(outPut)[len(outPut)-2] = leftNum + rightNum
	case top == "-":
		(outPut)[len(outPut)-2] = leftNum - rightNum
	case top == "*":
		outPut[len(outPut)-2] = leftNum * rightNum
	case top == "/":
		if rightNum == 0 {
			return outPut[0], fmt.Errorf("Ошибка, нельзя делить на 0")
		}
		outPut[len(outPut)-2] = leftNum / rightNum
	}

	outPut = outPut[:len(outPut)-1]
	stack = stack[:len(stack)-1]

	return outPut[0], nil
}

///////////////////////////////////////////////////////////////////////////////////

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
				return fmt.Errorf("Ошибка, вы ввели не верный последний символ")
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

		fmt.Println(bracketStackForCheck)

		if countLeftTypeSampethizes != countRightTypeSampethizes { // проверка количества скобок
			return fmt.Errorf("Ошибка, неверное количество скобок")
		}
	}

	fmt.Println("Количество пробелов", countSpaces, "\nКоличество операций", countOperations)

	if countOperations*2-(countLeftTypeSampethizes+countRightTypeSampethizes) != countSpaces {
		return fmt.Errorf("Ошибка, пиши нормально, с пробелами строку")
	}

	for i, char := range enterString {
		switch char {
		case '+', '-', '*', '/':
			if enterString[i+1] == '*' || enterString[i+1] == '+' || enterString[i+1] == '/' { //проверка на верность следующей операции сразу после текущей
				return fmt.Errorf("Ошибка, вы ввели не правильную операцию")
			}
		}

		if !(char >= '0' && char <= '9') && !(char == '-' || char == '/' || char == '*' || char == '+' || char == '(' || char == ')') { // проверка на верность самих символов
			return fmt.Errorf("Ошибка, вы ввели не правильную операцию, наверно")
		}

		if char == ')' {
			bracketStack = bracketStack[:len(bracketStack)-1]
		}

		if char == '(' {
			bracketStack = append(bracketStack, char)
		}
	}

	if len(bracketStack) != 0 {
		return fmt.Errorf("Ошибка, не правильно расставлены скобки")
	}

	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите строку: ")
	enterString, _ := reader.ReadString('\n')

	// Удаление символа новой строки из строки
	enterString = enterString[:len(enterString)-1]

	err := checkForCorrectString(enterString)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := infixToRPN(enterString)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
