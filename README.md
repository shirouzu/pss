# What is pss:
ps command with PSS/USS information

# Usage:
pss [ps args(default: aux)...]

# Build:
go build pss.go

# Recommend to set setuid bit for reading /proc/xxx/smaps
1st step: chown root.root pss<br>
2nd step: chmod 4755 pss

# Sample
	% pss
	USER       PID %CPU %MEM    VSZ     PSS   USS   RSS TTY      STAT START   TIME COMMAND
	root         1  0.0  0.0   7900     207    64    24 ?        Ss   11月10   0:05 init [2]
	  :
	shirouzu 22084  0.0  0.1  16192     733   224  1744 ?        S    22:36   0:00 sshd: shirouzu@pts/0
	shirouzu 22085  0.0  0.4   6232    2884  2364  4140 pts/0    Ss   22:36   0:00 -bash
	www-data 22629  0.0  0.0  24768    1840    32   856 ?        S    03:00   0:01 /usr/sbin/apache2 -k
	www-data 22630  2.5  9.0 3200712 201943 89392  91660 ?       Sl   03:00  32:01 /usr/sbin/apache2 -k
	  :
	Total:VSZ PSS USS RSS 6,762,184 589,939 200,684 270,500

