pss: ps command with PSS/USS information

Usage:
 % pss [ps args(default: aux)...]

Build:
 # go build pss.go

Recommend to set setuid bit for reading /proc/xxx/smaps
 # chown root.root pss
 # chmod 4755 pss



