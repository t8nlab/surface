/**
 * @package @titanpl/surface
 * 
 * Surface is an ultra-optimized high level modules provider native extension for the TitanPl framework,
 * handle data-heavy, IO-bound, and system-level tasks outside the JavaScript runtime 
 * for maximum performance and stability.
 */

export namespace csv {
  /**
   * Mode for CSV data representation
   * - 'object': Array of objects with headers as keys
   * - 'row' | 'raw' | 'rows': Array of arrays for maximum speed
   * - 'column': Data organized by columns
   */
  export type Mode = 'object' | 'row' | 'raw' | 'rows' | 'column';

  export interface OpenOptions {
    /** Whether the CSV has a header row (default: true) */
    header?: boolean;
    /** Automatically infer types like numbers and booleans (default: false) */
    inferTypes?: boolean;
    /** Output format mode (default: 'object') */
    mode?: Mode;
    /** Custom delimiter character (default: ',') */
    delimiter?: string;
    /** List of specific column names to select */
    select?: string[];
  }

  export interface NextOptions {
    /** Number of records to read in this chunk (default: 100) */
    size?: number;
  }

  export interface CreateOptions {
    /** Column headers for the CSV file */
    headers: string[];
    /** Custom delimiter character (default: ',') */
    delimiter?: string;
  }

  export interface Chunk<T = any> {
    /** The actual records (type depends on the mode) */
    rows: T[];
    /** Whether the end of the file has been reached */
    done: boolean;
    /** The mode used for this chunk */
    mode: Mode;
    /** Headers present in the CSV (only for raw/rows mode) */
    headers?: string[];
  }
}

export namespace smtp {
  export interface SendOptions {
    host: string;
    port: number;
    username: string;
    password: string;
    from?: string;
    to?: string;
    /** Carbon Copy recipients */
    cc?: string;
    /** Blind Carbon Copy recipients */
    bcc?: string;
    subject?: string;
    body: string;
    /** Force Implicit TLS (SSL) for port 465 (default: false) */
    ssl?: boolean;
    /** Enable Raw Mode: Scrapes From/To/Subject directly from the body headers (default: false) */
    raw?: boolean;
  }

  export interface BulkSendOptions extends Partial<SendOptions> {
    /** Array of email objects to send concurrently */
    emails: Partial<SendOptions>[];
    /** Max concurrent connections (default: 5) */
    concurrency?: number;
  }

  export interface Result {
    /** Success status */
    success: boolean;
    /** Message ID or info */
    message?: string;
    /** Error details if any */
    error?: string;
  }
}

export namespace image {
  export interface Step {
    action: 'resize' | 'crop' | 'grayscale' | 'blur';
    width?: number;
    height?: number;
    sigma?: number;
  }

  export interface ProcessOptions {
    /** Source: path, URL, or data:image/base64 */
    src: string;
    /** Destination file path. If omitted, returns Base64. */
    out?: string;
    /** Output format: 'jpg', 'png', 'webp' */
    format?: 'jpg' | 'png' | 'webp';
    /** Output quality (1-100) */
    quality?: number;
    /** Array of operations to perform in order */
    steps?: Step[];
  }

  export interface BatchOptions {
    /** Array of process options for each image */
    items: ProcessOptions[];
    /** Max parallel workers (default 4) */
    concurrency?: number;
  }

  export interface ResizeOptions {
    src: string;
    out?: string;
    width?: number;
    height?: number;
    quality?: number;
    format?: 'jpg' | 'png' | 'webp';
  }

  export interface CropOptions {
    src: string;
    out?: string;
    width: number;
    height: number;
    format?: 'jpg' | 'png' | 'webp';
  }

  export interface ProcessResult {
    status: 'success' | 'error';
    path?: string;
    base64?: string;
    error?: string;
  }
}

export namespace json {
  export interface OpenOptions {
    /** format: 'auto' | 'json' | 'jsonl' (default: 'auto') */
    format?: 'auto' | 'json' | 'jsonl';
    /** extraction path (e.g. "users[*].email") */
    fpath?: string;
  }

  export interface NextOptions {
    /** Number of records to read in this chunk (default: 100) */
    size?: number;
  }

  export interface WriteOptions {
    /** Output format: 'json' (array) | 'jsonl' (line-delimited) */
    format?: 'json' | 'jsonl';
  }

