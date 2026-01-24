export type JobStatus = 'pending' | 'running' | 'completed' | 'failed';

export interface CheckRequest {
  url: string;
  languages?: string[];
  profile?: string;
  target?: string;
  skip_checks?: string[];
  timeout_secs?: number;
}

export interface CheckResponse {
  job_id: string;
  status: JobStatus;
  message: string;
}

export interface JSONResult {
  name: string;
  id: string;
  passed: boolean;
  status: 'pass' | 'warn' | 'fail' | 'info';
  message?: string;
  language?: string;
  duration_ms: number;
  raw_output?: string;
}

export interface JSONSummary {
  total: number;
  passed: number;
  warnings: number;
  failed: number;
  info: number;
  score: number;
  total_duration_ms: number;
}

export interface JSONMaturity {
  level: string;
  description: string;
  suggestions?: string[];
}

export interface JSONOutput {
  languages: string[];
  results: JSONResult[];
  summary: JSONSummary;
  maturity: JSONMaturity;
  aborted: boolean;
  success: boolean;
}

export interface JobResponse {
  job_id: string;
  status: JobStatus;
  github_url: string;
  submitted_at: string;
  started_at?: string;
  completed_at?: string;
  request: CheckRequest;
  result?: JSONOutput;
  error?: string;
}

export interface HealthResponse {
  status: string;
}
