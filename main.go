package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const historyFile = "data/history.txt"

func main() {
	reader := bufio.NewReader(os.Stdin)

	calc := NewCalculator()
	conv := NewConverter()
	hm := NewHistoryManager(historyFile)

	for { // the TUI's main loop: keep showing the menu until the user exits
		printMenu()
		choice := readLine(reader)

		switch choice {
		case "1":
			runCalculator(reader, calc)
		case "2":
			runConverter(reader, conv)
		case "3":
			runBatchDemo(conv)
		case "4":
			showHistory(calc, conv)
		case "5":
			saveAllHistory(hm, calc, conv)
		case "6":
			loadAndPrintHistory(hm)
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func printMenu() {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 42))
	fmt.Println("   GO TUI — CALCULATOR & UNIT CONVERTER")
	fmt.Println(strings.Repeat("=", 42))
	fmt.Println("1. Calculator")
	fmt.Println("2. Unit Converter")
	fmt.Println("3. Batch Convert (goroutines demo)")
	fmt.Println("4. View Session History")
	fmt.Println("5. Save History to File")
	fmt.Println("6. Load History from File")
	fmt.Println("0. Exit")
	fmt.Print("Choose an option: ")
}

// --- input helpers -------------------------------------------------------

func readLine(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// readFloat loops until the user provides a value strconv can parse.
func readFloat(reader *bufio.Reader, prompt string) float64 {
	for {
		fmt.Print(prompt)
		input := readLine(reader)
		val, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Please enter a valid number.")
			continue
		}
		return val
	}
}

// --- menu actions ----------------------------------------------------------

func runCalculator(reader *bufio.Reader, calc *Calculator) {
	a := readFloat(reader, "Enter first number: ")
	fmt.Print("Enter operator (+, -, *, /): ")
	op := readLine(reader)
	b := readFloat(reader, "Enter second number: ")

	result, err := calc.Calculate(a, b, op)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Result: %.2f\n", result)
}

func runConverter(reader *bufio.Reader, conv *Converter) {
	fmt.Println("1. Length   2. Weight   3. Temperature")
	fmt.Print("Choose category: ")
	category := readLine(reader)

	var catName string
	switch category {
	case "1":
		catName = "length"
	case "2":
		catName = "weight"
	case "3":
		catName = "temperature"
	default:
		fmt.Println("Invalid category")
		return
	}
	fmt.Printf("Available units: %s\n", SupportedUnits(catName))

	value := readFloat(reader, "Enter value: ")
	fmt.Print("From unit: ")
	from := readLine(reader)
	fmt.Print("To unit: ")
	to := readLine(reader)

	var result float64
	var err error

	switch catName {
	case "length":
		result, err = conv.ConvertLength(value, from, to)
	case "weight":
		result, err = conv.ConvertWeight(value, from, to)
	case "temperature":
		result, err = conv.ConvertTemperature(value, from, to)
	}

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Result: %.2f\n", result)
}

// runBatchDemo shows several conversions happening concurrently via goroutines.
func runBatchDemo(conv *Converter) {
	jobs := []ConversionJob{
		{Value: 100, FromUnit: "cm", ToUnit: "m", Category: "length"},
		{Value: 1, FromUnit: "kg", ToUnit: "lb", Category: "weight"},
		{Value: 100, FromUnit: "C", ToUnit: "F", Category: "temperature"},
		{Value: 5, FromUnit: "mile", ToUnit: "km", Category: "length"},
		{Value: 32, FromUnit: "F", ToUnit: "K", Category: "temperature"},
	}

	fmt.Printf("Running %d conversions concurrently...\n", len(jobs))
	results := RunBatchConversion(conv, jobs)

	for _, r := range results { // loop over results collected from all goroutines
		if r.Err != nil {
			fmt.Printf("  [ERROR] %s -> %s: %v\n", r.Job.FromUnit, r.Job.ToUnit, r.Err)
			continue
		}
		fmt.Printf("  %.2f %s -> %.2f %s\n", r.Job.Value, r.Job.FromUnit, r.Output, r.Job.ToUnit)
	}
}

func showHistory(calc *Calculator, conv *Converter) {
	fmt.Println("--- Calculator History ---")
	if len(calc.History) == 0 {
		fmt.Println("(empty)")
	}
	for i, h := range calc.History {
		fmt.Printf("%d. %s\n", i+1, h)
	}

	fmt.Println("--- Converter History ---")
	if len(conv.History) == 0 {
		fmt.Println("(empty)")
	}
	for i, h := range conv.History {
		fmt.Printf("%d. %s\n", i+1, h)
	}
}

func saveAllHistory(hm *HistoryManager, calc *Calculator, conv *Converter) {
	all := append([]string{}, calc.History...)
	all = append(all, conv.History...)

	if len(all) == 0 {
		fmt.Println("Nothing to save yet.")
		return
	}

	if err := hm.SaveHistory(all); err != nil {
		fmt.Println("Failed to save history:", err)
		return
	}
	fmt.Println("History saved to", hm.FilePath)
}

func loadAndPrintHistory(hm *HistoryManager) {
	lines, err := hm.LoadHistory()
	if err != nil {
		fmt.Println("Failed to load history:", err)
		return
	}
	if len(lines) == 0 {
		fmt.Println("No saved history found.")
		return
	}
	fmt.Println("--- Saved History (from file) ---")
	for _, l := range lines {
		fmt.Println(l)
	}
}
