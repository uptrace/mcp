package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status      int
	written     int
	buf         bytes.Buffer
	maxBodySize int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n := len(p)
	w.written += n
	if w.buf.Len() < w.maxBodySize {
		remain := w.maxBodySize - w.buf.Len()
		if remain > 0 {
			w.buf.Write(p[:min(len(p), remain)])
		}
	}
	return w.ResponseWriter.Write(p)
}

func extractMethod(body string) string {
	var req struct {
		Method string `json:"method"`
	}
	json.Unmarshal([]byte(body), &req)
	return req.Method
}

func unescapeJSON(data string) any {
	var parsed any
	if err := json.Unmarshal([]byte(data), &parsed); err != nil {
		return data
	}

	if str, ok := parsed.(string); ok {
		return unescapeJSON(str)
	}

	return parsed
}

func prettyJSON(data string) string {
	if strings.HasPrefix(data, "event:") || strings.Contains(data, "\nevent:") {
		return prettySSE(data)
	}

	parsed := unescapeJSON(data)
	formatted, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return data
	}
	return string(formatted)
}

func prettySSE(data string) string {
	lines := strings.Split(data, "\n")
	result := make(map[string]any)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if key == "data" {
				result[key] = unescapeJSON(value)
			} else {
				result[key] = value
			}
		}
	}

	formatted, _ := json.MarshalIndent(result, "", "  ")
	return string(formatted)
}

func loggingMiddleware(
	next http.Handler,
	maxBody int,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqBody, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(reqBody))

		lw := &loggingResponseWriter{
			ResponseWriter: w,
			maxBodySize:    maxBody,
		}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)
		method := extractMethod(string(reqBody))

		timeStr := color.HiBlackString(time.Now().Format("15:04:05.000"))
		methodStr := color.CyanString("%-20s", method)
		statusStr := color.GreenString("%d", lw.status)
		if lw.status >= 400 {
			statusStr = color.RedString("%d", lw.status)
		}
		durationStr := color.MagentaString("%6s", duration.Round(time.Microsecond))

		fmt.Fprintf(os.Stderr, "%s %s %s %s\n",
			timeStr, methodStr, statusStr, durationStr)

		if len(reqBody) > 0 {
			prettyReq := prettyJSON(string(reqBody))
			fmt.Fprintf(os.Stderr, "%s\n%s\n\n",
				color.CyanString("Request:"),
				indent(prettyReq, "  "))
		}

		if lw.buf.Len() > 0 {
			prettyResp := prettyJSON(lw.buf.String())
			fmt.Fprintf(os.Stderr, "%s\n%s\n\n",
				color.MagentaString("Response:"),
				indent(prettyResp, "  "))
		}
	})
}

func indent(s string, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}
