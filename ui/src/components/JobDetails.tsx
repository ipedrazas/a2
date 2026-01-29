import { useEffect, useRef, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import type { JobResponse } from '../types';
import { api } from '../api/client';
import { JobStatus } from './JobStatus';
import { Results } from './Results';

export function JobDetails() {
  const { jobId } = useParams<{ jobId: string }>();
  const navigate = useNavigate();
  const [job, setJob] = useState<JobResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [polling, setPolling] = useState(false);
  const pollingRef = useRef(true); // Use ref to avoid stale closure

  useEffect(() => {
    if (!jobId) return;

    // Initial load
    loadJob();

    // Poll if job is not complete
    const interval = setInterval(async () => {
      if (!pollingRef.current) {
        clearInterval(interval);
        return;
      }

      setPolling(true);
      try {
        const updated = await api.getJob(jobId);
        setJob(updated);
        if (updated.status === 'completed' || updated.status === 'failed') {
          setPolling(false);
          pollingRef.current = false;
          clearInterval(interval);
        }
      } catch (err) {
        // Don't clear error on poll, only on initial load
      }
    }, 2000);

    return () => {
      pollingRef.current = false;
      clearInterval(interval);
    };
  }, [jobId]);

  async function loadJob() {
    if (!jobId) return;

    try {
      setLoading(true);
      setError(null);
      const data = await api.getJob(jobId);
      setJob(data);

      // Stop polling if already complete
      if (data.status === 'completed' || data.status === 'failed') {
        pollingRef.current = false;
        setPolling(false);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load job');
    } finally {
      setLoading(false);
    }
  }

  if (loading) {
    return (
      <div className="max-w-6xl mx-auto px-4 py-8">
        <div className="text-center py-12">
          <p className="text-gray-500">Loading job details...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-6xl mx-auto px-4 py-8">
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-800">{error}</p>
        </div>
        <button
          onClick={() => navigate('/jobs')}
          className="text-blue-600 hover:underline"
        >
          ← Back to Jobs
        </button>
      </div>
    );
  }

  if (!job) {
    return (
      <div className="max-w-6xl mx-auto px-4 py-8">
        <div className="text-center py-12">
          <p className="text-gray-500">Job not found</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <button
        onClick={() => navigate('/jobs')}
        className="mb-6 text-blue-600 hover:underline flex items-center gap-2"
      >
        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
        </svg>
        Back to Jobs
      </button>

      {polling && (
        <div className="mb-4 p-3 bg-blue-50 border border-blue-200 rounded-lg flex items-center gap-2">
          <div className="animate-spin rounded-full h-4 w-4 border-2 border-blue-600 border-t-transparent" />
          <span className="text-sm text-blue-800">Checking for updates...</span>
        </div>
      )}

      <JobStatus job={job} />

      {(job.status === 'completed' || job.status === 'failed') && job.result && (
        <div className="mt-6">
          <Results result={job.result} />
        </div>
      )}

      {job.error && (
        <div className="mt-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <h3 className="text-sm font-medium text-red-800 mb-2">Error</h3>
          <p className="text-sm text-red-700">{job.error}</p>
        </div>
      )}
    </div>
  );
}
