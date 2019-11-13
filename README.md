# What is pss:
ps command with PSS/USS information

# Usage:
pss [ps args(default: aux)...]

# Build:
go build pss.go

# Recommend to set setuid bit for reading /proc/xxx/smaps
1st step: chown root.root pss
2nd step: chmod 4755 pss