  export interface ToCsvOptions {
    /** JSON extraction path (e.g. "data.items[*]") */
    fpath?: string;
    /** CSV delimiter (default: ',') */
    delimiter?: string;
    /** Whether to include header row (default: true) */
    header?: boolean;
    /** Map of JSON keys to CSV column names or simple array of keys */
    columns?: Record<string, string> | string[];
  }
}


  /** Native Data Validation & Cleaning */
  export const clean: {
    /**
     * Natively validates emails in a file using regex and structural checks.
     */
    validateEmails(path: string): { valid: number; invalid: number };

    /**
     * Normalizes phone numbers to E.164 format.
     * Supports array of strings OR a file transformation task.
     */
    normalizePhones(input: string[] | { path: string; out: string }): any;

    /**
     * Natively removes duplicate rows from a file.
     */
    removeDuplicates(src: string, out: string): { processed: number; duplicates: number; success: boolean };

    /**
     * Complete data cleaning in a single native pass.
     */
    process(opts: { 
      src: string; 
      out: string; 
      normalize?: boolean; 
      dedup?: boolean; 
      concurrency?: number;
      phoneFields?: number[]; // Specify column indices to treat as phones
    }): { processed: number; duplicates: number; invalidEmails: number; workers: number; success: boolean };
  };

/** Native Web Extraction */
export declare const extract: {
  /**
   * Fetches raw HTML from a URL natively.
   */
  html(url: string): string;

  /**
   * Extracts all links (hrefs) from a URL natively.
   */
  links(url: string): string[];

  /**
   * Extracts SEO/OpenGraph metadata from a URL natively.
   */
  meta(url: string): Record<string, string>;
};

/** Standard object containing all CSV utilities */
export declare const csv: {
  /**
   * Opens a CSV file for reading and starts the native pre-fetching process.
   * Supports local file paths and cloud URLs (HTTP/HTTPS).
   * @param path Absolute path, relative path, or a public URL to the .csv file
   * @param opts Configuration for the reader (mode, delimiter, etc.)
   * @returns A native handler string to be used with other csv functions
   */
  open(path: string, opts?: csv.OpenOptions): string;

  /**
   * Reads a chunk of records from the native buffer.
   * @param handler The handler returned by csv.open()
   * @param opts Chunking configuration (chunk size, etc.)
   * @returns A chunk object containing rows and completion status
   */
  next<T = any>(handler: string, opts?: csv.NextOptions): csv.Chunk<T>;

  /**
   * Reads the entire remaining contents of the CSV file in one Go call.
   * @param handler The handler returned by csv.open()
   * @returns All records in the format specified by the mode
   */
  readAll<T = any>(handler: string): T[];

  /**
   * Closes the file handler and stops the native pre-fetching goroutine.
   * @param handler The handler to close
   */
  close(handler: string): void;

  /**
   * Creates or overwrites a CSV file for writing.
   * @param path Target file path
   * @param opts Configuration with headers and delimiter
   * @returns A native handler string for writing
   */
  create(path: string, opts: csv.CreateOptions): string;

  /**
   * Writes rows to an open CSV file.
   * @param handler The handler returned by csv.create()
   * @param rows Array of arrays or array of objects to write
   */
  write(handler: string, rows: any[]): void;
};

/** High-speed native SMTP utilities */
export declare const smtp: {
  /**
   * Sends a single email using the native Go engine.
   * @param opts Mail configuration (host, credentials, to, from, body, etc.)
   * @returns Result object with status
   */
  send(opts: smtp.SendOptions): smtp.Result;

  /**
   * Sends multiple emails concurrently using a native worker pool.
   * @param opts Bulk configuration including an array of emails and concurrency
   * @returns Array of result objects
   */
  bulk(opts: smtp.BulkSendOptions): smtp.Result[];

  /**
   * Renders a Go HTML template natively.
   * @param template The template string content
   * @param data Optional context data for the template
   * @returns The rendered HTML string
   */
  render(template: string, data?: any): string;

  /**
   * Reads and renders a Go template file directly from disk.
   * @param path Path to the .html template file
   * @param data Optional context data for the template
   * @returns The rendered HTML string
   */
  renderFile(path: string, data?: any): string;
};

