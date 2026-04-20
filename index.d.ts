/**
 * @package @titanpl/surface
 * 
 * Surface is an ultra-optimized high level modules provider native extension for the TitanPl framework,
 * handle data-heavy, IO-bound, and system-level tasks outside the JavaScript runtime 
 * for maximum performance and stability.
 */

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
}

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
}

export namespace json {
  export interface OpenOptions {
    /** format: 'auto' | 'json' | 'jsonl' (default: 'auto') */
    format?: 'auto' | 'json' | 'jsonl';
    /** extraction path (e.g. "users[*].email") */
    fpath?: string;
  }

  export interface NextOptions {
    size?: number;
  }

  export interface WriteOptions {
    /** Output format: 'json' (array) | 'jsonl' (line-delimited) */
    format?: 'json' | 'jsonl';
  }
}

/** Standard object containing all CSV utilities */
export declare const csv: {
  /**
   * Opens a CSV file for reading and starts the native pre-fetching process.
   * @param path Absolute or relative path to the .csv file
   * @param opts Configuration for the reader
   * @returns A native handler string to be used with other csv functions
   */
  open(path: string, opts?: csv.OpenOptions): string;

  /**
   * Reads a chunk of records from the native buffer.
   * @param handler The handler returned by csv.open()
   * @param opts Chunking configuration
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
   */
  create(path: string, opts: { headers: string[] }): string;

  /**
   * Writes rows to the CSV file.
   */
  write(handler: string, rows: any[]): void;
};

/** High-speed native SMTP utilities */
export declare const smtp: {
  /**
   * Sends a single email using the native Go engine.
   */
  send(opts: smtp.SendOptions): any;

  /**
   * Sends multiple emails concurrently using a native worker pool.
   */
  bulk(opts: smtp.BulkSendOptions): any[];

  /**
   * Renders a Go HTML template natively.
   */
  render(template: string, data?: any): string;

  /**
   * Reads and renders a Go template file directly from disk.
   */
  renderFile(path: string, data?: any): string;
};

/** Native high-performance image processing */
export declare const image: {
  /** 
   * Resizes an image natively.
   */
  resize(opts: image.ResizeOptions): { status: string; path?: string; base64?: string };

  /** 
   * Crops and fills an image to exact dimensions from the center.
   */
  crop(opts: image.CropOptions): { status: string; path?: string; base64?: string };

  /**
   * Complex pipeline processing.
   */
  process(opts: image.ProcessOptions): { status: string; path?: string; base64?: string };

  /**
   * Concurrent batch processing.
   */
  batch(opts: image.BatchOptions): any[];
};

/** Native JSON Streaming Module */
export declare const json: {
  /** Opens a JSON file for native streaming. Returns a native handler. */
  open(path: string, opts?: json.OpenOptions): string;
  /** Fetches the next chunk of JSON records. */
  next(handler: string, opts?: json.NextOptions): { rows: any[], done: boolean };
  /** Closes the JSON stream handler. */
  close(handler: string): void;
  /** Creates a new JSON/JSONL file for streaming output. */
  create(path: string): string;
  /** Writes a record to the native JSON stream. */
  write(handler: string, data: any, opts?: json.WriteOptions): void;
  /** Ultra-fast native serialization. Returns a string. */
  stringify(data: any): string;
};

declare const _default: {
  csv: typeof csv;
  smtp: typeof smtp;
  image: typeof image;
  json: typeof json;
};

export default _default;
