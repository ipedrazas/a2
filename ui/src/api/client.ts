import type {
  CheckRequest,
  CheckResponse,
  JobResponse,
  HealthResponse,
} from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || '';

export const api = {
  /**
   * Submit a check job
   */
  async submitCheck(request: CheckRequest): Promise<CheckResponse> {
    const response = await fetch(`${API_BASE_URL}/api/check`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to submit check');
    }

    return response.json();
  },

  /**
   * Get job status and results
   */
  async getJob(jobId: string): Promise<JobResponse> {
    const response = await fetch(`${API_BASE_URL}/api/check/${jobId}`);

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get job status');
    }

    return response.json();
  },

  /**
   * Health check
   */
  async health(): Promise<HealthResponse> {
    const response = await fetch(`${API_BASE_URL}/health`);
    return response.json();
  },

  /**
   * Poll job status until completion
   */
  async pollJob(
    jobId: string,
    onUpdate: (job: JobResponse) => void,
    intervalMs = 2000
  ): Promise<JobResponse> {
    return new Promise((resolve, reject) => {
      const poll = async () => {
        try {
          const job = await api.getJob(jobId);
          onUpdate(job);

          if (job.status === 'completed' || job.status === 'failed') {
            resolve(job);
            return;
          }

          setTimeout(poll, intervalMs);
        } catch (error) {
          reject(error);
        }
      };

      poll();
    });
  },
};
