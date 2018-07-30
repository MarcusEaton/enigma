# Enigma Machine in Go
A go program replicating the 1938 Enigma Machine as defined in https://en.wikipedia.org/wiki/Enigma_rotor_details

Written as a CLI and a web service (to make interoperability easier: next on my list is to write the Bombe https://en.wikipedia.org/wiki/Bombe)

# How to use it
## CLI
Compile enigma.go for your OS/Arch ($ go build enigmacli.go). There's an executable in the git folder compiled for MacOS.

Run the executable ($ ./enigmacli) and follow the instructions.

NB. this relies on the settings.json file being valid and in the same directory as the executable. By default they're set to those used on 31st October 1944.

To decrypt: set the machine back to the same state as that in which the message was encrypted (rerun the CLI) and use the cipher as the message.

## Web service
Compile enigma.go for your OS/Arch ($ go build enigmaservice.go)

Run the executable ($ ./enigmaservice)

Set the machine by hitting the POST "set" method (http://localhost:8080/set). The request will need to provide a setting object in the body, like the one defined in examplesettings.json (those used on 31st October 1944).

Encrypt a message by hitting the POST "encrypt" method (http://localhost:8080/encrypt). The request should contain the message in the body ALL IN CAPS (spaces will be ignored). The encrpt method will return the cipher text.

To decrypt: set the machine back to the same state as that in which the message was encrypted. Hit the encrypt method as described above but use the cipher as the message.

# Support
Tested and working on Mac and RaspberryPi. No error handling, I know, sorry! Send me a message if you need a hand.
