package sfCsv

import (
	"bytes"
	"encoding/json"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"sync"

	sfInput "github.com/t8nlab/surface/input"
)

type CsvStream struct {
	mu         sync.Mutex
	reader     *csv.Reader
	headers    []string
	selectMap  map[int]string
	inferTypes bool
	mode       string // "object", "column", "raw"
	file       io.ReadCloser

	// Async pre-fetching (Stores pre-marshaled JSON rows)
	records chan []byte
	done    bool
	err     error
	stop    chan struct{}
}

type CsvWriter struct {
	mu      sync.Mutex
	writer  *csv.Writer
	headers []string
	file    *os.File
}

var (
	streams      = make(map[string]*CsvStream)
	writers      = make(map[string]*CsvWriter)
	mu           sync.Mutex
	handlerCount int
)

func generateHandler() string {
	handlerCount++
	return fmt.Sprintf("csv_%d", handlerCount)
}

// --- Reader Functions ---

func CsvOpen(input map[string]any) (any, error) {
	path, err := sfInput.GetString(input, "path")
	if err != nil {
		return nil, err
	}

	header := sfInput.GetBool(input, "header", true)
	inferTypes := sfInput.GetBool(input, "inferTypes", false)
	
	// Better delimiter handling
	delimiter := ','
	if d, ok := input["delimiter"].(string); ok && len(d) > 0 {
		delimiter = rune(d[0])
	}

	file, err := sfInput.GetReader(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.Comma = delimiter

	var headers []string
	if header {
		headers, err = reader.Read()
		if err != nil {
			file.Close()
			return nil, err
		}
	}

	// Column selection
	var selectMap map[int]string
	if sel, ok := input["select"].([]any); ok {
		selectMap = make(map[int]string)
		for _, s := range sel {
			colName, ok := s.(string)
			if !ok {
				continue
			}
			if header {
				for i, h := range headers {
					if h == colName {
						selectMap[i] = colName
						break
					}
				}
			}
		}
	}

	// Determine Mode
	mode, _ := sfInput.GetString(input, "mode")
	if mode == "" {
		mode = "object"
	}

	mu.Lock()
	handler := generateHandler()
	s := &CsvStream{
		reader:     reader,
		headers:    headers,
		selectMap:  selectMap,
		inferTypes: inferTypes,
		mode:       mode,
		file:       file,
		records:    make(chan []byte, 2000), 
		stop:		make(chan struct{}),
	}
	streams[handler] = s
	mu.Unlock()

	go s.preFetch()

	return map[string]string{"handler": handler}, nil
}

func (s *CsvStream) preFetch() {
	defer close(s.records)

	mode := s.mode
	inferTypes := s.inferTypes
	headers := s.headers
	selectMap := s.selectMap

	type field struct {
		index int
		key   []byte
	}
	var fields []field

	if selectMap == nil {
		if mode == "object" {
			for i, h := range headers {
				fields = append(fields, field{
					index: i,
					key:   []byte(fmt.Sprintf("%q:", h)),
				})
			}
		}
	} else {
		var indices []int
		for i := range selectMap {
			indices = append(indices, i)
		}
		sort.Ints(indices)

		for _, idx := range indices {
			fields = append(fields, field{
				index: idx,
				key:   []byte(fmt.Sprintf("%q:", selectMap[idx])),
			})
		}
	}

	for {
		record, err := s.reader.Read()
		if err != nil {
			s.mu.Lock()
			if err != io.EOF {
				s.err = err
			}
			s.done = true
			s.mu.Unlock()
			return
		}

		var rowBuf bytes.Buffer
		if mode == "object" {
			rowBuf.WriteByte('{')
			written := 0
			for _, f := range fields {
				if f.index >= len(record) { continue }
				if written > 0 { rowBuf.WriteByte(',') }
				rowBuf.Write(f.key)
				val := record[f.index]
				if inferTypes {
					v, _ := json.Marshal(infer(val))
					rowBuf.Write(v)
				} else {
					v, _ := json.Marshal(val)
					rowBuf.Write(v)
				}
				written++
			}
			rowBuf.WriteByte('}')
		} else {
			rowBuf.WriteByte('[')
			written := 0
			if selectMap == nil {
				for i, val := range record {
					if i > 0 { rowBuf.WriteByte(',') }
					if inferTypes {
						v, _ := json.Marshal(infer(val))
						rowBuf.Write(v)
					} else {
						v, _ := json.Marshal(val)
						rowBuf.Write(v)
					}
				}
			} else {
				for _, f := range fields {
					if f.index >= len(record) { continue }
					if written > 0 { rowBuf.WriteByte(',') }
					val := record[f.index]
					if inferTypes {
						v, _ := json.Marshal(infer(val))
						rowBuf.Write(v)
					} else {
						v, _ := json.Marshal(val)
						rowBuf.Write(v)
					}
					written++
				}
			}
			rowBuf.WriteByte(']')
		}

		select {
		case s.records <- rowBuf.Bytes():
		case <-s.stop:
			return
		}
	}
}

func CsvNext(input map[string]any) (any, error) {
	handler, err := sfInput.GetString(input, "handler")
	if err != nil { return nil, err }

	size, _ := sfInput.GetInt(input, "size")
	if size <= 0 { size = 100 }

	mu.Lock()
	s, ok := streams[handler]
	mu.Unlock()

	if !ok {
		return nil, errors.New("invalid or closed handler: " + handler)
	}

	var combinedBuf bytes.Buffer
	combinedBuf.WriteByte('[')
	count := 0
	done := false
	
	for i := 0; i < size; i++ {
		select {
		case item, ok := <-s.records:
			if !ok {
				done = true
				break
			}
			if count > 0 { combinedBuf.WriteByte(',') }
			combinedBuf.Write(item)
			count++
		default:
			s.mu.Lock()
			isDone := s.done
			s.mu.Unlock()

			if count > 0 || isDone {
				i = size 
			} else {
				item, ok := <-s.records
				if !ok {
					done = true
					i = size
				} else {
					if count > 0 { combinedBuf.WriteByte(',') }
					combinedBuf.Write(item)
					count++
				}
			}
		}
		if done { break }
	}
	combinedBuf.WriteByte(']')

	s.mu.Lock()
	if (s.done && len(s.records) == 0) || count == 0 {
		done = s.done && len(s.records) == 0
	}
	s.mu.Unlock()

	return map[string]any{
		"done": done,
		"mode": s.mode,
		"rows": json.RawMessage(combinedBuf.Bytes()),
	}, nil
}

func CsvReadAll(input map[string]any) (any, error) {
	handler, err := sfInput.GetString(input, "handler")
	if err != nil { return nil, err }

	mu.Lock()
	s, ok := streams[handler]
	mu.Unlock()

	if !ok {
		return nil, errors.New("invalid or closed handler: " + handler)
	}

	var combinedBuf bytes.Buffer
	combinedBuf.WriteByte('[')
	count := 0
	for item := range s.records {
		if count > 0 { combinedBuf.WriteByte(',') }
		combinedBuf.Write(item)
		count++
	}
	combinedBuf.WriteByte(']')

	s.mu.Lock()
	s.done = true
	s.mu.Unlock()

	return json.RawMessage(combinedBuf.Bytes()), nil
}

func infer(v string) any {
	if i, err := strconv.Atoi(v); err == nil { return i }
	if f, err := strconv.ParseFloat(v, 64); err == nil { return f }
	if v == "true" || v == "false" { return v == "true" }
	return v
}

// --- Writer Functions ---

func CsvCreate(input map[string]any) (any, error) {
	path, err := sfInput.GetString(input, "path")
	if err != nil { return nil, err }

	headersAny, ok := input["headers"].([]any)
	if !ok { return nil, errors.New("headers must be an array of strings") }

	var headers []string
	for _, h := range headersAny {
		if s, ok := h.(string); ok { headers = append(headers, s) }
	}

	file, err := os.Create(path)
	if err != nil { return nil, err }

	writer := csv.NewWriter(file)
	if err := writer.Write(headers); err != nil {
		file.Close()
		return nil, err
	}

	mu.Lock()
	handler := generateHandler()
	writers[handler] = &CsvWriter{
		writer:  writer,
		headers: headers,
		file:    file,
	}
	mu.Unlock()

	return map[string]string{"handler": handler}, nil
}

func CsvWrite(input map[string]any) (any, error) {
	handler, err := sfInput.GetString(input, "handler")
	if err != nil { return nil, err }

	rows, ok := input["rows"].([]any)
	if !ok { return nil, errors.New("rows must be an array") }

	mu.Lock()
	w, ok := writers[handler]
	mu.Unlock()

	if !ok {
		return nil, errors.New("invalid or closed writer handler: " + handler)
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, r := range rows {
		rowMap, ok := r.(map[string]any)
		if !ok { continue }

		var record []string
		for _, h := range w.headers {
			if val, exists := rowMap[h]; exists {
				record = append(record, fmt.Sprintf("%v", val))
			} else {
				record = append(record, "")
			}
		}
		if err := w.writer.Write(record); err != nil { return nil, err }
	}

	w.writer.Flush()
	return true, nil
}

func CsvClose(input map[string]any) (any, error) {
	handler, err := sfInput.GetString(input, "handler")
	if err != nil { return nil, err }

	mu.Lock()
	defer mu.Unlock()

	if s, ok := streams[handler]; ok {
		s.mu.Lock()
		close(s.stop)
		s.file.Close()
		s.mu.Unlock()
		delete(streams, handler)
	}

	if w, ok := writers[handler]; ok {
		w.mu.Lock()
		w.writer.Flush()
		w.file.Close()
		w.mu.Unlock()
		delete(writers, handler)
	}

	return true, nil
}
