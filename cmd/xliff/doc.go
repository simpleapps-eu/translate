/*
Usage:
	#Translate

	Read the strings from -in <file> .strings file, then translate them using the translation
	from -xlf <file> XLIFF file and write out to -out <file> .strings file.

	e.g. xliff -out fr.strings -in en.strings -xlf fr.xlf

	#Normalize

	Read the strings from -in <file> .strings file and then write them out again to
	the -out <file> .strings file. This detects any errors in the .strings file,
	normalizes the strings and cleans up any formatting issues while writing them back out.

	e.g. xliff -out normalized.strings -in en.strings

	#Convert

	Read the en.strings and write out a fresh .xlf file to be send on to translators.

	e.g. xliff -xlf fr.xlf -in en.strings

*/
package main