/** Native high-performance image processing */
export declare const image: {
  /** 
   * Resizes an image natively.
   * @param opts Source, output path, dimensions and quality
   */
  resize(opts: image.ResizeOptions): image.ProcessResult;

  /** 
   * Crops and fills an image to exact dimensions from the center.
   * @param opts Source, output path, and target dimensions
   */
  crop(opts: image.CropOptions): image.ProcessResult;

  /**
   * Executes multiple image operations (resize, crop, blur, etc.) in a single native pass.
   * @param opts Source, output, and ordered list of processing steps
   */
  process(opts: image.ProcessOptions): image.ProcessResult;

  /**
   * Processes multiple images concurrently using a native worker pool.
   * @param opts Array of processing items and concurrency level
   */
  batch(opts: image.BatchOptions): image.ProcessResult[];
};

/** Native JSON Streaming Module */
export declare const json: {
  /** 
   * Opens a JSON/JSONL file for native streaming.
   * @param path Path to the JSON file
   * @param opts Configuration (format, extraction path)
   * @returns A native handler string
   */
  open(path: string, opts?: json.OpenOptions): string;

  /** 
   * Fetches the next chunk of JSON records from a stream.
   * @param handler The handler returned by json.open()
   * @param opts Chunking size
   * @returns Object with rows and done status
   */
  next(handler: string, opts?: json.NextOptions): { rows: any[], done: boolean };
  
  /** 
   * Reads the entire remaining contents of the JSON stream in one Go call.
   * @param handler The handler returned by json.open()
   * @returns All records in an array
   */
  readAll<T = any>(handler: string): T[];

  /** 
   * Closes the JSON stream handler and releases resources.
   * @param handler The handler to close
   */
  close(handler: string): void;

  /** 
   * Creates a new JSON or JSONL file for streaming output.
   * @param path Target file path
   * @returns A native handler string
   */
  create(path: string): string;

  /** 
   * Writes a single record to the native JSON stream.
   * @param handler The handler returned by json.create()
   * @param data The object or data to serialize
   * @param opts Output formatting (json vs jsonl)
   */
  write(handler: string, data: any, opts?: json.WriteOptions): void;

  /** 
   * Ultra-fast native serialization of JavaScript objects to JSON strings.
   * @param data The object to stringify
   * @returns Valid JSON string
   */
  stringify(data: any): string;

  /**
   * Native Cross-Engine Bridge: Streams JSON data directly into a CSV file.
   * Highly efficient conversion with zero JavaScript memory overhead for large files.
   * @param jsonPath Source JSON file path
   * @param csvPath Destination CSV file path
   * @param opts Conversion options (extraction path, delimiter, etc.)
   */
  toCSV(jsonPath: string, csvPath: string, opts?: json.ToCsvOptions): any;
};

export namespace http {
  export interface Config {
    url?: string;
    method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | string;
    headers?: Record<string, any>;
    params?: Record<string, any>;
    data?: any;
    /** Timeout in milliseconds (default: 30000) */
    timeout?: number;
  }

  export interface Response<T = any> {
    status: number;
    statusText: string;
    headers: Record<string, string>;
    data: T;
    url: string;
    ok: boolean;
  }
}

/** Native HTTP Client */
export declare const http: {
  /** Perform a custom HTTP request */
  request<T = any>(config: http.Config): http.Response<T>;
  request<T = any>(url: string, config?: http.Config): http.Response<T>;
  
  /** Shorthand for GET request */
  get<T = any>(url: string, config?: http.Config): http.Response<T>;
  
  /** Shorthand for POST request */
  post<T = any>(url: string, data?: any, config?: http.Config): http.Response<T>;
  
  /** Shorthand for PUT request */
  put<T = any>(url: string, data?: any, config?: http.Config): http.Response<T>;
  
  /** Shorthand for DELETE request */
  delete<T = any>(url: string, config?: http.Config): http.Response<T>;
  
  /** Shorthand for PATCH request */
  patch<T = any>(url: string, data?: any, config?: http.Config): http.Response<T>;
};

declare const _default: {
  csv: typeof csv;
  smtp: typeof smtp;
  image: typeof image;
  json: typeof json;
  clean: typeof clean;
  extract: typeof extract;
  http: typeof http;
};

export default _default;
