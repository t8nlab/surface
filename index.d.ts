/**
 * @package @titanpl/surface
 * 
 * Surface is a high-performance and top-level utils provider native extension for the TitanPl framework,
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

/** Standard object containing all CSV utilities */
export declare const csv: {
  /**
   * Opens a CSV file for reading and starts the native pre-fetching process.
   * @param path Absolute or relative path to the .csv file
   * @param opts Configuration for the reader
   * @returns A native handle string to be used with other csv functions
   * 
   * @example
   * const h = csv.open("data.csv", { mode: "object" });
   * try {
   *   let done = false;
   *   while (!done) {
   *     const chunk = csv.next(h, { size: 1000 });
   *     // process chunk.rows here
   *     done = chunk.done;
   *   }
   * } finally {
   *   csv.close(h);
   * }
   */
  open(path: string, opts?: csv.OpenOptions): string;

  /**
   * Reads a chunk of records from the native buffer.
   * This is extremely fast as Go pre-fetches records in the background.
   * @param handle The handle returned by csv.open()
   * @param opts Chunking configuration
   */
  next<T = any>(handle: string, opts?: csv.NextOptions): csv.Chunk<T>;

  /**
   * Reads the entire remaining contents of the CSV file in one Go call.
   * Moves the iteration loop into the native layer for zero JS overhead.
   * @param handle The handle returned by csv.open()
   * @returns All records in the format specified by the mode
   */
  readAll<T = any>(handle: string): T[];

  /**
   * Closes the file handle and stops the native pre-fetching goroutine.
   * Always call this in a finally block to prevent memory leaks.
   * @param handle The handle to close
   */
  close(handle: string): void;

  /**
   * Creates or overwrites a CSV file for writing.
   * @param path Destination path
   * @param opts Configuration including headers
   */
  create(path: string, opts: { headers: string[] }): string;

  /**
   * Writes rows to the CSV file.
   * @param handle The handle returned by csv.create()
   * @param rows Array of objects matching the headers
   */
  write(handle: string, rows: any[]): void;
};

/** High-speed native SMTP utilities */
export declare const smtp: {
  /**
   * Sends a single email using the native Go engine.
   * Supports both 587 (STARTTLS) and 465 (Direct SSL) ports.
   */
  send(opts: smtp.SendOptions): any;

  /**
   * Sends multiple emails concurrently using a native worker pool.
   * Perfect for high-speed newsletters or bulk notifications.
   */
  bulk(opts: smtp.BulkSendOptions): any[];

  /**
   * Renders a Go HTML template natively.
   * Uses standard Go {{.field}} syntax for data injection.
   */
  render(template: string, data?: any): string;

  /**
   * Reads and renders a Go template file directly from disk.
   * Much faster as it bypasses the JS file system layer.
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
   * Executes multiple steps (resize, crop, grayscale, etc.) in a single pass.
   */
  process(opts: image.ProcessOptions): { status: string; path?: string; base64?: string };

  /**
   * Concurrent batch processing.
   * Processes multiple images tasks in parallel using a native worker pool.
   */
  batch(opts: image.BatchOptions): any[];
};

declare const _default: {
  csv: typeof csv;
  smtp: typeof smtp;
  image: typeof image;
};

export default _default;
