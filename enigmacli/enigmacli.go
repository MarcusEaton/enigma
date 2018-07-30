package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// EnigmaMachine is the live state of the machine
var EnigmaMachine MachineState

// Letter lookups for converting runes into integers for arithmetic
var letterLookup map[rune]int
var reverseLetterLookup map[int]rune

func main() {
	// Set the machine using the settings file
	EnigmaMachine = loadSettingsFromFile("settings.json")
	fmt.Println("Machine set using settings.json")

	// read in the message
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter message (all in CAPS): ")
	message, _ := reader.ReadString('\n')

	// Encrypt the message
	cipher := EnigmaMachine.Encrypt(message)

	// Add spaces to the ciphertext so it prints formally
	pirintCipher := ""
	for i, r := range cipher {
		pirintCipher += string(r)
		if (i+1)%5 == 0 {
			pirintCipher += " "
		}
	}
	fmt.Println(pirintCipher)
}

// encryptHandler takes the message from the body of the http request, encrypts it using the current state
// of the machine, then prints the output.
func encryptHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	message := string(bodyBytes)
	cipher := EnigmaMachine.Encrypt(message)

	// Add spaces to the ciphertext so it prints formally
	pirintCipher := ""
	for i, r := range cipher {
		pirintCipher += string(r)
		if (i+1)%5 == 0 {
			pirintCipher += " "
		}
	}
	fmt.Fprintf(w, pirintCipher)
}

// Encrypt the message using the EnigmaMachine object
func (machine *MachineState) Encrypt(message string) string {

	cipher := ""
	for _, letter := range message {

		// Check that the rune is valid (any capital letter), otherwise skip it
		_, validRune := letterLookup[letter]
		if !validRune {
			continue
		}

		// Run the letter through the plug board, only change it if the letter is one of the ones with a plug
		newLetter, hasPlug := machine.Plugs[letter]
		if hasPlug {
			letter = newLetter
		}

		// Run the letter through the rotors
		i := letterLookup[letter]
		i = machine.Rotors[0].EncryptForward(i)
		i = machine.Rotors[1].EncryptForward(i)
		i = machine.Rotors[2].EncryptForward(i)

		// Run the letter through the reflector
		i = machine.Reflector.Map[i]

		// Run the letter back through the rotors
		i = machine.Rotors[2].EncryptBackward(i)
		i = machine.Rotors[1].EncryptBackward(i)
		i = machine.Rotors[0].EncryptBackward(i)

		// Run the output back through the plug board
		letter = reverseLetterLookup[i]
		newLetter, hasPlug = machine.Plugs[letter]
		if hasPlug {
			letter = newLetter
		}

		// Append the output letter to the ciphertext
		cipher = (cipher + string(letter))

		// Rotate the rotors, (mod(26)!)
		machine.Rotors[0].Position = mod(machine.Rotors[0].Position+1, 26)
		if machine.Rotors[0].Position == machine.Rotors[0].TurnPosition {
			machine.Rotors[1].Position = mod(machine.Rotors[1].Position+1, 26)
			if machine.Rotors[1].Position == machine.Rotors[1].TurnPosition {
				machine.Rotors[2].Position = mod(machine.Rotors[2].Position+1, 26)
			}
		}
	}

	return cipher
}

// MachineState object represents a state of an enigma machine.
type MachineState struct {
	Rotors    [3]Rotor
	Reflector Reflector
	Plugs     map[rune]rune
}

// Rotor object containing letter map
type Rotor struct {
	Number       int
	Position     int
	TurnPosition int
	Map          map[int]int
}

// EncryptForward runs a single letter through a single rotor
func (r *Rotor) EncryptForward(input int) (output int) {
	i := mod((input + r.Position), 26)
	i = r.Map[i]
	return mod((i - r.Position), 26)
}

// EncryptBackward runs a single letter through a single rotor backwards
func (r *Rotor) EncryptBackward(input int) (output int) {
	i := mod((input + r.Position), 26)
	for key, value := range r.Map {
		if i == value {
			i = key
			break
		}
	}
	return mod((i - r.Position), 26)
}

// LoadMap loads the map for the rotor based on its number
func (r *Rotor) LoadMap() {
	r.Map = map[int]int{}

	var mapString string

	// Enigma rotor definitions based on https://en.wikipedia.org/wiki/Enigma_rotor_details
	switch r.Number {
	case 1:
		mapString = "EKMFLGDQVZNTOWYHXUSPAIBRCJ"
		r.TurnPosition = letterLookup[[]rune("R")[0]]
	case 2:
		mapString = "AJDKSIRUXBLHWTMCQGZNPYFVOE"
		r.TurnPosition = letterLookup[[]rune("F")[0]]
	case 3:
		mapString = "BDFHJLCPRTXVZNYEIWGAKMUSQO"
		r.TurnPosition = letterLookup[[]rune("W")[0]]
	case 4:
		mapString = "ESOVPZJAYQUIRHXLNFTGKDCMWB"
		r.TurnPosition = letterLookup[[]rune("K")[0]]
	case 5:
		mapString = "VZBRGITYUPSDNHLXAWMJQOFECK"
		r.TurnPosition = letterLookup[[]rune("A")[0]]
	}

	// Loop through the mapString, populating r.map
	for i, char := range mapString {
		key := letterLookup[char]
		r.Map[i] = key
	}
}

// Reflector object
type Reflector struct {
	Letter string
	Map    map[int]int
}

