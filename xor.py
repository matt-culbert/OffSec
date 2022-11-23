#!/usr/bin/python3
from itertools import cycle
import os
import chardet
from anyascii import anyascii

data = "hello this is a test"
key = "fkghjw78hgw5w478ghw5h-245-t24=56=-sfdg=-"

def xor(data, key):
	new = [chr(ord(a) ^ ord(b)) for a,b in zip(data, cycle(key))] # Cycle here is really important, without it the program will XOR first x where x is the length of the key
                                                                  # Zip is equally important.
	result = "".join(new)
	return new

encoded = xor(data, key)
print(encoded)
print(xor(encoded, key))
test = ['\x0e', '\x0e', '\x0b', '\x04', '\x05', 'W', 'C', 'P', '\x01', '\x14', 'W', '\\', '\x04', '\x14', 'V', '\x18', '\x13', '\r', '\x04', 'A']
'''
Ok, here's the deal
In the XOR function, we get back either byte code like this or we get back the string joined together which doesn't render because it's unicode??? I don't know why. Don't ask
So here's the solution!
We get the byte code seperated into this list and do a string.join before decrypting. And it works!
'''
temp = ''.join(test)
print(xor(temp, key))