package sfJson

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	sfCsv "github.com/t8nlab/surface/csv"
	sfInput "github.com/t8nlab/surface/input"
)

type JsonStream struct {
	mu      sync.Mutex
	decoder *json.Decoder
	scanner *bufio.Scanner
	file    *os.File
	format  string
	path    string
	mode    string
	records chan []byte
	stop    chan struct{}
	done    bool
}

type JsonWriter struct {
	mu     sync.Mutex
	file   *os.File
	writer *bufio.Writer
	count  int
	format string
}

var (
	streams      = make(map[string]*JsonStream)
	writers      = make(map[string]*JsonWriter)
	mu           sync.Mutex
	handlerCount int
)

func generateHandler() string {
	handlerCount++
	return fmt.Sprintf("json_%d", handlerCount)
}

func JsonOpen(input map[string]any) (any, error) {
	path, err := sfInput.GetString(input, "path")
	if err != nil { return nil, err }

	fPath, _ := sfInput.GetString(input, "fpath")
	format, _ := sfInput.GetString(input, "format")
	mode, _ := sfInput.GetString(input, "mode")
	if mode == "" { mode = "object" }
	if format == "" { format = "auto" }

	file, err := os.Open(path)
	if err != nil { return nil, err }

	if format == "auto" || format == "" {
		if strings.HasSuffix(path, ".jsonl") || strings.HasSuffix(path, ".ndjson") {
			format = "jsonl"
		} else {
			format = "json"
		}
	}

	mu.Lock()
	handler := generateHandler()
	s := &JsonStream{
		file:    file,
		format:  format,
		path:    fPath,
		mode:    mode,
		records: make(chan []byte, 1000),
		stop:    make(chan struct{}),
	}

	if format == "jsonl" {
		s.scanner = bufio.NewScanner(file)
	} else {
		s.decoder = json.NewDecoder(file)
	}

	streams[handler] = s
	mu.Unlock()

	go s.preFetch()
	return map[string]string{"handler": handler}, nil
}

func (s *JsonStream) preFetch() {
	defer close(s.records)
	defer s.file.Close()

	if s.format == "jsonl" {
		for s.scanner.Scan() {
			line := s.scanner.Bytes()
			if len(line) == 0 { continue }
			select {
			case s.records <- line:
			case <-s.stop:
				return
			}
		}
	} else {
		if s.path != "" {
			parts := strings.Split(s.path, ".")
			s.walkTo(parts)
		} else {
			// Find array root
			for {
				token, err := s.decoder.Token()
				if err != nil { return }
				if delim, ok := token.(json.Delim); ok && delim == '[' {
					s.streamArray()
					return
				}
			}
		}
	}

	s.mu.Lock()
	s.done = true
	s.mu.Unlock()
}

func (s *JsonStream) walkTo(parts []string) {
	for i, part := range parts {
		isWildcard := strings.Contains(part, "[*]")
		cleanPart := strings.ReplaceAll(part, "[*]", "")
		for s.decoder.More() {
			t, err := s.decoder.Token()
			if err != nil { return }
			if t == cleanPart {
				if isWildcard || i == len(parts)-1 {
					nextT, _ := s.decoder.Token()
					if delim, ok := nextT.(json.Delim); ok && delim == '[' {
						s.streamArray()
						return
					} else {
						var raw json.RawMessage
						s.decoder.Decode(&raw)
						s.records <- raw
						return
					}
				}
				break
			}
		}
	}
}

func (s *JsonStream) streamArray() {
	for s.decoder.More() {
		var raw json.RawMessage
		if err := s.decoder.Decode(&raw); err != nil { break }
		select {
		case s.records <- raw:
		case <-s.stop:
			return
		}
	}
}

func JsonNext(input map[string]any) (any, error) {
	handler, _ := sfInput.GetString(input, "handler")
	size, _ := sfInput.GetInt(input, "size")
	if size <= 0 { size = 100 }

	mu.Lock()
	s, ok := streams[handler]
	mu.Unlock()
	if !ok { return nil, errors.New("invalid handler") }

	var buf bytes.Buffer
	buf.WriteByte('[')
	count := 0
	
	timeout := time.After(3 * time.Second)

	for i := 0; i < size; i++ {
		select {
		case item, ok := <-s.records:
			if !ok {
				buf.WriteByte(']')
				return map[string]any{"rows": json.RawMessage(buf.Bytes()), "done": true}, nil
			}
			if count > 0 { buf.WriteByte(',') }
			buf.Write(item)
			count++
		case <-timeout:
			// If we still have nothing after 3s, return what we have (empty array)
			i = size 
		}
	}
	buf.WriteByte(']')

	return map[string]any{
		"rows": json.RawMessage(buf.Bytes()),
		"done": false,
	}, nil
}

