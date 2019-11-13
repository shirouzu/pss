// Copyright (C) 2019 H.Shirouzu
// pss: ps command with PSS/USS information
//
// pss [ps args(default: aux)...]
//
// Build:
// # go build pss.go
//
// Recommend to set setuid bit for reading /proc/xxx/smaps
// # chown root.root pss
// # chmod 4755 pss

package main

import (
	"fmt"
	"strings"
	"os/exec"
	"os"
	"regexp"
	"strconv"
	"io/ioutil"
	"errors"
)

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return	v
}

func get_size(pid string) (int, int, int, int, error) {
	buf, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/smaps", pid))
	if err != nil {
		return 0, 0, 0, 0, errors.New("can't open smap")
	}
	sml := strings.Split(string(buf), "\n")

	vss, rss, pss, uss := 0, 0, 0, 0
	for _, l := range sml {
		w := strings.Fields(l)
		if len(w) == 3 {
			if w[0] == "Size:" {
				vss += atoi(w[1])
			} else if w[0] == "Rss:" {
				rss += atoi(w[1])
			} else if w[0] == "Pss:" || w[0] == "SwapPss:" {
				pss += atoi(w[1])
			} else if w[0] == "Private_Clean:" || w[0] == "Private_Dirty:" {
				uss += atoi(w[1])
			}
		}
	}
	return	vss, rss, pss, uss, nil
}

func get_termxy() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	s := string(out[:])

	w := regexp.MustCompile("[ \n]+").Split(s, 3)

	cx, _ := strconv.Atoi(w[1])
	cy, _ := strconv.Atoi(w[0])

	return	cx, cy
}

func get_wide_num(s string) int {
	num := 0
	idx := 0
	for _, c := range s {
		if (c & 0xc0) == 0x80 {
			idx++
		} else {
			idx = 0
		}
		if (idx == 2) {
			num++
		}
	}
	return	num
}

func num2strex(v int) string {
	s := fmt.Sprintf("%d", v)
	ret := ""
	for i, _ := range s {
		if i != 0 && i != len(s) -1 && ((len(s)-i) % 3) == 0 {
			ret += ","
		}
		ret += s[i:i+1]
	}
	return ret
}

type CalcPsRes struct {
	psl		[]string
	usl		[]string
	max_psw	int
	max_usw int
	pss_sum int
	uss_sum int
	vss_sum int
	rss_sum int
}

func calc_ps(title string, body []string) CalcPsRes {
	reg := regexp.MustCompile("[^ \t]+")

	cr := CalcPsRes{}
	cr.psl = make([]string, len(body))
	cr.usl = make([]string, len(body))
	cr.max_psw = 3
	cr.max_usw = 3

	pid_i := -1
	for i, w := range reg.FindAllString(title, 10) {
		if w == "PID" {
			pid_i = i
			break
		}
	}

	for i, l := range body {
		w := reg.FindAllString(l, 10)

		if len(w) < 2 {
			continue
		}
		vss, rss, pss, uss, err := get_size(w[pid_i])
		if err != nil {
			cr.psl[i] = "-"
			cr.usl[i] = "-"
			continue
		}
		cr.psl[i] = fmt.Sprintf("%d", pss)
		if len(cr.psl[i]) > cr.max_psw {
			cr.max_psw = len(cr.psl[i])
		}
		cr.usl[i] = fmt.Sprintf("%d", uss)
		if len(cr.usl[i]) > cr.max_usw {
			cr.max_usw = len(cr.usl[i])
		}

		cr.vss_sum += vss
		cr.rss_sum += rss
		cr.pss_sum += pss
		cr.uss_sum += uss
	}
	return	cr
}

func get_ins_idx(title string) int {
	ins_i := strings.Index(title, "VSZ")
	if ins_i != -1 {
		ins_i += 4
	} else {
		ins_i = strings.Index(title, "TTY")
		if ins_i > 0 {
			ins_i -= 1
		} else {
			ins_i = strings.LastIndex(title, " ")
			if ins_i == -1 {
				ins_i = 0
			}
		}
	}
	return	ins_i
}

func main() {
	args := strings.Join(os.Args[1:], " ")
	if args == "" {
		args = "aux"
	}

	out, _ := exec.Command("ps", args).Output()
	ll := strings.Split(string(out[:]), "\n")
	title, body := ll[0], ll[1:]

	cr := calc_ps(title, body)

	cx, _ := get_termxy()
	cx--

	ins_i := get_ins_idx(title)

	fmt.Printf("%.*s\n", cx,
		fmt.Sprintf("%s %*s %*s %s",
			title[:ins_i], cr.max_psw, "PSS", cr.max_usw, "USS", title[ins_i:]))

	for i, l := range body {
		if cr.psl[i] == "" {
			continue
		}
		s := fmt.Sprintf("%s %*s %*s %s",
			l[:ins_i], cr.max_psw, cr.psl[i], cr.max_usw, cr.usl[i], l[ins_i:])
		fmt.Printf("%.*s\n", cx + get_wide_num(s), s);
	}

	fmt.Printf("Total:VSZ PSS USS RSS ")
	fmt.Printf("%s ", num2strex(cr.vss_sum))
	fmt.Printf("%s ", num2strex(cr.pss_sum))
	fmt.Printf("%s ", num2strex(cr.uss_sum))
	fmt.Printf("%s ", num2strex(cr.rss_sum))
	fmt.Printf("\n")
}

