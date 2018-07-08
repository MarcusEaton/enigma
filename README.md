## Enigma Machine in Go
A web service written in Go. Replicating the 1938 Enigma Machine as defined in https://en.wikipedia.org/wiki/Enigma_rotor_details

Originally designed as a command line application, I converted it to a web service to make interoperability easier. Next on my list is to write the Bombe (https://en.wikipedia.org/wiki/Bombe)

This means you'll need an http request builder like Postman (https://www.getpostman.com/) if you want to use it without coding an interface.

## How to use it:
Compile enigma.go for your OS/Arch ($ go build enigma.go)
Run the executable (./enigma)
Set the machine by hitting the POST "set" method (http://localhost:8080/set)
The request will need to provide a setting object in the body, like the one defined in examplesettings.json
Encrypt a message by hitting the POST "encrypt" method (http://localhost:8080/encrypt)
The request should contain the message in the body ALL IN CAPS (spaces will be ignored)
The encrpt method will return the cipher text

To decrypt: set the machine back to the same state as that in which the message was encrypted
Hit the encrypt method as describe above but use the cipher as a message.

## Support
Tested and working on Mac and RaspberryPi. No error handling, I know, sorry! Send me a message if you need a hand.

I'm new to Go so constructive critisism/hints welcome!