// LoadMap loads the map for the reflector, using the letter.
func (r *Reflector) LoadMap() {
	r.Map = map[int]int{}

	var mapString string

	// Enigma reflector definitions based on https://en.wikipedia.org/wiki/Enigma_rotor_details
	switch r.Letter {
	case "A":
		mapString = "EJMZALYXVBWFCRQUONTSPIKHGD"
	case "B":
		mapString = "YRUHQSLDPXNGOKMIEBFZCWVJAT"
	case "C":
		mapString = "FVPJIAOYEDRZXWGCTKUQSBNMHL"
	}

	// Loop through the mapString, populating the actual reflector map
	for i, char := range mapString {
		key := letterLookup[char]
		r.Map[i] = key
	}
}

// Homemade Mod function because Go Mod doesn't handle negatives properly
func mod(d, m int) int {
	res := d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

// loadSettings pulls the data from settings.json and unmarshals it into a Setting object
func loadSettingsFromFile(file string) MachineState {

	// Grab the settings.json file
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return loadSettingsFromBytes(raw)
}

// loadSettings marshals a settings object from a json byte array
func loadSettingsFromBytes(settingsString []byte) MachineState {

	// Set the letter lookup (used by the LoadMap functions)
	makeLetterLookup()
	makeReverseLetterLookup()

	// Marshal it into a setting object
	var s setting
	err2 := json.Unmarshal(settingsString, &s)
	if err2 != nil {
		fmt.Println(err2.Error())
	}

	// Create our machine state
	var ms MachineState

	// Set the rotors based on the settings, and populate the rotor maps
	for i, rotor := range s.Rotors {
		ms.Rotors[i] = Rotor{rotor.Number, rotor.Position - 1, 0, nil}
		ms.Rotors[i].LoadMap()
	}

	// Set the reflector based on the settings
	ms.Reflector.Letter = s.Reflector
	ms.Reflector.LoadMap()

	// Populate the plugs
	ms.Plugs = map[rune]rune{}
	for _, plug := range s.Plugs {
		firstLetter := []rune(plug)[0]
		secondLetter := []rune(plug)[1]
		ms.Plugs[firstLetter] = secondLetter
		ms.Plugs[secondLetter] = firstLetter
	}

	return ms
}

type setting struct {
	Rotors    [3]rotorSetting `json:"Rotors"`
	Reflector string          `json:"Reflector"`
	Plugs     []string        `json:"Plugs"`
}

type rotorSetting struct {
	Number   int `json:"Number"`
	Position int `json:"Position"`
}

func makeLetterLookup() {
	letterLookup = map[rune]int{}
	letterLookup[[]rune("A")[0]] = 0
	letterLookup[[]rune("B")[0]] = 1
	letterLookup[[]rune("C")[0]] = 2
	letterLookup[[]rune("D")[0]] = 3
	letterLookup[[]rune("E")[0]] = 4
	letterLookup[[]rune("F")[0]] = 5
	letterLookup[[]rune("G")[0]] = 6
	letterLookup[[]rune("H")[0]] = 7
	letterLookup[[]rune("I")[0]] = 8
	letterLookup[[]rune("J")[0]] = 9
	letterLookup[[]rune("K")[0]] = 10
	letterLookup[[]rune("L")[0]] = 11
	letterLookup[[]rune("M")[0]] = 12
	letterLookup[[]rune("N")[0]] = 13
	letterLookup[[]rune("O")[0]] = 14
	letterLookup[[]rune("P")[0]] = 15
	letterLookup[[]rune("Q")[0]] = 16
	letterLookup[[]rune("R")[0]] = 17
	letterLookup[[]rune("S")[0]] = 18
	letterLookup[[]rune("T")[0]] = 19
	letterLookup[[]rune("U")[0]] = 20
	letterLookup[[]rune("V")[0]] = 21
	letterLookup[[]rune("W")[0]] = 22
	letterLookup[[]rune("X")[0]] = 23
	letterLookup[[]rune("Y")[0]] = 24
	letterLookup[[]rune("Z")[0]] = 25
}

func makeReverseLetterLookup() {
	reverseLetterLookup = map[int]rune{}
	reverseLetterLookup[0] = []rune("A")[0]
	reverseLetterLookup[1] = []rune("B")[0]
	reverseLetterLookup[2] = []rune("C")[0]
	reverseLetterLookup[3] = []rune("D")[0]
	reverseLetterLookup[4] = []rune("E")[0]
	reverseLetterLookup[5] = []rune("F")[0]
	reverseLetterLookup[6] = []rune("G")[0]
	reverseLetterLookup[7] = []rune("H")[0]
	reverseLetterLookup[8] = []rune("I")[0]
	reverseLetterLookup[9] = []rune("J")[0]
	reverseLetterLookup[10] = []rune("K")[0]
	reverseLetterLookup[11] = []rune("L")[0]
	reverseLetterLookup[12] = []rune("M")[0]
	reverseLetterLookup[13] = []rune("N")[0]
	reverseLetterLookup[14] = []rune("O")[0]
	reverseLetterLookup[15] = []rune("P")[0]
	reverseLetterLookup[16] = []rune("Q")[0]
	reverseLetterLookup[17] = []rune("R")[0]
	reverseLetterLookup[18] = []rune("S")[0]
	reverseLetterLookup[19] = []rune("T")[0]
	reverseLetterLookup[20] = []rune("U")[0]
	reverseLetterLookup[21] = []rune("V")[0]
	reverseLetterLookup[22] = []rune("W")[0]
	reverseLetterLookup[23] = []rune("X")[0]
	reverseLetterLookup[24] = []rune("Y")[0]
	reverseLetterLookup[25] = []rune("Z")[0]
}
