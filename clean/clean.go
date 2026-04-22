package sfClean

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	sfInput "github.com/t8nlab/surface/input"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2}\.?[a-z]{0,4}$`)
	phoneRegex = regexp.MustCompile(`[^\d+]`)
)

type CleanJob struct {
	Record []string
	Object map[string]any
	IsJSON bool
}

type CleanResult struct {
	Record []string
	Object map[string]any
	IsDup  bool
}

// Process handles massive multi-threaded data scrubbing for CSV and JSON
func Process(input map[string]any) (any, error) {
	srcPath, _ := input["src"].(string)
	if srcPath == "" {
		srcPath, _ = input["path"].(string)
	}
	outPath, _ := input["out"].(string)
	doNormalize, _ := input["normalize"].(bool)
	doDedup, _ := input["dedup"].(bool)

	if srcPath == "" { return nil, errors.New("source path is required") }

	isJSON := strings.ToLower(filepath.Ext(srcPath)) == ".json"

	sf, err := os.Open(srcPath)
	if err != nil { return nil, err }
	defer sf.Close()

	var df *os.File
	if outPath != "" {
		df, err = os.Create(outPath)
		if err != nil { return nil, err }
		defer df.Close()
	}

	// Concurrency
	numWorkers, _ := sfInput.GetInt(input, "concurrency")
	if numWorkers <= 0 { numWorkers = runtime.NumCPU() }

	// Targeted Fields (Experimental)
	phoneFields, _ := input["phoneFields"].([]any)
	isFieldTargeted := func(idx int) bool {
		if len(phoneFields) == 0 { return true }
		for _, f := range phoneFields {
			val, _ := f.(float64)
			if int(val) == idx { return true }
		}
		return false
	}

	jobs := make(chan CleanJob, 100)
	results := make(chan CleanResult, 100)
	var wg sync.WaitGroup
	var invalidEmails, duplicates, total int64
	seen := sync.Map{}

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				isDup := false
				if isJSON {
					processObject(job.Object, doNormalize, &invalidEmails)
					if doDedup {
						sig := fmt.Sprintf("%v", job.Object)
						if _, exists := seen.LoadOrStore(sig, true); exists { isDup = true }
					}
				} else {
					processRecord(job.Record, doNormalize, &invalidEmails, isFieldTargeted)
					if doDedup {
						sig := strings.Join(job.Record, "|")
						if _, exists := seen.LoadOrStore(sig, true); exists { isDup = true }
					}
				}
				results <- CleanResult{Record: job.Record, Object: job.Object, IsDup: isDup}
			}
		}()
	}

	// Collector
	done := make(chan bool)
	go func() {
		var csvWriter *csv.Writer
		if df != nil && !isJSON {
			csvWriter = csv.NewWriter(df)
			defer csvWriter.Flush()
		}
		
		isFirst := true
		if isJSON && df != nil { df.WriteString("[\n") }

		for res := range results {
			atomic.AddInt64(&total, 1)
			if res.IsDup {
				atomic.AddInt64(&duplicates, 1)
				continue
			}
			if df != nil {
				if isJSON {
					if !isFirst { df.WriteString(",\n") }
					b, _ := json.Marshal(res.Object)
					df.Write(b)
					isFirst = false
				} else if csvWriter != nil {
					csvWriter.Write(res.Record)
				}
			}
		}
		if isJSON && df != nil { df.WriteString("\n]") }
		done <- true
	}()

	// Feeder
	if isJSON {
		dec := json.NewDecoder(sf)
		if t, _ := dec.Token(); t == json.Delim('[') {
			for dec.More() {
				var obj map[string]any
				if err := dec.Decode(&obj); err == nil {
					jobs <- CleanJob{Object: obj, IsJSON: true}
				}
			}
		}
	} else {
		reader := csv.NewReader(sf)
		for {
			rec, err := reader.Read()
			if err == io.EOF { break }
			if err != nil { continue }
			jobs <- CleanJob{Record: rec, IsJSON: false}
		}
	}

	close(jobs)
	wg.Wait()
	close(results)
	<-done

	return map[string]any{
		"processed": total,
		"duplicates": duplicates,
		"invalidEmails": invalidEmails,
		"workers": numWorkers,
		"format": strings.TrimPrefix(filepath.Ext(srcPath), "."),
		"success": true,
	}, nil
}

func processRecord(rec []string, normalizeFields bool, invalidEmails *int64, isTargeted func(int) bool) {
	for i, field := range rec {
		// Basic trim always, but deep normalization collapses internal spaces
		var cleaned string
		if normalizeFields {
			cleaned = strings.Join(strings.Fields(field), " ")
		} else {
			cleaned = strings.TrimSpace(field)
		}

		if strings.Contains(cleaned, "@") {
			cleaned = strings.ToLower(cleaned)
			if !emailRegex.MatchString(cleaned) {
				atomic.AddInt64(invalidEmails, 1)
			}
		} else if isTargeted(i) && isLikelyPhone(cleaned) {
			cleaned = normalize(cleaned)
		}

		if normalizeFields {
			rec[i] = cleaned
		}
	}
}

func processObject(obj map[string]any, normalizeFields bool, invalidEmails *int64) {
	for k, v := range obj {
		if s, ok := v.(string); ok {
			var cleaned string
			if normalizeFields {
				cleaned = strings.Join(strings.Fields(s), " ")
			} else {
				cleaned = strings.TrimSpace(s)
			}

			if strings.Contains(cleaned, "@") {
				cleaned = strings.ToLower(cleaned)
				if !emailRegex.MatchString(cleaned) {
					atomic.AddInt64(invalidEmails, 1)
				}
			} else if isLikelyPhone(cleaned, k) {
				cleaned = normalize(cleaned)
			}

			if normalizeFields {
				obj[k] = cleaned
			}
		}
	}
}

func normalize(s string) string {
	clean := phoneRegex.ReplaceAllString(s, "")
	if !strings.HasPrefix(clean, "+") && len(clean) > 0 {
		clean = "+" + clean
	}
	return clean
}

func isLikelyPhone(s string, key ...string) bool {
	// If it's a JSON key, check for name
	if len(key) > 0 {
		k := strings.ToLower(key[0])
		if strings.Contains(k, "phone") || strings.Contains(k, "mobile") || strings.Contains(k, "tel") {
			return true
		}
	}

	// Heuristic: must have between 7-15 digits AND some phone symbols
	digits := 0
	hasSymbols := false
	for _, c := range s {
		if c >= '0' && c <= '9' {
			digits++
		} else if c == '+' || c == '(' || c == ')' || c == '-' || c == ' ' {
			hasSymbols = true
		}
	}
	
	// If it has a '+', it's very likely a phone
	if strings.HasPrefix(s, "+") { return digits >= 7 }
	
	// If it's just a raw number with no symbols, it might be an ID.
	// Only treat as phone if it's broad enough but risky.
	// Let's require at least one symbol if it's less than 10 digits
	if !hasSymbols && digits < 10 { return false }

	return digits >= 7 && digits <= 15
}

func ValidateEmails(input map[string]any) (any, error) {
	res, err := Process(input)
	if err != nil { return nil, err }
	m := res.(map[string]any)
	return map[string]any{
		"processed": m["processed"],
		"invalid":   m["invalidEmails"],
		"valid":     m["processed"].(int64) - m["invalidEmails"].(int64),
	}, nil
}
func NormalizePhones(input map[string]any) (any, error) {
	if raw, ok := input["phones"].([]any); ok {
		results := make([]string, 0, len(raw))
		for _, p := range raw {
			if s, ok := p.(string); ok { results = append(results, normalize(s)) }
		}
		return results, nil
	}
	input["normalize"] = true
	res, err := Process(input)
	if err != nil { return nil, err }
	m := res.(map[string]any)
	return map[string]any{"success": true, "processed": m["processed"]}, nil
}
func RemoveDuplicates(input map[string]any) (any, error) {
	input["dedup"] = true
	res, err := Process(input)
	if err != nil { return nil, err }
	m := res.(map[string]any)
	return map[string]any{
		"processed": m["processed"],
		"duplicates": m["duplicates"],
		"saved": m["processed"].(int64) - m["duplicates"].(int64),
	}, nil
}
