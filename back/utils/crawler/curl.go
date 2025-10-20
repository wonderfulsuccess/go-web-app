package crawler

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func RunCURL(cURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/bash", "-lc", cURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := stdout.Bytes()
	ct := ""
	_, name, _ := charset.DetermineEncoding(out, ct)
	switch strings.ToLower(name) {
	case "gbk", "gb2312", "gb18030":
		ub, e := simplifiedchinese.GBK.NewDecoder().Bytes(out)
		if e == nil {
			return string(ub), err
		}
	}
	return string(out), err
}