func JsonClose(input map[string]any) (any, error) {
	handler, _ := sfInput.GetString(input, "handler")
	mu.Lock()
	defer mu.Unlock()
	if s, ok := streams[handler]; ok {
		close(s.stop)
		delete(streams, handler)
	}
	if w, ok := writers[handler]; ok {
		w.mu.Lock()
		if w.format != "jsonl" && w.count > 0 {
			w.writer.WriteByte(']')
		}
		w.writer.Flush()
		w.file.Close()
		w.mu.Unlock()
		delete(writers, handler)
	}
	return map[string]any{"success": true}, nil
}

func JsonStringifyFast(input map[string]any) (any, error) {
	data, _ := input["data"]
	res, _ := json.Marshal(data)
	return map[string]any{"json": string(res)}, nil
}

func JsonCreate(input map[string]any) (any, error) {
	path, _ := sfInput.GetString(input, "path")
	format, _ := sfInput.GetString(input, "format")
	file, err := os.Create(path)
	if err != nil { return nil, err }
	mu.Lock()
	handler := generateHandler()
	writers[handler] = &JsonWriter{
		file:   file, 
		writer: bufio.NewWriter(file),
		format: format,
	}
	mu.Unlock()
	return map[string]any{"handler": handler}, nil
}

func JsonWrite(input map[string]any) (any, error) {
	handler, _ := sfInput.GetString(input, "handler")
	data, _ := input["data"]
	
	mu.Lock()
	w, ok := writers[handler]
	mu.Unlock()
	if !ok { return nil, errors.New("invalid writer") }

	w.mu.Lock()
	defer w.mu.Unlock()
	encoded, _ := json.Marshal(data)
	if w.format == "jsonl" {
		w.writer.Write(encoded)
		w.writer.WriteByte('\n')
	} else {
		if w.count == 0 { w.writer.WriteByte('[') } else { w.writer.WriteByte(',') }
		w.writer.Write(encoded)
	}
	
	w.count++
	return map[string]any{"success": true}, nil
}

func JsonToCsv(input map[string]any) (any, error) {
	jsonPath, _ := sfInput.GetString(input, "path")
	csvPath, _ := sfInput.GetString(input, "out")
	fPath, _ := sfInput.GetString(input, "fpath")

	// 1. Open JSON Stream
	res, err := JsonOpen(map[string]any{"path": jsonPath, "fpath": fPath})
	if err != nil { return nil, err }
	
	resMap, ok := res.(map[string]string)
	if !ok { return nil, errors.New("failed to initialize json reader") }
	handler := resMap["handler"]

	// 2. PATIENT READ
	chunk, err := JsonNext(map[string]any{"handler": handler, "size": 1})
	if err != nil { return nil, err }
	
	chunkMap, ok := chunk.(map[string]any)
	if !ok { return nil, errors.New("failed to read json chunk") }
	
	rowsRaw := chunkMap["rows"].(json.RawMessage)
	
	var firstRows []map[string]any
	json.Unmarshal(rowsRaw, &firstRows)
	if len(firstRows) == 0 { 
		JsonClose(map[string]any{"handler": handler})
		return nil, errors.New("empty json") 
	}

	var headers []string
	for k := range firstRows[0] { headers = append(headers, k) }

	// Convert []string to []any for sfCsv.CsvCreate compatibility
	headersAny := make([]any, len(headers))
	for i, h := range headers { headersAny[i] = h }

	// 3. Open CSV Writer
	csvRes, err := sfCsv.CsvCreate(map[string]any{"path": csvPath, "headers": headersAny})
	if err != nil { 
		JsonClose(map[string]any{"handler": handler})
		return nil, fmt.Errorf("failed to create csv file: %v", err) 
	}
	
	csvResMap, ok := csvRes.(map[string]string)
	if !ok { return nil, errors.New("failed to initialize csv writer") }
	ch := csvResMap["handler"]

	// 4. Stream Loop: Write the header row(s) first
	rowsAny := make([]any, len(firstRows))
	for i, r := range firstRows { rowsAny[i] = r }
	
	_, err = sfCsv.CsvWrite(map[string]any{"handler": ch, "rows": rowsAny})
	if err != nil { return nil, fmt.Errorf("failed to write first row to csv: %v", err) }

	done := false
	for !done {
		nextChunk, err := JsonNext(map[string]any{"handler": handler, "size": 1000})
		if err != nil { break }
		
		ncMap := nextChunk.(map[string]any)
		done = ncMap["done"].(bool)
		
		var rows []map[string]any
		json.Unmarshal(ncMap["rows"].(json.RawMessage), &rows)
		if len(rows) > 0 {
			// Convert to []any
			batchAny := make([]any, len(rows))
			for i, r := range rows { batchAny[i] = r }
			
			_, err = sfCsv.CsvWrite(map[string]any{"handler": ch, "rows": batchAny})
			if err != nil { return nil, fmt.Errorf("failed to write batch to csv: %v", err) }
		}
	}

	JsonClose(map[string]any{"handler": handler})
	sfCsv.CsvClose(map[string]any{"handler": ch})

	return map[string]any{"success": true, "out": csvPath}, nil
}
